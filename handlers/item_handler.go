package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"go-postgres-api/config"
	"go-postgres-api/models"

	"github.com/gin-gonic/gin"
)

// GET /api/items
func GetAllItems(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, name, price, stock, created_at, updated_at FROM items ORDER BY id ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	defer rows.Close()

	items := []models.Item{}
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Stock, &item.CreatedAt, &item.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": items})
}

// GET /api/items/:id
func GetItemByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid item ID"})
		return
	}

	var item models.Item
	query := "SELECT id, name, price, stock, created_at, updated_at FROM items WHERE id = $1"
	err = config.DB.QueryRow(query, id).Scan(&item.ID, &item.Name, &item.Price, &item.Stock, &item.CreatedAt, &item.UpdatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Item not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": item})
}

// POST /api/items
func CreateItem(c *gin.Context) {
	var input models.Item
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": err.Error()})
		return
	}

	query := `INSERT INTO items (name, price, stock) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := config.DB.QueryRow(query, input.Name, input.Price, input.Stock).Scan(&input.ID, &input.CreatedAt, &input.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Item created", "data": input})
}

// PUT /api/items/:id
func UpdateItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid item ID"})
		return
	}

	var input models.UpdateItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": err.Error()})
		return
	}

	query := `UPDATE items 
              SET name = COALESCE($1, name),
                  price = COALESCE($2, price),
                  stock = COALESCE($3, stock),
                  updated_at = $4
              WHERE id = $5`

	res, err := config.DB.Exec(query, input.Name, input.Price, input.Stock, time.Now(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Item updated"})
}

// DELETE /api/items/:id
func DeleteItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid item ID"})
		return
	}

	res, err := config.DB.Exec("DELETE FROM items WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Item deleted"})
}
