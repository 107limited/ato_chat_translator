package chat

import (
	"ato_chat/dbAto"
	"ato_chat/models"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(username, email, password string) error {
	// Implement logic to insert a new user into the database
	return nil
}

// IsEmailTaken checks if an email address is already registered in the database.
func IsEmailTaken(email string) bool {
	db, err := dbAto.OpenDB()
	if err != nil {
		return true // Assume email is taken if database connection fails
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return true // Assume email is taken if there's an error in the query
	}

	return count > 0
}

// HashPassword hashes a password using bcrypt.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CreateUser creates a new user record in the database.
// func CreateUser(user models.User) error {
// 	db, err := dbAto.OpenDB()
// 	if err != nil {
// 		return err
// 	}
// 	defer db.Close()

// 	hashedPassword, err := HashPassword(user.Password)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = db.Exec("INSERT INTO users (email, password, nama, company_id, role_id) VALUES (?, ?, ?, ?, ?)",
// 		user.Email, hashedPassword, user.Name, user.CompanyID, user.RoleID)

// 	return nil
// }

// AuthenticateUser authenticates a user based on email and password.
func AuthenticateUser(email, password string) (*models.User, error) {
	db, err := dbAto.OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var user models.User
	err = db.QueryRow("SELECT id, email, password FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		return nil, err // User with the given email not found
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err // Password does not match
	}

	return &user, nil
}
