package main

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type account struct {
	id       int64
	login    string
	password string
	enabled  bool
}

func (i account) Table() string {
	return "accounts"
}

type Product struct {
	ID      string `column:"id"`
	Name    string `column:"name"`
	Price   int64  `column:"price"`
	InStock bool   `column:"in_stock"`
}

func (p Product) Table() string {
	return "products"
}

func insertQuery(in interface{}) (string, error) {
	t := reflect.TypeOf(in)
	v := reflect.ValueOf(in)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() != reflect.Struct {
		return "", errors.New("in must be struct or pointer to struct")
	}

	m := v.MethodByName("Table")
	if !m.IsValid() {
		return "", errors.New(t.Name() + " does not have method Table")
	}

	if !m.CanConvert(reflect.TypeOf(func() string { return "" })) {
		return "", errors.New("invalid method `Table`")
	}
	table := m.Interface().(func() string)()

	numField := v.NumField()

	cols := make([]string, numField)
	vals := make([]string, numField)

	for i := 0; i < numField; i++ {
		if columnTag := t.Field(i).Tag.Get("column"); columnTag != "" {
			cols[i] = columnTag
		} else {
			cols[i] = t.Field(i).Name
		}
		switch v.Field(i).Kind() {
		case reflect.Int64:
			vals[i] = strconv.FormatInt(v.Field(i).Int(), 10)
		case reflect.String:
			vals[i] = "'" + v.Field(i).String() + "'"
		case reflect.Bool:
			vals[i] = strconv.FormatBool(v.Field(i).Bool())
		default:
			return "", errors.New("unsupported type" + v.Field(i).Type().Name())
		}
	}

	query := "INSERT INTO " + table + " (" + strings.Join(cols, ", ") + ") VALUES (" + strings.Join(vals, ", ") + ")"

	return query, nil
}
