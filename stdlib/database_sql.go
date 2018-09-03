package stdlib

// Generated by 'goexports database/sql'. Do not edit!

import (
	"database/sql"
	"reflect"
)

func init() {
	Value["database/sql"] = map[string]reflect.Value{
		"Drivers":              reflect.ValueOf(sql.Drivers),
		"ErrConnDone":          reflect.ValueOf(sql.ErrConnDone),
		"ErrNoRows":            reflect.ValueOf(sql.ErrNoRows),
		"ErrTxDone":            reflect.ValueOf(sql.ErrTxDone),
		"LevelDefault":         reflect.ValueOf(sql.LevelDefault),
		"LevelLinearizable":    reflect.ValueOf(sql.LevelLinearizable),
		"LevelReadCommitted":   reflect.ValueOf(sql.LevelReadCommitted),
		"LevelReadUncommitted": reflect.ValueOf(sql.LevelReadUncommitted),
		"LevelRepeatableRead":  reflect.ValueOf(sql.LevelRepeatableRead),
		"LevelSerializable":    reflect.ValueOf(sql.LevelSerializable),
		"LevelSnapshot":        reflect.ValueOf(sql.LevelSnapshot),
		"LevelWriteCommitted":  reflect.ValueOf(sql.LevelWriteCommitted),
		"Named":                reflect.ValueOf(sql.Named),
		"Open":                 reflect.ValueOf(sql.Open),
		"OpenDB":               reflect.ValueOf(sql.OpenDB),
		"Register":             reflect.ValueOf(sql.Register),
	}
	Type["database/sql"] = map[string]reflect.Type{
		"ColumnType":     reflect.TypeOf((*sql.ColumnType)(nil)).Elem(),
		"Conn":           reflect.TypeOf((*sql.Conn)(nil)).Elem(),
		"DB":             reflect.TypeOf((*sql.DB)(nil)).Elem(),
		"DBStats":        reflect.TypeOf((*sql.DBStats)(nil)).Elem(),
		"IsolationLevel": reflect.TypeOf((*sql.IsolationLevel)(nil)).Elem(),
		"NamedArg":       reflect.TypeOf((*sql.NamedArg)(nil)).Elem(),
		"NullBool":       reflect.TypeOf((*sql.NullBool)(nil)).Elem(),
		"NullFloat64":    reflect.TypeOf((*sql.NullFloat64)(nil)).Elem(),
		"NullInt64":      reflect.TypeOf((*sql.NullInt64)(nil)).Elem(),
		"NullString":     reflect.TypeOf((*sql.NullString)(nil)).Elem(),
		"Out":            reflect.TypeOf((*sql.Out)(nil)).Elem(),
		"RawBytes":       reflect.TypeOf((*sql.RawBytes)(nil)).Elem(),
		"Result":         reflect.TypeOf((*sql.Result)(nil)).Elem(),
		"Row":            reflect.TypeOf((*sql.Row)(nil)).Elem(),
		"Rows":           reflect.TypeOf((*sql.Rows)(nil)).Elem(),
		"Scanner":        reflect.TypeOf((*sql.Scanner)(nil)).Elem(),
		"Stmt":           reflect.TypeOf((*sql.Stmt)(nil)).Elem(),
		"Tx":             reflect.TypeOf((*sql.Tx)(nil)).Elem(),
		"TxOptions":      reflect.TypeOf((*sql.TxOptions)(nil)).Elem(),
	}
}