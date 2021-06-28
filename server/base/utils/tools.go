package utils

import (
	"fmt"
	"gameserver-997/server/base/logger"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"time"
)

func HttpRequestWrap(uri string, targat func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
				logger.Info("===================http server panic recover=============== uri: %s err:%v", uri, err)
			}
		}()
		st := time.Now()
		logger.Debug("User-Agent: ", request.Header["User-Agent"])
		targat(response, request)
		logger.Debug(fmt.Sprintf("%s cost total time: %f ms", uri, time.Now().Sub(st).Seconds()*1000))
	}
}

func ReSettingLog() {
	// --------------------------------------------init log start
	// logger.SetConsole(GlobalObject.SetToConsole)
	// if GlobalObject.LogFileType == logger.ROLLINGFILE {
	// 	logger.SetRollingFile(GlobalObject.LogPath, GlobalObject.LogName,
	// 		GlobalObject.MaxLogNum, GlobalObject.MaxFileSize, GlobalObject.LogFileUnit)
	// } else {
	// 	logger.SetRollingDaily(GlobalObject.LogPath, GlobalObject.LogName)
	// 	logger.SetLevel(GlobalObject.LogLevel)
	// }
	logger.ResetLogger(GlobalObject.LogLevel, GlobalObject.LogPath, GlobalObject.LogName)
	// --------------------------------------------init log end
}

func XingoTry(f reflect.Value, args []reflect.Value) (res []reflect.Value, returnError error) {
	defer func() {
		if err := recover(); err != nil {
			returnError = err.(error)
		}
	}()
	return f.Call(args), nil
}

//读取配置文件
func ReadConf(filePath string) ([]byte, error) {
	env := os.Getenv("environmentName")
	if env == "" {
		env = "local"
	}
	filePath = path.Join("conf", env, filePath)
	return ReadFile(filePath)
}

func ReadFile(filePath string) ([]byte, error) {
	fBytes, err := ioutil.ReadFile(filePath)
	if err == nil {
		return fBytes, nil
	}
	dir, err1 := filepath.Abs(filepath.Dir(os.Args[0]))
	if err1 != nil {
		return nil, err
	}
	fBytes, err = ioutil.ReadFile(filepath.Join(dir, filePath))
	return fBytes, err
}

func GetRouteName(obj interface{}) (string, string) {
	typeString := reflect.TypeOf(obj).String()
	typeString = typeString[5:]
	pfStart := -1
	serverNameStart := -1
	interNameStart := -1
	for i := 0; i < len(typeString); i++ {
		if typeString[i] >= 65 && typeString[i] <= 90 {
			//首字母大写
			if pfStart < 0 {
				pfStart = i
				continue
			}
			if serverNameStart < 0 {
				serverNameStart = i
				continue
			}

			if interNameStart < 0 {
				interNameStart = i
				break
			}
		}
	}

	pfStr := typeString[pfStart:serverNameStart]
	serverStr := typeString[serverNameStart:interNameStart]
	interStr := typeString[interNameStart:]
	routeName := serverStr + "_" + interStr
	fmt.Println("[GetRouteName]: ", routeName)

	return pfStr, routeName
}
