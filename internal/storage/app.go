package storage

import "github.com/pedroxer/auth-service/internal/models"

func (s *Storage) GetApp(appId int) (models.App, error) {
	query := `SELECT app_id, name, secret FROM auth.apps WHERE app_id = $1`

	var app models.App

	if err := s.dbConn.QueryRow(query, appId).Scan(&app.ID, &app.Name, &app.Secret); err != nil {
		return models.App{}, err
	}
	return app, nil
}

func (s *Storage) AddApp(app models.App) error {
	query := `INSERT INTO auth.apps (app_id, name, secret) VALUES ($1, $2, $3)`

	_, err := s.dbConn.Exec(query, app.ID, app.Name, app.Secret)

	return err

}

func (s *Storage) GetAppSecret(appId int) (string, error) {
	query := `SELECT secret FROM auth.apps WHERE app_id = $1`

	var secret string

	if err := s.dbConn.QueryRow(query, appId).Scan(&secret); err != nil {
		return "", err
	}
	return secret, nil
}
