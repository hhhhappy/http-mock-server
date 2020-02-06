package httpServer

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"http-mock-server/manager/config"
	"net/http"
	"strings"
)

var errMethodNotSupported = errors.New("Method not supported: ")

func Run() error {
	router := gin.Default()
	if len(config.GetConf().DefaultHeaders) != 0 {
		router.Use(func(context *gin.Context) {
			for key, value := range config.GetConf().DefaultHeaders {
				context.Header(key, value)
			}
		})
	}

	mockGroup := router.Group("/mock_http")
	{
		for _, request := range config.GetConf().Requests {
			switch strings.ToUpper(request.Type) {
			case http.MethodPost:
				mockGroup.POST(request.Url, callback)
				break
			case http.MethodGet:
				mockGroup.GET(request.Url, callback)
				break
			default:
				fmt.Println(errMethodNotSupported.Error() + request.Type)
				return errMethodNotSupported
			}
		}
	}

	server := &http.Server{
		Addr:    ":" + config.GetConf().Port,
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
