package controller

import (
	"FinalGo/dto"
	"FinalGo/model"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
)

func getConnection() (*sql.DB, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(viper.Get("mysql.dsn"))
	dsn := viper.GetString("mysql.dsn")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	return db, err
}

func Customer(router *gin.Engine) {
	routes := router.Group("/customer")
	{
		routes.POST("/auth/login", Login)
		routes.PUT("/update-address", UpdateAddress)
		routes.PUT("/update-password", UpdatePassword)
	}
}
func Login(c *gin.Context) {
	customerDto := dto.LoginDto{}

	if err := c.ShouldBindJSON(&customerDto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db, err := getConnection()
	if err != nil {
		c.JSON(500, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	var customer model.Customer
	query := "SELECT customer_id, first_name, last_name, email, phone_number, address, password, created_at, updated_at FROM customer WHERE email = ?"
	err = db.QueryRow(query, customerDto.Email).Scan(
		&customer.CustomerID,
		&customer.FirstName,
		&customer.LastName,
		&customer.Email,
		&customer.PhoneNumber,
		&customer.Address,
		&customer.Password,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(401, gin.H{"message": "Invalid email or password"})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	if customerDto.Password != customer.Password {
		c.JSON(401, gin.H{"message": "Invalid email or password"})
		return
	}

	var customerData dto.CustomerDto
	copier.Copy(&customerData, &customer)

	c.JSON(200, gin.H{"customer": customerData})
}
func UpdateAddress(c *gin.Context) {
	request := dto.UpdateAddressDto{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db, err := getConnection()
	if err != nil {
		c.JSON(500, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	// อัปเดตที่อยู่ของลูกค้า
	query := "UPDATE customer SET address = ? WHERE email = ?"
	result, err := db.Exec(query, request.Address, request.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// ตรวจสอบว่ามีแถวที่ถูกอัปเดตหรือไม่
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"message": "Customer not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Address updated successfully"})
}

func UpdatePassword(c *gin.Context) {
	// รับค่า email, oldPassword, newPassword จากผู้ใช้ และเก็บไว้ในตัวแปร request
	var request = dto.UpdatePasswordDto{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db, err := getConnection()
	if err != nil {
		c.JSON(500, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	// ดึงข้อมูลรหัสผ่านปัจจุบันจากฐานข้อมูล
	var storedPassword string
	query := "SELECT password FROM customer WHERE email = ?"
	err = db.QueryRow(query, request.Email).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"message": "Customer not found"})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	// ตรวจสอบว่ารหัสผ่านเก่าถูกต้องหรือไม่
	if request.OldPassword != storedPassword {
		c.JSON(401, gin.H{"message": "Old password is incorrect"})
		return
	}

	// อัปเดตรหัสผ่านใหม่
	updateQuery := "UPDATE customer SET password = ? WHERE email = ?"
	_, err = db.Exec(updateQuery, request.NewPassword, request.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Password updated successfully"})
}
