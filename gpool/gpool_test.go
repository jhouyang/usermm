package gpool

import (
    "time"
    "context"
    "testing"
    "google.golang.org/grpc"
)

func Test_NewPool(t *testing.T) {
    pool, err := NewPool(func() (*grpc.ClientConn, error){
        return grpc.Dial("www.baidu.com", grpc.WithInsecure())
}, 1, 3, time.Second * 1)
    if err != nil {
        t.Error("create failed:", err.Error())
    }
    defer pool.Close()

    if pool.Available() != 3 {
        t.Error("Available should be 3, now: ", pool.Available())
    }

    if pool.Capacity() != 3 {
        t.Error("pool.Capacity should be 3, now: ", pool.Capacity())
    }
}

func Test_Put(t* testing.T) {
    pool, err := NewPool(func() (*grpc.ClientConn, error){
        return grpc.Dial("www.baidu.com", grpc.WithInsecure())
}, 1, 3, time.Second * 1)
    if err != nil {
        t.Error("create failed:", err.Error())
    }
    defer pool.Close()

    var conn Conn
    conn.C, _ = pool.factory()
    // put to full pool
    err = pool.Put(&conn)
    if err == nil || err.Error() != ErrPoolFulled.Error() {
        t.Error("Put to full pool should return error")
    }

    // get and put
    ctx, cancel := context.WithDeadline(context.Background(),  time.Now().Add(10 * time.Millisecond))
    conn2, err := pool.Get(ctx)
    cancel()
    if err != nil {
        t.Error("get from pool failed, ", err.Error())
    }
    err = pool.Put(conn2)
    if err != nil {
        t.Error("Put failed, ", err.Error())
    }
}

func Test_Get(t* testing.T) {
    pool, err := NewPool(func() (*grpc.ClientConn, error){
        return grpc.Dial("www.baidu.com", grpc.WithInsecure())
}, 1, 1, time.Second * 1)
    if err != nil {
        t.Error("create failed:", err.Error())
    }
    defer pool.Close()

    // normal get
    ctx, cancel := context.WithDeadline(context.Background(),  time.Now().Add(10 * time.Millisecond))
    conn1, err := pool.Get(ctx)
    cancel()
    if err != nil {
        t.Error("get from pool failed, ", err.Error())
    }

    // get from empty pool
    ctx2, cancel2 := context.WithDeadline(context.Background(),  time.Now().Add(10 * time.Millisecond))
    _, err = pool.Get(ctx2)
    cancel2()

    pool.Put(conn1)
    if err == nil || err.Error() != ErrTimeout.Error() {
        t.Error("get from empty pool succ")
    }
}
