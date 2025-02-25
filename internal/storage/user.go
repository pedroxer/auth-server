package storage

import "github.com/pedroxer/auth-service/internal/models"

func (s *Storage) GetUser(email string) (models.User, error) {
	query := `SELECT user_id, 
       email, 
       password, 
       team, 
       role, 
       position FROM auth.users WHERE email = $1`

	var user models.User
	if err := s.dbConn.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Team, &user.Role, &user.Position); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *Storage) GetUserByID(id int) (models.User, error) {
	query := `SELECT user_id, 
	   email, 
	   password, 
	   team, 
	   role, 
	   position FROM auth.users WHERE user_id = $1`

	var user models.User
	if err := s.dbConn.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Password, &user.Team, &user.Role, &user.Position); err != nil {
		return models.User{}, err
	}
	return user, nil
}
