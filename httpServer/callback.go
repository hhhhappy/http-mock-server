package httpServer

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"http-mock-server/manager/config"
	"http-mock-server/manager/log"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func callback(context *gin.Context) {
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

	requestDef := config.GetConf().GetRequestDefinition(baseUrl)
	if requestDef == nil {
		log.Log("Can't find url's definition. Please check your configure file. Calling: " + baseUrl)
		context.String(http.StatusInternalServerError, ``)
		return
	}

	// Set custom header
	for key, value := range requestDef.Header {
		context.Header(key, value)
	}

	if requestDef.Code == 0 {
		context.Writer.WriteHeader(http.StatusOK)
	}else{
		context.Writer.WriteHeader(requestDef.Code)
	}

	if len(requestDef.ReturnBodyFile) != 0 {
		f, err := os.Open(requestDef.ReturnBodyFile)
		if err != nil {
			log.Log(err.Error())
			log.Log("Can't find url's definition. Please check your configure file. Calling: " + baseUrl)
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
