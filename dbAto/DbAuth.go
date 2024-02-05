package dbAto

import (
	jwtforreg "ato_chat/JwtForReg"
	"ato_chat/config"
	"errors"
	"log"

	"ato_chat/models"
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// GetAllUsers mengembalikan semua pengguna dari database
func GetAllUsers(db *sql.DB) ([]models.User, error) {
    var users []models.User

    query := `
    SELECT users.id, users.email, users.password, users.company_id, users.role_id, users.name,
           companies.company_name AS company_name, roles.role_name AS role_name
    FROM users
    LEFT JOIN companies ON users.company_id = companies.id
    LEFT JOIN roles ON users.role_id = roles.id`
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var u models.User
        // Sesuaikan pemanggilan rows.Scan untuk memasukkan CompanyName dan RoleName
        if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.CompanyID, &u.RoleID, &u.Name, &u.CompanyName, &u.RoleName); err != nil {
            return nil, err
        }
        // Pertimbangkan untuk mengosongkan Password sebelum menambahkannya ke slice jika Anda tidak ingin mengirimkan password kembali dalam respons
        u.Password = ""
        users = append(users, u)
    }

    return users, nil
}

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
func GenerateTokenAndLogin(db *sql.DB, email string, companyID int) (string, error) {
    // Jika diperlukan, lakukan validasi tambahan atau langkah-langkah pra-login di sini

    // Membuat token JWT untuk pengguna yang berhasil login
    token, err := jwtforreg.CreateToken(email, companyID)
    if err != nil {
        return "", fmt.Errorf("failed to create authentication token: %v", err)
    }

    return token, nil
}


// AuthenticateUser melakukan autentikasi pengguna berdasarkan email dan password
func AuthenticateUser(email, password string) (*models.User, error) {
	db, err := config.OpenDB()
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}
	defer db.Close()

	var user models.User
	err = db.QueryRow("SELECT email, password FROM users WHERE email = ?", email).
		Scan(&user.Email, &user.Password)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return nil, err
	}

	// Log email and password yang diambil dari database
	log.Printf("Authenticating user: Email: %s, DB Password: %s", user.Email, user.Password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Password comparison error: %v", err)
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}

// GetUserIDByEmail mengambil ID pengguna dari database berdasarkan email.
func GetUserIDByEmail(db *sql.DB, email string) (int, error) {
    var id int
    query := `SELECT id FROM users WHERE email = ?`
    err := db.QueryRow(query, email).Scan(&id)
    if err != nil {
        return 0, err
    }
    return id, nil
}

// GetUserById mengambil user berdasarkan id dengan informasi company dan role.
func GetUserById(db *sql.DB, userID int) (*models.User, error) {
    var user models.User
    query := `
SELECT users.id, users.email, users.password, users.company_id, users.role_id, users.name,
       companies.company_name AS company_name, roles.role_name AS role_name
FROM users
LEFT JOIN companies ON users.company_id = companies.id
LEFT JOIN roles ON users.role_id = roles.id
WHERE users.id = ?
`
    err := db.QueryRow(query, userID).Scan(
        &user.ID, &user.Email, &user.Password, &user.CompanyID, &user.RoleID, &user.Name,
        &user.CompanyName, &user.RoleName,
    )
    if err != nil {
        return nil, fmt.Errorf("error querying user by ID: %v", err)
    }

    // Hapus nilai password untuk keamanan.
    user.Password = ""

    return &user, nil
}

// Dalam package dbAto atau package yang sesuai untuk akses database

func GetUsersByCompanyId(db *sql.DB, companyId int) ([]models.User, error) {
    var users []models.User

    query := `
    SELECT users.id, users.email, users.company_id, users.role_id, users.name,
           companies.company_name AS company_name, roles.role_name AS role_name
    FROM users
    LEFT JOIN companies ON users.company_id = companies.id
    LEFT JOIN roles ON users.role_id = roles.id
    WHERE users.company_id = ?`
    rows, err := db.Query(query, companyId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var u models.User
        if err := rows.Scan(&u.ID, &u.Email, &u.CompanyID, &u.RoleID, &u.Name, &u.CompanyName, &u.RoleName); err != nil {
            return nil, err
        }
        u.Password = "" // Kosongkan password untuk keamanan
        users = append(users, u)
    }

    return users, nil
}

func GetUsersByCompanyName(db *sql.DB, companyName string) ([]models.User, error) {
    var users []models.User

    // Perbarui query untuk menggunakan `company_name` sebagai kolom untuk join.
    query := `
    SELECT u.id, u.email, u.password, u.company_id, u.role_id, u.name, c.company_name, r.role_name
FROM users u
JOIN companies c ON u.company_id = c.id
LEFT JOIN roles r ON u.role_id = r.id
WHERE c.company_name = ?
`

    rows, err := db.Query(query, companyName)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var u models.User
        // Sesuaikan rows.Scan untuk mencakup company_name dan role_name.
        if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.CompanyID, &u.RoleID, &u.Name, &u.CompanyName, &u.RoleName); err != nil {
            return nil, err
        }
        u.Password = "" // Kosongkan password untuk keamanan
        users = append(users, u)
    }

    return users, nil
}
