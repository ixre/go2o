package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"ops/cf/log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var _ Orm = new(DBMap)

type TableMapMeta struct {
	tableName     string
	pkFieldName   string
	pkIsAuto      bool
	FieldNames    []string //预留，可能会用到
	FieldMapNames []string
}

//it's a IOrm Implements for mysql
type DBMap struct {
	tableMap map[string]*TableMapMeta
	*sql.DB
	useTrace bool
}

func NewDBMap(db *sql.DB) *DBMap {
	return &DBMap{
		DB:       db,
		tableMap: make(map[string]*TableMapMeta),
	}
}

func (this *DBMap) err(err error) {
	if this.useTrace {
		log.Println("[Error]:%s ", err.Error())
	}
}

func (this *DBMap) getTableMapMeta(t reflect.Type) *TableMapMeta {
	m, exists := this.tableMap[t.String()]
	if exists {
		return m
	}

	names, maps := this.getFields(t)
	pkName, pkIsAuto := this.getPKName(t)
	m = &TableMapMeta{
		tableName:     t.Name(),
		pkFieldName:   pkName,
		pkIsAuto:      pkIsAuto,
		FieldNames:    names,
		FieldMapNames: maps,
	}

	this.tableMap[t.String()] = m

	if this.useTrace {
		log.Println("[DbMap Meta]:", m)
	}

	return m
}

