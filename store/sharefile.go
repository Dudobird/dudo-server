package store

import (
	"encoding/base64"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/Dudobird/dudo-server/config"
	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/utils"
	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

type fileToken struct {
	FileID string
	UserID string
	jwt.StandardClaims
}

// CreateShareToken create a new share file token
func (store *FileStore) CreateShareToken(fileID string, days int) (string, error) {
	exist := store.StorageFileExistCheck(fileID)
	if exist != true {
		return "", &utils.ErrResourceNotFound
	}
	tokenSecret := config.GetConfig().Application.Token
	token := jwt.NewWithClaims(
		jwt.GetSigningMethod("HS256"),
		&fileToken{
			FileID: fileID,
			UserID: store.userID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * time.Duration(days)).Unix(),
			},
		},
	)
	shareFile := &models.ShareFiles{
		FileID: fileID,
		Expire: time.Now().AddDate(0, 0, days),
		UserID: store.userID,
	}
	err := store.DB.Save(shareFile).Error
	if err != nil {
		log.Errorf("save share file info fail : %s", err)
		return "", &utils.ErrInternalServerError
	}
	t, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Errorf("save share file info fail : %s", err)
		return "", &utils.ErrInternalServerError
	}
	return t, nil
}

// VerifyShareToken check token and return file id if success
func (store *FileStore) VerifyShareToken(token string) (string, string, error) {
	if token == "" {
		return "", "", &utils.ErrTokenIsNotValid
	}
	pathDecodeBase64, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", "", &utils.ErrPostDataNotCorrect
	}

	tokenSecret := config.GetConfig().Application.Token
	pathDecodeBase64Str := string(pathDecodeBase64)
	fileTokenObject := &fileToken{}
	parseToken, err := jwt.ParseWithClaims(pathDecodeBase64Str, fileTokenObject, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil && !parseToken.Valid {
		return "", "", &utils.ErrTokenIsNotValid
	}
	return fileTokenObject.FileID, fileTokenObject.UserID, nil
}

// GetAllSharedFiles get all shared files
func (store *FileStore) GetAllSharedFiles() ([]models.ShareFiles, error) {
	files := []models.ShareFiles{}
	err := store.DB.Model(&models.ShareFiles{}).Where("user_id = ?", store.userID).Find(&files).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return files, err
	}
	return files, nil
}
