package dbAto

import (
	jwtforreg "ato_chat/JwtForReg"
	
	"ato_chat/models"
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func IsValidEmail(email string) bool {
	// Definisikan ekspresi reguler untuk validasi format email
	// Ekspresi reguler ini memeriksa apakah email memiliki format yang benar
	// Sesuaikan ekspresi reguler sesuai kebutuhan Anda
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Buat objek regex
	re := regexp.MustCompile(regex)

	// Gunakan metode MatchString untuk memeriksa apakah email cocok dengan pola ekspresi reguler
	return re.MatchString(email)
}

// IsEmailTaken checks if an email address is already registered in the database.
func IsEmailTaken(db *sql.DB, email string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
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
func CreateUserAndLogin(db *sql.DB, user models.User) (string, error) {
    // Gunakan db yang diberikan untuk menjalankan query
    _, err := db.Exec("INSERT INTO users (email, password, company_id, name, role_id) VALUES (?, ?, ?, ?, ?)",
        user.Email, user.Password, user.CompanyID, user.Name, user.RoleID)

    if err != nil {
        // Menambahkan detail error untuk membantu mendiagnosis masalah
        return "", fmt.Errorf("failed to create user: %v", err) // Memperbaiki jumlah nilai yang dikembalikan
    }

    // Jika tidak ada error dan pengguna berhasil disimpan:
    // Memperbaiki jumlah argumen yang diperlukan oleh CreateTokenOrSession
    // asumsikan password telah di-hash dan disimpan di user.Password
    token, err := jwtforreg.CreateTokenOrSession(user.Email, user.Password, user.CompanyID)
    if err != nil {
        return "", fmt.Errorf("failed to create authentication token: %v", err)
    }

    // Memperbaiki jumlah nilai yang dikembalikan
    return token, nil
}


// AuthenticateUser authenticates a user based on username and password.
func AuthenticateUser(username, password string) (*models.User, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var user models.User
	err = db.QueryRow("SELECT id, username, email, password FROM users WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err // Password does not match
	}

	return &user, nil
}
