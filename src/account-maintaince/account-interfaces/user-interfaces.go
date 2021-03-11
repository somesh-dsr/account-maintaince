package userinterface

type Credit struct {
	Activity string
	Payload struct{
		UserId string
		Amount int
		Priority int
		Expiry  int
		Type string
	}
}
type Debit struct {
	Activity string
	payload struct{
		UserId string
		Amount int
	}
}
type Account struct {
	Name string
	MailId string
	MobileNo string
	Address struct{
		Pincode string
		Street  string
		City    string
		State   string
	}
}
type Logs struct {
	Date string
	Priority int
	Amount   int
}
type createAccount interface {
	createUserAcount() error
}
type creditAccount interface {
	creditAmount() error

}
type debitAccount interface {
	debitAmount() error
}
type transitionLogs interface {
	getDebitLogs()
	getCreditLogs()
}

