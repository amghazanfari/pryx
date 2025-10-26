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

type EndpointDeleteResponse struct {
	Message string `json:"object"`
}

func (e Endpoint) List(w http.ResponseWriter, r *http.Request) {
	var endpoints *[]models.Endpoint

	w.Header().Set("Content-Type", "application/json")

	endpoints, err := e.EndpointService.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	endpointResponse := EndpointListResponse{
		Object: "list",
		Data:   endpoints,
	}

	endpintResponseBytes, err := json.Marshal(endpointResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(endpintResponseBytes)

}

func (e Endpoint) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	err := e.EndpointService.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deleteResponse := EndpointDeleteResponse{
		Message: "the endpoint deleted successfully",
	}

	endpointResponseBytes, err := json.Marshal(deleteResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(endpointResponseBytes)

}
