package accdebit

import (
	dboperations "account-maintaince/db-operations"
	"account-maintaince/generic-operations"
	"errors"
	"time"
)

type Debit struct {
	Activity string
	Payload struct{
		UserId string
		Amount int
	}
}
func (debit *Debit)FillData(data []byte,debAccs interface{})error{
	err := generic.UnMarshalData(data, debAccs)
	return err
}
func (debit *Debit)AddTransitionRecord()error{
	dbConnection := dboperations.GetDBref()
	err := dbConnection.CustomQuery("insert into '"+debit.Payload.UserId+"_transition'(date,amount,transition_type,priority,reason,timestamp) values(?,?,?,?,?,?)",time.Now().String(),debit.Payload.Amount,debit.Activity,0,debit.Activity,0)
	return err
}
func (debit *Debit)UpdateBal()error{
	dbConnection := dboperations.GetDBref()
	rows,err := dbConnection.GetuserData("select bal,id from usersaccount where id = '"+debit.Payload.UserId+"'")
	avabal := 0
	id := ""
	for rows.Next(){
		rows.Scan(&avabal,&id)
	}
	rows.Close()
	if id !=  debit.Payload.UserId || id == "" {
		return errors.New("Account not found.")
	}
	if avabal < debit.Payload.Amount {
		return errors.New("Insufficient Funds")
	}
	err = dbConnection.CustomQuery("update usersaccount set bal = ? where id = ?",(avabal-debit.Payload.Amount),debit.Payload.UserId)
	return err
}
func (debit *Debit)UpdatePriorityTable()error{
	dbConnection := dboperations.GetDBref()
	rows,err := dbConnection.GetuserData("select priority,bal from '"+debit.Payload.UserId+"_priority' order by priority desc")
	priorityMap := make(map[int]int)
	var priority,bal int
	for rows.Next(){
		rows.Scan(&priority,&bal)
		priorityMap[priority] = bal
	}
	rows.Close()
	debitAmount := debit.Payload.Amount
	for key,value := range priorityMap{
		if debitAmount > value {
			if value <=0 {
				continue
			}
			err = dbConnection.CustomQuery("update '"+debit.Payload.UserId+ "_priority' set bal = ? where priority = ?",0,key)
			debitAmount = debitAmount - value
		}else {
			err = dbConnection.CustomQuery("update '"+debit.Payload.UserId+ "_priority' set bal = ? where priority = ?",value-debit.Payload.Amount,key)
			break
		}

	}
	return err
}
