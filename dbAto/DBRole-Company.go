package dbAto

import (
	"database/sql"
	"log"
	"strings"
)

// Role mendefinisikan struktur untuk role
type Role struct {
	ID       int    `json:"id"`
	RoleName string `json:"role_name"`
}

// GetAllRoles mengambil semua role dari tabel roles
func GetAllRoles(db *sql.DB) ([]Role, error) {
	var roles []Role

	// Query untuk mengambil semua role
	rows, err := db.Query("SELECT id, role_name FROM roles")
	if err != nil {
		log.Printf("Error when querying roles: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Iterasi melalui hasil query
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role.ID, &role.RoleName); err != nil {
			log.Printf("Error when scanning roles: %v", err)
			return nil, err
		}
		roles = append(roles, role)
	}

	// Periksa error saat iterasi
	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
		return nil, err
	}

	return roles, nil
}

// CreateRole inserts a new role into the database.
func CreateRole(db *sql.DB, roleName string) (int, error) {
	var roleID int
	query := "INSERT INTO roles (role_name) VALUES (?) RETURNING id"
	err := db.QueryRow(query, roleName).Scan(&roleID)
	if err != nil {
		log.Printf("Error creating new role: %v", err)
		return 0, err
	}
	return roleID, nil
}

// Di dalam paket dbAto atau paket yang sesuai
func GetCompanyNameByID(db *sql.DB, companyID int) (string, error) {
	var companyName string
	query := "SELECT company_name FROM companies WHERE id = ?"
	err := db.QueryRow(query, companyID).Scan(&companyName)
	if err != nil {
		return "", err
	}
	return companyName, nil
}

// Di dalam paket dbAto atau paket yang sesuai
func GetCompanyIDByEmail(db *sql.DB, email string) (int, error) {
	var companyID int
	query := "SELECT company_id FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&companyID)
	if err != nil {
		return 0, err
	}
	return companyID, nil
}

// Di dalam paket dbAto atau paket yang sesuai
func GetAdditionalUserInfo(db *sql.DB, email string) (string, string, string, error) {
	var name, companyName, roleName string
	query := `
        SELECT u.name, c.company_name, r.role_name 
        FROM users u 
        LEFT JOIN companies c ON u.company_id = c.id 
        LEFT JOIN roles r ON u.role_id = r.id 
        WHERE u.email = ?`
	err := db.QueryRow(query, email).Scan(&name, &companyName, &roleName)
	if err != nil {
		return "", "", "", err
	}
	return name, companyName, roleName, nil
}

// ToTitleCase konversi string ke title case (huruf pertama kapital setiap kata).
func ToTitleCase(str string) string {
	// Pisahkan string menjadi slice berdasarkan spasi
	words := strings.Fields(str)
	for i, word := range words {
		// Konversi setiap kata: huruf pertama kapital, sisanya huruf kecil
		words[i] = strings.Title(strings.ToLower(word))
	}
	// Gabungkan kembali slice menjadi satu string
	return strings.Join(words, " ")
}

func CheckRoleExistsOrCreate(db *sql.DB, roleName string) (int64, error) {
    var roleID int64
    // Format roleName ke title case
    formattedRoleName := ToTitleCase(roleName)
    
    err := db.QueryRow("SELECT id FROM roles WHERE LOWER(role_name) = LOWER(?)", formattedRoleName).Scan(&roleID)

    if err == sql.ErrNoRows {
        // Jika role tidak ada, masukkan dengan nama yang sudah diformat
        result, err := db.Exec("INSERT INTO roles (role_name) VALUES (?)", formattedRoleName)
        if err != nil {
            return 0, err
        }
        roleID, err = result.LastInsertId()
        if err != nil {
            return 0, err
        }
    } else if err != nil {
        return 0, err
    }

    return roleID, nil
}

