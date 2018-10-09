package main

import (
    "fmt"
    "math/rand"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/astaxie/beego/logs"

    "rpcclient"
    codeModule "code"
)

// login
func randomLoginHandler(c *gin.Context) {
    // check params
    uid := rand.Int63n(10000000)
    username := fmt.Sprintf("username%d", uid)
    passwd := "e10adc3949ba59abbe56e057f20f883e"

    if len(passwd) != 32 {
        logs.Error("Invalid passwd:", passwd)
        c.JSON(http.StatusBadRequest, rpcclient.FormatResponse(codeModule.CodeInvalidPasswd, "", nil))
        return
    }

    // communicate with rcp server
    ret, token, rsp := rpcclient.Login(map[string]string{"username":username, "passwd":passwd})
    // set cookie
    logs.Debug("set cookie with expire:", token)
    if ret == http.StatusOK && token != "" {
        c.SetCookie("token", token, config.Logic.Tokenexpire, "/", config.Server.IP, false, true)
        logs.Debug("set cookie with expire:", config.Logic.Tokenexpire)
    }

    logs.Debug("succ get response from backend with", rsp["code"], " and msg:", rsp["msg"])
    c.JSON(ret, rsp)
}
