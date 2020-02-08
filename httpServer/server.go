package httpServer

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"http-mock-server/manager/config"
	"http-mock-server/manager/log"
	"io"
	"net/http"
	"path"
	"strings"
	"time"
)

const urlPrefix = "/mock_http"
const requestDefLog = "[%s]		\"%s\"	"
const withBodyLog = "Response body: \"%s\"\n"
const withoutBodyLog = "Without response body\n"

var errMethodNotSupported = errors.New("Method not supported: ")

func Run() error {
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Log access
	if config.GetConf().LogAccessSummary {
		gin.DefaultWriter = io.MultiWriter(log.GetAccessWriter())
		router.Use(gin.LoggerWithFormatter(logFormat))
	}

	// Set default headers
	if len(config.GetConf().DefaultHeaders) != 0 {
		router.Use(func(context *gin.Context) {
			for key, value := range config.GetConf().DefaultHeaders {
				context.Header(key, value)
			}
		})
	}

	// Set router
	mockGroup := router.Group(urlPrefix)
	if err := setRouter(mockGroup);err != nil{
		return err
	}

	server := &http.Server{
		Addr:    ":" + config.GetConf().Port,
		Handler: router,
	}

	fmt.Println("Server is Listening to port:", config.GetConf().Port)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func setRouter(g *gin.RouterGroup) error {
	for _, request := range config.GetConf().Requests {
		requestMethod := strings.ToUpper(request.Type)
		switch requestMethod {
		case http.MethodPost:
			g.POST(request.Url, callback)
			break
		case http.MethodGet:
			g.GET(request.Url, callback)
			break
		case http.MethodDelete:
			g.DELETE(request.Url, callback)
			break
		case http.MethodHead:
			g.HEAD(request.Url, callback)
			break
		case http.MethodOptions:
			g.OPTIONS(request.Url, callback)
			break
		case http.MethodPatch:
			g.PATCH(request.Url, callback)
			break
		case http.MethodPut:
			g.PUT(request.Url, callback)
			break
		default:
			fmt.Println(errMethodNotSupported.Error() + request.Type)
			return errMethodNotSupported
		}
		if len(request.ReturnBodyFile) == 0{
			fmt.Printf(requestDefLog + withoutBodyLog, requestMethod, path.Join(urlPrefix, request.Url))
		}else{
			fmt.Printf(requestDefLog + withBodyLog, requestMethod, path.Join(urlPrefix, request.Url), request.ReturnBodyFile)
		}
	}

	fmt.Printf("%d APIs were set!\n", len(config.GetConf().Requests))
	return nil
}

func logFormat(param gin.LogFormatterParams) string {
	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
		param.ClientIP,
		param.TimeStamp.Format(time.RFC1123),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}
