package main

import (
    "os"
    "fmt"
    "time"
    "flag"
    "conf"
    "io/ioutil"
    "rpcclient"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/astaxie/beego/logs"
)

var config conf.HTTPConf

func init() {
    // parser config
    var confFile string
    flag.StringVar(&confFile, "c", "../conf/httpserver.yaml", "config file")
    flag.Parse()

    err := conf.ConfParser(confFile, &config)
    if err != nil {
        logs.Critical("Parser config failed, err:", err.Error())
        os.Exit(-1)
    }

    // init log
    logConfig := fmt.Sprintf(`{"filename":"%s","level":%s,"maxlines":0,"maxsize":0,"daily":true,"maxdays":%s}`,
                             config.Log.Logfile, config.Log.Loglevel, config.Log.Maxdays)
    logs.SetLogger(logs.AdapterFile, logConfig)
    logs.EnableFuncCallDepth(true)
    logs.SetLogFuncCallDepth(3)
    logs.Async()

    // init userclient (pool)
    err = rpcclient.InitPool(config.Rpcserver.Addr, config.Pool.Initsize, config.Pool.Capacity, time.Duration(config.Pool.Maxidle) * time.Second)
    if err != nil {
        logs.Critical("InitPool failed, err:", err.Error())
        os.Exit(-2)
    }
}

// cleanup global objects
func finalize() {
    rpcclient.DestoryPool()
}

func main() {
    defer finalize()

    gin.SetMode(gin.ReleaseMode)
    gin.DefaultWriter = ioutil.Discard

    engine := gin.Default()
    engine.Any("/welcome", webRoot)
    engine.POST("/login", loginHandler)
    engine.POST("/logout", logoutHandler)
    engine.GET("/getuserinfo", getUserinfoHandler)
    engine.POST("/editnickname", editNicknameHandler)
    engine.POST("/uploadpic", uploadHeadurlHandler)

    engine.POST("/randlogin", randomLoginHandler)
    engine.Static("/static/", "./static/")
    engine.Static("/upload/images/", "./upload/images/")

    engine.Run(fmt.Sprintf(":%d", config.Server.Port))
}

func webRoot(context *gin.Context) {
    context.String(http.StatusOK, "hello, world")
}
