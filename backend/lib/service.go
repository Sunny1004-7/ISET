package lib

import (
	"embed"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddCRUD[T any](router gin.IRouter, path string, db *gorm.DB) *gin.RouterGroup {
	return APIBuilder(func(group *gin.RouterGroup) *gin.RouterGroup {
		group.GET("", GetAll[T](db, nil))
		group.GET("/:id", Get[T](db, nil))
		group.POST("", Create[T](db, nil))
		group.PUT("/:id", Update[T](db))
		group.DELETE("/:id", Delete[T](db))
		return group
	})(router, path)
}
func AddCRUDNew[T any](router gin.IRouter, path string, db *gorm.DB, processGet func(*gorm.DB, *gin.Context) *gorm.DB, processGetAll func(*gorm.DB, *gin.Context) *gorm.DB, processCreate func(*gorm.DB, *gin.Context) *gorm.DB) *gin.RouterGroup {
	return APIBuilder(func(group *gin.RouterGroup) *gin.RouterGroup {
		group.GET("", GetAll[T](db, nil))
		group.GET("/:id", Get[T](db, nil))
		group.POST("", Create[T](db, nil))
		group.PUT("/:id", Update[T](db))
		group.DELETE("/:id", Delete[T](db))
		return group
	})(router, path)
}
func AddStatic(router *gin.Engine, staticFileDir []string) {
	for _, dir := range staticFileDir {
		router.NoRoute(gin.WrapH(http.FileServer(http.Dir(dir))))
	}
}
func AddStaticFS(router *gin.Engine, fs embed.FS) {
	router.NoRoute(gin.WrapH(http.FileServer(http.FS(fs))))
}

// handlers for gorm
func Create[T any](db *gorm.DB, process func(*gorm.DB, *T) *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		d := new(T)
		if err := c.ShouldBindJSON(d); err != nil {
			c.AbortWithStatus(404)
			log.Println("[gorm]parse creation data failed: ", err)
		} else {
			if process != nil {
				if process(db, d).Error != nil {
					c.AbortWithStatus(404)
					log.Println("[gorm] create data process failed: ", err)
				}
			} else if err := db.Create(d).Error; err != nil {
				c.AbortWithStatus(404)
				log.Println("[gorm] create data failed: ", err)
			}
			c.JSON(200, d)
		}
	}
}
func Get[T any](db *gorm.DB, process func(*gorm.DB, *gin.Context) *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, d := c.Param("id"), new(T)
		if process != nil {
			if process(db, c).First(d).Error != nil {
				c.AbortWithStatus(404)
				log.Println("[gorm] db query process failed")
			}
		} else if err := db.Where("id = ?", id).First(d).Error; err != nil {
			c.AbortWithStatus(404)
			log.Println(err)
		}
		c.JSON(200, d)
	}
}
func GetAll[T any](db *gorm.DB, process func(*gorm.DB, *gin.Context) *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		d := new([]T)
		if process != nil {
			if process(db, c).Find(d).Error != nil {
				c.AbortWithStatus(404)
				log.Println("[gorm] db query all process failed")
			}
		} else if err := db.Find(d).Error; err != nil {
			c.AbortWithStatus(404)
			log.Println(err)
		}
		c.JSON(200, d)
	}
}
func Update[T any](db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		d := new(T)
		if err := c.ShouldBindJSON(d); err != nil {
			c.AbortWithStatus(404)
			log.Println("[gorm]parse update data failed: ", err)
		} else {
			if err := db.Save(&d).Error; err != nil {
				c.AbortWithStatus(404)
				log.Println(err)
			} else {
				c.JSON(200, d)
			}
		}
	}
}
func Delete[T any](db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		var d T
		if err := db.Where("id = ?", id).Delete(&d).Error; err != nil {
			c.AbortWithStatus(404)
			log.Println(err)
		} else {
			c.JSON(200, d)
		}
	}
}
