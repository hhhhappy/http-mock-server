package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"http-mock-server/manager/config"
	"http-mock-server/manager/log"
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
		for _, url := range config.GetConf().UrlList {
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
		Addr:    ":" + config.GetConf().Port,
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
	log.LogRequest(context.Request.Method, string(query), string(header), string(body), baseUrl)

	urlDefine := config.GetConf().GetUrlDefinition(baseUrl)
	if urlDefine == nil {
		log.Log("Can't find url's definition. Please check your configure file. Calling: " + baseUrl)
		context.String(http.StatusInternalServerError, ``)
		return
	}

	if len(urlDefine.ReturnBodyFile) == 0 {		// Don't have custom return body, return default body
		context.String(http.StatusOK, `Hi, world. This is a greeting from Hhhhappy !`)
	}else{
		context.File(urlDefine.ReturnBodyFile)

	}

	// Set custom header
	for key, value := range urlDefine.Header {
		context.Header(key, value)
	}

	return
}
