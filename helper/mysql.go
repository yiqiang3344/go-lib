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

func ConnectDB(cfg string) (*sqlx.DB, error) {
	if len(connMap) == 0 {
		connMap = make(map[string]*sqlx.DB)
	}

	if _, ok := connMap[cfg]; ok {
		if err := connMap[cfg].Ping(); err == nil {
			//DebugLog("mysql connect ping success:"+cfg, "")
			return connMap[cfg], nil
		}
		//DebugLog("mysql connect ping failed:"+cfg, "")
	}

	cfgMap := config.Get(cfg).StringMap(map[string]string{})
	conn, err := sqlx.Open("mysql", cfgMap["user"]+":"+cfgMap["password"]+"@tcp("+cfgMap["host"]+":"+cfgMap["port"]+")/"+cfgMap["database"])
	if err != nil {
		ErrorLog("open mysql failed:"+err.Error(), "")
		return nil, err
	}
	//DebugLog("redis pool create:"+cfg, "")
	connMap[cfg] = conn
	return conn, nil
}
