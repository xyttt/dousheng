package config

const VideoNum = 5 // != 0

// SQL相关
const DBusername = "root"
const DBpassword = "root"
const DBhost = "127.0.0.1"
const DBport = 3306
const DBname = "test3"

const MaxFreeConn = 20
const MaxConn = 100
const MaxQueryNumber = 100

//const DBtimeout = "10s"

// MinIO
const Endpoint = "172.27.229.7:9000" //需要改为服务器IP
const AccessKeyID = "minioadmin"
const SecretAccessKey = "minioadmin"
const UseSSL = false

const BucketName = "doushengbuc"
const Location = "cn-north-1"

// MinIO URL

// 时间常量
var OneDayOfHours = 60 * 60 * 24

// Secret 密钥
var Secret = "dousheng"
