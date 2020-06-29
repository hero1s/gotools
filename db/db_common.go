package db

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/hero1s/gotools/i18n"
)

func NewOrmWithDB(db string) orm.Ormer {
	o := orm.NewOrm()
	o.Using(db)
	return o
}

func NewMasterOrmWithDB(db string) orm.Ormer {
	o := orm.NewOrm()
	o.Using(GetMasterAliasName(db))
	return o
}

func NewOrm(multiOrm []orm.Ormer, db ...string) (o orm.Ormer) {
	if len(multiOrm) == 0 {
		o = orm.NewOrm()
		if len(db) == 1 {
			o.Using(db[0])
		}
	} else if len(multiOrm) == 1 {
		o = multiOrm[0]
	} else {
		panic("只能传一个Ormer")
	}
	return
}

//构造查询
func NewQueryBuilder() (qb orm.QueryBuilder, err error) {
	qb, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		return nil, i18n.WrapDatabaseError(err)
	}

	return qb, nil
}

type Table struct {
	TableName string
	DbName    string
	// TableField map[string]struct{}
}

/*
func filterData(data map[string]interface{}, tableField map[string]struct{}) {
	for k := range data {
		if _, ok := tableField[k]; !ok {
			delete(data, k)
		}
	}
}
*/

func WrapDatabaseError(err error) error {
	if err == orm.ErrNoRows {
		return i18n.RecordNotFound
	}
	return i18n.WrapDatabaseError(err)
}

func NewTableRecord(dbName, tableName string, data map[string]interface{}, multiOrm ...orm.Ormer) (int64, error) {
	// filterData(data, t.TableField)
	o := NewOrm(multiOrm, dbName)
	values, sql := InsertSql(tableName, data)
	result, err := o.Raw(sql, values).Exec()
	if err != nil {
		return 0, i18n.WrapDatabaseError(err)
	}
	n, err := result.LastInsertId()
	return n, i18n.WrapDatabaseError(err)
}

// 一次性插入多条值
func NewMultiTableRecord(dbName, tableName string, data map[string][]interface{}, multiOrm ...orm.Ormer) error {
	// filterData(data, t.TableField)
	values, sql := MultiInsertSql(tableName, data)
	if len(values) == 0 {
		return i18n.WrapDatabaseError(errors.New("insert multi的value数据不齐"))
	}
	o := NewOrm(multiOrm, dbName)
	rp, err := o.Raw(sql).Prepare()
	defer rp.Close()
	if err != nil {
		return i18n.WrapDatabaseError(err)
	}
	for _, value := range values {
		_, err := rp.Exec(value...)
		if err != nil {
			return i18n.WrapDatabaseError(err)
		}
	}
	return rp.Close()
}

func NewOrUpdateRecord(dbName, tableName string, data map[string]interface{}, multiOrm ...orm.Ormer) (int64, error) {
	// filterData(data, t.TableField)
	o := NewOrm(multiOrm, dbName)
	values, sql := InsertOrUpdateSql(tableName, data)
	result, err := o.Raw(sql, values).Exec()
	if err != nil {
		return 0, i18n.WrapDatabaseError(err)
	}
	n, err := result.LastInsertId()
	return n, i18n.WrapDatabaseError(err)
}

func UpdateTableRecord(dbName, tableName string, data map[string]interface{}, conditionSql string, multiOrm ...orm.Ormer) (int64, error) {
	// filterData(data, t.TableField)
	o := NewOrm(multiOrm, dbName)
	values, sql := UpdateSql(tableName, data, conditionSql)
	result, err := o.Raw(sql, values).Exec()
	if err != nil {
		return 0, i18n.WrapDatabaseError(err)
	}
	n, err := result.RowsAffected()
	return n, i18n.WrapDatabaseError(err)
}

func UpdateAndAddTableRecord(dbName, tableName string, data map[string]interface{}, conditionSql string, multiOrm ...orm.Ormer) (int64, error) {
	// filterData(data, t.TableField)
	o := NewOrm(multiOrm, dbName)
	values, sql := UpdateAndAddSql(tableName, data, conditionSql)
	result, err := o.Raw(sql, values).Exec()
	if err != nil {
		return 0, i18n.WrapDatabaseError(err)
	}
	n, err := result.RowsAffected()
	return n, i18n.WrapDatabaseError(err)
}

func DeleteTableRecord(dbName, tableName string, conditionSql string, multiOrm ...orm.Ormer) (int64, error) {
	o := NewOrm(multiOrm, dbName)
	sql := DeleteSql(tableName, conditionSql)
	result, err := o.Raw(sql).Exec()
	if err != nil {
		return 0, i18n.WrapDatabaseError(err)
	}
	n, err := result.RowsAffected()
	return n, i18n.WrapDatabaseError(err)
}

