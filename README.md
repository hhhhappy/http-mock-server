# http-mock-server
 A simple http mock server, based on [Gin](https://github.com/gin-gonic/gin) framework.

Can be used in **Linux** , **Windows** and **Mac**.

Method supported: GET, POST, DELETE, OPTIONS, HEAD, PUT, PATCH

Notice: do not support file uploading !

So if there is any issues, please report it, I will resolve it as soon as possible.

## How to use?

### Download

Download zip file according to your operation system: 

- [Linux](https://github.com/hhhhappy/http-mock-server/releases/download/v0.0.1/http-mock-server-linux.zip)
- [Windows](https://github.com/hhhhappy/http-mock-server/releases/download/v0.0.1/http-mock-server-win.zip)
- [Mac](https://github.com/hhhhappy/http-mock-server/releases/download/v0.0.1/http-mock-server-mac.zip)

### Install

Unzip the zip file into your custom directory.

### Run

#### Windows

Double click on *http-mock-server.exe* in GUI.

Or execute `http-mock-server.exe` in CLI.

The execution output will be:

```powershell
[POST]		"post_example"	Response body: "./config/requests/post_example.html"
[GET]		"get_example"	Response body: "./config/requests/get_example.jpg"
[DELETE]		"delete_example"	Response body: "./config/requests/delete_example"
[OPTIONS]		"options_example"	Response body: "./config/requests/options_example"
[HEAD]		"head_example"	Without response body
[PUT]		"put_example"	Response body: "./config/requests/put_example"
[PATCH]		"/patch/patch_example"	Response body: "./config/requests/patch_example"
7 APIs were set!
Server is Listening to port: 8588
```

#### Linux or Mac

Execute `http-mock-server` in CLI.

The execution output will be same as Windows.


## Server Configure

You can find configure file in `config/config.yml`

Its content will be:

```yaml
port: 8588    # Define the listen Port
logPath: ./output   # Log output directory
logAccessSummary: true   # If need to log access's summary
defaultHeaders:  # Set common header for all the requests
  Access-Control-Allow-Origin: "*"
  Access-Control-Allow-Methods: "*"
  Access-Control-Allow-Headers: "*"
```

Explanation:

- ***port***: http server listen port
- ***logPath***: where should it put log files
- ***logAccessSummary***: if need to log access log
- ***defaultHeaders***: you can define common headers by writing a key-value map here.

## Request definition
You can define your own requests in `config/requests.yml`

````yaml
  # Custom Url
  # With this value, final url will be http://<ip>:<port>/mock_http/post_example
- url: post_example

  # HTTP Method:
  # Only support GET and POST for now
  type: post

  # Custom Return Body's Filepath:
  # You can DIY body of return request by creating a body file
  returnBodyFile: ./config/requests/post_example.html

  # Custom return Header:
  # You can DIY header of return request
  header:
    version: 1.220
    session: ab2b1aab2cce31111
  # Custom response's status code
  # If not set, return 200 OK by default
  code: 201

- url: get_example
  type: get
  returnBodyFile: ./config/requests/get_example.jpg
  code: 209
````
Explanation:
  - ***url***: base of URL
    - For example:  `/mock/get_example`  => `http://<ip>:<port>/mock/get_example`
  - ***type***: HTTP Method of URL, only support **POST** and **GET** for now
  - ***returnBodyFile***: Custom Return Body's filepath, You can DIY body of return request by creating a body file
  - ***header***: Custom return Header
  - ***code***: Custom response's status code, if not set, return 200 OK by default

## Log Request

Content of all requests will be saved in files (*.request). Filename is based on your custom url. 

For example, I defined URL `/mock_http/callback2`, after sending requests, a file `mock_http.callback2.request` will be created in log directory. 

It's content will be:

```bash
2020/02/11 11:36:53 PATCH
[Query] 
{
    "i": ["1"],
    "q": ["9"]
}

[Header]
{
    "Accept-Encoding": ["gzip, deflate, br"],
    "Connection": ["keep-alive"],
    "Postman-Token": ["df95c946-aede-47fd-96af-e89a97df5c57"],
    "Content-Length": ["25"],
    "Content-Type": ["application/x-www-form-urlencoded"],
    "User-Agent": ["PostmanRuntime/7.22.0"],
    "Accept": ["*/*"],
    "Cache-Control": ["no-cache"]
}

[Body]
admin=admin&date=1980-7-1
```

## Examples

Examples are included in directory config/requests, you can test it by yourself.
