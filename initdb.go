package main

import (
    "fmt"
    "os"
    "flag"
    "time"
    "utils"
    "conf"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

var config conf.TCPConf
var db *gorm.DB

// User gorm db struct
type User struct {
    ID       int32       `gorm:"type:int(11);primary key"`
    Username string      `gorm:"type:varchar(64);unique;not null"`
    Nickname string      `gorm:"type:varchar(128)"`
    Passwd   string      `gorm:"type:varchar(32);not null"`
    Skey     string      `gorm:"type:varchar(16);not null"`
    Headurl  string      `gorm:"type:varchar(128);unique;not null"`
    Uptime   int64      `gorm:"type:datetime"`
}

// TableName generate tablename
func (u User) TableName() string {
    var value int
    for _, c := range []rune(u.Username) {
        value = value + int(c)
    }
    return fmt.Sprintf("userinfo_tab_%d", value % 20)
}

func getTableName(username string) string {
    var value int
    for _, c := range []rune(username) {
        value = value + int(c)
    }
    return fmt.Sprintf("userinfo_tab_%d", value % 20)
}


func init() {
    // parser config
    var confFile string
    flag.StringVar(&confFile, "c", "conf/tcpserver.yaml", "config file")
    flag.Parse()

    err := conf.ConfParser(confFile, &config)
    if err != nil {
        fmt.Println("parser config failed:", err.Error())
        os.Exit(-1)
    }

    // init db
    conninfo := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", config.Db.User, config.Db.Passwd, config.Db.Host, config.Db.Db)
    db, err = gorm.Open("mysql", conninfo)
    if err != nil {
        fmt.Println("connect to db failed:", err.Error())
        os.Exit(-1)
    }
    db.DB().SetMaxIdleConns(config.Db.Conn.Maxidle)
    db.DB().SetMaxOpenConns(config.Db.Conn.Maxopen)
    db.LogMode(true)
}

func generateSkey() string {
    myTime := fmt.Sprintf("%d", time.Now().Unix())
    str := utils.Md5String(myTime)
    return str[0:6]
}

func createTable() {
    sql := `CREATE TABLE IF NOT EXISTS userinfo_tab_0 (
id INT(11) NOT NULL AUTO_INCREMENT COMMENT 'primary key',
username VARCHAR(64) NOT NULL COMMENT 'unique id',
nickname VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'user nickname, can be empty',
passwd VARCHAR(32) NOT NULL COMMENT 'md5 result of real password and key',
skey VARCHAR(16) NOT NULL COMMENT 'secure key of each user',
headurl VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'user headurl, can be empty',
uptime int(64) NOT NULL DEFAULT 0 COMMENT 'update time: unix timestamp',
PRIMARY KEY(id),
UNIQUE KEY username_unique (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='user info table';`
    db.Exec(sql)
    for i := 1; i < 20; i++ {
        tableName := fmt.Sprintf("userinfo_tab_%d", i)
        db.Exec(fmt.Sprintf("create table if not exists %s like userinfo_tab_0", tableName))
    }
}


func insertRecord() {
    db.LogMode(false)

    var ch chan int64
    var cnum chan int
    maxProcs := 50
    ch = make(chan int64, maxProcs)
    cnum = make(chan int, maxProcs)
    var startTime = time.Now().Unix()
    for i := 0; i < maxProcs; i++ {
        go func(ch chan int64, cnum chan int) {
            var uid int64
            var tableName, username, nickname, skey, password string
            for {
                uid = <-ch
                if uid == 0 {
                    cnum <- 1
                    break
                }
                username = fmt.Sprintf("username%d", uid)
                tableName = getTableName(username)
                nickname = fmt.Sprintf("nickname%d", uid)
                skey = generateSkey()
                password = utils.Md5String(utils.Md5String("123456") + skey)
                db.Table(tableName).Create(&User{Username: username, Nickname: nickname, Passwd: password, Skey: skey, Uptime: time.Now().Unix()})
            }
        }(ch, cnum)
    }
    fmt.Println("Start to create user data,Please wait...")
    totalNum := 10000000
    for i := 1; i <= totalNum; i++ {
        if int64(i)%20000 == 0 {
            fmt.Println(time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("Completed %.1f%%", float64(i*100)/float64(totalNum)))
        }
        ch <- int64(i)
    }

    for i := 0; i < maxProcs; i++ {
        ch <- int64(0)
    }
    for i := 0; i < maxProcs; i++ {
        <-cnum
    }
    var endTime = time.Now().Unix()
    fmt.Println("Done.Cost", endTime-startTime, "s.")
}

func main() {
    createTable()
    insertRecord()
}

