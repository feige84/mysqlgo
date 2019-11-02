package mysqlgo

import (
	"fmt"
	"testing"
)

func TestExecute(t *testing.T) {
	//fmt.Println("return:", data7, apiErr7)
	var err error
	MyDb, err = NewDbLib("mysql", fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8", "root", "xxxxxxx", "tcp", "127.0.0.1", 3306, "doudashi"))
	if err != nil {
		panic(err)
	}
	MyDb.Debug = true

	//data := []DbRow{
	//	{
	//		"q_id":       11111111,
	//		"q_dateline": time.Now().Unix(),
	//		"q_status":   1,
	//	},
	//	{
	//		"q_id":       22222222,
	//		"q_dateline": time.Now().Unix(),
	//		"q_status":   1,
	//	},
	//	{
	//		"q_id":       33333333,
	//		"q_dateline": time.Now().Unix(),
	//		"q_status":   1,
	//	},
	//}

	exists := MyDb.Exists("dy_goods", "product_id", 111111111)
	fmt.Println(exists)
}
