package dto

// ProductDto ใช้เก็บข้อมูลของสินค้า
type ProductDto struct {
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name"`
	Description   string `json:"description"`
	Price         string `json:"price"`
	StockQuantity int    `json:"stock_quantity"`
}
