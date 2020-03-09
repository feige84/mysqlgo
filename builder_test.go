package mysqlgo

import (
	"fmt"
	"testing"
)

func TestBuilder(t *testing.T) {
	s := &SelectSql{}
	str, args := s.Table("v_user").Select("*").Where("user_nickname LIKE ?", "%哈哈%").Order("user_id DESC, user_reg_time DESC").Limit(10).Offset(10).BuildSql()
	fmt.Println(str)
	fmt.Println(args)
}
