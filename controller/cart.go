package controller

import (
	"FinalGo/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ฟังก์ชันสำหรับตั้งค่าเส้นทาง (routes) ของรถเข็น
func Cart(router *gin.Engine) {
	routes := router.Group("/cart")
	{
		routes.GET("/all", GetCustomerCarts) // API สำหรับดูรถเข็นทั้งหมดของลูกค้า
	}
}

// ฟังก์ชันสำหรับดึงข้อมูลรถเข็นทั้งหมดของลูกค้า พร้อมสินค้าในแต่ละคัน
func GetCustomerCarts(c *gin.Context) {
	// รับ `customer_id` จาก query parameter
	customerID := c.Query("customer_id")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_id is required"})
		return
	}

	// เชื่อมต่อฐานข้อมูล
	db, err := getConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	// ดึงข้อมูลรถเข็นทั้งหมดของลูกค้าตาม `customer_id`
	query := `
		SELECT c.cart_id, c.cart_name, p.product_id, p.product_name, p.description, 
		       ci.quantity, p.price, (ci.quantity * p.price) AS total_price
		FROM cart c
		JOIN cart_item ci ON c.cart_id = ci.cart_id
		JOIN product p ON ci.product_id = p.product_id
		WHERE c.customer_id = ?
		ORDER BY c.cart_id
	`
	rows, err := db.Query(query, customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// ใช้ map เพื่อจัดกลุ่มสินค้าตามรถเข็นแต่ละคัน
	cartMap := make(map[int]*dto.CartDetails)

	for rows.Next() {
		var cartID int
		var cartName string
		var item dto.CartItemDetails

		err := rows.Scan(&cartID, &cartName, &item.ProductID, &item.ProductName, &item.Description, &item.Quantity, &item.Price, &item.TotalPrice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ถ้ารถเข็นยังไม่มีใน map ให้สร้างใหม่
		if _, exists := cartMap[cartID]; !exists {
			cartMap[cartID] = &dto.CartDetails{
				CartID:   cartID,
				CartName: cartName,
				Items:    []dto.CartItemDetails{},
			}
		}

		// เพิ่มสินค้าเข้าไปในรายการของรถเข็น
		cartMap[cartID].Items = append(cartMap[cartID].Items, item)
	}

	// แปลง `map` เป็น `slice`
	var carts []dto.CartDetails
	for _, cart := range cartMap {
		carts = append(carts, *cart)
	}

	// ส่งข้อมูล JSON กลับไป
	c.JSON(http.StatusOK, gin.H{"carts": carts})
}
