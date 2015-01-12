package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/atnet/gof/log"
	"reflect"
	"strings"
)

var _ Orm = new(simpleOrm)

//it's a IOrm Implements for mysql
type simpleOrm struct {
	tableMap map[string]*TableMapMeta
	*sql.DB
	useTrace bool
}

func NewOrm(db *sql.DB) Orm {
	return &simpleOrm{
		DB:       db,
		tableMap: make(map[string]*TableMapMeta),
	}
}

func (this *simpleOrm) err(err error) {
	if this.useTrace {
		log.Println("[ORM][Error]:%s ", err.Error())
	}
}

func (this *simpleOrm) getTableMapMeta(t reflect.Type) *TableMapMeta {
	m, exists := this.tableMap[t.String()]
	if exists {
		return m
	}

	m = GetTableMapMeta(t)
	this.tableMap[t.String()] = m

	if this.useTrace {
		log.Println("[ORM][META]:", m)
	}

	return m
}

func (this *simpleOrm) getTName(t reflect.Type) string {
	//todo: 用int做键
	v, exists := this.tableMap[t.String()]
	if exists {
		return v.TableName
	}
	return t.Name()
}

//if not defined primary key.the first key will as primary key
func (this *simpleOrm) getPKName(t reflect.Type) (pkName string, pkIsAuto bool) {
	v, exists := this.tableMap[t.String()]
	if exists {
		return v.PkFieldName, v.PkIsAuto
	}
	return GetPKName(t)
}

func (this *simpleOrm) SetTrace(b bool) {
	this.useTrace = b
}

//create a fixed table map
func (this *simpleOrm) CreateTableMap(v interface{}, tableName string) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	meta := this.getTableMapMeta(t)
	meta.TableName = tableName
	this.tableMap[t.String()] = meta
}

func (this *simpleOrm) Get(primaryVal interface{}, entity interface{}) error {
	var sql string
	var fieldLen int
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	val := reflect.ValueOf(entity)
	if val.Kind() != reflect.Ptr {
		return errors.New("Unaddressable of entity ,it must be a ptr")
	}
	val = val.Elem()

	/* build sql */
	meta := this.getTableMapMeta(t)
	fieldLen = len(meta.FieldNames)
	fieldArr := make([]string, fieldLen)
	var scanVal []interface{} = make([]interface{}, fieldLen)
	var rawBytes [][]byte = make([][]byte, fieldLen)

	for i, v := range meta.FieldMapNames {
		fieldArr[i] = v
		scanVal[i] = &rawBytes[i]
	}

	sql = fmt.Sprintf("SELECT %s FROM %s WHERE %s=?",
		strings.Join(fieldArr, ","),
		meta.TableName,
		meta.PkFieldName,
	)

	if this.useTrace {
		log.Println(fmt.Sprintf("[ORM][SQL]:%s , [Params]:%s", sql, primaryVal))
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
	err = row.Scan(scanVal...)
	if err != nil {
		this.err(err)
		return err
	}
	for i := 0; i < fieldLen; i++ {
		field := val.Field(i)
		SetField(field, rawBytes[i])
	}
	return nil
}

func (this *simpleOrm) GetBy(entity interface{}, where string,
	args ...interface{}) error {
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
	var scanVal []interface{} = make([]interface{}, fieldLen)
	var rawBytes [][]byte = make([][]byte, fieldLen)

	for i, v := range meta.FieldMapNames {
		fieldArr[i] = v
		scanVal[i] = &rawBytes[i]
	}

	sql = fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(fieldArr, ","),
		meta.TableName,
		where,
	)

	if this.useTrace {
		log.Println(fmt.Sprintf("[ORM][SQL]:%s , [Params]:%s", sql, where))
	}

	/* query */
	stmt, err := this.DB.Prepare(sql)
	if err != nil {
		if this.useTrace {
			log.Println("[ORM][Error]:", err.Error(), " [SQL]:", sql)
		}
		return errors.New(err.Error() + "\n[SQL]:" + sql)
	}
	defer stmt.Close()

	row := stmt.QueryRow(args...)
	err = row.Scan(scanVal...)

	if err != nil {
		return err
	}

	for i := 0; i < fieldLen; i++ {
		field := val.Field(i)
		SetField(field, rawBytes[i])
	}
	return nil
}

