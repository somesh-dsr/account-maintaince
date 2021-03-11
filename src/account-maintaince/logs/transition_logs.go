package transitions

import (
	"account-maintaince/db-operations"
)
type Logs struct {
	Date string
	Priority int
	Amount   int
}


func (tLogs *Logs)GetLogs(tType string,logsSlice []Logs,userId string)[]Logs{
	dbConnection := dboperations.GetDBref()
	rows,_ := dbConnection.GetuserData("select date,amount,priority from '"+userId+"_transition' where transition_type = '"+tType+"'"+" or reason = '"+tType+"'")
	var date string
	var amount,priority int
	for rows.Next() {
		rows.Scan(&date,&amount,&priority)
		logsSlice = append(logsSlice,Logs{
			Date: date,
			Amount: amount,
			Priority: priority,
		})
	}
	return logsSlice
}
/*func (tLogs *LogsReq)FillData(data []byte,logReq interface{}){
	generic.UnMarshalData(data,logReq)
}*/