package orders

import (
	"gotrack/modules/users"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	EmployeeID  int    `json:"employee_id" gorm:"column:employee_id"`
	Customer    string `json:"customer"`
	Location    string `json:"location"`
	Status      string `json:"status"` // "pending" or "completed"
	Description string `json:"description"`

	Employee     users.User    `gorm:"foreignKey:EmployeeID; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	OrderDetails []OrderDetail `json:"order_details" gorm:"foreignKey:OrderID"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderRequest struct {
	EmployeeID   int           `json:"employee_id"`
	Customer     string        `json:"customer"`
	Location     string        `json:"location"`
	Status       string        `json:"status"`
	Description  string        `json:"description"`
	OrderDetails []OrderDetail `json:"order_details"`
}

type OrderDetail struct {
	gorm.Model
	OrderID int    `json:"order_id" gorm:"column:order_id"`
	Item    string `json:"item"`
	Qty     int    `json:"qty"`

	// Order Order `gorm:"foreignKey:OrderID; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (OrderDetail) TableName() string {
	return "order_detail"
}

type OrderHistory struct {
	gorm.Model
	OrderID        int    `json:"order_id" gorm:"column:order_id"`
	DetailLocation int    `json:"detail_location"`
	Note           string `json:"note"`
	Status         string `json:"status"` // "completed"

	Order    Order                `gorm:"foreignKey:OrderID; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Location users.DetailLocation `gorm:"foreignKey:DetailLocation; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (OrderHistory) TableName() string {
	return "order_history"
}
