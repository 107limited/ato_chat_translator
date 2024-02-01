package dbAto

import (
    "database/sql"
    "log"

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

