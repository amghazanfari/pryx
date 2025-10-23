package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Endpoint struct {
	ID        string    `json:"-"`
	Name      string    `json:"id"`
	APIKey    string    `json:"-"`
	URLAdress string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	Object    string    `json:"object"`
	Timestamp int64     `json:"created_at"`
	OwnedBy   string    `json:"owned_by"`
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
	SELECT name, api_key, url_address FROM endpoint
	WHERE model_id = $1
	`, model.ID)

	if err != nil {
		return nil, fmt.Errorf("error getting list of endpoints: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		ep := Endpoint{}
		err = rows.Scan(&ep.Name, &ep.APIKey, &ep.URLAdress)
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
	SELECT name, created_at, model_id, url_address, object FROM endpoint
	`)

	if err != nil {
		return nil, fmt.Errorf("error getting list of models: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		ep := Endpoint{}
		var modelID int
		err = rows.Scan(&ep.Name, &ep.CreatedAt, &modelID, &ep.URLAdress, &ep.Object)
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
