package main

import (
    "fmt"
    "conf"
    "time"
    "errors"
    "encoding/json"
    "github.com/go-redis/redis"
)

var redisConn *redis.Client

const userInfoPrefix = "userinfo_"
const tokenKeyPrefix = "token_"

// init redis connection pool
func initRedisConn(conf *conf.TCPConf) error {
    redisConn = redis.NewClient(&redis.Options{
        Addr : conf.Redis.Addr,
        Password : conf.Redis.Passwd,
        DB : conf.Redis.Db,
        PoolSize : conf.Redis.Poolsize,
    })
    if redisConn == nil {
        return errors.New("Failed to call redis.NewClient")
    }

    _, err := redisConn.Ping().Result()
    if err != nil {
        msg := fmt.Sprintf("Failed to ping redis, err:%s", err.Error())
       return errors.New(msg)
    }
    return nil
}

// cleanup
func closeCache() {
    redisConn.Close()
}

// get cached userinfo
func getUserCacheInfo(username string) ( User, error) {
    redisKey := userInfoPrefix + username
    val, err := redisConn.Get(redisKey).Result()
    var user User
    if err != nil {
        return user, err
    }
    err = json.Unmarshal([]byte(val), &user)
    return user, err
}

// set cached userinfo
func setUserCacheInfo(user User) error {
    redisKey := userInfoPrefix + user.Username
    val, err := json.Marshal(user)
    if err != nil {
        return err
    }
    expired := time.Second * time.Duration(config.Redis.Cache.Userexpired)
    _, err = redisConn.Set(redisKey, val, expired).Result()
    return err
}

// get token info
func getTokenInfo(token string) ( User, error) {
    redisKey := tokenKeyPrefix + token
    val, err := redisConn.Get(redisKey).Result()
    var user User
    if err != nil {
        return user, err
    }
    err = json.Unmarshal([]byte(val), &user)
    return user, err
}

// set cached userinfo
func setTokenInfo(user User, token string) error {
    redisKey := tokenKeyPrefix + token
    val, err := json.Marshal(user)
    if err != nil {
        return err
    }
    expired := time.Second * time.Duration(config.Redis.Cache.Tokenexpired)
    _, err = redisConn.Set(redisKey, val, expired).Result()
    return err
}

// update cached userinfo, if failed, try to delete it from cache
func updateCachedUserinfo(user User) error {
    err := setUserCacheInfo(user)
    if err != nil {
        redisKey := userInfoPrefix + user.Username
        redisConn.Del(redisKey).Result()
    }
    return err
}

// delete token cache info
func delTokenInfo(token string) error {
    redisKey := tokenKeyPrefix + token
    _, err := redisConn.Del(redisKey).Result()
    return err
}

