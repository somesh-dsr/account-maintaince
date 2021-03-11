package accscheduler

import (
	"fmt"
	"strconv"
	"time"
	"account-maintaince/db-operations"
)

var Notify chan string
type RunningTimers struct {
	Priority int
	Expiry int64
	UserId string
	Reschedule chan struct{}
	Amount int
}
var UserMap map[string]RunningTimers
func CreateNewUserTImer(){
	UserMap = make(map[string]RunningTimers)
	Notify = make(chan string)
	for {
		select {
			case id := <-Notify:
				fmt.Println("Got new acc notification:",id)
				go CreateNewTimer(id)
		}
	}
}
func CreateNewTimer(id string){

	select {
		case <- UserMap[id].Reschedule:
				fmt.Println("Rescheduling")
		case <- time.After(time.Duration(UserMap[id].Expiry)):
				if removeOldEntries(id,UserMap[id].Expiry,UserMap[id].Priority,UserMap[id].Amount) {
					return
				}
	}
	CreateNewTimer(id)
}

func removeOldEntries(id string, expiry int64,priority int,amount int) bool {
	dbConnection := dboperations.GetDBref()
	rows, _ := dbConnection.GetuserData("select amount,priority from '"+id+"_transition' where timestamp <= "+strconv.Itoa(int(time.Now().Unix()))+" and reason='credit'")
	autodebitMap := make(map[int]int)

	var pri,priAmount int
	if rows == nil {
		return notifyWithNewTimer(id,expiry)
	}
	for rows.Next() {
		rows.Scan(&pri,&priAmount)
		if _,status := autodebitMap[pri]; !status {
			autodebitMap[pri] = priAmount
		}else{
			autodebitMap[pri] = autodebitMap[pri]+priAmount
		}
	}
	rows.Close()



	rows, _ = dbConnection.GetuserData("select bal from '"+id+"_priority' where priority = "+strconv.Itoa(priority))
	priorityTable := make(map[int]int)

	for rows.Next() {
		rows.Scan(&pri,&priAmount)
		if priAmount == 0 {
			continue
		}
		priorityTable[pri]=priAmount
	}
	rows.Close()


	sum := 0
	for key,value := range autodebitMap {
		if value >= priorityTable[key] {
			priorityTable[key] = 0
			sum += value
		}else  {
			sum  += priorityTable[key] - value
			priorityTable[key] -= value
		}
	}

	rows, _ = dbConnection.GetuserData("select bal from usersaccount where id='"+id+"'")
	avabal := 0
	for rows.Next(){
		rows.Scan(&avabal)
	}
	rows.Close()
	if avabal <= sum {
		avabal = 0
	}else {
		avabal -= sum
	}
	_ = dbConnection.CustomQuery("update usersaccount set bal = ? where id = ?",(avabal),id)
	_ = dbConnection.CustomQuery("update '"+id+"_transition' set reason = ? where timestamp <= ? and reason='credit'","auto debit",time.Now().Unix())
	//fmt.Println("Updated acc with auto debit transition ",err)
	return notifyWithNewTimer(id,expiry)
}
func notifyWithNewTimer(id string,expiry int64)bool{
	dbConnection := dboperations.GetDBref()
	rows ,_ := dbConnection.GetuserData("select distinct timestamp from '"+id+"_transition' where timestamp >= "+strconv.Itoa(int(expiry))+" and reason='credit' and order by timestamp asc limit 1")
	var timestamp int64
	if rows != nil {
		for rows.Next() {
			rows.Scan(&timestamp)
		}
		rows.Close()
	}
	if timestamp == 0 {
		return true
	}
	userObj := UserMap[id]
	userObj.Expiry = timestamp
	userObj.Reschedule <- struct{}{}
	UserMap[id] = userObj
	return false
}
