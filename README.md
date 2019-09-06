# http-mock-server
 A simple http mock server, based on [Gin](https://github.com/gin-gonic/gin) framework.

Can be used in **Linux** , **Windows** and **Mac**. 

*For Mac version, I  don't have a Mac yet (It's freaking expensive...), so I never test it , sorry for that. Hopefully it will work.*

Bug life is a lifestyle of developer. 

So if there is any issues, please report it, I will resolve it as soon as possible.

## How to use?

### Download

Download zip file according to your operation system: 

- Linux
- Windows
- Mac

### Install

Unzip the zip file into your custom directory.

### Run

#### Windows

Double click on *http-mock-server.exe* in GUI.

Or execute `http-mock-server.exe` in CLI.

The execution output will be:

```powershell
C:\Users\hhhhappy\Documents\http-mock-server\http-mock-server>http-mock-server.exe
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /mock_http/callback       --> main.callBackAction (4 handlers)
[GIN-debug] GET    /mock_http/callback2      --> main.callBackAction (4 handlers)
```

#### Linux

Execute `http-mock-server` in CLI.

The execution output will be:

```
[hhhhappy@linux http-mock-server]$ ./http-mock-server 
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /mock_http/callback       --> main.callBackAction (4 handlers)
[GIN-debug] GET    /mock_http/callback2      --> main.callBackAction (4 handlers)
```

## Configure

You can find configure file in `config/config.yml`

Its content will be:

```yaml
port: 8588    # Define the listen Port
logPath: ./   # Log output directory
urlList:          # Define your listen URL
    # Custom Url
    # With this value, final url will be http://<ip>:<port>/mock_http/callback
  - url: callback_html

    # HTTP Method:
    # Only support GET and POST for now
    type: post

    # Custom Return Body's Filepath:
    # You can DIY body of return request by creating a body file
    returnBodyFile: ./config/callback.html

    # Custom return Header:
    # You can DIY header of return request
    header:
      version: 1.220
      session: ab2b1aab2cce31111

  - url: callback_string
    type: get
    returnBodyFile: ./config/callback_string

  - url: callback_file
    type: get
    returnBodyFile: ./config/callback_file.jpg
```

Explanation:

- ***port***: http server listen port
- ***logPath***: where should it put log files
- ***urlList***: definition of mock http URL, you can define multiple URL here.
  - ***url***: base of URL
    - `mock_http` is fixed string, can not be changed for now.
    - For example:  `callback`  => `http://<ip>:<port>/mock_http/callback`
  - ***type***: HTTP Method of URL, only support **POST** and **GET** for now
  - ***returnBodyFile***: Custom Return Body's filepath, You can DIY body of return request by creating a body file
  - ***header***: Custom return Header

## Log Request

Content of all requests will be saved in files (*.request). Filename is based on your custom url. 

For example, I defined URL `/mock_http/callback2`, after sending requests, a file `callback2.request` will be created in current directory. 

It's content will be:

```bash
2019/09/06 14:01:05 
[Method]
POST

[Query] 
{}

[Header]
{
    "Accept": [
        "*/*"
    ],
    "Accept-Encoding": [
        "gzip, deflate"
    ],
    "Cache-Control": [
        "no-cache"
    ],
    "Connection": [
        "keep-alive"
    ],
    "Content-Length": [
        "52"
    ],
    "Content-Type": [
        "text/plain"
    ],
    "Postman-Token": [
        "5178aac3-d1fa-4738-af63-3c69cd74bed7"
    ],
    "User-Agent": [
        "PostmanRuntime/7.6.0"
    ]
}

[Body]
Hhhhappy est un beau garçon. Hhhhappy is handsome.
```

## Examples

Examples are included in source, you can test it by yourself.

- Return HTML: `http://<ip>:<port>//mock_http/callback_html`
- Return  Raw text: `http://<ip>:<port>//mock_http/callback_string`
- Return File: `http://<ip>:<port>//mock_http/callback_file`