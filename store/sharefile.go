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
	ShareID string
	FileID  string
	UserID  string
	jwt.StandardClaims
}

// CreateShareToken create a new share file token
func (store *FileStore) CreateShareToken(fileID string, days int) (string, error) {
	exist := store.StorageFileExistCheck(fileID)
	if exist != true {
		return "", &utils.ErrResourceNotFound
	}
	tokenSecret := config.GetConfig().Application.Token
	id := utils.GenRandomID("share", 10)
	token := jwt.NewWithClaims(
		jwt.GetSigningMethod("HS256"),
		&fileToken{
			ShareID: id,
			FileID:  fileID,
			UserID:  store.userID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * time.Duration(days)).Unix(),
			},
		},
	)
	shareFile := &models.ShareFiles{
		ID:     id,
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

// ShareFileExistCheck  return true when file exist or false if not exist
func (store *FileStore) ShareFileExistCheck(shareID string) bool {
	existCheckStorage := &models.ShareFiles{}
	notFoundChecker := store.DB.Model(&models.ShareFiles{}).Where("id = ?", shareID).First(&existCheckStorage).RecordNotFound()
	if notFoundChecker == false {
		return true
	}
	return false
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
	// check shareid exist or not
	if exist := store.ShareFileExistCheck(fileTokenObject.ShareID); exist == true {
		return fileTokenObject.FileID, fileTokenObject.UserID, nil
	}
	return "", "", &utils.ErrResourceNotFound

}

// GetAllSharedFiles get all shared files
func (store *FileStore) GetAllSharedFiles() ([]models.ShareFiles, error) {
	files := []models.ShareFiles{}
	err := store.DB.Preload("StorageFile").Model(&models.ShareFiles{}).Where("user_id = ?", store.userID).Find(&files).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return files, err
	}
	return files, nil
}

// DeleteShareFilesRef delete share file with id
func (store *FileStore) DeleteShareFilesRef(id string) error {
	err := store.DB.Unscoped().Where("id = ? and user_id = ?", id, store.userID).Delete(&models.ShareFiles{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &utils.ErrResourceNotFound
		}
		return &utils.ErrInternalServerError
	}
	return nil
}
