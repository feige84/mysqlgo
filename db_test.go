package mysqlgo

import (
	"fmt"
	"testing"
)

func TestExecute(t *testing.T) {
	//fmt.Println("return:", data7, apiErr7)
	//MyDb = GetDB()
	var err error
	MyDb, err = NewDbLib("mysql", fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8", "root", "123456", "tcp", "127.0.0.1", 3306, "dbname"))
	if err != nil {
		panic(err)
	}
	MyDb.Debug = true
	rows, err := MyDb.GetAll("SELECT VERSION() AS version")
	fmt.Println(rows, err)
}
