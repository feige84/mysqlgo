package mysqlgo

import (
	"fmt"
	"strings"
)

type SelectSql struct {
	andConditions []map[string]interface{}
	group         string
	fields        interface{}
	orders        []string
	offset        interface{}
	limit         interface{}
	tableName     string
	raw           string
}

func (s *SelectSql) Table(name string) *SelectSql {
	s.tableName = name
	return s
}

func (s *SelectSql) Select(query interface{}) *SelectSql {
	s.fields = query
	return s
}

func (s *SelectSql) Limit(limit interface{}) *SelectSql {
	s.limit = limit
	return s
}

func (s *SelectSql) Offset(offset interface{}) *SelectSql {
	s.offset = offset
	return s
}

func (s *SelectSql) Order(value string) *SelectSql {
	if value != "" {
		s.orders = append(s.orders, value)
	}
	return s
}

func (s *SelectSql) Where(query interface{}, values ...interface{}) *SelectSql {
	s.andConditions = append(s.andConditions, map[string]interface{}{"query": query, "args": values})
	return s
}

func (s *SelectSql) Group(query string) *SelectSql {
	s.group = query
	return s
}

func (s *SelectSql) Raw(b string) *SelectSql {
	s.raw = b
	return s
}

func (s *SelectSql) BuildSql() (str string, args []interface{}) {
	if s.raw != "" {
		str = s.raw
		return
	}

	if s.tableName == "" {
		panic(fmt.Errorf("table is null"))
	}

	fields := "*"
	if s.fields != nil {
		switch value := s.fields.(type) {
		case string:
			fields = value
		case []string:
			fields = strings.Join(value, ", ")
		}
	}

	var andConditions []string
	where := ""
	if s.andConditions != nil {
		for _, conds := range s.andConditions {
			//conds = map[string]interface{}{"query": query, "args": values}
			andConditions = append(andConditions, conds["query"].(string))

			for _, v := range conds["args"].([]interface{}) {
				args = append(args, v)
			}
		}
		if len(andConditions) > 0 {
			where = " WHERE " + strings.Join(andConditions, " AND ")
		}
	}

	if s.group != "" {
		where += " GROUP BY " + s.group
	}

	var orders []string
	if s.orders != nil {
		for _, order := range s.orders {
			orders = append(orders, order)
		}
		where += " ORDER BY " + strings.Join(orders, ",")
	}

	if s.limit != nil {
		switch value := s.limit.(type) {
		case int, int64:
			where += fmt.Sprintf(" LIMIT %d ", value)
		}
	}

	if s.offset != nil {
		switch value := s.offset.(type) {
		case int, int64:
			where += fmt.Sprintf(" OFFSET %d ", value)
		}
	}

	query := strings.Replace(fmt.Sprintf("SELECT %v FROM %v %v", fields, s.tableName, where), "  ", " ", -1)

	return query, args
}
