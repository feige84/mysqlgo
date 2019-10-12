package mysqlgo

import (
	"database/sql"
	"fmt"
	"math"
	"reflect"
	"runtime/debug"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var MyDb *DbLib

type DbRow map[string]interface{}

type DbLib struct {
	Db    *sql.DB
	Debug bool
}

func NewDbLib(driver, dsn string) (*DbLib, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("Sql open error: %s\n%s", err, debug.Stack())
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(43200 * time.Second)
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("Db ping error: %s\n%s", err, debug.Stack())
	}

	p := new(DbLib)
	p.Db = db
	return p, nil
}

func scanRow(rows *sql.Rows) (DbRow, error) {
	columns, _ := rows.Columns()
	vals := make([]interface{}, len(columns))
	valsPtr := make([]interface{}, len(columns))

	for i := range vals {
		valsPtr[i] = &vals[i]
	}

	err := rows.Scan(valsPtr...)

	if err != nil {
		return nil, fmt.Errorf("rows scan error: %s\n%s", err, debug.Stack())
	}

	r := make(DbRow)

	for i, v := range columns {
		if va, ok := vals[i].([]byte); ok {
			r[v] = string(va)
		} else {
			r[v] = vals[i]
		}
	}

	return r, nil

}

// 获取多行记录 不怎么好用。暂时留在这里。
func (d *DbLib) GetOne(sql string, args ...interface{}) (DbRow, error) {
	rows, err := d.Db.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %s\n%s", err, debug.Stack())
	}
	defer rows.Close()
	rows.Next()
	result, err := scanRow(rows)
	return result, err
}

// 获取多行记录 不怎么好用。暂时留在这里。
func (d *DbLib) GetAll(sql string, args ...interface{}) ([]DbRow, error) {
	rows, err := d.Db.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %s\n%s", err, debug.Stack())
	}
	defer rows.Close()

	result := make([]DbRow, 0)

	for rows.Next() {
		r, err := scanRow(rows)
		if err != nil {
			continue
		}

		result = append(result, r)
	}

	return result, nil

}

// 写入记录
/*
data := DbRow{}
data["creative_id"] = 4444
data["web_url"] = "yyyyy"
result, err := MyDb.Insert("dy_ad", data)
*/
func (d *DbLib) Insert(table string, data DbRow) (int64, error) {
	fields := make([]string, 0)
	vals := make([]interface{}, 0)
	placeHolder := make([]string, 0)

	for f, v := range data {
		fields = append(fields, f)
		vals = append(vals, v)
		placeHolder = append(placeHolder, "?")
	}

	sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ", table, strings.Join(fields, ","), strings.Join(placeHolder, ","))
	if d.Debug {
		fmt.Println(sqlStr, vals)
	}
	result, err := d.Db.Exec(sqlStr, vals...)
	if err != nil {
		return 0, fmt.Errorf("insert error: %s\n%s", err, debug.Stack())
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get insert lastid error: %s\n%s", err, debug.Stack())
	}

	//这里是有些表没有自增主键。获取不到insertLastId。就获取影响行数。
	if lastInsertId == 0 {
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf("get insert rows affected error: %s\n%s", err, debug.Stack())
		}
		if rowsAffected > 0 {
			return rowsAffected, nil
		}
	}

	return lastInsertId, nil
}

