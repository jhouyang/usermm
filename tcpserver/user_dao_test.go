package main

import (
    "testing"
)

func Test_getUserInfo(t *testing.T) {
    existUname := "test"
    _, err := getUserInfo(existUname)
    if err != nil {
        t.Error("getUserInfo failed, err:", err.Error())
    }
}

func Test_editUserInfo(t *testing.T) {
    // mode 1
    cnt := editUserInfo("test", "testusername", "", "", 1)
    if cnt != 0 && cnt != 1 {
        t.Error("editUserInfo nickname failed, cnt=", cnt)
    }
    // mode 2
    cnt = editUserInfo("test", "", "www.helloworld.com", "", 2)
    if cnt != 0 && cnt != 1 {
        t.Error("editUserInfo nickname failed, cnt=", cnt)
    }
    // mode 3
    cnt = editUserInfo("test", "testusername1", "www.google.com", "", 3)
    if cnt != 0 && cnt != 1 {
        t.Error("editUserInfo nickname failed, cnt=", cnt)
    }
}