func (this *simpleOrm) GetByQuery(entity interface{}, sql string,
	args ...interface{}) error {
	var fieldLen int
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	val := reflect.ValueOf(entity)
	if val.Kind() != reflect.Ptr {
		return errors.New("Unaddressable of entity ,it must be a ptr")
	}

	val = val.Elem()

	/* build sql */
	meta := this.getTableMapMeta(t)
	fieldLen = len(meta.FieldNames)
	fieldArr := make([]string, fieldLen)
	var scanVal []interface{} = make([]interface{}, fieldLen)
	var rawBytes [][]byte = make([][]byte, fieldLen)

	for i, v := range meta.FieldMapNames {
		fieldArr[i] = meta.TableName + "." + v
		scanVal[i] = &rawBytes[i]
	}

	if strings.Index(sql, "*") != -1 {
		sql = strings.Replace(sql, "*", strings.Join(fieldArr, ","), 1)
	}

	if this.useTrace {
		log.Println(fmt.Sprintf("[ORM][SQL]:%s", sql))
	}

	/* query */
	stmt, err := this.DB.Prepare(sql)
	if err != nil {
		if this.useTrace {
			log.Println("[ORM][Error]:", err.Error(), " [SQL]:", sql)
		}
		return errors.New(err.Error() + "\n[SQL]:" + sql)
	}
	defer stmt.Close()

	row := stmt.QueryRow(args...)
	err = row.Scan(scanVal...)

	if err != nil {
		return err
	}

	for i := 0; i < fieldLen; i++ {
		field := val.Field(i)
		SetField(field, rawBytes[i])
	}
	return nil
}

//Select more than 1 entity list
//@to : refrence to queryed entity list
//@entity : query condition
//@where : other condition
func (this *simpleOrm) Select(to interface{}, where string, args ...interface{}) error {
	return this.selectBy(to, where, false, args...)
}

func (this *simpleOrm) SelectByQuery(to interface{}, sql string, args ...interface{}) error {
	return this.selectBy(to, sql, true, args...)
}

// query rows
func (this *simpleOrm) selectBy(to interface{}, sql string, fullSql bool, args ...interface{}) error {
	var fieldLen int
	var eleIsPtr bool // 元素是否为指针

	toVal := reflect.Indirect(reflect.ValueOf(to))
	toTyp := reflect.TypeOf(to).Elem()

	if toTyp.Kind() == reflect.Ptr {
		toTyp = toTyp.Elem()
	}

	if toTyp.Kind() != reflect.Slice {
		return errors.New("to must be slice")
	}

	baseTyp := toTyp.Elem()
	if baseTyp.Kind() == reflect.Ptr {
		eleIsPtr = true
		baseTyp = baseTyp.Elem()
	}

	/* build sql */
	meta := this.getTableMapMeta(baseTyp)
	fieldLen = len(meta.FieldNames)

	fieldArr := make([]string, fieldLen)
	var scanVal []interface{} = make([]interface{}, fieldLen)
	var rawBytes [][]byte = make([][]byte, fieldLen)

	for i, v := range meta.FieldMapNames {
		fieldArr[i] = meta.TableName + "." + v
		scanVal[i] = &rawBytes[i]
	}

	if fullSql {
		if strings.Index(sql, "*") != -1 {
			sql = strings.Replace(sql, "*", strings.Join(fieldArr, ","), 1)
		}
	} else {
		where := sql
		if len(where) == 0 {
			sql = fmt.Sprintf("SELECT %s FROM %s",
				strings.Join(fieldArr, ","),
				meta.TableName)
		} else {
			// 此时,sql为查询条件
			sql = fmt.Sprintf("SELECT %s FROM %s WHERE %s",
				strings.Join(fieldArr, ","),
				meta.TableName,
				where)
		}
	}

	if this.useTrace {
		log.Println(fmt.Sprintf("[ORM][SQL]:%s", sql))
	}

	/* query */
	stmt, err := this.DB.Prepare(sql)
	if err != nil {
		err = errors.New(fmt.Sprintf("%s - [SQL]: %s- [Args]:%+v", err.Error(), sql, args))
		this.err(err)
		return err
	}

	defer stmt.Close()
	rows, err := stmt.Query(args...)

	if err != nil {
		err = errors.New(err.Error() + "\n[SQL]:" + sql)
		this.err(err)
		return err
	}

	defer rows.Close()

	/* 用反射来对输出结果复制 */

	toArr := toVal

	for rows.Next() {
		e := reflect.New(baseTyp)
		v := e.Elem()

		rows.Scan(scanVal...)
		for i := 0; i < fieldLen; i++ {
			SetField(v.Field(i), rawBytes[i])
		}
		if eleIsPtr {
			toArr = reflect.Append(toArr, e)
		} else {
			toArr = reflect.Append(toArr, v)
		}
	}
	toVal.Set(toArr)

	return nil
}

