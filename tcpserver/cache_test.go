package main

import (
    "time"
    "testing"
)

func Test_setTokenInfo(t* testing.T) {
    user := User{ID: 1, Username: "test", Nickname: "nickname", Passwd: "47ec2dd791e31e2ef2076caf64ed9b3d", Skey: "123456", Headurl: "www.google.com"}
    token := "4f38176c59172644e5ae0718d39b59d7"

    err := setTokenInfo(user, token)
    if err != nil {
        t.Error("setTokenInfo failed, ", err.Error())
    }
}

func Test_getTokenInfo(t* testing.T) {
    token := "4f38176c59172644e5ae0718d39b59d7"
    user, err := getTokenInfo(token)
    if err != nil {
        t.Error("getTokenInfo failed, ", err.Error())
    }

    if user.Username != "test" {
        t.Error("getTokenInfo not match, expect: test, get: ", user.Username)
    }
}

func Test_delTokenInfo(t* testing.T) {
    token := "4f38176c59172644e5ae0718d39b59d7"
    err := delTokenInfo(token)
    if err != nil {
        t.Error("delTokenInfo failed, ", err.Error())
    }

    // try get should return nil
    _, err = getTokenInfo(token)
    if err == nil {
        t.Error("delTokenInfo succ, but still get userinfo")
    }
}

func Test_tokenCacheInfoExpired(t* testing.T) {
    config.Redis.Cache.Tokenexpired = 2
    // set TokenInfo
    user := User{ID: 1, Username: "test", Nickname: "nickname", Passwd: "47ec2dd791e31e2ef2076caf64ed9b3d", Skey: "123456", Headurl: "www.google.com"}
    token := "4f38176c59172644e5ae0718d39b59d7"

    err := setTokenInfo(user, token)
    if err != nil {
        t.Error("setTokenInfo failed, ", err.Error())
    }

    // sleep to let token info expired
    time.Sleep(time.Duration(3) * time.Second)

    // try get should return nil
    _, err = getTokenInfo(token)
    if err == nil {
        t.Error("delTokenInfo succ, but still get userinfo")
    }
}

func Test_setUserCacheInfo(t* testing.T) {
    user := User{ID: 1, Username: "test", Nickname: "nickname", Passwd: "47ec2dd791e31e2ef2076caf64ed9b3d", Skey: "123456", Headurl: "www.google.com"}
    err := setUserCacheInfo(user)
    if err != nil {
        t.Error("setUserCacheInfo failed, ", err.Error())
    }
}

func Test_getUserCacheInfo(t* testing.T) {
    user, err := getUserCacheInfo("test")
    if err != nil {
        t.Error("getUserCacheInfo failed, ", err.Error())
    }

    if user.Username != "test" {
        t.Error("getUserCacheInfo not match, expect: test, get: ", user.Username)
    }
}
