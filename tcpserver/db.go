package main

import (
    "fmt"
    "conf"
    "time"
    "errors"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

// User gorm user object
type User struct {
    ID       int32       `gorm:"type:int(11);primary key"`
    Username string      `gorm:"type:varchar(64);unique;not null"`
    Nickname string      `gorm:"type:varchar(128)"`
    Passwd   string      `gorm:"type:varchar(32);not null"`
    Skey     string      `gorm:"type:varchar(16);not null"`
    Headurl  string      `gorm:"type:varchar(128);unique;not null"`
    Uptime   int64      `gorm:"type:datetime"`
}

// TableName gorm use this to get tablename
// NOTE : it only works int where caulse
func (u User) TableName() string {
    var value int
    for _, c := range []rune(u.Username) {
        value = value + int(c)
    }
    return fmt.Sprintf("userinfo_tab_%d", value % 20)
}

// for table
func getTableName(username string) string {
    var value int
    for _, c := range []rune(username) {
        value = value + int(c)
    }
    return fmt.Sprintf("userinfo_tab_%d", value % 20)
}

// init conn
func initDbConn(config *conf.TCPConf) error {
    conninfo := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", config.Db.User, config.Db.Passwd, config.Db.Host, config.Db.Db)
    var err error
    db, err = gorm.Open("mysql", conninfo)
    if err != nil {
        msg := fmt.Sprintf("Failed to connect to db '%s', err: %s", conninfo, err.Error())
        return errors.New(msg)
    }

    db.DB().SetMaxIdleConns(config.Db.Conn.Maxidle)
    db.DB().SetMaxOpenConns(config.Db.Conn.Maxopen)
    //db.LogMode(true)
    return nil
}

// cleanup
func closeDB() {
    db.Close()
}

// query
func getDbUserInfo(username string) (User, error) {
   var quser User
   db.Table(getTableName(username)).Where(&User{Username:username}).First(&quser)
   if quser.Username == "" {
       return quser, fmt.Errorf("user(%s) not exists", username)
   }
   return quser, nil
}

// update nickname
func updateDbNickname(username, nickname string) int64 {
   return db.Table(getTableName(username)).Model(&User{}).Where("`username` = ?", username).Updates(User{Nickname: nickname, Uptime: time.Now().Unix()}).RowsAffected
}

// update headurl
func updateDbHeadurl(username, url string) int64 {
   return db.Table(getTableName(username)).Model(&User{}).Where("`username` = ?", username).Updates(User{Headurl: url, Uptime: time.Now().Unix()}).RowsAffected
}

// update nickname and headurl
func updateDbUserinfo(username, nickname, url string) int64 {
   return db.Table(getTableName(username)).Model(&User{}).Where("`username` = ?", username).Updates(User{Nickname: nickname, Headurl: url, Uptime: time.Now().Unix()}).RowsAffected
}