func (this *simpleOrm) Delete(entity interface{}, where string,
	args ...interface{}) (effect int64, err error) {
	var sql string

	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	/* build sql */
	meta := this.getTableMapMeta(t)

	if where == "" {
		return 0, errors.New("Unknown condition")
	}

	sql = fmt.Sprintf("DELETE FROM %s WHERE %s",
		meta.TableName,
		where,
	)

	if this.useTrace {
		log.Println(fmt.Sprintf("[ORM][SQL]:%s , [Params]:%s", sql, args))
	}

	/* query */
	stmt, err := this.DB.Prepare(sql)
	if err != nil {
		if this.useTrace {
			log.Println("[ORM][Error]:", err.Error(), " [SQL]:", sql)
		}
		return 0, errors.New(err.Error() + "\n[SQL]" + sql)
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	var rowNum int64 = 0
	if err == nil {
		rowNum, err = result.RowsAffected()
	}
	if err != nil {
		return rowNum, errors.New(err.Error() + "\n[SQL]" + sql)
	}
	return rowNum, nil
}

func (this *simpleOrm) DeleteByPk(entity interface{}, primary interface{}) (err error) {
	var sql string
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	/* build sql */
	meta := this.getTableMapMeta(t)

	sql = fmt.Sprintf("DELETE FROM %s WHERE %s=?",
		meta.TableName,
		meta.PkFieldName,
	)

	if this.useTrace {
		log.Println(fmt.Sprintf("[ORM][SQL]:%s , [Params]:%s", sql, primary))
	}

	/* query */
	stmt, err := this.DB.Prepare(sql)
	if err != nil {
		if this.useTrace {
			log.Println("[ORM][Error]:", err.Error(), " [SQL]:", sql)
		}
		return errors.New(err.Error() + "\n[SQL]" + sql)
	}
	defer stmt.Close()

	_, err = stmt.Exec(primary)
	if err != nil {
		return errors.New(err.Error() + "\n[SQL]" + sql)
	}
	return nil
}

func (this *simpleOrm) Save(primaryKey interface{}, entity interface{}) (rows int64, lastInsertId int64, err error) {
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
	params, fieldArr := ItrFieldForSave(meta, &val, false)

	//insert
	if primaryKey == nil {
		var pArr = make([]string, len(fieldArr))
		for i, _ := range pArr {
			pArr[i] = "?"
		}

		sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", meta.TableName,
			strings.Join(fieldArr, ","),
			strings.Join(pArr, ","),
		)

		if this.useTrace {
			log.Println(fmt.Sprintf("[ORM][SQL]:%s , [Params]:%s", sql, params))
		}

		/* query */
		stmt, err := this.DB.Prepare(sql)
		if err != nil {
			if this.useTrace {
				log.Println("[ORM][Error]:", err.Error(), " [SQL]:", sql)
			}
			return 0, 0, errors.New(err.Error() + "\n[SQL]" + sql)
		}
		defer stmt.Close()

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

		sql = fmt.Sprintf("UPDATE %s SET %s WHERE %s=?", meta.TableName,
			setCond,
			meta.PkFieldName,
		)

		/* query */
		stmt, err := this.DB.Prepare(sql)
		if err != nil {
			if this.useTrace {
				log.Println("[ORM][Error]:", err.Error(), " [SQL]:", sql)
			}
			return 0, 0, errors.New(err.Error() + "\n[SQL]" + sql)
		}
		defer stmt.Close()

		params = append(params, primaryKey)

		if this.useTrace {
			log.Println(fmt.Sprintf("[ORM][SQL]:%s , [Params]:%s", sql, params))
		}

		result, err := stmt.Exec(params...)
		var rowNum int64 = 0
		if err == nil {
			rowNum, err = result.RowsAffected()
			return rowNum, 0, err
		}
		return rowNum, 0, errors.New(err.Error() + "\n[SQL]" + sql)
	}
}
