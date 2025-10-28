package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/amghazanfari/pryx/models"
)

type Model struct {
	ModelService *models.ModelService
}

type ModelCreateRequest struct {
	ModelName    string  `json:"model_name"`
	EndpointName string  `json:"endpoint_name"`
	APIKey       string  `json:"api_key"`
	URLAdress    string  `json:"url_address"`
	InputPrice   float32 `json:"input_price,omitempty"`
	OutputPrice  float32 `json:"output_price,omitempty"`
	Active       bool    `json:"active,omitempty"`
}

type ModelCreateResponse struct {
	Message string
	Model   *models.Model
}

type ModelListResponse struct {
	Object string          `json:"object"`
	Data   *[]models.Model `json:"data"`
}

type ModelRetrieveResponse struct {
	*models.Model
}

func (e Model) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var ep ModelCreateRequest
	err = json.Unmarshal(body, &ep)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if ep.ModelName == "" || ep.URLAdress == "" || ep.EndpointName == "" {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	model, err := e.ModelService.Create(ep.ModelName, ep.EndpointName, ep.APIKey, ep.URLAdress, ep.InputPrice, ep.OutputPrice, ep.Active)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	modelResponse := ModelCreateResponse{
		Message: "the model added succesfully",
		Model:   model,
	}

	modelResponseBytes, err := json.Marshal(modelResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(modelResponseBytes)

}

func (e Model) List(w http.ResponseWriter, r *http.Request) {
	var models *[]models.Model

	w.Header().Set("Content-Type", "application/json")

	models, err := e.ModelService.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	modelResponse := ModelListResponse{
		Object: "list",
		Data:   models,
	}

	modelResponseBytes, err := json.Marshal(modelResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(modelResponseBytes)

}

func (e Model) Retrieve(w http.ResponseWriter, r *http.Request) {
	var model *models.Model

	w.Header().Set("Content-Type", "application/json")
	modelName := r.URL.Query().Get("model")

	model, err := e.ModelService.Retrieve(modelName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	modelResponse := ModelRetrieveResponse{
		model,
	}

	modelResponseBytes, err := json.Marshal(modelResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(modelResponseBytes)

}
