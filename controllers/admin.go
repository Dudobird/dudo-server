package controllers

import (
	"net/http"
	"strconv"

	"github.com/Dudobird/dudo-server/utils"

	"github.com/Dudobird/dudo-server/models"
)

const (
	DefaultPageSize = 20
)

type pagination struct {
	Page int
	Size int
}

func getPaginationInfoFromURL(r *http.Request) (*pagination, error) {
	queryValues := r.URL.Query()
	pageFromQuery := queryValues.Get("page")
	sizeFromQuery := queryValues.Get("size")
	var page int = 0
	var size int = DefaultPageSize
	var err error
	if pageFromQuery != "" {
		page, err = strconv.Atoi(pageFromQuery)
		if err != nil {
			return nil, utils.ErrPostDataNotCorrect
		}
	}
	if sizeFromQuery != "" {
		size, err = strconv.Atoi(sizeFromQuery)
		if err != nil {
			return nil, utils.ErrPostDataNotCorrect
		}
	}
	return &pagination{Page: page, Size: size}, nil
}

// GetAdminUsers get a list of all users
func GetAdminUsers(w http.ResponseWriter, r *http.Request) {
	// /admin/users?page=xxx&size=xxx
	p, err := getPaginationInfoFromURL(r)
	if err != nil {
		utils.JSONRespnseWithErr(w, err)
		return
	}
	users, err := models.GetUsers(p.Page, p.Size)
	if err != nil {
		utils.JSONRespnseWithErr(w, err)
		return
	}
	utils.JSONMessageWithData(w, http.StatusOK, "", users)
}
