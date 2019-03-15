/**
* Copyright 2015 Huawei Technologies Co., Ltd. All rights reserved.
* eSDK is licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

//日志级别
const (
	OFF     = 0
	ERROR   = 1
	WARNING = 2
	INFO    = 3
	DEBUG   = 4
)

//日志输出区域
const (
	RECORD_CONSOLE             = 1
	RECORD_LOGFILE             = 2
	RECORD_CONSOLE_AND_LOGFILE = 3
)

type logMsg struct {
	logPath      string //日志文件路径,请确定文件路径存在
	logName      string //日志文件前缀名
	logLevel     uint   //日志级别，取值为logger.OFF,logger.ERROR,logger.WARNING,logger.INFO,logger.DEBUG。
	nowDay       string //当前日期
	consolePrint uint   //日志打印区域，1：只打印控制台，2：只记录日志文件，3：打印控制台并记录日志文件，其他值无效
}

var logInfo *logMsg = nil
var logHander *log.Logger = nil
var logFile *os.File = nil

func InitLog(logPath, logName string, logLevel, consolePrint uint) {
	if consolePrint < RECORD_CONSOLE || consolePrint > RECORD_CONSOLE_AND_LOGFILE {
		panic("the parameter of consolePrint input out of range,the value of level are " +
			"logger.RECORD_CONSOLE:only record on console," +
			"logger.RECORD_LOGFILE:only record on log file," +
			"logger.RECORD_CONSOLE_AND_LOGFILE:both record on console and log file")
	}
	if consolePrint > RECORD_CONSOLE && logLevel > OFF && (logPath == "" || logName == "") {
		panic("input logPath or LogName is empty")
	}
	if logLevel > DEBUG || logLevel < OFF {
		panic("the parameter of logLevel is out of range,the value of level are:" +
			"logger.ERROR,logger.WARNING,logger.INFO or logger.DEBUG")
	}

	var name string = logName
	if strings.Contains(logName, ".") {
		if strings.Split(logName, ".")[1] == "log" {
			name = strings.Split(logName, ".")[0]
		}
	}
	logInfo = &logMsg{logPath, name, logLevel, getNowDay(), consolePrint}
}

func LOG(level uint, format string, v ...interface{}) {
	if logInfo != nil && logInfo.logLevel != OFF {
		if level <= logInfo.logLevel {
			openLogFile()
			pc, _, line, _ := runtime.Caller(1)
			fm := fmt.Sprintf("[%s:%d]", runtime.FuncForPC(pc).Name(), line)
			str := fmt.Sprintf(format, v...)
			var msg string
			switch level {
			case ERROR:
				msg = fmt.Sprintf("%s[ERROR]:%s", fm, str)
			case WARNING:
				msg = fmt.Sprintf("%s[WARNING]:%s", fm, str)
			case INFO:
				msg = fmt.Sprintf("%s[INFO]:%s", fm, str)
			case DEBUG:
				msg = fmt.Sprintf("%s[DEBUG]:%s", fm, str)
			}
			recordLog(msg)
		}
	}

}

func recordLog(msg string) {
	if logInfo.consolePrint == RECORD_CONSOLE {
		fmt.Println(getNowTime(), msg)
	}
	if logInfo.consolePrint == RECORD_LOGFILE {
		logHander.Println(msg)
	}
	if logInfo.consolePrint == RECORD_CONSOLE_AND_LOGFILE {
		fmt.Println(getNowTime(), msg)
		logHander.Println(msg)
	}
}

func openLogFile() {
	if logInfo.consolePrint > RECORD_CONSOLE {
		filePath := filepath.Join(logInfo.logPath, logInfo.logName+"_"+logInfo.nowDay+".log")
		_, err := os.Stat(filePath)
		if getNowDay() != logInfo.nowDay || os.IsNotExist(err) || logFile == nil {
			logInfo.nowDay = getNowDay()
			closeLogFile()
			filePath = filepath.Join(logInfo.logPath, logInfo.logName+"_"+logInfo.nowDay+".log")
			logFile, err = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
			if err != nil {
				panic(err)
			}
			logHander = log.New(logFile, "", log.LstdFlags)
		}
	}
}

func closeLogFile() {
	if logFile != nil {
		logFile.Sync()
		defer logFile.Close()
	}
}

//当前日期
func getNowDay() string {
	//2017-05-24 16:41:57.5772647 +0800 WST
	return strings.Split(time.Now().String(), " ")[0]
}

func getNowTime() string {
	return time.Now().Format("2006/01/02 15:04:05")
}
