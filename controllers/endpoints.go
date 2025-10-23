package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/amghazanfari/pryx/models"
)

type Endpoint struct {
	EndpointService *models.EndpointService
}

type EndpointListResponse struct {
	Object string             `json:"object"`
	Data   *[]models.Endpoint `json:"data"`
}

func (e Endpoint) List(w http.ResponseWriter, r *http.Request) {
	var endpoints *[]models.Endpoint

	w.Header().Set("Content-Type", "application/json")

	endpoints, err := e.EndpointService.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	modelResponse := EndpointListResponse{
		Object: "list",
		Data:   endpoints,
	}

	modelResponseBytes, err := json.Marshal(modelResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(modelResponseBytes)

}
