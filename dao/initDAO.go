package dao

import (
	"dousheng/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

type users struct {
	Id       int64  `gorm:"column:id; not null; type:bigint; primaryKey; autoIncrement; comment:'用户id'"`
	Name     string `gorm:"column:name; not null; type:varchar(255); comment:'用户名'"`
	Password string `gorm:"column:password; not null; type:varchar(255); comment:'用户密码'"`
}

type videos struct {
	Id          int64     `gorm:"column:id; not null; type:bigint; primaryKey; autoIncrement; comment:'视频id'"`
	AuthorId    int64     `gorm:"column:author_id; not null; type:bigint; comment:'作者id'"`
	PlayUrl     string    `gorm:"column:play_url; not null; type:varchar(255); comment:'播放url'"`
	CoverUrl    string    `gorm:"column:cover_url; not null; type:varchar(255); comment:'封面url'"`
	PublishTime time.Time `gorm:"column:publish_time; not null; default:current_timestamp; type:timestamp; comment:'发布时间'"`
	Title       string    `gorm:"column:title; not null; type:varchar(255); comment:'视频标题'"`
}

type favorites struct {
	Id         int64 `gorm:"column:id; not null; type:bigint; primaryKey; autoIncrement; comment:'点赞id'"`
	UserId     int64 `gorm:"column:user_id; not null; type:bigint; comment:'用户id'"`
	VideoId    int64 `gorm:"column:video_id; not null; type:bigint; comment:'用户点赞视频id'"`
	IsFavorite int8  `gorm:"column:is_favorite; not null; type:tinyint; default:1 ;comment:'是否点赞，默认为1'"`
}

type follows struct {
	Id       int64 `gorm:"column:id; not null; type:bigint; primaryKey; autoIncrement; comment:'关注id'"`
	UserId   int64 `gorm:"column:user_id; not null; type:bigint; comment:'用户id'"`
	FollowId int64 `gorm:"column:followed_id; not null; type:bigint; comment:'被关注者id'"`
	IsFollow int8  `gorm:"column:is_follow; not null; type:tinyint; default:1 ;comment:'是否关注，默认为1'"`
}

type comments struct {
	Id          int64     `gorm:"column:id; not null; type:bigint; primaryKey; autoIncrement; comment:'评论id'"`
	UserId      int64     `gorm:"column:user_id; not null; type:bigint; comment:'评论用户id'"`
	VideoId     int64     `gorm:"column:video_id; not null; type:bigint; comment:'用户评论视频id'"`
	CommentText string    `gorm:"column:comment_text; not null; type:varchar(255); comment:'评论内容'"`
	PublishTime time.Time `gorm:"column:comment_time; not null; default:current_timestamp; type:timestamp; comment:'评论时间'"`
	IsComment   int8      `gorm:"column:is_comment; not null; type:tinyint; default:1 ;comment:'是否关注，默认为1'"`
}

func CreateTables() {
	if !DB.Migrator().HasTable(&users{}) {
		DB.AutoMigrate(&users{})
		log.Print("successfully created table-users")
	} else {
		log.Print("table-users existed")
	}

	if !DB.Migrator().HasTable(&videos{}) {
		DB.AutoMigrate(&videos{})
		log.Print("successfully created table-videos")
	} else {
		log.Print("table-videos existed")
	}

	if !DB.Migrator().HasTable(&favorites{}) {
		DB.AutoMigrate(&favorites{})
		log.Print("successfully created table-favorites")
	} else {
		log.Print("table-favorites existed")
	}

	if !DB.Migrator().HasTable(&follows{}) {
		DB.AutoMigrate(&follows{})
		log.Print("successfully created table-follows")
	} else {
		log.Print("table-follows existed")
	}

	if !DB.Migrator().HasTable(&comments{}) {
		DB.AutoMigrate(&comments{})
		log.Print("successfully created table-comments")
	} else {
		log.Print("table-comments existed")
	}
}

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBusername, config.DBpassword, config.DBhost, config.DBport, config.DBname)
	var err error

	slowLogger := logger.New(
		//将标准输出作为Writer
		log.New(os.Stdout, "\r\n", log.LstdFlags),

		logger.Config{
			//设定慢查询时间阈值为1ms
			SlowThreshold: 1 * time.Microsecond,
			//设置日志级别，只有Warn和Info级别会输出慢查询日志
			LogLevel: logger.Warn,
		},
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: slowLogger,
	})

	if err != nil {
		log.Panicln("failed with DB connecting ,error: ", err.Error())
	}

	sqlDB, _ := DB.DB()

	//设置数据库连接池参数
	sqlDB.SetMaxOpenConns(config.MaxConn)     //设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(config.MaxFreeConn) //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。

}

// 获取gorm db对象，其他包需要执行数据库查询的时候，只要通过tools.getDB()获取db对象即可。
// 不用担心协程并发使用同样的db对象会共用同一个连接，db对象在调用他的方法的时候会从数据库连接池中获取新的连接
func GetDB() *gorm.DB {
	return DB
}
