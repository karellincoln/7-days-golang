package session

import (
	"errors"
	"github.com/karellincoln/7-day-golang/gee-orm/clause"
	"reflect"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	// recordValues 其实是[][]interface{}.interface{}的灵活，也带来了一些理解的成本
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		s.CallMethod(BeforeInsert, value)
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	s.CallMethod(AfterInsert, nil)
	return result.RowsAffected()
}

func (s *Session) Find(values interface{}) error {
	s.CallMethod(BeforeQuery, nil)
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			// 获取dest中的成员对象的地址，并作为一个值保存在values中，在下面的scan中用于填充其中的值。
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		s.CallMethod(AfterQuery, dest.Addr().Interface())
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}

// support map[string]interface{}
// also support kv list: "Name", "Tom", "Age", 18, ....
func (s *Session) Update(kv ...interface{}) (int64, error) {
	s.CallMethod(BeforeUpdate, nil)
	m, ok := kv[0].(map[string]interface{})
	if !ok { // 如果不是kv，则将参数构建成kv的键值对
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterUpdate, nil)
	return result.RowsAffected()
}

// Delete records with where clause
func (s *Session) Delete() (int64, error) {
	s.CallMethod(BeforeDelete, nil)
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterDelete, nil)
	return result.RowsAffected()
}

// Count records with where clause
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

// chain 调用
// Limit adds limit condition to clause
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

// Where adds limit condition to clause
func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

// OrderBy adds order by condition to clause
func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}
