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
