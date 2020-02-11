package httpServer

import (
	"github.com/gin-gonic/gin"
	"http-mock-server/manager/config"
	"http-mock-server/manager/log"
	"io"
	"net/http"
	"os"
)

func callback(context *gin.Context) {
	header, err := arrStringMap(context.Request.Header).MarshalJSON()
	if err != nil {
		log.Log(err)
	}

	query, err := arrStringMap(context.Request.URL.Query()).MarshalJSON()
	if err != nil {
		log.Log(err)
	}

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

	urlPath := context.Request.URL.Path

	// Log received request
	log.LogRequest(context.Request.Method, string(query), string(header), string(body), urlPath)

	requestDef := config.GetConf().GetRequestDefinition(urlPath)
	if requestDef == nil {
		log.Log("Can't find url's definition. Please check your configure file. Calling: " + urlPath)
		context.String(http.StatusInternalServerError, ``)
		return
	}

	// Set custom header
	for key, value := range requestDef.Header {
		context.Header(key, value)
	}

	if requestDef.Code == 0 {
		context.Writer.WriteHeader(http.StatusOK)
	} else {
		context.Writer.WriteHeader(requestDef.Code)
	}

	if len(requestDef.ReturnBodyFile) != 0 {
		f, err := os.Open(requestDef.ReturnBodyFile)
		if err != nil {
			log.Log(err.Error())
			log.Log("Can't find url's definition. Please check your configure file. Calling: " + urlPath)
			context.String(http.StatusInternalServerError, ``)
			return
		}

		defer f.Close()

		for {
			n, errRead := f.Read(buf)
			_, errWrite := context.Writer.Write(buf[0:n])
			if errWrite != nil {
				log.Log(errWrite.Error())
				context.String(http.StatusInternalServerError, ``)
				return
			}
			if errRead == io.EOF {
				break
			}
		}
	}

	return
}
