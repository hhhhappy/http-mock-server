package logManager

import (
	"fmt"
	"http-mock-server/manager/configManager"
	"log"
	"os"
	"sync"
)

var (
	/* File writer */
	fileWriterMap = map[string]*os.File{}
	mutex *sync.RWMutex
)

func init(){
	basePath := configManager.GetConf().LogPath

	err := Mkdir(basePath)
	if err != nil {
		panic(err)
	}

	mutex = new(sync.RWMutex)
}

func Mkdir(dir string) error {
	_, err := os.Stat(dir)

	// dir exists
	if err == nil {
		return nil
	}

	err = os.MkdirAll(dir, 0777)

	return err
}

// LogInfo logs info-level info
func Info(msg string, fileName string) {
	logToFile(msg, fileName)
}

// log to fileManager
func logToFile(msg string, fileName string) {
	mutex.Lock()
	defer mutex.Unlock()
	var writer *os.File
	writer, ok := fileWriterMap[fileName]
	if !ok {
		var err error
		writer, err = os.OpenFile(GetLogFilePath(fileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fileWriterMap[fileName] = writer
	}

	log.SetOutput(writer)
	log.SetFlags(log.Ldate | log.Ltime)

	log.Println(fmt.Sprintf("%s", msg))
}

func GetLogFilePath(fileName string) string {
	basePath := configManager.GetConf().LogPath
	return fmt.Sprintf("%s%s%s", basePath, "/", fileName + ".log")
}