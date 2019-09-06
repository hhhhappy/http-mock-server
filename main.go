package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"http-mock-server/manager/configManager"
	"http-mock-server/manager/logManager"
	"io"
	"net/http"
	"path"
	"strings"
)

const HttpMethodPost = "POST"
const HttpMethodGet = "GET"

func main() {
	router := gin.Default()

	router.Use(func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Methods", "*")
		context.Header("Access-Control-Allow-Headers", "*")
	})

	testGroup := router.Group("/mock_http")
	{
		for _, url := range configManager.GetConf().UrlList {
			switch strings.ToUpper(url.Type) {
			case HttpMethodPost:
				testGroup.POST(url.Url, callBackAction)
				break
			case HttpMethodGet:
				testGroup.GET(url.Url, callBackAction)
				break
			default:
				fmt.Println("Method not supported: " + url.Type)
				return
			}
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
	header,_ := json.MarshalIndent(context.Request.Header, "", "    ")
	query,_ := json.MarshalIndent(context.Request.URL.Query(), "", "    ")

	// Retrieve Body
	buf := make([]byte, 128)
	var body []byte
	for {
		n, err := context.Request.Body.Read(buf)
		body = append(body, buf[0:n]...)

		if err == io.EOF {
			break
		}
	}

	baseUrl:= path.Base(strings.Split(context.Request.RequestURI, "?")[0])

	// Log received request
	logManager.LogRequest(context.Request.Method, string(query), string(header), string(body), baseUrl)

	urlDefine := configManager.GetConf().GetUrlDefinition(baseUrl)
	if urlDefine == nil {
		logManager.Log("Can't find url's definition. Please check your configure file. Calling: " + baseUrl)
		context.String(http.StatusInternalServerError, ``)
		return
	}

	if urlDefine.Status == 0 {
		urlDefine.Status = http.StatusOK
	}

	if len(urlDefine.ReturnBodyFile) == 0 {		// Don't have custom return body, return default body
		context.String(urlDefine.Status, `Hi, world. This is a greeting from Hhhhappy !`)
	}else{
		context.File(urlDefine.ReturnBodyFile)
		context.Status(urlDefine.Status)
	}

	// Set custom header
	for key, value := range urlDefine.Header {
		context.Header(key, value)
	}

	return
}
