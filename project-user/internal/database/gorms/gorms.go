package gorms

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"project-user/config"
)

var _db *gorm.DB

func init() {
	//配置MySQL连接参数
	username := config.AppConf.MysqlConfig.Username //账号
	password := config.AppConf.MysqlConfig.Password //密码
	host := config.AppConf.MysqlConfig.Host         //数据库地址，可以是Ip或者域名
	port := config.AppConf.MysqlConfig.Port         //数据库端口
	Dbname := config.AppConf.MysqlConfig.Db         //数据库名
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
	zap.L().Debug("mysql dsn", zap.String("dsn", dsn))
	var err error
	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
}
func GetDB() *gorm.DB {
	return _db
}

type GormConn struct {
	db *gorm.DB
}

func (g *GormConn) Begin() {
	g.db = GetDB().Begin()
}

func New() *GormConn {
	return &GormConn{db: GetDB()}
}
func NewTran() *GormConn {
	return &GormConn{db: GetDB().Begin()}
}
func (g *GormConn) Session(ctx context.Context) *gorm.DB {
	return g.db.Session(&gorm.Session{Context: ctx})
}

func (g *GormConn) RoolBack() {
	g.db.Rollback()
}
func (g *GormConn) Commit() {
	g.db.Commit()
}

func (g *GormConn) Tx(ctx context.Context) *gorm.DB {
	return g.db.WithContext(ctx)
}
