package main

import (
	"github.com/gin-gonic/gin"
	"http-mock-server/manager/configManager"
	"http-mock-server/manager/logManager"
	"io"
	"net/http"
	"path"
)

func main() {
	router := gin.Default()

	router.Use(func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Methods", "GET, POST")
		context.Header("Access-Control-Allow-Headers", "x-requested-with,content-type")
	})

	testGroup := router.Group("/test")
	{
		for _, url := range configManager.GetConf().Url {
			testGroup.POST(url, callBackAction)
			testGroup.GET(url, callBackAction)
		}
	}

	server := &http.Server{
		Addr:    ":" + configManager.GetConf().Port,
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		println(err.Error())
	}
}

func callBackAction(context *gin.Context) {
	buf := make([]byte, 1024)
	var result []byte
	for {
		n, err := context.Request.Body.Read(buf)
		result = append(result, buf[0:n]...)

		if err == io.EOF {
			break
		}
	}
	logManager.Info(string(result), path.Base(context.Request.RequestURI))
	context.JSON(http.StatusOK, gin.H{
		"result": 0,
	})
	return
}
