package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Endpoint struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	APIKey      string    `json:"-"`
	URLAdress   string    `json:"-"`
	CreatedAt   time.Time `json:"-"`
	Object      string    `json:"object"`
	Timestamp   int64     `json:"created_at"`
	OwnedBy     string    `json:"owned_by"`
	InputPrice  float32   `json:"input_price,omitempty"`
	OutputPrice float32   `json:"output_price,omitempty"`
	Active      bool      `json:"active,omitempty"`
}

type EndpointService struct {
	DB *sql.DB
}

func (es *EndpointService) ListByModel(modelName string) (*[]Endpoint, error) {
	var model Model
	var endpoints []Endpoint

	row := es.DB.QueryRow(`
	SELECT id from model
	WHERE name = $1
	`, modelName)

	err := row.Scan(&model.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting id of models: %w", err)
	}

	rows, err := es.DB.Query(`
	SELECT name, api_key, url_address, input_price, output_price, active FROM endpoint
	WHERE model_id = $1
	`, model.ID)

	if err != nil {
		return nil, fmt.Errorf("error getting list of endpoints: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		ep := Endpoint{}
		err = rows.Scan(&ep.Name, &ep.APIKey, &ep.URLAdress, &ep.InputPrice, &ep.OutputPrice, &ep.Active)
		if err != nil {
			return nil, fmt.Errorf("error getting list of models: %w", err)
		}
		ep.Timestamp = ep.CreatedAt.Unix()
		endpoints = append(endpoints, ep)
	}

	return &endpoints, nil
}

func (es *EndpointService) List() (*[]Endpoint, error) {
	var endpoints []Endpoint

	rows, err := es.DB.Query(`
	SELECT id, name, created_at, model_id, url_address, object, input_price, output_price, active FROM endpoint
	`)

	if err != nil {
		return nil, fmt.Errorf("error getting list of models: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		ep := Endpoint{}
		var modelID int
		err = rows.Scan(&ep.ID, &ep.Name, &ep.CreatedAt, &modelID, &ep.URLAdress, &ep.Object, &ep.InputPrice, &ep.OutputPrice, &ep.Active)
		if err != nil {
			return nil, fmt.Errorf("error getting list of endpoints: %w", err)
		}
		ep.Timestamp = ep.CreatedAt.Unix()
		row := es.DB.QueryRow(`
		SELECT name FROM model
		WHERE id = $1
		`, modelID)
		err = row.Scan(&ep.OwnedBy)
		if err != nil {
			return nil, fmt.Errorf("error getting list of endpoints: %w", err)
		}
		endpoints = append(endpoints, ep)
	}

	return &endpoints, nil
}

func (es *EndpointService) Delete(id string) error {
	_, err := es.DB.Exec(`
	DELETE FROM endpoint
	WHERE id = $1
	`, id)

	if err != nil {
		return fmt.Errorf("error deleting an endpoint: %w", err)
	}

	return nil
}
