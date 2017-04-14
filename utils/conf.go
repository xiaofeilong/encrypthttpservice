package utils

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	LogInfo        loginfo
	Httpserverinfo httpserver
	Databaseinfo   database
	Encrytionkey   encryptkey
}

type httpserver struct {
	Serveraddress string
}

type loginfo struct {
	LogFile  string
	LogLevel int
}

type database struct {
	MysqlServer   string
	Mysqlusername string
	MysqlPassword string
	Mysqldbname   string
	MysqlMaxconn  int
	MysqlMaxidle  int
}

type encryptkey struct {
	Key  string
	Mode []string
}

var (
	Conf *Config
)

func init() {
	var (
		fp       *os.File
		fcontent []byte
	)
	Conf = new(Config)
	var err error
	if fp, err = os.Open("conf.toml"); err != nil {
		fmt.Println("open error ", err)
		return
	}

	if fcontent, err = ioutil.ReadAll(fp); err != nil {
		fmt.Println("ReadAll error ", err)
		return
	}

	if err = toml.Unmarshal(fcontent, Conf); err != nil {
		fmt.Println("toml.Unmarshal error ", err)
		return
	}

	return

}
