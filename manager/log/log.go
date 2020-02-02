package log

import (
	"fmt"
	"http-mock-server/manager/config"
	"log"
	"os"
	"sync"
)

var (
	/* File writer */
	logFileWriter        *os.File
	requestFileWriterMap = map[string]*os.File{}
	mutex                *sync.RWMutex
)

const messagePattern = `
[Method]
%s

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

	logFileWriter, err = os.OpenFile(getLogFilePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
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
func Log(msg string) {
	mutex.Lock()
	defer mutex.Unlock()
	if logFileWriter == nil {
		var err error
		logFileWriter, err = os.OpenFile(getLogFilePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	logToFile(msg, logFileWriter)
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
	return fmt.Sprintf("%s%s%s", basePath, "/", fileName+".request")
}

// make log file path
func getLogFilePath() string {
	basePath := config.GetConf().LogPath
	return fmt.Sprintf("%s%s%s", basePath, "/", "error.log")
}
