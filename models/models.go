package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Model struct {
	ID        string     `json:"-"`
	Name      string     `json:"id"`
	CreatedAt time.Time  `json:"-"`
	Object    string     `json:"object"`
	Timestamp int64      `json:"created_at"`
	OwnedBy   string     `json:"owned_by"`
	Endpoints []Endpoint `json:"-"`
}

type ModelService struct {
	DB *sql.DB
}

func (es *ModelService) Create(modelName, endpointName, apiKey, urlAddress string) (*Model, error) {
	model := Model{
		Name: modelName,
	}

	row := es.DB.QueryRow(`
	WITH ins AS (
		INSERT INTO model (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
		RETURNING id, created_at, owned_by, object
	)
	SELECT * FROM ins
	UNION ALL
	SELECT id, created_at, owned_by, object
	FROM model
	WHERE name = $1
	AND NOT EXISTS (SELECT 1 FROM ins)
	;`, model.Name)

	err := row.Scan(&model.ID, &model.CreatedAt, &model.OwnedBy, &model.Object)
	model.Timestamp = model.CreatedAt.Unix()

	if err != nil {
		return nil, fmt.Errorf("error while create model: %w", err)
	}

	endpoint := Endpoint{
		Name:      endpointName,
		APIKey:    apiKey,
		URLAdress: urlAddress,
	}

	_, err = es.DB.Query(`
	INSERT INTO endpoint (name, api_key, url_address, model_id)
	VALUES ($1, $2, $3, $4)
	`, endpoint.Name, endpoint.APIKey, endpoint.URLAdress, model.ID)

	if err != nil {
		return nil, fmt.Errorf("error while create endpoint: %w", err)
	}

	return &model, nil
}

func (es *ModelService) List() (*[]Model, error) {
	var models []Model

	rows, err := es.DB.Query(`
	SELECT name, created_at, owned_by, object FROM model
	`)

	if err != nil {
		return nil, fmt.Errorf("error getting list of models: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		ep := Model{}
		err = rows.Scan(&ep.Name, &ep.CreatedAt, &ep.OwnedBy, &ep.Object)
		if err != nil {
			return nil, fmt.Errorf("error getting list of models: %w", err)
		}
		ep.Timestamp = ep.CreatedAt.Unix()
		models = append(models, ep)
	}

	return &models, nil
}

func (es *ModelService) Retrieve(modelName string) (*Model, error) {
	var model Model

	row := es.DB.QueryRow(`
	SELECT name, created_at, owned_by, object FROM model
	WHERE name = $1
	`, modelName)
	err := row.Scan(&model.Name, &model.CreatedAt, &model.OwnedBy, &model.Object)
	model.Timestamp = model.CreatedAt.Unix()

	if err != nil {
		return nil, fmt.Errorf("error getting list of models: %w", err)
	}

	return &model, nil
}
