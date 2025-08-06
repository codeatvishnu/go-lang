package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// Item represents a single data entry.
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var (
	items = make(map[string]Item)
	mu    sync.RWMutex
)

func main() {
	// Create a new Gin router.
	router := gin.Default()
  router.GET("/health", func(c *gin.Context){
    c.JSON(200,gin.H{"message":"ok"})
  })

	// Define the API routes.
	router.GET("/items", getItems)
	router.POST("/items", postItem)
	router.DELETE("/items/:id", deleteItem)

	// Run the server on port 8080.
	router.Run(":8080")
}

// getItems retrieves all items.
func getItems(c *gin.Context) {
	mu.RLock()
	defer mu.RUnlock()

	var allItems []Item
	for _, item := range items {
		allItems = append(allItems, item)
	}

	c.JSON(http.StatusOK, allItems)
}

// postItem creates a new item.
func postItem(c *gin.Context) {
	var newItem Item
	if err := c.BindJSON(&newItem); err != nil {
		return // Gin will automatically send a 400 Bad Request error.
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := items[newItem.ID]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Item with this ID already exists"})
		return
	}

	items[newItem.ID] = newItem
	c.JSON(http.StatusCreated, newItem)
}

// deleteItem deletes an item by its ID.
func deleteItem(c *gin.Context) {
	id := c.Param("id")

	mu.Lock()
	defer mu.Unlock()

	if _, exists := items[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	delete(items, id)
	c.Status(http.StatusNoContent)
}