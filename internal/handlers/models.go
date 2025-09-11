package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"pryx/internal/models"
	"strings"
)

type createModelRequest struct {
	Name      string `json:"name"`
	ModelName string `json:"model_name"`
	Endpoint  string `json:"endpoint"`
	APIKey    string `json:"api_key"`
}

func (h *Handler) AddModelHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Info("incoming request")

		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "method not supported",
			})
			return
		}
		var model createModelRequest

		if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
			log.Error("invalid request body: ", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid json payload",
			})
			return
		}

		model.Name = strings.TrimSpace(model.Name)
		model.ModelName = strings.TrimSpace(model.ModelName)
		model.Endpoint = strings.TrimSpace(model.Endpoint)

		if model.Name == "" || model.ModelName == "" || model.Endpoint == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid values, name, model name, and endpoint should not be empty"})
			return
		}

		item := models.Model{Name: model.Name, ModelName: model.ModelName, Endpoint: model.Endpoint, APIKey: model.APIKey}
		if err := h.DB.WithContext(r.Context()).Create(&item).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "database error"})
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(item)
	}
}
