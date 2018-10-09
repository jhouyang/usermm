package conf

import (
    "testing"
)

func Test_ConfParser(t *testing.T) {
    // http
    var conf HTTPConf
    err := ConfParser("http_test.yaml", &conf)
    if err == nil {
        t.Error("Unstrict failed!")
    }

    err = ConfParser("httpserver.yaml", &conf)
    if err != nil {
        t.Error(err.Error())
    }
    // tcp
    var tconf TCPConf
    err = ConfParser("tcp_test.yaml", &tconf)
    if err == nil {
        t.Error("Unstrict failed!")
    }

    err = ConfParser("tcpserver.yaml", &tconf)
    if err != nil {
        t.Error(err.Error())
    }
}
