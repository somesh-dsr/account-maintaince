package dboperations

import (
	"database/sql"
	_ "github.com/mutecomm/go-sqlcipher"
	"sync"
)

var Isfree bool
var dbconnection *Sqlitedb
type Sqlitedb struct {
	DB        *sql.DB
	dbRWMutex sync.RWMutex
}
func InitiateDatabase(dbname string){
	createDatabase(dbname)

	dbconnection.CreateTable("create table if not exists usersaccount(id varchar(100) primary key not null,name varchar(100) not null,bal integer not null,mail varchar(100),mobile varchar(12) unique not null)")
	dbconnection.CreateTable("create table if not exists address(id varchar(100) primary key not null,street varchar(50),city varchar(50) not null,pincode integer not null,state varchar(50) not null,district varchar(50) not null)")
}
func (db *Sqlitedb)pending(query string)error{
	db.dbRWMutex.Lock()
	defer db.dbRWMutex.Unlock()
	_,err := db.DB.Exec(query)
	return err
}
func createDatabase(dbName string)(*Sqlitedb,error){
	dbconnection = new(Sqlitedb)
	var err error
	dbconnection.DB,err = sql.Open("sqlite3", dbName)
	dbconnection.DB.SetMaxOpenConns(1)
	return  dbconnection,err
}
func (db *Sqlitedb)GetuserData(query string)(*sql.Rows,error){
	defer func() {
		 recover()
		//fmt.Println(err)
	}()
	db.dbRWMutex.Lock()
	defer db.dbRWMutex.Unlock()
	rows,err := db.DB.Query(query)
	return rows,err
}
func (db *Sqlitedb)CreateTable(query string)error{
	db.dbRWMutex.Lock()
	defer db.dbRWMutex.Unlock()
	_,err := db.DB.Exec(query)
	return err
}
func (db *Sqlitedb)CustomQuery(query string,data ...interface{})error{
	defer func() {
		 recover()
	}()
	db.dbRWMutex.Lock()
	defer db.dbRWMutex.Unlock()
	//fmt.Println(query,data)
	pRef, err := db.DB.Prepare(query)
	defer pRef.Close()
	_,err = pRef.Exec(data...)
	return err

}
func GetDBref()*Sqlitedb{
	return dbconnection
}
