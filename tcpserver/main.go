package main

import (
    "os"
    "fmt"
    "conf"
    "time"
    "flag"
    "math/rand"
    "github.com/astaxie/beego/logs"
)

var config conf.TCPConf

func init() {
    // parser config
    var confFile string
    flag.StringVar(&confFile, "c", "../conf/tcpserver.yaml", "config file")
    flag.Parse()

    err := conf.ConfParser(confFile, &config)
    if err != nil {
        logs.Critical("parser config failed:", err.Error())
        os.Exit(-1)
    }

    // init log
    logConfig := fmt.Sprintf(`{"filename":"%s","level":%s,"maxlines":0,"maxsize":0,"daily":true,"maxdays":%s}`,
                             config.Log.Logfile, config.Log.Loglevel, config.Log.Maxdays)
    logs.SetLogger(logs.AdapterFile, logConfig)
    logs.EnableFuncCallDepth(true)
    logs.SetLogFuncCallDepth(3)
    // logs.Async()

    // init redis
    err = initRedisConn(&config)
    if err != nil {
        logs.Critical("initRedisConn failed:", err.Error())
        os.Exit(-1)
    }

    // init db
    err = initDbConn(&config)
    if err != nil {
        logs.Critical("initDbConn failed:", err.Error())
        os.Exit(-1)
    }
    logs.Info("init successfully!")
}

// cleanup
func finalize() {
    closeCache()
    closeDB()
}

func main() {
    defer finalize()
    // generate random seed global
    rand.Seed(time.Now().UTC().UnixNano())
    // start event loop
    start(&config)
}

