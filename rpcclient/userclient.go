package rpcclient

import (
    "time"
    "context"
    "strconv"
    "net/http"

    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
    "github.com/gin-gonic/gin"
    "github.com/astaxie/beego/logs"

    codeModule "code"
    pb "proto"
    gpool "gpool"
)

var pool *gpool.GPool

// InitPool  init grpc client connection pool
func InitPool(addr string, init, capacity uint32, maxIdle time.Duration) error {
    // init grpc client pool
    var err error
    pool, err = gpool.NewPool(func () (*grpc.ClientConn, error) {
                     conn, err := grpc.Dial(addr, grpc.WithInsecure())
                     if err != nil {
                         return nil, err
                     }
                     // return pb.NewUserServiceClient(conn), nil
                     return conn, nil
                 },
                 init, capacity, maxIdle)
    return err
}

// DestoryPool destroy connection pool
func DestoryPool() {
    pool.Close()
}

// clientWrap
type clientWrap struct {
    conn   *gpool.Conn
    client pb.UserServiceClient
}

// get client
func getRPCClient() (*clientWrap, error) {
    // get conn
    ctx, cancel := context.WithDeadline(context.Background(),  time.Now().Add(10 * time.Millisecond))
    conn, err := pool.Get(ctx)
    // call cancel to avoid leak
    cancel()

    if err != nil {
        return nil, err
    }
    client := pb.NewUserServiceClient(conn.C)
    return &clientWrap{conn, client}, nil
}

// freeclient
func freeRPCClient(wrap* clientWrap) {
    err := pool.Put(wrap.conn)
    if err != nil {
        logs.Error("Failed to reclaime conn, err:", err.Error())
    }
}

// FormatResponse : suppress golint error
/* {
 *   code: int   // error code, 0 for succ
 *   msg: string // succ or error msg
 *   data: {     // response data object
 *   }
 * }
 */
func FormatResponse(code int, msg string, data map[string]string) map[string]interface{} {
    if msg == "" {
        msg = codeModule.CodeMsg[code]
    }
    return gin.H{"code": code, "msg": msg, "data": data}
}

// Login : userlogin handler
func Login(args map[string]string) (int, string, map[string]interface{}) {
    // get uuid
    uuid := args["uuid"]
    // communicate with rcp server
    client, err := getRPCClient()
    if err != nil {
        logs.Error(uuid, " -- Failed to getRPCClient, err:", err.Error())
        return http.StatusInternalServerError, "", FormatResponse(codeModule.CodeInternalErr, "", nil)
    }
    defer freeRPCClient(client)

    ctx := metadata.AppendToOutgoingContext(context.Background(), "uuid", uuid)
    rsp, err := client.client.Login(ctx, &pb.LoginRequest{Username: args["username"], Passwd: args["passwd"]})
    if err != nil {
        logs.Error(uuid, " -- Failed to communicate with TCP server, err:", err.Error())
        return http.StatusOK, "", FormatResponse(codeModule.CodeErrBackend, "", nil)
    }

    var token string
    logs.Debug(uuid, " -- Succ get token:", rsp.Token, " code:", rsp.Code)
    if rsp.Code == codeModule.CodeSucc && rsp.Token != "" {
        token = rsp.Token
    }
    return http.StatusOK, token, FormatResponse(int(rsp.Code), rsp.Msg, map[string]string{"username":rsp.Username, "nickname":rsp.Nickname, "headurl":rsp.Headurl})
}

// Logout : user logout
func Logout(args map[string]string) (int, map[string]interface{}) {
    // get uuid
    uuid := args["uuid"]
    // communicate with rcp server
    client, err := getRPCClient()
    if err != nil {
        logs.Error(uuid, " -- Failed to getRPCClient, err:", err.Error())
        return http.StatusInternalServerError, FormatResponse(codeModule.CodeInternalErr, "", nil)
    }
    defer freeRPCClient(client)

    ctx := metadata.AppendToOutgoingContext(context.Background(), "uuid", uuid)
    rsp, err := client.client.Logout(ctx, &pb.CommRequest{Token: args["token"], Username: args["username"]})
    if err != nil {
        logs.Error(uuid, " -- Failed to communicate with TCP server, err:", err.Error())
        return http.StatusOK, FormatResponse(codeModule.CodeErrBackend, "", nil)
    }
    logs.Debug(uuid, "Succ to get response from backend with ", rsp.Code, " and msg:", rsp.Msg)
    return http.StatusOK, FormatResponse(int(rsp.Code), rsp.Msg, nil)
}

