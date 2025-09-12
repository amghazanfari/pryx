package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"gorm.io/gorm"
	"pryx/internal/auth"
	"pryx/internal/models"
)

type createUserReq struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type createKeyReq struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Scopes string `json:"scopes"` // e.g. "completion:invoke,model:write"
}

type createKeyResp struct {
	ID     uint   `json:"id"`
	Prefix string `json:"prefix"`
	Key    string `json:"key"` // only returned once
}

func (h *Handler) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createUserReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"bad json"}`, http.StatusBadRequest); return
		}
		req.Email = strings.TrimSpace(req.Email)
		if req.Email == "" {
			http.Error(w, `{"error":"email required"}`, http.StatusBadRequest); return
		}
		u := models.User{Email: req.Email, Name: strings.TrimSpace(req.Name), IsActive: true}
		if err := h.DB.WithContext(r.Context()).Create(&u).Error; err != nil {
			http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError); return
		}
		w.Header().Set("Content-Type","application/json")
		_ = json.NewEncoder(w).Encode(u)
	}
}

func (h *Handler) CreateAPIKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createKeyReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"bad json"}`, http.StatusBadRequest); return
		}

		var u models.User
		if err := h.DB.WithContext(r.Context()).First(&u, req.UserID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, `{"error":"user not found"}`, http.StatusNotFound); return
			}
			http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError); return
		}
		plain, prefix, hash, err := auth.GenerateAPIKey()
		if err != nil {
			http.Error(w, `{"error":"key gen failed"}`, http.StatusInternalServerError); return
		}
		key := models.APIKey{
			UserID: u.ID, Name: strings.TrimSpace(req.Name),
			Prefix: prefix, Hash: hash, Scopes: strings.TrimSpace(req.Scopes),
		}
		if err := h.DB.WithContext(r.Context()).Create(&key).Error; err != nil {
			http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError); return
		}
		resp := createKeyResp{ID: key.ID, Prefix: key.Prefix, Key: plain}
		w.Header().Set("Content-Type","application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}
