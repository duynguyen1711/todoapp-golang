package main

import (
	// "encoding/json"

	"log"
	"strconv"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TodoItem struct {
	Id          int        `json:"id" gorm:"id;"`
	Status      string     `json:"status" gorm:"status;"`
	Description string     `json:"description" gorm:"description;"`
	CreatedAt   *time.Time `json:"created_at" gorm:"created_at;"`
	UpdatedAt   *time.Time `json:"updated_at"  gorm:"updated_at;"`
}

func (TodoItem) TableName() string {
	return "todo_lists"
}

type TodoItemCreation struct {
	Id          int    `json:"-" gorm:"column:id;"`
	Title       string `json:"title" gorm:"column:title;"`
	Description string `json:"description" gorm:"column:description;"`
}

func (TodoItemCreation) TableName() string {
	return TodoItem{}.TableName()
}

type TodoItemUpdate struct {
	Title       *string    `json:"title" gorm:"column:title;"`
	Description *string    `json:"description" gorm:"column:description;"`
	Status      *string    `json:"status" gorm:"status;"`
	UpdatedAt   *time.Time `json:"updated_at"  gorm:"updated_at;"`
}

func (TodoItemUpdate) TableName() string {
	return TodoItem{}.TableName()
}
func main() {
	dsn := "root:@tcp(127.0.0.1:3307)/todolist?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.POST("/items", CreateItem(db))
		v1.GET("/items", GetListItem(db))
		v1.GET("/items/:id", GetItemById(db))
		v1.PATCH("/items/:id", UpdateItem(db))
		v1.DELETE("/items/:id", DeleteItem(db))
	}
	r.Run(":3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
func CreateItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data TodoItemCreation
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err := db.Create(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": data.Id,
		})
	}

}

func GetItemById(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data TodoItem
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		data.Id = id
		if err := db.First(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}
func UpdateItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data TodoItemUpdate

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		if err := db.Where("id=?", id).Updates(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})

			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": true,
			"time": data.UpdatedAt,
		})
	}
}
func DeleteItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}
		if err := db.Table(TodoItem{}.TableName()).Where("id = ?", id).Delete(nil).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "deleted",
		})
	}
}
func GetListItem(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var data []TodoItem
		if err := db.Order("id desc").Find(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}
