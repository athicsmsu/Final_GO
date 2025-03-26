package dto

// CartItemDto ใช้เก็บข้อมูลของสินค้าที่อยู่ในรถเข็น
type CartItemDto struct {
	CustomerID int    `json:"customer_id"` // ID ของลูกค้าที่เจ้าของรถเข็น
	CartName   string `json:"cart_name"`   // ชื่อของรถเข็น
	ProductID  int    `json:"product_id"`  // ID ของสินค้า
	Quantity   int    `json:"quantity"`    // จำนวนสินค้าในรถเข็น
}

type CartItemDetails struct {
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	TotalPrice  float64 `json:"total_price"`
}

type CartDetails struct {
	CartID   int               `json:"cart_id"`
	CartName string            `json:"cart_name"`
	Items    []CartItemDetails `json:"items"`
}
