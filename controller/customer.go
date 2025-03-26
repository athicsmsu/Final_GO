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
		// routes.GET("/", getAllCustomer)
		routes.POST("/auth/login", Login)
		// routes.GET("/:name", getCustomerName)
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
