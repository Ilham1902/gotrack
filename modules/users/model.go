package users

import (
	"errors"
	"gotrack/helpers/common"
	"regexp"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string `json:"password"`
	Role     string `json:"role" gorm:"type:varchar(20)"` // "owner" or "employee"
	IP       string `json:"ip"`
}

func (User) TableName() string {
	return "users"
}

type IPInfo struct {
	gorm.Model
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"` // Format: "latitude,longitude"
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
	UserID   uint   `json:"user_id"`
}

func (IPInfo) TableName() string {
	return "ip_info"
}

type DetailLocation struct {
	gorm.Model
	IpID int `json:"ip_id" gorm:"column:ip_id"`

	Location IPInfo `gorm:"foreignKey:IpID; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (DetailLocation) TableName() string {
	return "detail_location"
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l *LoginRequest) ValidateLogin() (err error) {
	if common.IsEmptyField(l.Username) {
		return errors.New("username required")
	}

	if common.IsEmptyField(l.Password) {
		return errors.New("password required")
	}

	return
}

type LoginResponse struct {
	Token string `json:"token"`
}

type SignUpRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	ReTypePassword string `json:"re_type_password"`
	Role           string `json:"role"`
}

func (s *SignUpRequest) ValidateSignUp() (err error) {
	if common.IsEmptyField(s.Username) {
		return errors.New("username required")
	}

	if common.IsEmptyField(s.Password) {
		return errors.New("password required")
	}

	if common.IsEmptyField(s.ReTypePassword) {
		return errors.New("retype password required")
	}

	if common.IsEmptyField(s.Role) {
		return errors.New("role required")
	}

	if s.ReTypePassword != s.Password {
		return errors.New("password mismatch")
	}

	re := regexp.MustCompile(`^(.{8,})$`)
	if !re.MatchString(s.Password) {
		return errors.New("please make sure that the password contains at least 8 characters")
	}

	return nil
}

func (s *SignUpRequest) ConvertToModelForSignUp() (user User, err error) {
	hashedPassword, err := common.HashPassword(s.Password)
	if err != nil {
		err = errors.New("hashing password failed")
		return
	}

	return User{
		Username: s.Username,
		Password: hashedPassword,
		Role:     s.Role,
	}, nil
}
