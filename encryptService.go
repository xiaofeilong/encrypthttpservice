package main

import (
	"encoding/json"
	"encryptService/dataservice"
	"encryptService/utils"
	"math/rand"
	"net"
	"net/http"
	"runtime"
	"time"
)

const (
	SUCCESS        = 0
	PARA_ERROR     = 1
	DATABASE_ERROR = 2
)

type EncrypyInfo struct {
	Ret         int    //返回值 0--success
	Encryptmode string //加密方式
	Encryptkey  string //密钥
}

func GetRandomString(keylen int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < keylen; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//生成加密密钥并存入数据库中
func generatekey(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	var result EncrypyInfo
	result.Ret = PARA_ERROR
	if len(req.Form["videoid"]) > 0 && len(req.Form["encryptmode"]) > 0 {
		for _, mode := range utils.Conf.Encrytionkey.Mode {
			if mode == req.Form["encryptmode"][0] {
				randkey := GetRandomString(128)
				result.Ret = SUCCESS
				encryptkey, encryptmode := dataservice.Insertkey(req.Form["videoid"][0], req.Form["encryptmode"][0], randkey)
				result.Encryptmode = encryptmode
				result.Encryptkey = utils.AESCBCEncrypter([]byte(utils.Conf.Encrytionkey.Key), []byte(encryptkey))
				break
			}

		}

	}
	backresult, _ := json.Marshal(result)
	w.Header().Add("Connection", "close")
	//返回值加密AES
	w.Write(backresult)
}

//获取密钥
func getencryptkey(w http.ResponseWriter, req *http.Request) {
	var encryptkey, encryptmode string
	var ret bool
	req.ParseForm()
	var result EncrypyInfo
	result.Ret = PARA_ERROR
	if len(req.Form["videoid"]) > 0 {
		encryptkey, encryptmode, ret = dataservice.Getkey(req.Form["videoid"][0])
		if ret {
			result.Ret = SUCCESS
			result.Encryptmode = encryptmode
			result.Encryptkey = utils.AESCBCEncrypter([]byte(utils.Conf.Encrytionkey.Key), []byte(encryptkey))
		} else {
			result.Ret = DATABASE_ERROR
		}

	}
	backresult, _ := json.Marshal(result)
	w.Header().Add("Connection", "close")
	//返回值加密AES
	w.Write(backresult)
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	utils.SetConfig(utils.Conf.LogInfo.LogLevel, utils.Conf.LogInfo.LogFile)
	var server string = utils.Conf.Httpserverinfo.Serveraddress
	http.HandleFunc("/generatekey", generatekey)
	http.HandleFunc("/getencryptkey", getencryptkey)
	var listener net.Listener
	listener, _ = net.Listen("tcp", server)
	go func() {
		err := http.Serve(listener, nil)
		if err != nil {
			//utils.Logger.Println("test")
		}
	}()

	t := time.NewTimer(time.Second * 5)

	for {
		select {
		case <-t.C:
			t.Reset(time.Second * 5)
		}
	}

}
