package users

import (
	"gotrack/helpers/common"

	"gorm.io/gorm"
)

type Repository interface {
	Login(user LoginRequest) (result User, err error)
	SignUp(user User) (err error)
	Update(user User, id int, role string) (err error)
	Delete(id int) (err error)
	GetAll(search string, page int, limit int) (users []User, err error)
	FindByID(id uint) (User, error)
	UpdateIPEmployee(userID uint, ipAddress string) error
	TrackEmployeeLocation(userID uint, ipAddress string) (geolocation IPInfo, err error)
}

type userRepository struct {
	db *gorm.DB
}

func NewRepository(database *gorm.DB) Repository {
	return &userRepository{
		db: database,
	}
}

func (r *userRepository) Login(user LoginRequest) (User, error) {
	var result User

	err := r.db.Where("username = ?", user.Username).First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return result, nil // or return a specific error for user not found
		}
		return result, err
	}

	return result, nil
}

func (r *userRepository) SignUp(user User) (err error) {
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Update(user User, id int, role string) (err error) {
	updateData := map[string]interface{}{
		"Username": user.Username,
	}

	if user.Password != "" {
		hashedPassword, err := common.HashPassword(user.Password)
		if err != nil {
			return err
		}
		updateData["Password"] = hashedPassword
	}

	// Only add Role to updateData if the user is an owner
	if role == "owner" {
		updateData["Role"] = user.Role
	}

	if err := r.db.Model(&User{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Delete(id int) (err error) {
	if err := r.db.Where("id = ?", id).Delete(&User{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetAll(search string, page int, limit int) (users []User, err error) {
	var data []User
	query := r.db.Model(&User{})

	if search != "" {
		query = query.Where("username LIKE ? OR role LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	// Paginasi
	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	}

	if err = query.Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *userRepository) FindByID(id uint) (user User, err error) {

	if err = r.db.First(&user, id).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *userRepository) UpdateIPEmployee(userID uint, ipAddress string) error {
	return r.db.Model(&User{}).Where("id = ?", userID).Update("ip", ipAddress).Error
}

// TrackEmployeeLocation implements Repository.
func (r *userRepository) TrackEmployeeLocation(userID uint, ipAddress string) (geolocation IPInfo, err error) {
	panic("unimplemented")
}
