package main

import (
    "path"
    "strings"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/astaxie/beego/logs"

    "utils"
    "rpcclient"
    codeModule "code"
)

// generate upload image file name
func generateImgName(fname, postfix string) string {
    ext := path.Ext(fname)
    fileName := strings.TrimSuffix(fname, ext)
    fileName = utils.Md5String(fileName + postfix)

    return fileName + ext
}

// login
func loginHandler(c *gin.Context) {
    // check params
    username := c.PostForm("username")
    passwd := c.PostForm("passwd")

    if len(passwd) != 32 {
        logs.Error("Invalid passwd:", passwd)
        c.JSON(http.StatusBadRequest, rpcclient.FormatResponse(codeModule.CodeInvalidPasswd, "", nil))
        return
    }

    uuid := utils.GenerateToken(username)
    logs.Debug(uuid, " -- loginHandler access from:", username, "@", passwd)

    // communicate with rcp server
    ret, token, rsp := rpcclient.Login(map[string]string{"username":username, "passwd":passwd, "uuid":uuid})
    // set cookie
    if ret == http.StatusOK && token != "" {
        c.SetCookie("token", token, config.Logic.Tokenexpire, "/", config.Server.IP, false, true)
        logs.Debug(uuid, " -- Set token ", token, "with expire:", config.Logic.Tokenexpire)
    }

    logs.Debug(uuid, " -- Succ get response from backend with", rsp["code"], " and msg:", rsp["msg"])
    c.JSON(ret, rsp)
}

// logout
func logoutHandler(c* gin.Context) {
    // check params
    username := c.PostForm("username")
    token, err := c.Cookie("token")
    if err != nil {
        logs.Error("Failed to get token from cookie, err:", err.Error())
        c.JSON(http.StatusBadRequest, rpcclient.FormatResponse(codeModule.CodeTokenNotFound, "", nil))
        return
    }

    if len(token) != 32 {
        logs.Error("Invalid token :", token)
        c.JSON(http.StatusBadRequest, rpcclient.FormatResponse(codeModule.CodeInvalidToken, "", nil))
        return
    }
    uuid := utils.GenerateToken(username)
    logs.Debug(uuid, " -- logoutHandler access from:", username, " with token:", token)

    // communicate with rcp server
    ret, rsp := rpcclient.Logout(map[string]string{"username":username, "token":token, "uuid":uuid})

    logs.Debug(uuid, " -- Succ to get response from backend with ", rsp["code"], " and msg:", rsp["msg"])
    c.JSON(ret, rsp)
}

// edit nickname
func editNicknameHandler(c* gin.Context) {
    // check params
    username := c.PostForm("username")
    nickname := c.PostForm("newnickname")
    token, err := c.Cookie("token")
    logs.Debug("access from:", username, " with token:", token, " and newname:", nickname)
    if err != nil {
        logs.Error("Failed to get token from cookie, err:", err.Error())
        c.JSON(http.StatusBadRequest, rpcclient.FormatResponse(codeModule.CodeTokenNotFound, "", nil))
        return
    }

    if len(token) != 32 {
        logs.Error("Invalid token :", token)
        c.JSON(http.StatusBadRequest, rpcclient.FormatResponse(codeModule.CodeInvalidToken, "", nil))
        return
    }

    uuid := utils.GenerateToken(username)
    logs.Debug(uuid, " -- editNicknameHandler access from:", username, " with token:", token, " new nickname:", nickname)
    // communicate with rcp server
    ret, rsp := rpcclient.EditUserinfo(map[string]string{"username":username, "token":token, "nickname":nickname, "headurl":"", "mode": "1", "uuid":uuid})

    logs.Debug(uuid, " -- Succ to get response from backend with ", rsp["code"], " and msg:", rsp["msg"])
    c.JSON(ret, rsp)
}

