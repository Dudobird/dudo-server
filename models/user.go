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
	UserID uint
	jwt.StandardClaims
}

// User include user authenticate information
type User struct {
	gorm.Model
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Token    string    `json:"token" sql:"-"`
	Storages []Storage `json:"-"`
}

// ToJSONBytes will format the accout information to json []byte
func (account *User) ToJSONBytes() []byte {
	return []byte(fmt.Sprintf(`{"email":"%s","password":"%s"}`, account.Email, account.Password))
}

// CheckIfEmailExist return true if email already exist in database
// return error != nil when sever query fail
func (account *User) CheckIfEmailExist() (bool, error) {
	temp := &User{}
	err := GetDB().Table("users").Where("email=?", account.Email).First(temp).Error
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
func (account *User) Validate() (bool, string) {
	validate := validator.New()
	err := accountValidate(validate, "email", account.Email)
	if err != nil {
		return false, "email address is require"
	}
	err = accountValidate(validate, "password", account.Password)
	if err != nil {
		return false, "password is required"
	}

	return true, "validate success"
}

func (account *User) createToken() (string, error) {
	tokenSecret := config.GetConfig().Application.Token
	token := jwt.NewWithClaims(
		jwt.GetSigningMethod("HS256"),
		&Token{
			UserID: account.ID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			},
		},
	)
	return token.SignedString([]byte(tokenSecret))

}

// Create will valid user infomation and create it
func (account *User) Create() *utils.Message {
	if status, message := account.Validate(); status != true {
		return utils.NewMessage(http.StatusBadRequest, message)
	}
	if status, err := account.CheckIfEmailExist(); status == true || err != nil {
		return utils.NewMessage(http.StatusBadRequest, err.Error())
	}
	hashedPasswd, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPasswd)

	GetDB().Create(account)
	if account.ID <= 0 {
		log.Errorf("server create account fail for %s", account.Email)
		return utils.NewMessage(http.StatusInternalServerError, "server create account fail")
	}
	token, err := account.createToken()
	if err != nil {
		return utils.NewMessage(http.StatusInternalServerError, "server create account fail")
	}
	account.Token = token
	account.Password = ""
	message := utils.NewMessage(http.StatusCreated, "account create success")
	message.Data = account
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
func Logout(userID uint) *utils.Message {
	account := &User{}
	GetDB().Table("users").Where("id = ?", userID).First(account)
	if account.Email == "" {
		return utils.NewMessage(http.StatusNotFound, "user not found")
	}
	account.Token = ""
	return utils.NewMessage(http.StatusOK, "logout user success")
}

// UpdatePassword update the user password
func UpdatePassword(userID uint, password, newPassword string) *utils.Message {
	validate := validator.New()
	if err := accountValidate(validate, "password", newPassword); err != nil {
		return utils.NewMessage(http.StatusBadRequest, "new password format error")
	}

	account := &User{}
	err := GetDB().Table("users").Where("id = ?", userID).First(account).Error
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
	return utils.NewMessage(http.StatusOK, "update password success")
}

// GetUser return user infomation based userid
func GetUser(userID uint) *User {
	account := &User{}
	err := GetDB().Table("accounts").Where("id = ?", userID).First(account).Error
	if err != nil {
		return nil
	}
	account.Password = ""
	return account
}