// 获取单行记录
func SingleRecordByAny(dbName, tableName string, conditionSql string, record interface{}, multiOrm ...orm.Ormer) error {
	o := NewOrm(multiOrm, dbName)
	sql := QuerySingleSql(tableName, conditionSql)
	err := o.Raw(sql).QueryRow(record)
	if err == orm.ErrNoRows {
		return i18n.RecordNotFound
	}
	return i18n.WrapDatabaseError(err)
}

//获取多行记录
func MultiRecordByAny(dbName, tableName string, conditionSql string, record interface{}, multiOrm ...orm.Ormer) error {
	o := NewOrm(multiOrm, dbName)
	sql := QueryMultiSql(tableName, conditionSql)
	_, err := o.Raw(sql).QueryRows(record)
	return i18n.WrapDatabaseError(err)
}

//多行记录带分页
func MultiRecordByAnyOrderLimit(dbName, tableName string, conditionSql string, orderbyCondition string, pageIndex, pageSize int64, record interface{}, multiOrm ...orm.Ormer) error {
	o := NewOrm(multiOrm, dbName)
	sql := QueryMultiSqlOderByLimit(tableName, conditionSql, orderbyCondition, pageIndex, pageSize)
	_, err := o.Raw(sql).QueryRows(record)
	return i18n.WrapDatabaseError(err)
}

func ALLRecord(dbName, tableName string, record interface{}, multiOrm ...orm.Ormer) error {
	return MultiRecordByAny(dbName, tableName, "", record, multiOrm...)
}

//返回true就表示不存在
func CheckNoExist(err error) bool {
	if err != nil && err == orm.ErrNoRows {
		return true
	}
	return false
}

func (t *Table) FullTableName() string {
	return t.DbName + "." + t.TableName
}

func (t *Table) NewOrm(multiOrm ...orm.Ormer) orm.Ormer {
	o := NewOrm(multiOrm, t.DbName)
	return o
}

func (t *Table) NewMasterOrm() orm.Ormer {
	o := NewMasterOrmWithDB(t.DbName)
	return o
}

func (t *Table) NewTableRecord(data map[string]interface{}, multiOrm ...orm.Ormer) (int64, error) {
	return NewTableRecord(t.DbName, t.TableName, data, multiOrm...)
}

func (t *Table) NewOrUpdateRecord(data map[string]interface{}, multiOrm ...orm.Ormer) (int64, error) {
	return NewOrUpdateRecord(t.DbName, t.TableName, data, multiOrm...)
}

func (t *Table) UpdateTableRecord(data map[string]interface{}, conditionSql string, multiOrm ...orm.Ormer) (int64, error) {
	return UpdateTableRecord(t.DbName, t.TableName, data, conditionSql, multiOrm...)
}

func (t *Table) UpdateAndAddTableRecord(data map[string]interface{}, conditionSql string, multiOrm ...orm.Ormer) (int64, error) {
	return UpdateAndAddTableRecord(t.DbName, t.TableName, data, conditionSql, multiOrm...)
}

func (t *Table) DeleteTableRecord(conditionSql string, multiOrm ...orm.Ormer) (int64, error) {
	return DeleteTableRecord(t.DbName, t.TableName, conditionSql, multiOrm...)
}

// 获取单行记录
func (t *Table) SingleRecordByAny(conditionSql string, record interface{}, multiOrm ...orm.Ormer) error {
	if len(multiOrm) == 0 {
		return SingleRecordByAny(t.DbName, t.TableName, conditionSql, record, t.NewMasterOrm())
	}
	return SingleRecordByAny(t.DbName, t.TableName, conditionSql, record, multiOrm...)
}

// 获取多行记录
func (t *Table) MultiRecordByAny(conditionSql string, record interface{}, multiOrm ...orm.Ormer) error {
	return MultiRecordByAny(t.DbName, t.TableName, conditionSql, record, multiOrm...)
}

// 带分页的多行记录
func (t *Table) MultiRecordByAnyOrderLimit(conditionSql string, orderbyCondition string, pageIndex, pageSize int64, record interface{}, multiOrm ...orm.Ormer) error {
	return MultiRecordByAnyOrderLimit(t.DbName, t.TableName, conditionSql, orderbyCondition, pageIndex, pageSize, record)
}

func (t *Table) ALLRecord(record interface{}, multiOrm ...orm.Ormer) error {
	return t.MultiRecordByAny("", record, multiOrm...)
}

// 统计条件数量
func (t *Table) CountRecord(conditionSql string, multiOrm ...orm.Ormer) int64 {
	sql := fmt.Sprintf(`SELECT  COUNT(1) FROM %v WHERE %v`, t.TableName, conditionSql)
	var count int64
	err := t.NewOrm().Raw(sql).QueryRow(&count)
	if err != nil {
		return 0
	}
	return count
}
