package main

import (
    "testing"
    pb "proto"
    "golang.org/x/net/context"
)

func Test_Login(t *testing.T) {
    server := UserServer{}
    req := pb.LoginRequest{Username: "test", Passwd:"098f6bcd4621d373cade4e832627b4f6"}
    _, err := server.Login(context.Background(), &req)
    if err != nil {
        t.Error("Login failed, err:", err.Error())
    }
}

func Test_GetUserInfo(t *testing.T) {
    // login
    server := UserServer{}
    req := pb.LoginRequest{Username: "test", Passwd:"098f6bcd4621d373cade4e832627b4f6"}
    rsp, err := server.Login(context.Background(), &req)
    if err != nil {
        t.Error("Login failed, err:", err.Error())
        return
    }

    // and query userinfo
    creq := pb.CommRequest{Username: "test", Token: rsp.Token}
    _, err = server.GetUserInfo(context.Background(), &creq)
    if err != nil {
        t.Error("GetUserInfo failed, err:", err.Error())
        return
    }
}

func Test_EditUserInfo(t *testing.T) {
    // login
    server := UserServer{}
    req := pb.LoginRequest{Username: "test", Passwd:"098f6bcd4621d373cade4e832627b4f6"}
    rsp, err := server.Login(context.Background(), &req)
    if err != nil {
        t.Error("Login failed, err:", err.Error())
        return
    }

    // and query userinfo
    ereq := pb.EditRequest{Username: "test", Token: rsp.Token, Nickname: "hello", Mode: 1}
    _, err = server.EditUserInfo(context.Background(), &ereq)
    if err != nil {
        t.Error("GetUserInfo failed, err:", err.Error())
        return
    }
}

func Test_Logout(t *testing.T) {
    // login
    server := UserServer{}
    req := pb.LoginRequest{Username: "test", Passwd:"098f6bcd4621d373cade4e832627b4f6"}
    rsp, err := server.Login(context.Background(), &req)
    if err != nil {
        t.Error("Login failed, err:", err.Error())
        return
    }

    // and query userinfo
    creq := pb.CommRequest{Username: "test", Token: rsp.Token}
    _, err = server.Logout(context.Background(), &creq)
    if err != nil {
        t.Error("GetUserInfo failed, err:", err.Error())
        return
    }
}
