package main

import (
    "github.com/astaxie/beego/logs"
)

const (
    editUsername = 1
    editHeadurl  = 2
    editBoth     = 3
)

func getUserInfo(username string) (User, error) {
    // try cache
    user, err := getUserCacheInfo(username)
    if err == nil && user.Username == username {
        return user, err
    }

    // get from db
    user, err = getDbUserInfo(username)
    if err != nil {
        return user, err
    }

    // update cache
    serr := setUserCacheInfo(user)
    if serr != nil {
        logs.Error("cache userinfo failed for user:", user.Username, " with err:", serr.Error())
    }

    return user, err
}

// edit userinfo
func editUserInfo(username, nickname, headurl, token string, mode uint32) int64 {
    // update db info
    var affectedRows int64
    switch mode {
        case editUsername:
            affectedRows = updateDbNickname(username, nickname)
        case editHeadurl:
            affectedRows = updateDbHeadurl(username, headurl)
        case editBoth:
            affectedRows = updateDbUserinfo(username, nickname, headurl)
        default:
            // do nothing
            break
    }

    // on successing, update cache or delete it if updating failed
    if affectedRows == 1 {
        user, err := getDbUserInfo(username)
        if err == nil {
            updateCachedUserinfo(user)
            if token != "" {
                err = setTokenInfo(user, token)
                if err != nil {
                    logs.Error("update token failed:", err.Error())
                    delTokenInfo(token)
                }
            }
        } else {
            logs.Error("Failed to get dbUserInfo for cache, username:", username, " with err:", err.Error())
        }
    }
    return affectedRows
}

// auth
func auth(username, token string) bool {
    user, err := getTokenInfo(token)
    if err != nil {
        logs.Error("failed to getTokenInfo, token:", token)
        return false
    }
    if user.Username != username {
        logs.Error("invalid token info, username not match!")
        return false
    }
    return true
}

