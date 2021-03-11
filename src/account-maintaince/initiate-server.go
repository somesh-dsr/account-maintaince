package main

import (
	"account-maintaince/api-processor"
	"account-maintaince/db-operations"
	accscheduler "account-maintaince/scheduler"
)

func main() {
    go dboperations.InitiateDatabase("/tmp/users.db")
    go accscheduler.CreateNewUserTImer()
    processapis.InitiateServer()
}
