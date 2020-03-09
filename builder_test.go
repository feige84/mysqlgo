package mysqlgo

import (
	"fmt"
	"testing"
)

func TestBuilder(t *testing.T) {
	s := &SelectSql{}
	s.Table("v_user").Select("*")
	s.Where("user_id=?", 1)
	s.Where("user_nickname LIKE ?", "%哈哈%").Order("user_id DESC, user_reg_time DESC").Limit(10).Offset(10)

	str2, args2 := s.BuildSql()
	fmt.Println(str2)
	fmt.Println(args2)
	str, args := s.Count("user_id").BuildSql()
	fmt.Println(str)
	fmt.Println(args)
}
