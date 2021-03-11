package processapis

import (
	"github.com/gin-gonic/gin"
	"account-maintaince/account-operations/create"
	//"account-maintaince/account-interfaces"
	"account-maintaince/account-operations/credit"
	"account-maintaince/account-operations/debit"
	"account-maintaince/logs"
	"io/ioutil"
)

func InitiateServer(){
	eng := gin.Default()
	eng.POST("/users/account/credit",creditAmount)
	eng.POST("/users/account/debit",debitAmount)
	eng.POST("/users/create",createNewAcc)
	eng.GET("/users/debit/logs",getDebitLogs)
	eng.GET("/users/credit/logs",getCreditLogs)
	eng.GET("/users/expired/amount",getExpiredLogs)
	eng.Run("127.0.0.1:100")
}

func getExpiredLogs(eLogs *gin.Context) {
	var LogResp []transitions.Logs
	refRes := transitions.Logs{}
	LogResp = refRes.GetLogs("auto debit",LogResp,eLogs.Query("userId"))
	sendLogs(LogResp,eLogs)
}
func creditAmount(credit *gin.Context){
	data,err := ioutil.ReadAll(credit.Request.Body)

	var creditAccs acccredit.CreditData
	//fmt.Println(creditAccs)
	if err = creditAccs.FillData(data,&creditAccs); err != nil {
		sendAPiResponse(err.Error(),409,"error",credit)
		return
	}
	if err = creditAccs.AddTransitionRecord(); err != nil {
		sendAPiResponse(err.Error(),409,"error",credit)
		return
	}
	if err = creditAccs.UpdateBal(); err != nil {
		sendAPiResponse(err.Error(),409,"error",credit)
		return
	}
	sendAPiResponse("Amount Credited successfully",201,"success",credit)
}
func debitAmount(debit *gin.Context){
	data,err := ioutil.ReadAll(debit.Request.Body)

	var debitAccs accdebit.Debit
	if err = debitAccs.FillData(data,&debitAccs);err != nil {
		sendAPiResponse(err.Error(),409,"error",debit)
		return
	}
	if err = debitAccs.UpdateBal(); err != nil {
		sendAPiResponse(err.Error(),409,"error",debit)
		return
	}
	if err = debitAccs.AddTransitionRecord(); err != nil {
		sendAPiResponse(err.Error(),409,"error",debit)
		return
	}
	if err = debitAccs.UpdatePriorityTable(); err != nil {
		sendAPiResponse(err.Error(),409,"error",debit)
		return
	}
	sendAPiResponse("Amount Debited successfully",201,"success",debit)
}
func createNewAcc(naccount *gin.Context){
	data,err := ioutil.ReadAll(naccount.Request.Body)

	var naccounts nacccreate.Account
	if err = naccounts.FillData(data,&naccounts); err != nil {
		sendAPiResponse(err.Error(),409,"error",naccount)
		return
	}
	key := ""
	if err,key = naccounts.CreateNewAccount(); err != nil {
		sendAPiResponse(err.Error(),409,"error",naccount)
		return
	}
	sendAPiResponse("Account Created Successfully: key is "+key,201,"success",naccount)
}

func sendAPiResponse(apiResult string,responseCode int,result string, c *gin.Context) {
	go func() {
		 recover()
	}()
	c.JSON(responseCode,gin.H{
		"statusCode" : responseCode,
		 result: apiResult,
	})
}
func getDebitLogs(dLogs *gin.Context){
	/*data,err := ioutil.ReadAll(dLogs.Request.Body)
	if err != nil {

	}*/
	var LogResp []transitions.Logs
	refRes := transitions.Logs{}
	LogResp = refRes.GetLogs("debit",LogResp,dLogs.Query("userId"))
	sendLogs(LogResp,dLogs)
}
func getCreditLogs(cLogs *gin.Context){
	//data,err := ioutil.ReadAll(cLogs.Request.Body)
	/*if err != nil {

	}*/
	var LogResp []transitions.Logs
	refRes := transitions.Logs{}
	LogResp = refRes.GetLogs("credit",LogResp,cLogs.Query("userId"))
	sendLogs(LogResp,cLogs)
}
func sendLogs(logs []transitions.Logs,c *gin.Context){
	var respHeader []gin.H
    for _,logObj := range  logs{
    	respHeader = append(respHeader,gin.H{
    		"transitionDate" : logObj.Date,
    		"amount" : logObj.Amount,
    		"amountPriority" : logObj.Priority,
		})
	}

    c.JSON(200,respHeader)
}

