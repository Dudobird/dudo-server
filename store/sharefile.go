package store

import (
	"time"

	"github.com/Dudobird/dudo-server/config"
	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/utils"
	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

type fileToken struct {
	FileID string
	jwt.StandardClaims
}

// CreateShareToken create a new share file token
func (store *FileStore) CreateShareToken(fileID string, days int) (string, error) {
	tokenSecret := config.GetConfig().Application.Token
	token := jwt.NewWithClaims(
		jwt.GetSigningMethod("HS256"),
		&fileToken{
			FileID: fileID,
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
	}
	return token.SignedString([]byte(tokenSecret))
}

// VerifyShareToken check token and return file id if success
func (store *FileStore) VerifyShareToken(token string) (string, error) {
	tokenSecret := config.GetConfig().Application.Token
	if token == "" {
		return "", &utils.ErrTokenIsNotValid
	}
	fileTokenObject := &fileToken{}
	parseToken, err := jwt.ParseWithClaims(token, fileTokenObject, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil && !parseToken.Valid {
		return "", &utils.ErrTokenIsNotValid
	}
	return fileTokenObject.FileID, nil
}