// EditUserinfo  edit user nickname/headurl
func EditUserinfo(args map[string]string) (int, map[string]interface{}) {
    // get uuid
    uuid := args["uuid"]

    headurl := args["headurl"]
    // get connection
    client, err := getRPCClient()
    if err != nil {
        logs.Error(uuid, " -- Failed to getRPCClient, err:", err.Error())
        return http.StatusInternalServerError, FormatResponse(codeModule.CodeInternalErr, "", nil)
    }
    defer freeRPCClient(client)

    // update userinfo
    mode, _ := strconv.Atoi(args["mode"])
    ctx := metadata.AppendToOutgoingContext(context.Background(), "uuid", uuid)
    editRsp, err := client.client.EditUserInfo(ctx,
                          &pb.EditRequest{Username: args["username"], Token: args["token"], Nickname: args["nickname"], Headurl: headurl, Mode: uint32(mode)})
    if err != nil {
        logs.Error(uuid, " -- Failed to communicate with TCP server, err:", err.Error())
        return http.StatusOK, FormatResponse(codeModule.CodeErrBackend, "", nil)
    }
    data := map[string]string{}
    if editRsp.Code == 0 && headurl != "" {
        data["headurl"] = headurl
    }
    return http.StatusOK, FormatResponse(int(editRsp.Code), editRsp.Msg, data)
}

// GetUserinfo get userinfo handler
func GetUserinfo(args map[string]string) (int, map[string]interface{}) {
    // get uuid
    uuid := args["uuid"]
    // communicate with rcp server
    client, err := getRPCClient()
    if err != nil {
        logs.Error(uuid, " -- Failed to getRPCClient, err:", err.Error())
        return http.StatusInternalServerError, FormatResponse(codeModule.CodeInternalErr, "", nil)
    }
    defer freeRPCClient(client)

    ctx := metadata.AppendToOutgoingContext(context.Background(), "uuid", uuid)
    rsp, err := client.client.GetUserInfo(ctx, &pb.CommRequest{Token: args["token"], Username: args["username"]})
    if err != nil {
        logs.Error(uuid, " -- Failed to communicate with TCP server, err:", err.Error())
        return http.StatusOK, FormatResponse(codeModule.CodeErrBackend, "", nil)
    }
    response := FormatResponse(int(rsp.Code), rsp.Msg, map[string]string{"username":rsp.Username, "nickname":rsp.Nickname, "headurl":rsp.Headurl})
    return http.StatusOK, response
}

// Auth user getUserInfo to auth
func Auth(args map[string]string) (int, int, string) {
    // get uuid
    uuid := args["uuid"]
    // communicate with rcp server
    client, err := getRPCClient()
    if err != nil {
        logs.Error(uuid, " -- Failed to getRPCClient, err:", err.Error())
        return http.StatusInternalServerError, codeModule.CodeInternalErr, codeModule.CodeMsg[codeModule.CodeInternalErr]
    }
    defer freeRPCClient(client)

    ctx := metadata.AppendToOutgoingContext(context.Background(), "uuid", uuid)
    rsp, err := client.client.GetUserInfo(ctx, &pb.CommRequest{Token: args["token"], Username: args["username"]})
    if err != nil {
        logs.Error(uuid, " -- Failed to communicate with TCP server, err:", err.Error())
        return http.StatusOK, codeModule.CodeErrBackend, codeModule.CodeMsg[codeModule.CodeErrBackend]
    }
    if rsp.Code == 0 {
        return http.StatusOK, codeModule.CodeSucc, codeModule.CodeMsg[codeModule.CodeSucc]
    }
    return http.StatusOK, int(rsp.Code), rsp.Msg
}
