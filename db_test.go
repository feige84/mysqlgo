package mysqlgo

import (
	"fmt"
	"testing"
	"time"
)

func TestExecute(t *testing.T) {
	//fmt.Println("return:", data7, apiErr7)
	var err error
	MyDb, err = NewDbLib("mysql", fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8", "root", "xxxx", "tcp", "127.0.0.1", 3306, "doudashi"))
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

	data := DbRow{
		"q_id":       44444444,
		"q_dateline": time.Now().Unix(),
		"q_status":   1,
	}
	rows, err := MyDb.ReplaceInto("dy_goods_queue", data)
	fmt.Println(rows, err)
}
