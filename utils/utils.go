package utils

import (
    "fmt"
    "io/ioutil"
    "math/rand"
    "crypto/md5"
    "encoding/hex"
    "mime/multipart"
)

// Md5String return md5 value of source string
func Md5String(s string) string {
    h := md5.New()
    h.Write([]byte(s))
    str := hex.EncodeToString(h.Sum(nil))
    return str
}

// GenerateToken return encrypt token of string
func GenerateToken(uname string) string {
    return Md5String(fmt.Sprintf("%s:%d", uname, rand.Intn(999999)))
}

// GetFileSize get size of file
func GetFileSize(f multipart.File) (int, error) {
    content, err := ioutil.ReadAll(f)
    return len(content), err
}

