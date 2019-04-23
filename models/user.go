package models

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Dudobird/dudo-server/config"
	log "github.com/sirupsen/logrus"

	"github.com/Dudobird/dudo-server/utils"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	validator "gopkg.in/go-playground/validator.v9"
)

// Token contains the user authenticate information
type Token struct {
	UserID string
	jwt.StandardClaims
}

// User include user authenticate information
type User struct {
	ID        string    `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `gorm:"DEFAULT:current_timestamp"`
	UpdatedAt time.Time `gorm:"DEFAULT:current_timestamp"`
	DeletedAt *time.Time
	Email     string `json:"email" gorm:"not null;type:varchar(100);unique_index"`
	Password  string `json:"password" gorm:"not null"`
	RoleID    uint
	Token     string `json:"token" sql:"-"`
	// some relation to other modals
	Files   []StorageFile `json:"-"`
	Profile Profile       `json:"-"`
}

// ToJSONBytes will format the accout information to json []byte
// use only in test
func (u *User) ToJSONBytes() []byte {
	return []byte(fmt.Sprintf(`{"email":"%s","password":"%s"}`, u.Email, u.Password))
}

// CheckIfEmailExist return true if email already exist in database
// return error != nil when sever query fail
func (u *User) CheckIfEmailExist() (bool, error) {
	temp := &User{}
	err := GetDB().Table("users").Where("email=?", u.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorln(err)
		return false, errors.New("server unavailable")
	}
	if temp.Email != "" {
		return true, errors.New("email is already in use")
	}
	return false, nil
}

func accountValidate(validate *validator.Validate, field string, value string) error {
	switch field {
	case "email":
		return validate.Var(value, "required,email")
	case "password":
		return validate.Var(value, "required")
	default:
		return nil
	}
}

// Validate will check if the user registe info is correct
func (u *User) Validate() (bool, string) {
	validate := validator.New()
	err := accountValidate(validate, "email", u.Email)
	if err != nil {
		return false, "email address is require"
	}
	err = accountValidate(validate, "password", u.Password)
	if err != nil {
		return false, "password is required"
	}

	return true, "validate success"
}

func (u *User) createToken() (string, error) {
	tokenSecret := config.GetConfig().Application.Token
	token := jwt.NewWithClaims(
		jwt.GetSigningMethod("HS256"),
		&Token{
			UserID: u.ID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			},
		},
	)
	return token.SignedString([]byte(tokenSecret))
}

// Create will valid user infomation and create it
func (u *User) Create() *utils.Message {
	if status, message := u.Validate(); status != true {
		return utils.NewMessage(http.StatusBadRequest, message)
	}
	if status, err := u.CheckIfEmailExist(); status == true || err != nil {
		return utils.NewMessage(http.StatusBadRequest, err.Error())
	}
	// default user level
	if u.RoleID == 0 {
		u.RoleID = UserRoleID
	}
	hashedPasswd, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashedPasswd)
	u.ID = utils.GenRandomID("user", 12)
	// set defaute for profile
	u.Profile = Profile{
		DiskLimit:     utils.GetFileSizeFromReadable(config.GetConfig().Application.DefaultDiskLimit),
		ProfileImage:  config.GetConfig().Application.DefaultProfileImage,
		UsageDiskSize: uint64(0),
		Name:          u.ID,
	}
	err := GetDB().Create(u).Error
	if err != nil {
		log.Errorf("server sql fail for %+v:%s", u, err)
		return utils.NewMessage(http.StatusInternalServerError, "server create account fail")
	}

	token, err := u.createToken()
	if err != nil {
		return utils.NewMessage(http.StatusInternalServerError, "server create account fail")
	}
	u.Token = token
	u.Password = ""
	message := utils.NewMessage(http.StatusCreated, "account create success")
	message.Data = u
	return message
}

// Login will login the user with email and password
// if success, it will save the token and return success message
// or return forbidden etc message
func Login(email, password string) *utils.Message {
	account := &User{}
	tokenSecret := config.GetConfig().Application.Token
	tempAccout := &User{
		Email:    email,
		Password: password,
	}
	// check the input again, if not correct no need for sql query
	if status, message := tempAccout.Validate(); status != true {
		return utils.NewMessage(http.StatusBadRequest, message)
	}
	err := GetDB().Table("users").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.NewMessage(http.StatusNotFound, "email not found")
		}
		log.Errorf("user login fail for %s", email)
		return utils.NewMessage(http.StatusInternalServerError, "server unavailable")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		return utils.NewMessage(http.StatusForbidden, "password not correct")
	}

	account.Password = ""

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &Token{UserID: account.ID})
	tokenString, _ := token.SignedString([]byte(tokenSecret))
	account.Token = tokenString

	message := utils.NewMessage(http.StatusOK, "login success")
	message.Data = account
	return message
}

// Logout user will delete the user token from database
func Logout(userID string) *utils.Message {
	account := &User{}
	err := GetDB().Table("users").Where("id = ?", userID).First(account).Error
	if err == gorm.ErrRecordNotFound {
		return utils.NewMessage(http.StatusNotFound, "user not found")
	}
	account.Token = ""
	return utils.NewMessage(http.StatusOK, "logout user success")
}

// UpdatePassword update the user password
func UpdatePassword(userID string, password, newPassword string) *utils.Message {
	validate := validator.New()
	if err := accountValidate(validate, "password", newPassword); err != nil {
		return utils.NewMessage(http.StatusBadRequest, "new password format error")
	}

	account := &User{}
	err := GetDB().Model(&User{}).Where("id = ?", userID).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.NewMessage(http.StatusNotFound, "user not found")
		}
		log.Errorf("user login fail for %v", err)
		return utils.NewMessage(http.StatusInternalServerError, "server unavailable")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		return utils.NewMessage(http.StatusForbidden, "password not correct")
	}
	hashedPasswd, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	account.Password = string(hashedPasswd)
	err = GetDB().Save(account).Error
	if err != nil {
		log.Errorf("update password fail: %v", err)
		return utils.NewMessage(http.StatusInternalServerError, "server unavailable")
	}
	return utils.NewMessage(http.StatusOK, "update password success")
}

// GetUser return user infomation based userid
func GetUser(userID string) (*User, *utils.CustomError) {
	account := &User{}
	err := GetDB().Table("users").Where("id = ?", userID).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &utils.ErrUserNotFound
		}
		return nil, &utils.ErrInternalServerError
	}
	account.Password = ""
	return account, nil
}

// GetUserWithEmail return user infomation based user email
func GetUserWithEmail(email string) (*User, *utils.CustomError) {
	account := &User{}
	err := GetDB().Table("users").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &utils.ErrUserNotFound
		}
		return nil, &utils.ErrInternalServerError
	}
	account.Password = ""
	return account, nil
}

// InsertAdminUser insert a new admin account
func InsertAdminUser(email string, password string) error {
	account := &User{Email: email, Password: password, RoleID: AdminRoleID}
	message := account.Create()
	if message.Status != http.StatusCreated {
		return errors.New(message.Message)
	}
	return nil
}
