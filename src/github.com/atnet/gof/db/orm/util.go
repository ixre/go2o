/**
 * Copyright 2013 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-10-22 21:43
 * description :
 * history :
 */

package orm

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

// 获取表元数据
func GetTableMapMeta(t reflect.Type) *TableMapMeta {
	names, maps := GetFields(t)
	pkName, pkIsAuto := GetPKName(t)
	m := &TableMapMeta{
		TableName:     t.Name(),
		PkFieldName:   pkName,
		PkIsAuto:      pkIsAuto,
		FieldNames:    names,
		FieldMapNames: maps,
	}
	return m
}

//if not defined primary key.the first key will as primary key
func GetPKName(t reflect.Type) (pkName string, pkIsAuto bool) {
	var ti int = t.NumField()

	ffc := func(f reflect.StructField) (string, bool) {
		if f.Tag != "" {
			var iauto bool
			var fname string

			if ia := f.Tag.Get("auto"); ia == "yes" || ia == "1" {
				iauto = true
			}

			if fname = f.Tag.Get("db"); fname != "" {
				return fname, iauto
			}
			return f.Name, iauto
		}
		return f.Name, false
	}

	for i := 0; i < ti; i++ {
		f := t.Field(i)
		if f.Tag != "" {
			pk := f.Tag.Get("pk")
			if pk == "1" || pk == "yes" {
				return ffc(f)
			}
		}
	}

	return ffc(t.Field(0))
}

// 获取实体的字段
func GetFields(t reflect.Type) (names []string, mapNames []string) {
	names = []string{}
	mapNames = []string{}

	fnum := t.NumField()
	var fmn string

	for i := 0; i < fnum; i++ {
		f := t.Field(i)
		if f.Tag != "" {
			fmn = f.Tag.Get("db")
			if fmn == "-" || fmn == "_" || len(fmn) == 0 {
				break
			}
		}
		if fmn == "" {
			fmn = f.Name
		}
		mapNames = append(mapNames, fmn)
		names = append(names, f.Name)
		fmn = ""
	}

	return names, mapNames
}

func SetField(field reflect.Value, d []byte) {
	if field.IsValid() {
		//fmt.Println(field.String(), "==>", field.Type().Kind())
		switch field.Type().Kind() {
		case reflect.String:
			field.Set(reflect.ValueOf(string(d)))
			return

		case reflect.Int:
			val, err := strconv.ParseInt(string(d), 10, 0)
			if err == nil {
				field.Set(reflect.ValueOf(int(val)))
			}
		case reflect.Int32:
			val, err := strconv.ParseInt(string(d), 10, 32)
			if err == nil {
				field.Set(reflect.ValueOf(val))
			}
		case reflect.Int64:
			val, err := strconv.ParseInt(string(d), 10, 64)
			if err == nil {
				field.Set(reflect.ValueOf(val))
			}

		case reflect.Float32:
			val, err := strconv.ParseFloat(string(d), 32)
			if err == nil {
				field.Set(reflect.ValueOf(float32(val)))
			}

		case reflect.Float64:
			val, err := strconv.ParseFloat(string(d), 64)
			if err == nil {
				field.Set(reflect.ValueOf(val))
			}

		case reflect.Bool:
			strVal := string(d)
			val := strings.ToLower(strVal) == "true" || strVal == "1"
			field.Set(reflect.ValueOf(val))
			return

			//接口类型
		case reflect.Struct:
			//fmt.Println(reflect.TypeOf(time.Now()), field.Type())
			if reflect.TypeOf(time.Now()) == field.Type() {
				t, err := time.Parse("2006-01-02 15:04:05", string(d))
				if err == nil {
					field.Set(reflect.ValueOf(t.Local()))
				}
			}
			return
		}

	}
}

//遍历所有列，并得到参数及列名
func ItrFieldForSave(meta *TableMapMeta, val *reflect.Value, includePk bool) (params []interface{}, fieldArr []string) {
	var isSet bool
	for i, k := range meta.FieldMapNames {

		if !includePk && meta.PkIsAuto &&
			meta.FieldMapNames[i] == meta.PkFieldName {
			continue
		}

		field := val.Field(i)
		isSet = false

		switch field.Type().Kind() {
		case reflect.String:
			if field.String() != "" {
				isSet = true
				if val.Kind() == reflect.Ptr {
					params = append(params, field.String())
				} else {
					params = append(params, field.String())
				}
			}
			break
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			//if field.Int() != 0 {
			isSet = true
			params = append(params, field.Int())
			//}
		case reflect.Float32, reflect.Float64:
			//if v := field.Float(); v != 0 {
			isSet = true
			params = append(params, field.Float())
			//}

		case reflect.Bool:
			strVal := field.String()
			val := strings.ToLower(strVal) == "true" || strVal == "1"
			field.Set(reflect.ValueOf(val))
			break

		case reflect.Struct:
			v := field.Interface()
			switch v.(type) {
			case time.Time:
				if v.(time.Time).Year() > 1 {
					isSet = true
					params = append(params, v.(time.Time))
				}
			}
		}

		if isSet {
			fieldArr = append(fieldArr, k)
		}
	}
	return params, fieldArr
}

func ItrField(meta *TableMapMeta, val *reflect.Value, includePk bool) (params []interface{}, fieldArr []string) {
	var isSet bool
	for i, k := range meta.FieldMapNames {

		if !includePk && meta.PkIsAuto &&
			meta.FieldMapNames[i] == meta.PkFieldName {
			continue
		}

		field := val.Field(i)
		isSet = false

		switch field.Type().Kind() {
		case reflect.String:
			if field.String() != "" {
				isSet = true
				if val.Kind() == reflect.Ptr {
					params = append(params, field.String())
				} else {
					params = append(params, field.String())
				}
			}
			break
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if field.Int() != 0 {
				isSet = true
				params = append(params, field.Int())
			}
		case reflect.Float32, reflect.Float64:
			if v := field.Float(); v != 0 {
				isSet = true
				params = append(params, field.Float())
			}

			//		case reflect.Bool:
			//			val := strings.ToLower(strVal) == "true" || strVal == "1"
			//			field.Set(reflect.ValueOf(val))
			//			break

		case reflect.Struct:
			v := field.Interface()
			switch v.(type) {
			case time.Time:
				if v.(time.Time).Year() > 1 {
					isSet = true
					params = append(params, v.(time.Time))
				}
			}
		}

		if isSet {
			fieldArr = append(fieldArr, k)
		}
	}
	return params, fieldArr
}
