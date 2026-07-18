package dao

import (
	"fmt"
	"sftpbackend/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DB *gorm.DB
)

func InitMariaDB() (err error) {
	// 从配置文件中获取数据库连接信息
	config := config.GlobalConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.Charset,
		config.Query,
	)
	var db *gorm.DB
	// 使用gorm.io/gorm库来连接数据库，并配置命名策略（这里可按需设置其他配置项）
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 禁用外键约束，避免在迁移时触发错误
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: config.Prefix, // 表名前缀，对于“User”表，其表名应为“t_users”。
			// SingularTable: true, // 使用单数表名
		},
	})
	if err != nil {
		return err
	}

	// 测试数据库连接是否正常
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	// 设置连接池参数
	// 最大打开的连接数，根据实际需求调整
	sqlDB.SetMaxOpenConns(100)
	// 最大空闲连接数，根据实际需求调整
	sqlDB.SetMaxIdleConns(20)
	// 连接的最大存活时间，根据实际需求调整，这里设置为1小时
	sqlDB.SetConnMaxLifetime(time.Hour)
	// 空闲连接的最大存活时间，根据实际需求调整，这里设置为30分钟
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	DB = db
	// 测试数据库连接是否正常
	err = sqlDB.Ping()
	return err
}

func Close() {
	// 获取底层的数据库连接对象并关闭连接
	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Println("获取数据库连接对象失败:", err)
		return
	}
	err = sqlDB.Close()
	if err != nil {
		fmt.Println("关闭数据库连接失败:", err)
	}
}
