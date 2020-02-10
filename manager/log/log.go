package log

import (
	"encoding/json"
	"fmt"
	"http-mock-server/manager/config"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	/* File writer */
	errLogFileWriter     *os.File
	accessLogFileWriter  *os.File
	requestFileWriterMap = map[string]*os.File{}
	mutex                *sync.RWMutex
)

const messagePattern = `%s
[Query] 
%s

[Header]
%s

[Body]
%s
`

func init() {
	basePath := config.GetConf().LogPath

	err := mkdir(basePath)
	if err != nil {
		panic(err)
	}

	mutex = new(sync.RWMutex)

	errLogFileWriter, err = os.OpenFile(getErrLogFilePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	accessLogFileWriter, err = os.OpenFile(getAccessLogFilePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
}

func GetAccessWriter() *os.File {
	return accessLogFileWriter
}

// Log received requests
func LogRequest(method, query, header, body string, fileName string) {
	msg := fmt.Sprintf(messagePattern, method, query, header, body)

	mutex.Lock()
	defer mutex.Unlock()
	var writer *os.File
	writer, ok := requestFileWriterMap[fileName]
	if !ok {
		var err error
		writer, err = os.OpenFile(getRequestFilePath(fileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		requestFileWriterMap[fileName] = writer
	}

	logToFile(msg, writer)
}

// Log system error
func Log(msg interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	if errLogFileWriter == nil {
		var err error
		errLogFileWriter, err = os.OpenFile(getErrLogFilePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	if m, ok := msg.(string); ok {
		logToFile(m, errLogFileWriter)
	} else {
		m, err := json.Marshal(msg)
		if err != nil {
			fmt.Println(err)
			logToFile("Can not encode interface to json: "+err.Error(), errLogFileWriter)
		} else {
			logToFile(string(m), errLogFileWriter)
		}
	}
}

// log to fileManager
func logToFile(msg string, writer *os.File) {
	log.SetOutput(writer)
	log.SetFlags(log.Ldate | log.Ltime)

	log.Println(fmt.Sprintf("%s", msg))
}

// make request file path
func getRequestFilePath(fileName string) string {
	basePath := config.GetConf().LogPath
	if fileName[0] == '/' {
		fileName = fileName[1:]
	}
	fileName = strings.ReplaceAll(fileName, "/", ".")
	return fmt.Sprintf("%s%s%s", basePath, "/", fileName+".request")
}

// make error log file path
func getErrLogFilePath() string {
	basePath := config.GetConf().LogPath
	return fmt.Sprintf("%s%s%s", basePath, "/", "error.log")
}

// make access log file path
func getAccessLogFilePath() string {
	basePath := config.GetConf().LogPath
	return fmt.Sprintf("%s%s%s", basePath, "/", "access.log")
}

func mkdir(dir string) error {
	_, err := os.Stat(dir)

	// dir exists
	if err == nil {
		return nil
	}

	err = os.MkdirAll(dir, 0777)

	return err
}