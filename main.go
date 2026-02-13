package main

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Book struct {
	ID     uint    `json:"id" gorm:"primaryKey"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
	Img    string  `json:"img"`
	Desc   string  `json:"desc"`
}

var db *gorm.DB

func main() {
	// DB ulanish
	database, err := gorm.Open(sqlite.Open("books.db"), &gorm.Config{})
	if err != nil {
		panic("Database ulanmadı ❌")
	}
	db = database
	db.AutoMigrate(&Book{})

	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// Static uploads
	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	{
		// Book CRUD
		api.GET("/books", getBooks)
		api.GET("/books/:id", getBook)
		api.POST("/books", createBook)
		api.PUT("/books/:id", updateBook)
		api.DELETE("/books/:id", deleteBook)

		// Image upload
		api.POST("/upload", uploadImage)
	}

	r.Run(":8080")
}

/////////////////////
// HANDLERS
/////////////////////

func getBooks(c *gin.Context) {
	var books []Book
	db.Find(&books)
	c.JSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
	var book Book
	if err := db.First(&book, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Topilmadi"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func createBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&book)
	c.JSON(http.StatusCreated, book)
}

func updateBook(c *gin.Context) {
	var book Book
	if err := db.First(&book, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Topilmadi"})
		return
	}

	var input Book
	c.ShouldBindJSON(&input)

	book.Title = input.Title
	book.Author = input.Author
	book.Price = input.Price
	book.Img = input.Img
	book.Desc = input.Desc

	db.Save(&book)
	c.JSON(http.StatusOK, book)
}

func deleteBook(c *gin.Context) {
	db.Delete(&Book{}, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"message": "O'chirildi"})
}

// Upload Image
func uploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rasm topilmadi"})
		return
	}

	// Papka borligini tekshirish
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		os.Mkdir("./uploads", os.ModePerm)
	}

	savePath := "./uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rasm saqlanmadi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": "https://book-ejyl.onrender.com/uploads/" + file.Filename,
	})
}
