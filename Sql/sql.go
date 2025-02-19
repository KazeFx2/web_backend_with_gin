package Sql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"main/Logger"
	"main/Vars"
)

var Db *sql.DB

func ConnectPSQL() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		Vars.Host, Vars.Port, Vars.User, Vars.Password, Vars.Dbname)
	var err error
	Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		Logger.Log(Logger.LogMerge, fmt.Sprintf("Connect to sql '%s' failed.", psqlInfo))
		return err
	}
	err = Db.Ping()
	if err != nil {
		Logger.Log(Logger.LogMerge, fmt.Sprintf("Connect to sql '%s' failed.", psqlInfo))
		return err
	}
	_, err = Db.Exec("SET TIME ZONE 'Asia/Shanghai'")
	if err != nil {
		return err
	}
	Logger.Log(Logger.LogInfo, "SQL Connected.")
	return nil
}
