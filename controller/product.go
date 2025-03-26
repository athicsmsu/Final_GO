package controller

import (
	"FinalGo/dto"
	"FinalGo/model"
	"database/sql"

	"github.com/gin-gonic/gin"
)

// ฟังก์ชันที่ใช้ในการตั้งค่าเส้นทาง (routes) สำหรับสินค้า
func Product(router *gin.Engine) {
	routes := router.Group("/product") // สร้างกลุ่ม route สำหรับ /product
	{
		// เส้นทางสำหรับค้นหาสินค้า
		routes.GET("/search", SearchProducts)

		// เส้นทางสำหรับเพิ่มสินค้า
		routes.POST("/add-to-cart", AddToCart)
	}
}

// Product API: ค้นหาสินค้าจากรายละเอียดสินค้าและช่วงราคา
func SearchProducts(c *gin.Context) {
	// อ่านพารามิเตอร์จาก query string
	description := c.DefaultQuery("description", "")
	minPrice := c.DefaultQuery("min_price", "0")
	maxPrice := c.DefaultQuery("max_price", "10000")

	// สร้างคำสั่ง SQL
	query := `
		SELECT product_id, product_name, description, price, stock_quantity 
		FROM product 
		WHERE description LIKE ? AND price BETWEEN ? AND ?
	`
	db, err := getConnection()
	if err != nil {
		c.JSON(500, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	rows, err := db.Query(query, "%"+description+"%", minPrice, maxPrice)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// จัดเก็บผลลัพธ์สินค้า
	var products []dto.ProductDto
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ProductID, &product.ProductName, &product.Description, &product.Price, &product.StockQuantity); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		products = append(products, dto.ProductDto{
			ProductID:     product.ProductID,
			ProductName:   product.ProductName,
			Description:   product.Description,
			Price:         product.Price,
			StockQuantity: product.StockQuantity,
		})
	}

	c.JSON(200, gin.H{"products": products})
}

// Add to Cart: เพิ่มสินค้าลงในรถเข็น
func AddToCart(c *gin.Context) {
	var cartItem dto.CartItemDto
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// ตรวจสอบว่ามีรถเข็นที่ชื่อที่ต้องการหรือยัง
	db, err := getConnection()
	if err != nil {
		c.JSON(500, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	var cart model.Cart
	query := "SELECT cart_id FROM cart WHERE customer_id = ? AND cart_name = ?"
	err = db.QueryRow(query, cartItem.CustomerID, cartItem.CartName).Scan(&cart.CartID)

	// ถ้าไม่มีรถเข็นให้สร้างใหม่
	if err == sql.ErrNoRows {
		insertCartQuery := "INSERT INTO cart (customer_id, cart_name) VALUES (?, ?)"
		res, err := db.Exec(insertCartQuery, cartItem.CustomerID, cartItem.CartName)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		cartID, err := res.LastInsertId()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		cart.CartID = int(cartID)
	}

	// ตรวจสอบว่ามีสินค้านี้ในรถเข็นแล้วหรือไม่
	var existingItem model.CartItem
	checkQuery := "SELECT cart_item_id, quantity FROM cart_item WHERE cart_id = ? AND product_id = ?"
	err = db.QueryRow(checkQuery, cart.CartID, cartItem.ProductID).Scan(&existingItem.CartItemID, &existingItem.Quantity)

	// ถ้ามีสินค้าแล้วให้เพิ่มจำนวน
	if err == nil {
		updateQuery := "UPDATE cart_item SET quantity = ? WHERE cart_item_id = ?"
		_, err := db.Exec(updateQuery, existingItem.Quantity+cartItem.Quantity, existingItem.CartItemID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Quantity updated in cart"})
		return
	}

	// ถ้าไม่มีสินค้าให้เพิ่มสินค้าลงใน cart_item
	insertQuery := "INSERT INTO cart_item (cart_id, product_id, quantity) VALUES (?, ?, ?)"
	_, err = db.Exec(insertQuery, cart.CartID, cartItem.ProductID, cartItem.Quantity)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Product added to cart"})
}
