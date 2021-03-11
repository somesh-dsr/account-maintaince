package nacccreate

import (
	"account-maintaince/generic-operations"
	"account-maintaince/db-operations"
)
type Account struct {
	Name string
	MailId string
	MobileNo string
	Address struct{
		Pincode int
		Street  string
		City    string
		State   string
	}
}
func (acc *Account)FillData(data []byte,accs interface{})error{
	err := generic.UnMarshalData(data,accs)
	return err
}
func (acc *Account)CreateNewAccount()(error,string){
	key,_ := generic.GetKey()
	dbConnection := dboperations.GetDBref()
	err := dbConnection.CustomQuery("insert into usersaccount(id,name,bal,mail,mobile) values(?,?,?,?,?)",key,acc.Name,0,acc.MailId,acc.MobileNo)
	if err != nil {
		return err, ""
	}
	dbConnection.CustomQuery("insert into address(id,street,city,pincode,state,district) values(?,?,?,?,?,?)",key,acc.Name,0,acc.MailId,acc.MobileNo)
	err = dbConnection.CreateTable("create table if not exists '"+key+"_transition'(date varchar(20) not null,amount integer not null,transition_type varchar(10) not null,priority integer,reason varchar(20),timestamp integer);"+"create table if not exists '"+key+"_priority'(priority integer not null,bal integer not null)")
	return err,key
}