// uploadHeadurlHandle
func uploadHeadurlHandler(c* gin.Context) {
    // check params
    username := c.Query("username")
    token, err := c.Cookie("token")
    logs.Debug("access from:", username, " with token:", token)
    if err != nil {
        logs.Error("Failed to get token from cookie, err:", err.Error())
        c.JSON(http.StatusBadRequest, rpcclient.FormatResponse(codeModule.CodeTokenNotFound, "", nil))
        return
    }

    uuid := utils.GenerateToken(username)
    logs.Debug(uuid, " -- uploadHeadurlHandler access from:", username, " with token:", token)
    // step 1 : auth
    httpCode, tcpCode, msg := rpcclient.Auth(map[string]string{"username":username, "token":token, "uuid":uuid})
    if httpCode != http.StatusOK || tcpCode != 0 {
        logs.Error(uuid, " -- uploadHeadurlHandler Auth failed, msg:", msg)
        c.JSON(httpCode, rpcclient.FormatResponse(tcpCode, msg, nil))
        return
    }
    logs.Debug(uuid, " -- uploadHeadurlHandler Auth succ")
    // step 2 : save upload picture into file
    // save picture
    file, image, err := c.Request.FormFile("picture")
    if err != nil {
        logs.Error(uuid, " -- Failed to FormFile, err:", err.Error())
        c.JSON(http.StatusOK, rpcclient.FormatResponse(codeModule.CodeFormFileFailed, "", nil))
    }

    //// check image
    if image == nil {
        logs.Error(uuid, " -- Failed to get image from formfile!")
        c.JSON(http.StatusOK, rpcclient.FormatResponse(codeModule.CodeFormFileFailed, "", nil))
        return
    }
    //// check filesize
    size, err := utils.GetFileSize(file)
    if err != nil {
        logs.Error(uuid, " -- Failed to get filesize, err:", err.Error())
        c.JSON(http.StatusOK, rpcclient.FormatResponse(codeModule.CodeFileSizeErr, "", nil))
        return
    }
    if size == 0 || size > config.Image.Maxsize * 1024 * 1024 {
        logs.Error(uuid, " -- Filesize illegal, size:", size)
        c.JSON(http.StatusOK, rpcclient.FormatResponse(codeModule.CodeFileSizeErr, "", nil))
        return
    }
    logs.Debug(uuid, " -- uploadHeadurlHandler CheckImage succ")
    //// save
    imageName := generateImgName(image.Filename, username)
    fullPath  := config.Image.Savepath + imageName

    if err = c.SaveUploadedFile(image, fullPath); err != nil {
        logs.Error(uuid, " -- Failed to save file, err:", err.Error())
        c.JSON(http.StatusInternalServerError, rpcclient.FormatResponse(codeModule.CodeInternalErr, "", nil))
        return
    }
    logs.Debug(uuid, " -- Succ to save upload image, path:", fullPath)

    // step 3 : update picture info
    imageURL := config.Image.Prefixurl + "/" + fullPath
    ret, editRsp := rpcclient.EditUserinfo(map[string]string{"username": username, "token": token, "nickname": "", "headurl": imageURL, "mode": "2", "uuid":uuid})
    logs.Debug(uuid, " -- editUserInfo response:", ret)
    c.JSON(ret, editRsp)
}

// get user info
func getUserinfoHandler(c* gin.Context) {
    // check params
    username := c.Query("username")
    token, err := c.Cookie("token")
    logs.Debug("access from:", username, " with token:", token)
    if err != nil {
        logs.Error("Failed to get token from cookie, err:", err.Error())
        c.JSON(http.StatusBadRequest, rpcclient.FormatResponse(codeModule.CodeTokenNotFound, "", nil))
        return
    }

    if len(token) != 32 {
        logs.Error("Invalid token :", token)
        c.JSON(http.StatusBadRequest, rpcclient.FormatResponse(codeModule.CodeInvalidToken, "", nil))
        return
    }

    uuid := utils.GenerateToken(username)
    logs.Debug(uuid, " -- getUserinfoHandler access from:", username, " with token:", token)
    // communicate with rcp server
    ret, rsp := rpcclient.GetUserinfo(map[string]string{"username":username, "token":token, "uuid":uuid})
    logs.Debug(uuid, " -- Succ to get response from backend with ", rsp["code"], " and msg:", rsp["msg"])
    c.JSON(ret, rsp)
}

