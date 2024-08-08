package orders

import (
	"errors"
	"gotrack/modules/users"

	"gorm.io/gorm"
)

type Repository interface {
	Create(order *Order) error
	GetAll(role string, idUser int, search string, page int, limit int) (result []Order, err error)
	GetByID(id int) (Order, error)
	Delete(id int) error
	Update(order Order, id int, details []OrderDetail) error
	FindEmployee(id int) (*users.User, error)
	IsOrderExists(id int) (bool, error)
	CreateOrderDetails(details []OrderDetail) error
	Delivery(id int) error
	Success(id int, ip string, filename string) error
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

func (o *orderRepository) IsOrderExists(id int) (bool, error) {
	var order Order

	if err := o.db.First(&order, id).Error; err != nil {
		return false, err
	}

	return true, nil
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
	var order Order

	err := o.db.Model(&Order{}).Where("id = ?", id).
		Preload("OrderDetails").
		Preload("Employee").
		First(&order).Error
	if err != nil {
		return Order{}, err
	}

	// Jika status adalah "Success", preload DetailLocation dan Location
	if order.Status == "Success" {
		err = o.db.Preload("DetailLocation.Location").Where("id = ?", id).First(&order).Error
		if err != nil {
			return Order{}, err
		}
	} else {
		order.DetailLocation = nil
	}

	return order, nil
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

// Delivery implements Repository.
func (o *orderRepository) Delivery(id int) error {
	var order Order

	if err := o.db.Select("Status").First(&order, id).Error; err != nil {
		return errors.New("data order tidak ditemukan")
	}

	if order.Status != "Pending" {
		return errors.New("status harus pending")
	}

	if err := o.db.Model(&Order{}).Where("id = ?", id).Update("Status", "Delivery").Error; err != nil {
		return err
	}

	return nil
}

// Success implements Repository.
func (o *orderRepository) Success(id int, ip string, filename string) error {
	var order Order

	if err := o.db.Select("Status").First(&order, id).Error; err != nil {
		return errors.New("data order tidak ditemukan")
	}

	if order.Status != "Delivery" {
		return errors.New("status harus delivery")
	}

	ipInfo, err := getIPInfo(ip)
	if err != nil {
		return errors.New("unable to get IP info")
	}

	return o.db.Transaction(func(tx *gorm.DB) error {
		ipRecord := &users.IPInfo{
			IP:       ipInfo.IP,
			Hostname: ipInfo.Hostname,
			City:     ipInfo.City,
			Region:   ipInfo.Region,
			Country:  ipInfo.Country,
			Loc:      ipInfo.Loc,
			Org:      ipInfo.Org,
			Postal:   ipInfo.Postal,
			Timezone: ipInfo.Timezone,
		}
		if err := tx.Create(ipRecord).Error; err != nil {
			return errors.New("unable to save IP info")
		}

		detailLocation := &users.DetailLocation{
			IpID:    int(ipRecord.ID),
			OrderID: id,
			Pict:    filename,
		}
		if err := tx.Create(detailLocation).Error; err != nil {
			return errors.New("unable to save detail location")
		}

		if err := tx.Model(&Order{}).Where("id = ?", id).Update("Status", "Success").Error; err != nil {
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
