package acccredit

import (
	"account-maintaince/db-operations"
	"errors"
	"strconv"

	//"account-maintaince/account-interfaces"
	//userinterface "account-maintaince/account-interfaces"
	//"fmt"
	"account-maintaince/generic-operations"
	"account-maintaince/scheduler"
	"time"
)

type CreditData struct {
	Activity string
	Payload struct{
		UserId string
		Amount int
		Priority int
		Expiry  int64
		Type string
	}
}
func (creditData *CreditData)FillData(data []byte,creditAccs interface{}) error{
	err := generic.UnMarshalData(data,creditAccs)
	return err
}
func (creditData *CreditData)AddTransitionRecord()error{
	creditData.Payload.Expiry = time.Now().Unix()+creditData.Payload.Expiry
	dbConnection := dboperations.GetDBref()
	err := dbConnection.CustomQuery("insert into '"+creditData.Payload.UserId+"_transition' (date,amount,transition_type,priority,reason,timestamp) values(?,?,?,?,?,?)",time.Now().String(),creditData.Payload.Amount,creditData.Activity,creditData.Payload.Priority,creditData.Activity,creditData.Payload.Expiry)
	return err
}
func (creditData *CreditData)UpdateBal()error{
	dbConnection := dboperations.GetDBref()
	rows,err := dbConnection.GetuserData("select id,bal from usersaccount where id='"+creditData.Payload.UserId+"'")

	avabal := 0
	id := ""
	for rows.Next(){
		rows.Scan(&id,&avabal)
	}
	rows.Close()
	if id !=  creditData.Payload.UserId || id == "" {
		return errors.New("Account not found.")
	}

	err = dbConnection.CustomQuery("update usersaccount set bal = ? where id = ? ",creditData.Payload.Amount+avabal,creditData.Payload.UserId)
	rows,err = dbConnection.GetuserData("select bal from '"+creditData.Payload.UserId+"_priority' where priority = ? "+ strconv.Itoa(creditData.Payload.Priority))
	avabal = -1

	if rows != nil {
		for rows.Next(){
			rows.Scan(&avabal)
		}
		rows.Close()
	}

	if avabal == -1 {
		err = dbConnection.CustomQuery("insert into '"+creditData.Payload.UserId+"_priority' (priority,bal) values(?,?)",creditData.Payload.Priority,creditData.Payload.Amount)
	}else {
		err = dbConnection.CustomQuery("update '"+creditData.Payload.UserId+ "_priority' set bal = ? where priority = ?",(creditData.Payload.Amount+avabal),creditData.Payload.Priority)
	}
	if err == nil {
		if userObj,status := accscheduler.UserMap[creditData.Payload.UserId]; !status {
			accscheduler.UserMap[creditData.Payload.UserId] = accscheduler.RunningTimers{
				UserId: creditData.Payload.UserId,
				Amount: creditData.Payload.Amount,
				Expiry: creditData.Payload.Expiry,
				Priority: creditData.Payload.Priority,
				Reschedule: make(chan struct{}),
			}
			accscheduler.Notify<-creditData.Payload.UserId
		}else{
				if userObj.Expiry > creditData.Payload.Expiry {
					userObj.Expiry = creditData.Payload.Expiry
					userObj.Reschedule  <- struct{}{}
					accscheduler.UserMap[creditData.Payload.UserId] = userObj
				}
		}
	}
	return err
}
