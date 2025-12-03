package lib

import (
	"errors"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 为Gin router 添加handler
func APIBuilder(handlers ...func(*gin.RouterGroup) *gin.RouterGroup) func(gin.IRouter, string) *gin.RouterGroup {
	return func(router gin.IRouter, path string) *gin.RouterGroup {
		group := router.Group(path)
		for _, handler := range handlers {
			group = handler(group)
		}
		return group
	}
}

// 创建数据库连接
func NewDB(config *DatabaseConfig, preProcessor func(*gorm.DB) error) *gorm.DB {
	var db *gorm.DB
	var err error
	switch config.Type {
	case "mysql":
		dsn := config.User + ":" + config.Password + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.DB + "?charset=utf8mb4&parseTime=True&loc=Local"
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.DB), &gorm.Config{})
	}
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if config.Migrate {
		if preProcessor == nil {
			log.Fatalf("Migrator is nil")
		}
		if err = preProcessor(db); err != nil {
			log.Fatalf("Failed to migrate tables: %v", err)
		}
	}
	return db
}

// 配置管理
func LoadConfig[Config any]() *Config {
	if _, err := os.Stat("config.yaml"); err != nil {
		confTmpl := new(Config)
		data, _ := yaml.Marshal(confTmpl)
		os.WriteFile("config.yaml", []byte(data), 0644)
		log.Fatal(errors.New("config file not found, a template file has been created"))
	}
	if err := func() error {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		return viper.ReadInConfig()
	}(); err != nil {
		log.Fatal(errors.New("config file read failed"))
	}
	config := new(Config)
	if err := viper.Unmarshal(config); err != nil {
		log.Fatal(errors.New("config file parse failed"))
	}
	return config
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// finish this golang cors middleware
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Origin")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Next()
	}
}
