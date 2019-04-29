package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Dudobird/dudo-server/utils"
	"github.com/gorilla/mux"

	"github.com/Dudobird/dudo-server/models"
)

const (
	DefaultPageSize = 20
)

type pagination struct {
	Page   int
	Size   int
	Search string
}

func getPaginationInfoFromURL(r *http.Request) (*pagination, error) {
	queryValues := r.URL.Query()
	pageFromQuery := queryValues.Get("page")
	sizeFromQuery := queryValues.Get("size")
	searchFromQuery := queryValues.Get("q")
	var page = 0
	var size = DefaultPageSize
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
	return &pagination{Page: page, Size: size, Search: searchFromQuery}, nil
}

// AdminGetUsers get a list of all users
func AdminGetUsers(w http.ResponseWriter, r *http.Request) {
	// /admin/users?page=xxx&size=xxx
	p, err := getPaginationInfoFromURL(r)
	if err != nil {
		utils.JSONRespnseWithErr(w, err)
		return
	}
	users, err := models.GetUsers(p.Page, p.Size, p.Search)
	if err != nil {
		utils.JSONRespnseWithErr(w, err)
		return
	}
	utils.JSONMessageWithData(w, http.StatusOK, "", users)
}

// AdminDeleteUser delete user from database
func AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := models.DeleteUserWithID(id, true)
	if err != nil {
		utils.JSONRespnseWithErr(w, err)
		return
	}
	utils.JSONMessageWithData(w, http.StatusOK, "", id)
}

func AdminChangeUserStorageLimit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	type ReadableSize struct {
		ReadableSize string `json:"readableSize"`
	}
	readableSize := ReadableSize{}
	err := json.NewDecoder(r.Body).Decode(&readableSize)
	if err != nil {
		utils.JSONRespnseWithErr(w, utils.ErrPostDataNotCorrect)
		return
	}
	err = models.ChangeUserStorageSize(id, readableSize.ReadableSize)
	if err != nil {
		utils.JSONRespnseWithErr(w, err)
		return
	}
	utils.JSONMessageWithData(w, http.StatusOK, "", id)
}

// AdminChangeUserPassword change user password
func AdminChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	type NewPassword struct {
		Password string `json:"password"`
	}
	passwordData := NewPassword{}
	err := json.NewDecoder(r.Body).Decode(&passwordData)
	if err != nil {
		utils.JSONRespnseWithErr(w, utils.ErrPostDataNotCorrect)
		return
	}
	err = models.ChangeUserPassword(id, passwordData.Password)
	if err != nil {
		utils.JSONRespnseWithErr(w, err)
		return
	}
	utils.JSONMessageWithData(w, http.StatusOK, "", id)
}