func (this *DBMap) getFields(t reflect.Type) (names []string, mapNames []string) {
	names = []string{}
	mapNames = []string{}

	fnum := t.NumField()
	var fmn string

	for i := 0; i < fnum; i++ {
		f := t.Field(i)
		if f.Tag != "" {
			fmn = f.Tag.Get("db")
			if fmn == "_" || len(fmn) == 0 {
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

func (this *DBMap) getTName(t reflect.Type) string {
	//todo: 用int做键
	v, exists := this.tableMap[t.String()]
	if exists {
		return v.tableName
	}
	return t.Name()
}

//if not defined primary key.the first key will as primary key
func (this *DBMap) getPKName(t reflect.Type) (pkName string, pkIsAuto bool) {
	v, exists := this.tableMap[t.String()]
	if exists {
		return v.pkFieldName, v.pkIsAuto
	}

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

func (this *DBMap) SetTrace(b bool) {
	this.useTrace = b
}

//create a fixed table map
func (this *DBMap) CreateTableMap(v interface{}, tableName string) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	meta := this.getTableMapMeta(t)
	meta.tableName = tableName
	this.tableMap[t.String()] = meta
}

func (this *DBMap) Get(entity interface{}, primaryVal interface{}) error {
	var sql string
	var fieldLen int
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	val := reflect.ValueOf(entity)
	if val.Kind() != reflect.Ptr {
		return errors.New("unaddressable of entity ,it must be a ptr")
	}
	val = val.Elem()

	/* build sql */
	meta := this.getTableMapMeta(t)
	fieldLen = len(meta.FieldNames)
	fieldArr := make([]string, fieldLen)
	var scanVals []interface{} = make([]interface{}, fieldLen)
	var rawBytes [][]byte = make([][]byte, fieldLen)

	for i, v := range meta.FieldMapNames {
		fieldArr[i] = v
		scanVals[i] = &rawBytes[i]
	}

	sql = fmt.Sprintf("SELECT %s FROM %s WHERE %s=?",
		strings.Join(fieldArr, ","),
		meta.tableName,
		meta.pkFieldName,
	)

	if this.useTrace {
		log.Println(fmt.Sprintf("[SQL]:%s , [Params]:%s", sql, primaryVal))
	}

	/* query */
	stmt, err := this.DB.Prepare(sql)
	if err != nil {
		err = errors.New(err.Error() + "\n[SQL]:" + sql)
		this.err(err)
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRow(primaryVal)
	err = row.Scan(scanVals...)
	if err != nil {
		this.err(err)
		return err
	}
	for i := 0; i < fieldLen; i++ {
		field := val.Field(i)
		setField(field, rawBytes[i])
	}
	return nil
}

func (this *DBMap) GetBy(entity interface{}, where string) error {
	var sql string
	var fieldLen int
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	val := reflect.ValueOf(entity)
	if val.Kind() != reflect.Ptr {
		return errors.New("unaddressable of entity ,it must be a ptr")
	}

	if strings.Trim(where, "") == "" {
		return errors.New("param where can't be empty ")
	}

	val = val.Elem()

	if !val.IsValid() {
		return errors.New("not validate")
	}

	/* build sql */
	meta := this.getTableMapMeta(t)
	fieldLen = len(meta.FieldNames)
	fieldArr := make([]string, fieldLen)
	var scanVals []interface{} = make([]interface{}, fieldLen)
	var rawBytes [][]byte = make([][]byte, fieldLen)

	for i, v := range meta.FieldMapNames {
		fieldArr[i] = v
		scanVals[i] = &rawBytes[i]
	}

	sql = fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(fieldArr, ","),
		meta.tableName,
		where,
	)

	if this.useTrace {
		log.Println(fmt.Sprintf("[SQL]:%s , [Params]:%s", sql, where))
	}

	/* query */
	stmt, err := this.DB.Prepare(sql)
	if err != nil {
		if this.useTrace {
			log.Println("[SQL ERROR]:", err.Error(), " [SQL]:", sql)
		}
		return errors.New(err.Error() + "\n[SQL]:" + sql)
	}
	defer stmt.Close()

	row := stmt.QueryRow()
	err = row.Scan(scanVals...)

	if err != nil {
		return err
	}

	for i := 0; i < fieldLen; i++ {
		field := val.Field(i)
		setField(field, rawBytes[i])
	}
	return nil
}

func (this *DBMap) GetByQuery(entity interface{}, sql string) error {
	var fieldLen int
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	val := reflect.ValueOf(entity)
	if val.Kind() != reflect.Ptr {
		return errors.New("unaddressable of entity ,it must be a ptr")
	}

	val = val.Elem()

	/* build sql */
	meta := this.getTableMapMeta(t)
	fieldLen = len(meta.FieldNames)
	fieldArr := make([]string, fieldLen)
	var scanVals []interface{} = make([]interface{}, fieldLen)
	var rawBytes [][]byte = make([][]byte, fieldLen)

	for i, v := range meta.FieldMapNames {
		fieldArr[i] = meta.tableName + "." + v
		scanVals[i] = &rawBytes[i]
	}

	if strings.Index(sql, "*") != -1 {
		sql = strings.Replace(sql, "*", strings.Join(fieldArr, ","), 1)
	}

	if this.useTrace {
		log.Println(fmt.Sprintf("[SQL]:%s", sql))
	}

	/* query */
	stmt, err := this.DB.Prepare(sql)
	if err != nil {
		if this.useTrace {
			log.Println("[SQL ERROR]:", err.Error(), " [SQL]:", sql)
		}
		return errors.New(err.Error() + "\n[SQL]:" + sql)
	}
	defer stmt.Close()

	row := stmt.QueryRow()
	err = row.Scan(scanVals...)

	if err != nil {
		return err
	}

	for i := 0; i < fieldLen; i++ {
		field := val.Field(i)
		setField(field, rawBytes[i])
	}
	return nil
}

//Select more than 1 entity list
//@to : refrence to queryed entity list
//@entity : query condition
//@where : other condition
func (this *DBMap) Select(to interface{}, entity interface{}, where string) error {
	var sql string
	var condition string = where
	var fieldLen int

	if reflect.ValueOf(to).Kind() != reflect.Ptr {
		return errors.New("unaddressable of to ,it must be a ptr")
	}

	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	val := reflect.Indirect(reflect.ValueOf(entity))

	/* build sql */
	meta := this.getTableMapMeta(t)
	fieldLen = len(meta.FieldNames)
	params, pFieldArr := itrField(meta, &val, true)

	if len(pFieldArr) != 0 {
		if strings.Trim(condition, " ") != "" {
			condition = condition + " AND "
		}
		condition = condition + strings.Join(pFieldArr, "=? AND ") + "=?"
	}

	/*
		if isSet {
			if condition == "" {
				condition = fmt.Sprintf("%s=?", k)
			} else {
				condition = fmt.Sprintf("%s AND %s=?", condition, k)
			}

		}

		if where != "" {
			if condition != "" {
				condition = condition + " AND " + where
			} else {
				condition = where
			}
		}*/

	fieldArr := make([]string, fieldLen)
	var scanVals []interface{} = make([]interface{}, fieldLen)
	var rawBytes [][]byte = make([][]byte, fieldLen)

	for i, v := range meta.FieldMapNames {
		fieldArr[i] = v
		scanVals[i] = &rawBytes[i]
	}

	sql = fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(fieldArr, ","),
		meta.tableName,
		condition,
	)

	if this.useTrace {
		log.Println(fmt.Sprintf("[SQL]:%s , [Params]:%s", sql, params))
	}

	/* query */
	stmt, err := this.DB.Prepare(sql)
	if err != nil {
		if this.useTrace {
			log.Println("[SQL ERROR]:", err.Error(), " [SQL]:", sql)
		}
		return errors.New(err.Error() + "\n[SQL]:" + sql)
	}
	defer stmt.Close()
	rows, err := stmt.Query(params...)
	defer rows.Close()

	if err != nil {
		if this.useTrace {
			log.Println("[SQL ERROR]:", err.Error(), " [SQL]:", sql)
		}
		return errors.New(err.Error() + "\n[SQL]:" + sql)
	}

	/* 用反射来对输出结果复制 */
	toArr := reflect.ValueOf(to).Elem()

	for rows.Next() {
		e := reflect.Indirect(reflect.New(t))

		rows.Scan(scanVals...)
		for i := 0; i < fieldLen; i++ {
			field := e.Field(i)
			setField(field, rawBytes[i])
		}
		toArr = reflect.Append(toArr, e)
	}

	reflect.ValueOf(to).Elem().Set(toArr)
	return nil
}

func (this *DBMap) SelectByQuery(to interface{}, entity interface{}, sql string) error {
	var fieldLen int

	if reflect.ValueOf(to).Kind() != reflect.Ptr {
		return errors.New("unaddressable of to ,it must be a ptr")
	}

	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	/* build sql */
	meta := this.getTableMapMeta(t)
	fieldLen = len(meta.FieldNames)
	params := []interface{}{}

	fieldArr := make([]string, fieldLen)
	var scanVals []interface{} = make([]interface{}, fieldLen)
	var rawBytes [][]byte = make([][]byte, fieldLen)

	for i, v := range meta.FieldMapNames {
		fieldArr[i] = meta.tableName + "." + v
		scanVals[i] = &rawBytes[i]
	}

	if strings.Index(sql, "*") != -1 {
		sql = strings.Replace(sql, "*", strings.Join(fieldArr, ","), 1)
	}

	if this.useTrace {
		log.Println(fmt.Sprintf("[SQL]:%s", sql))
	}

	/* query */
	stmt, err := this.DB.Prepare(sql)
	if err != nil {
		err = errors.New(err.Error() + "\n[SQL]:" + sql)
		this.err(err)
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(params...)
	defer rows.Close()

	if err != nil {
		err = errors.New(err.Error() + "\n[SQL]:" + sql)
		this.err(err)
		return err
	}

	/* 用反射来对输出结果复制 */

	//todo:如果已经有指定数量的，则用索引赋值
	toArr := reflect.ValueOf(to).Elem()
	// _cap := cap(toArr)
	// var i *int = 0
	for rows.Next() {
		e := reflect.Indirect(reflect.New(t))

		rows.Scan(scanVals...)
		for i := 0; i < fieldLen; i++ {
			field := e.Field(i)
			setField(field, rawBytes[i])
		}

		//		*i = *i + 1
		//		if *i<_cap {
		//			continue
		//		}
		toArr = reflect.Append(toArr, e)

	}

	reflect.ValueOf(to).Elem().Set(toArr)
	return nil
}

func (this *DBMap) Delete(entity interface{}, where string) (effect int64, err error) {
	var sql string
	var condition string = where

	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	val := reflect.Indirect(reflect.ValueOf(entity))

	/* build sql */
	meta := this.getTableMapMeta(t)
	params, fieldArr := itrField(meta, &val, true)

	if len(fieldArr) != 0 {
		if strings.Trim(condition, " ") != "" {
			condition = condition + " AND "
		}
		condition = condition + strings.Join(fieldArr, "=? AND ") + "=?"
	}

	if condition == "" {
		return 0, errors.New("unknown condition")
	}

	sql = fmt.Sprintf("DELETE FROM %s WHERE %s",
		meta.tableName,
		condition,
	)

	if this.useTrace {
		log.Println(fmt.Sprintf("[SQL]:%s , [Params]:%s", sql, params))
	}

	/* query */
	stmt, err := this.DB.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		if this.useTrace {
			log.Println("[SQL ERROR]:", err.Error(), " [SQL]:", sql)
		}
		return 0, errors.New(err.Error() + "\n[SQL]" + sql)
	}

	result, err := stmt.Exec(params...)
	var rowNum int64 = 0
	if err == nil {
		rowNum, err = result.RowsAffected()
	}
	if err != nil {
		return rowNum, errors.New(err.Error() + "\n[SQL]" + sql)
	}
	return rowNum, nil
}

func (this *DBMap) DeleteByPk(entity interface{}, primary interface{}) (err error) {
	var sql string
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	/* build sql */
	meta := this.getTableMapMeta(t)

	sql = fmt.Sprintf("DELETE FROM %s WHERE %s=?",
		meta.tableName,
		meta.pkFieldName,
	)

	if this.useTrace {
		log.Println(fmt.Sprintf("[SQL]:%s , [Params]:%s", sql, primary))
	}

	/* query */
	stmt, err := this.DB.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		if this.useTrace {
			log.Println("[SQL ERROR]:", err.Error(), " [SQL]:", sql)
		}
		return errors.New(err.Error() + "\n[SQL]" + sql)
	}

	_, err = stmt.Exec(primary)
	if err != nil {
		return errors.New(err.Error() + "\n[SQL]" + sql)
	}
	return nil
}

func (this *DBMap) Save(primaryKey interface{}, entity interface{}) (rows int64, lastInsertId int64, err error) {
	var sql string
	//var condition string
	//var fieldLen int

	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	val := reflect.Indirect(reflect.ValueOf(entity))

	/* build sql */
	meta := this.getTableMapMeta(t)
	//fieldLen = len(meta.FieldNames)
	params, fieldArr := itrFieldForSave(meta, &val, false)

	//insert
	if primaryKey == nil {
		var pArr = make([]string, len(fieldArr))
		for i, _ := range pArr {
			pArr[i] = "?"
		}

		sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", meta.tableName,
			strings.Join(fieldArr, ","),
			strings.Join(pArr, ","),
		)

		if this.useTrace {
			log.Println(fmt.Sprintf("[SQL]:%s , [Params]:%s", sql, params))
		}

		/* query */
		stmt, err := this.DB.Prepare(sql)
		defer stmt.Close()
		if err != nil {
			if this.useTrace {
				log.Println("[SQL ERROR]:", err.Error(), " [SQL]:", sql)
			}
			return 0, 0, errors.New(err.Error() + "\n[SQL]" + sql)
		}

		result, err := stmt.Exec(params...)
		var rowNum int64 = 0
		var lastInsertId int64 = 0
		if err == nil {
			rowNum, err = result.RowsAffected()
			lastInsertId, _ = result.LastInsertId()
			return rowNum, lastInsertId, err
		}
		return rowNum, lastInsertId, errors.New(err.Error() + "\n[SQL]" + sql)
	} else {
		//update model

		var setCond string

		for i, k := range fieldArr {
			if i == 0 {
				setCond = fmt.Sprintf("%s = ?", k)
			} else {
				setCond = fmt.Sprintf("%s,%s = ?", setCond, k)
			}
		}

		sql = fmt.Sprintf("UPDATE %s SET %s WHERE %s=?", meta.tableName,
			setCond,
			meta.pkFieldName,
		)

		/* query */
		stmt, err := this.DB.Prepare(sql)
		defer stmt.Close()
		if err != nil {
			if this.useTrace {
				log.Println("[SQL ERROR]:", err.Error(), " [SQL]:", sql)
			}
			return 0, 0, errors.New(err.Error() + "\n[SQL]" + sql)
		}

		params = append(params, primaryKey)

		if this.useTrace {
			log.Println(fmt.Sprintf("[SQL]:%s , [Params]:%s", sql, params))
		}

		result, err := stmt.Exec(params...)
		var rowNum int64 = 0
		if err == nil {
			rowNum, err = result.RowsAffected()
			return rowNum, 0, err
		}
		return rowNum, 0, errors.New(err.Error() + "\n[SQL]" + sql)
	}

	//	for i, v := range fieldArr {
	//		if i == 0 {
	//			condition = fmt.Sprintf("%s=?", v)
	//		}else {
	//			condition = fmt.Sprintf(",%s=?", v)
	//		}
	//	}
}

func setField(field reflect.Value, d []byte) {
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
func itrFieldForSave(meta *TableMapMeta, val *reflect.Value, includePk bool) (params []interface{}, fieldArr []string) {
	var isSet bool
	for i, k := range meta.FieldMapNames {

		if !includePk && meta.pkIsAuto && meta.FieldMapNames[i] == meta.pkFieldName {
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

			//		case reflect.Bool:
			//			strVal := field.String()
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

func itrField(meta *TableMapMeta, val *reflect.Value, includePk bool) (params []interface{}, fieldArr []string) {
	var isSet bool
	for i, k := range meta.FieldMapNames {

		if !includePk && meta.pkIsAuto && meta.FieldMapNames[i] == meta.pkFieldName {
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
