package db

import (
	"fmt"
	"time"
)

// 使用者可能要在前面封装一层
func CreateTomorrowTable(dbName, oldTable, newTable string) func() error {
	return func() error {
		CreateTodayTable(dbName, oldTable, newTable)()
		t := time.Now().Add(24 * time.Hour).Format("20060102")
		sql := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %v LIKE %v
		`, newTable+t, oldTable)

		o := NewOrmWithDB(dbName)
		_, err := o.Raw(sql).Exec()
		return err
	}
}

func CreateTodayTable(dbName, oldTable, newTable string) func() error {
	return func() error {
		t := time.Now().Format("20060102")
		sql := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %v LIKE %v
		`, newTable+t, oldTable)
		o := NewOrmWithDB(dbName)
		_, err := o.Raw(sql).Exec()
		return err
	}
}

func CreateYesterdayTable(dbName, oldTable, newTable string) func() error {
	return func() error {
		CreateTwoDayBeforeTable(dbName, oldTable, newTable)()
		t := time.Now().Add(-24 * time.Hour).Format("20060102")
		sql := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %v LIKE %v
		`, newTable+t, oldTable)

		o := NewOrmWithDB(dbName)
		_, err := o.Raw(sql).Exec()
		return err
	}
}

func CreateTwoDayBeforeTable(dbName, oldTable, newTable string) func() error {
	return func() error {
		t := time.Now().Add(-48 * time.Hour).Format("20060102")
		sql := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %v LIKE %v
		`, newTable+t, oldTable)

		o := NewOrmWithDB(dbName)
		_, err := o.Raw(sql).Exec()
		return err
	}
}

// 使用者可能要在前面封装一层
func CreateNextMonthTable(dbName, oldTable, newTable string) func() error {
	return func() error {
		CreateThisMonthTable(dbName, oldTable, newTable)()
		thisMonth := time.Now()
		nextMonth := thisMonth.AddDate(0, 1, 0)
		t1 := nextMonth.Format("200601")
		sql := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %v LIKE %v
		`, newTable+t1, oldTable)

		o := NewOrmWithDB(dbName)
		_, err := o.Raw(sql).Exec()
		return err
	}
}

func CreateThisMonthTable(dbName, oldTable, newTable string) func() error {
	return func() error {
		t := time.Now().Format("200601")
		sql := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %v LIKE %v
		`, newTable+t, oldTable)

		o := NewOrmWithDB(dbName)
		_, err := o.Raw(sql).Exec()
		return err
	}
}
