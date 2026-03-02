package database

import (
	"log"
	"os"
	"path"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"sun-panel/models/datatype"
	"time"

	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	MYSQL  = "mysql"
	SQLITE = "sqlite"
)

type DbClient interface {
	Connect() (db *gorm.DB, err error)
}

type MySQLConfig struct {
	Username    string
	Password    string
	Host        string
	Port        string
	Database    string
	WaitTimeout int
}

type SQLiteConfig struct {
	Filename string
}

func DbInit(dbClient DbClient) (db *gorm.DB, dbErr error) {
	db, dbErr = dbClient.Connect()
	if dbErr != nil {
		return
	}
	return
}

// Connect mysql连接
func (d *MySQLConfig) Connect() (db *gorm.DB, err error) {
	dsn := d.Username + ":" + d.Password + "@tcp(" + d.Host + ":" + d.Port + ")/" + d.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: GetLogger(),
		NamingStrategy: schema.NamingStrategy{
			// TablePrefix:   "blog_",
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(10)  // SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDb.SetMaxOpenConns(100) // SetMaxOpenConns 设置打开数据库连接的最大数量。
	wait_timeout := d.WaitTimeout
	sqlDb.SetConnMaxLifetime(time.Duration(wait_timeout * int(time.Second))) // SetConnMaxLifetime 设置了连接可复用的最大时间。
	return
}

// Connect sqllite3连接
func (d *SQLiteConfig) Connect() (db *gorm.DB, err error) {
	filePath := d.Filename

	// If database file does not exist, try to seed it from ./seed/database/database.db
	// This provides a safe, deterministic “base template” without overwriting existing user data.
	if ok, _ := cmn.PathExists(filePath); !ok {
		seedDb := path.Join("seed", "database", "database.db")
		if seedOk, _ := cmn.PathExists(seedDb); seedOk {
			// ensure parent dir exists
			_ = os.MkdirAll(path.Dir(filePath), 0700)
			if errCopy := copyFile(seedDb, filePath); errCopy != nil {
				// do not hard fail; fallback to empty db creation
				log.Println("seed db copy failed:", errCopy)
			}
		}
	}

	exists := false
	if exists, err = cmn.PathExists(path.Dir(filePath)); err != nil {
		return
	} else {
		// 创建文件夹
		if !exists {
			if err = os.MkdirAll(path.Dir(filePath), 0700); err != nil {
				return
			}
		}

		db, err = gorm.Open(sqlite.Open(filePath), &gorm.Config{
			Logger: GetLogger(),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	}

	return
}

// 日志
func GetLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Warn, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 彩色打印
		},
	)

}

// 创建数据库
func CreateDatabase(driver string, db *gorm.DB) error {

	// mysql特殊处理
	if driver == MYSQL {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	// 创建数据表
	err := db.AutoMigrate(
		&models.User{},
		&models.SystemSetting{},
		&models.ItemIcon{},
		&models.UserConfig{},
		&models.File{},
		&models.ItemIconGroup{},
		&models.ModuleConfig{},
	)

	return err
}

// 初始化一个用户,一个用户都没有的时候创建一个
func NotFoundAndCreateUser(db *gorm.DB) error {
	fUser := models.User{}
	if err := db.First(&fUser).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		username := "admin"
		fUser.Mail = "admin@example.com"
		fUser.Username = username
		fUser.Name = "管理员"
		fUser.Status = 1
		fUser.Role = 1
		fUser.Password = cmn.PasswordEncryption("admin")

		if errCreate := db.Create(&fUser).Error; errCreate != nil {
			return errCreate
		}
	}

	return nil
}

// 初始化示例分组和网站数据（如果没有分组）
func NotFoundAndCreateExampleData(db *gorm.DB) error {
	var count int64
	db.Model(&models.ItemIconGroup{}).Count(&count)
	if count > 0 {
		return nil // 已有分组，不创建示例
	}

	userID := uint(1) // 默认第一个用户

	// 示例分组数据
	groups := []models.ItemIconGroup{
		{Title: "常用网站", Sort: 1, UserId: userID},
		{Title: "技术博客", Sort: 2, UserId: userID},
		{Title: "视频娱乐", Sort: 3, UserId: userID},
	}
	if err := db.Create(&groups).Error; err != nil {
		return err
	}

	// 示例网站数据（网站模式）
	sites := []models.ItemIcon{
		{Title: "Google", Url: "https://www.google.com", Icon: datatype.ItemIconIconInfo{ItemType: 1, Src: "https://www.google.com/favicon.ico"}, ItemIconGroupId: int(groups[0].ID), UserId: userID, Sort: 1},
		{Title: "GitHub", Url: "https://github.com", Icon: datatype.ItemIconIconInfo{ItemType: 1, Src: "https://github.com/favicon.ico"}, ItemIconGroupId: int(groups[0].ID), UserId: userID, Sort: 2},
		{Title: "Baidu", Url: "https://www.baidu.com", Icon: datatype.ItemIconIconInfo{ItemType: 1, Src: "https://www.baidu.com/favicon.ico"}, ItemIconGroupId: int(groups[0].ID), UserId: userID, Sort: 3},
		{Title: "知乎", Url: "https://www.zhihu.com", Icon: datatype.ItemIconIconInfo{ItemType: 1, Src: "https://www.zhihu.com/favicon.ico"}, ItemIconGroupId: int(groups[1].ID), UserId: userID, Sort: 1},
		{Title: "CSDN", Url: "https://blog.csdn.net", Icon: datatype.ItemIconIconInfo{ItemType: 1, Src: "https://blog.csdn.net/favicon.ico"}, ItemIconGroupId: int(groups[1].ID), UserId: userID, Sort: 2},
		{Title: "YouTube", Url: "https://www.youtube.com", Icon: datatype.ItemIconIconInfo{ItemType: 1, Src: "https://www.youtube.com/favicon.ico"}, ItemIconGroupId: int(groups[2].ID), UserId: userID, Sort: 1},
	}
	if err := db.Create(&sites).Error; err != nil {
		return err
	}

	// 示例网页数据（网页模式 - 信息流）
	webpages := []models.ItemIcon{
		{Title: "科技动态 - 今日头条", Url: "https://www.toutiao.com/", Description: "实时科技资讯", ItemIconGroupId: int(groups[1].ID), UserId: userID, Sort: 10},
		{Title: "编程教程 - 菜鸟教程", Url: "https://www.runoob.com/", Description: "入门编程教程", ItemIconGroupId: int(groups[1].ID), UserId: userID, Sort: 11},
		{Title: "开源项目 - Gitee", Url: "https://gitee.com/", Description: "代码托管平台", ItemIconGroupId: int(groups[0].ID), UserId: userID, Sort: 4},
		{Title: "技术社区 - V2EX", Url: "https://www.v2ex.com/", Description: "程序员讨论社区", ItemIconGroupId: int(groups[1].ID), UserId: userID, Sort: 12},
	}
	if err := db.Create(&webpages).Error; err != nil {
		return err
	}

	return nil
}
