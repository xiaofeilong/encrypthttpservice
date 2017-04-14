package dataservice

import (
	"database/sql"
	"encryptService/utils"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	mysqldb *sql.DB
)

func init() {
	var server string = utils.Conf.Databaseinfo.Mysqlusername + ":" + utils.Conf.Databaseinfo.MysqlPassword +
		"@tcp(" + utils.Conf.Databaseinfo.MysqlServer + ")/" + utils.Conf.Databaseinfo.Mysqldbname + "?charset=utf8"
	var err error
	mysqldb, err = sql.Open("mysql", server)
	if err != nil {
		fmt.Println(err.Error())
	}
	mysqldb.SetMaxOpenConns(utils.Conf.Databaseinfo.MysqlMaxconn)
	mysqldb.SetMaxIdleConns(utils.Conf.Databaseinfo.MysqlMaxidle)
}

func Insertkey(videoid string, encryptmode string, encryptkey string) (string, string) {
	var querystr string = "SELECT encryptkey, encryptmode FROM videoencryptinfo where videoid=\"" + videoid + "\";"
	rows, err := mysqldb.Query(querystr)
	if err != nil {
		fmt.Println(err.Error())
		return "", ""
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&encryptkey, &encryptmode)
		return encryptkey, encryptmode
	}

	stmt, err := mysqldb.Prepare(`INSERT videoencryptinfo (videoid, encryptkey, encryptmode) values (?,?,?)`)
	res, err := stmt.Exec(videoid, encryptkey, encryptmode)
	_, err = res.LastInsertId()
	return encryptkey, encryptmode
}

func Getkey(videoid string) (string, string, bool) {
	var ret bool
	ret = false
	var encryptkey string
	var encryptmode string
	var querystr string = "SELECT videoid, encryptkey, encryptmode FROM videoencryptinfo where videoid=\"" + videoid + "\";"
	fmt.Println(querystr)
	rows, err := mysqldb.Query(querystr)
	if err != nil {
		return encryptkey, encryptmode, ret
	}
	defer rows.Close()

	for rows.Next() {
		var videoid string
		rows.Scan(&videoid, &encryptkey, &encryptmode)
		ret = true
		break
	}
	fmt.Println(encryptkey)
	return encryptkey, encryptmode, ret
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