/*
//批量插入，未完成。
data := DbRow{}
data["creative_id"] = 4444
data["web_url"] = "yyyyy"
result, err := MyDb.Insert("dy_ad", data)
*/
func (d *DbLib) InsertMulti(table string, data interface{}) (int64, error) {

	sind := reflect.Indirect(reflect.ValueOf(data))

	switch sind.Kind() {
	case reflect.Array, reflect.Slice:
		if sind.Len() == 0 {
			return 0, fmt.Errorf("args error may be empty")
		}
	case reflect.Map:
		if oneData, ok := data.(DbRow); ok {
			return d.Insert(table, oneData)
		}
		return 0, fmt.Errorf("args is not DbRow")
	default:
		return 0, fmt.Errorf("args error may be empty")
	}

	fields := make([]string, 0)
	values := make([]interface{}, 0)
	placeHolder := make([]string, 0)

	i := 0
	for _, d := range data.([]DbRow) {
		for f, v := range d {
			if i == 0 {
				fields = append(fields, f) //第一次取值的时候取字段
			}
			values = append(values, v)
			placeHolder = append(placeHolder, "?")
		}
		fmt.Println("fields:", fields)
		i++
	}

	marks := make([]string, len(fields))
	for i := range marks {
		marks[i] = "?"
	}
	Q := "`"
	sep := fmt.Sprintf("%s, %s", Q, Q)
	qmarks := strings.Join(marks, ", ")
	columns := strings.Join(fields, sep)

	multi := len(values) / len(fields)
	qmarks = strings.Repeat(qmarks+"), (", multi-1) + qmarks

	query := fmt.Sprintf("INSERT INTO %s%s%s (%s%s%s) VALUES (%s)", Q, table, Q, Q, columns, Q, qmarks)
	fmt.Println("query:", query, values)
	//sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ", table, strings.Join(fields, ","), strings.Join(placeHolder, ","))
	/*
		fields := make([]string, 0)
		vals := make([]interface{}, 0)
		placeHolder := make([]string, 0)

		for f, v := range data {
			fields = append(fields, f)
			vals = append(vals, v)
			placeHolder = append(placeHolder, "?")
		}

		sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ", table, strings.Join(fields, ","), strings.Join(placeHolder, ","))
		if d.Debug {
			fmt.Println(sqlStr, vals)
		}
		result, err := d.Db.Exec(sqlStr, vals...)
		if err != nil {
			return 0, fmt.Errorf("insert error: %s\n%s", err, debug.Stack())
		}

		lastInsertId, err := result.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("get insert lastid error: %s\n%s", err, debug.Stack())
		}

		//这里是有些表没有自增主键。获取不到insertLastId。就获取影响行数。
		if lastInsertId == 0 {
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				return 0, fmt.Errorf("get insert rows affected error: %s\n%s", err, debug.Stack())
			}
			if rowsAffected > 0 {
				return rowsAffected, nil
			}
		}

		return lastInsertId, nil
	*/
	return 0, nil
}

// 更新记录
/*
data := DbRow{}
data["creative_id"] = 4444
data["web_url"] = "yyyyy"
result, err := MyDb.Update("dy_ad", "ad_id=?", data, 112)
*/
func (d *DbLib) Update(table, condition string, data DbRow, args ...interface{}) (int64, error) {
	params := make([]string, 0)
	vals := make([]interface{}, 0)

	for f, v := range data {
		params = append(params, f+"=?")
		vals = append(vals, v)
	}

	sqlStr := "UPDATE %s SET %s"
	if condition != "" {
		sqlStr += " WHERE %s"
		sqlStr = fmt.Sprintf(sqlStr, table, strings.Join(params, ","), condition)
		vals = append(vals, args...)
	} else {
		sqlStr = fmt.Sprintf(sqlStr, table, strings.Join(params, ","))
	}

	if d.Debug {
		fmt.Println(sqlStr, vals)
	}
	result, err := d.Db.Exec(sqlStr, vals...)
	if err != nil {
		return 0, fmt.Errorf("update error: %s\n%s", err, debug.Stack())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("get update rows affected error: %s\n%s", err, debug.Stack())
	}

	return rowsAffected, nil
}

// 删除记录
/*
result, err := MyDb.Delete("dy_ad", "ad_id=?", 111)
*/
func (d *DbLib) Delete(table, condition string, args ...interface{}) (int64, error) {
	sqlStr := "DELETE FROM %s "
	if condition != "" {
		sqlStr += "WHERE %s"
		sqlStr = fmt.Sprintf(sqlStr, table, condition)
	} else {
		sqlStr = fmt.Sprintf(sqlStr, table)
	}

	if d.Debug {
		fmt.Println(sqlStr, args)
	}
	result, err := d.Db.Exec(sqlStr, args...)
	if err != nil {
		return 0, fmt.Errorf("delete error: %s\n%s", err, debug.Stack())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("get delete rows affected error: %s\n%s", err, debug.Stack())
	}
	return rowsAffected, nil
}

//执行记录
func (d *DbLib) Execute(sqlStr string, args ...interface{}) (int64, error) {
	if d.Debug {
		fmt.Println(sqlStr, args)
	}
	res, err := d.Db.Exec(sqlStr, args...)
	if err != nil {
		return 0, fmt.Errorf("execute error: %s\n%s", err, debug.Stack())
	}

	var rowsAffected, lastInsertId int64
	if rowsAffected, err = res.RowsAffected(); err != nil {
		rowsAffected = 0
	}
	if lastInsertId, err = res.LastInsertId(); err != nil {
		lastInsertId = 0
	}
	return int64(math.Max(float64(rowsAffected), float64(lastInsertId))), nil
}

func Struct2Map(structData interface{}) map[string]interface{} {
	if structData != nil {
		result := make(map[string]interface{})
		object := reflect.ValueOf(structData)
		ref := object.Elem()
		typeOfType := ref.Type()
		for i := 0; i < ref.NumField(); i++ {
			field := ref.Field(i)
			result[typeOfType.Field(i).Name] = field.Interface()
			//fmt.Printf("%d. %s %s = %v \n", i, typeOfType.Field(i).Name, field.Type(), field.Interface())
		}
		//fmt.Println(result)
		return result
	}
	return nil
}
