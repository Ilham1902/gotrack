package orders

import (
	"gotrack/modules/users"
	"log"

	"gorm.io/gorm"
)

type Repository interface {
	Create(order *Order) error
	GetAll(role string, idUser int, search string, page int, limit int) (result []Order, err error)
	GetByID(id int) (Order, error)
	Delete(id int) error
	Update(order Order, id int, details []OrderDetail) error
	FindEmployee(id int) (*users.User, error)
	CreateOrderDetails(details []OrderDetail) error
	// GetBooks(id int) (Categories, error)
}

type orderRepository struct {
	db *gorm.DB
}

func (o *orderRepository) CreateOrderDetails(details []OrderDetail) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&details).Error; err != nil {
			return err
		}
		return nil
	})
}

func (o *orderRepository) DeleteOrderDetails(id int) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id = ?", id).Delete(&OrderDetail{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (o *orderRepository) FindEmployee(id int) (*users.User, error) {
	var user users.User

	if err := o.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Create implements Repository.
func (o *orderRepository) Create(order *Order) error {
	if err := o.db.Create(order).Error; err != nil {
		return err
	}
	return nil
}

// Delete implements Repository.
func (o *orderRepository) Delete(id int) error {
	var order Order

	if err := o.db.Delete(&order, id).Error; err != nil {
		return err
	}

	o.DeleteOrderDetails(id)

	return nil
}

// GetAll implements Repository.
func (o *orderRepository) GetAll(role string, idUser int, search string, page int, limit int) (result []Order, err error) {
	var data []Order
	query := o.db.Model(&Order{}).Preload("OrderDetails").Preload("Employee")

	if search != "" {
		query = query.Where("customer LIKE ? OR location LIKE ? OR status LIKE ? OR description LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Paginasi
	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	}

	log.Println("Generated SQL:", query.Statement.SQL.String())

	if role == "owner" {
		if err = query.Find(&data).Error; err != nil {
			return nil, err
		}
	} else if role == "employee" {
		if err = query.Where("employee_id = ?", idUser).Find(&data).Error; err != nil {
			return nil, err
		}
	}

	return data, nil
}

// GetByID implements Repository.
func (o *orderRepository) GetByID(id int) (Order, error) {
	panic("unimplemented")
}

// Update implements Repository.
func (o *orderRepository) Update(order Order, id int, details []OrderDetail) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		// Update order
		if err := tx.Model(&Order{}).Where("id = ?", id).Updates(order).Error; err != nil {
			return err
		}

		// Delete existing order details
		if err := tx.Where("order_id = ?", id).Delete(&OrderDetail{}).Error; err != nil {
			return err
		}

		if err := tx.Create(&details).Error; err != nil {
			return err
		}

		return nil
	})
}

func NewRepository(database *gorm.DB) Repository {
	return &orderRepository{
		db: database,
	}
}
