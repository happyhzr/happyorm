package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...any) (string, []any)

var generators map[Type]generator

func init() {
	generators = map[Type]generator{
		INSERT:  _insert,
		VALUES:  _values,
		SELECT:  _select,
		LIMIT:   _limit,
		WHERE:   _where,
		ORDERBY: _orderBy,
	}
}

func genBindVars(num int) string {
	vars := []string{}
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ",")
}

func _insert(values ...any) (string, []any) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []any{}
}

func _values(values ...any) (string, []any) {
	bindStr := ""
	sql := &strings.Builder{}
	vars := []any{}
	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]any)
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

func _select(values ...any) (string, []any) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []any{}
}

func _limit(values ...any) (string, []any) {
	return "LIMIT ?", values
}

func _where(values ...any) (string, []any) {
	desc := values[0]
	vars := values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func _orderBy(values ...any) (string, []any) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []any{}
}
