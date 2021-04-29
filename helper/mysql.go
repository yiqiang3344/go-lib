package helper

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/micro/go-micro/v2/config"
)

var connMap map[string]*sqlx.DB

func DefaultDB() (*sqlx.DB, error) {
	return ConnectDB("db")
}

func ConnectDB(name string) (*sqlx.DB, error) {
	if len(connMap) == 0 {
		connMap = make(map[string]*sqlx.DB)
	}

	if _, ok := connMap[name]; ok {
		if err := connMap[name].Ping(); err == nil {
			//DebugLog("mysql connect ping success:"+name, "")
			return connMap[name], nil
		}
		//DebugLog("mysql connect ping failed:"+name, "")
	}

	cfgMap := config.Get(name).StringMap(map[string]string{})
	conn, err := sqlx.Open("mysql", cfgMap["user"]+":"+cfgMap["password"]+"@tcp("+cfgMap["host"]+":"+cfgMap["port"]+")/"+cfgMap["database"])
	if err != nil {
		ErrorLog("open mysql failed:"+err.Error(), "")
		return nil, err
	}
	//DebugLog("redis pool create:"+name, "")
	connMap[name] = conn
	return conn, nil
}
