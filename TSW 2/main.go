package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

type Articles struct {
	Nickname string `json:"nick_name"`
	Title    string `json:"title"`
	Content  string `json:"Content"`
	Date     string `json:"date"`
}
type Comments struct {
	Nickname string `json:"nick_name"`
	Content  string `json:"content"`
	Comment  string ` json:"comment"`
	Date     string `json:"date"`
}
type CreateArticleInput struct {
	Nickname string `json:"nick_name" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"Content" binding:"required"`
}
type CreateCommentInput struct {
	Nickname string `json:"nick_name" binding:"required"`
	Comment  string `json:"comment" binding:"required"`
	Content  string `json:"Content" binding:"required"`
}

func ConnectDatabase() {
	database, err := gorm.Open("mysql", "root:1505@tcp(127.0.0.1:3306)/articles")

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&Articles{}, &Comments{})

	DB = database
}

// LIST ALL ARTICLES WITH EVERY FIELDS
func ListallArticles(c *gin.Context) {
	var article []Articles
	DB.Find(&article)
	c.JSON(http.StatusOK, gin.H{"data": article})
}

// Get a Article
func FindArticle(c *gin.Context) {
	var task Articles
	title := c.Request.URL.Query().Get("Title")
	if err := DB.Where("title = ?", title).First(&task).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})

}

// POST AN ARTICLE
func CreateArticle(c *gin.Context) {
	var input CreateArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inputs := Articles{Nickname: input.Nickname, Title: input.Title, Content: input.Content, Date: "2023-02-03 00:00:00"}
	DB.Create(&inputs)

	c.JSON(http.StatusOK, gin.H{"data": inputs})

}

// list comments from the article
func ListallComments(c *gin.Context) {
	var comment []Comments
	DB.Find(&comment)
	c.JSON(http.StatusOK, gin.H{"data": comment})
}

// comment on a article
func CreateComments(c *gin.Context) {
	var input CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inputs := Comments{Nickname: input.Nickname, Comment: input.Comment, Content: input.Content, Date: "2023-02-03 00:00:00"}
	DB.Create(&inputs)

	c.JSON(http.StatusOK, gin.H{"data": inputs})

}
func main() {
	ConnectDatabase()
	r := gin.Default()
	//ARTICLE ENDPOINTS
	r.GET("/api/articles", ListallArticles)
	r.GET("/api/articles/:title", FindArticle)
	r.POST("/api/articles/create", CreateArticle)
	//COMMENT ENDPOINTS
	r.GET("/api/comments/", ListallComments)
	r.POST("/api/comments/create", CreateComments)
	r.Run()
}
