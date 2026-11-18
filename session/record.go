package session

import (
	"reflect"

	"github.com/happyhzr/happyorm/clause"
)

func (s *Session) Insert(values ...any) (int64, error) {
	recordValues := []any{}
	for _, value := range values {
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
	return result.RowsAffected()
}

func (s *Session) Find(values any) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	elemType := destSlice.Type().Elem()
	modelType := elemType
	pointerElements := elemType.Kind() == reflect.Ptr
	if pointerElements {
		modelType = elemType.Elem()
	}
	table := s.Model(reflect.New(modelType).Interface()).RefTable()
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}
	for rows.Next() {
		record := reflect.New(modelType)
		dest := record.Elem()
		values := []any{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		err := rows.Scan(values...)
		if err != nil {
			return err
		}
		if pointerElements {
			dest = record
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}
