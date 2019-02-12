package models

import (
	"net/http"
	"os"
	"strings"

	"github.com/zhangmingkai4315/dudo-server/utils"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Token contains the user authenticate information
type Token struct {
	UserID uint
	jwt.StandardClaims
}

// Account include user authenticate information
type Account struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

// Validate will check if the user registe info is correct
func (account *Account) Validate() (int, string) {
	if !strings.Contains(account.Email, "@") {
		return http.StatusBadRequest, "email address is require"
	}

	if len(account.Password) < 6 {
		return http.StatusBadRequest, "password is required"
	}
	temp := &Account{}

	err := GetDB().Table("accounts").Where("email=?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return http.StatusServiceUnavailable, "server unavailable"
	}
	if temp.Email != "" {
		return http.StatusBadRequest, "email is already in use"
	}
	return http.StatusOK, "validate success"
}

// Create will valid user infomation and create it
func (account *Account) Create() *utils.Message {
	if status, message := account.Validate(); status != http.StatusOK {
		return utils.NewMessage(status, message)
	}

	hashedPasswd, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPasswd)

	GetDB().Create(account)
	if account.ID <= 0 {
		return utils.NewMessage(http.StatusInternalServerError, "server create account fail")
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &Token{UserID: account.ID})
	tokenString, err := token.SignedString([]byte(os.Getenv("token_password")))
	if err != nil {
		return utils.NewMessage(http.StatusInternalServerError, "server create account fail")
	}
	account.Token = tokenString
	account.Password = ""
	message := utils.NewMessage(http.StatusCreated, "account create success")
	message.Data = account
	return message
}

// Login will login the user with email and password
// if success, it will save the token and return success message
// or return forbidden etc message
func Login(email, password string) *utils.Message {
	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.NewMessage(http.StatusNotFound, "email not found")
		}
		return utils.NewMessage(http.StatusInternalServerError, "server unavailable")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		return utils.NewMessage(http.StatusForbidden, "password not correct")
	}

	account.Password = ""

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &Token{UserID: account.ID})
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	message := utils.NewMessage(http.StatusOK, "login success")
	message.Data = account
	return message
}

// GetUser return user infomation based userid
func GetUser(userID uint) *Account {
	account := &Account{}
	GetDB().Table("accounts").Where("id = ?", userID).First(account)
	if account.Email == "" {
		return nil
	}
	account.Password = ""
	return account
}
