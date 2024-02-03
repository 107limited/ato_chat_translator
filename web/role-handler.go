package web

import (
    "net/http"
    "encoding/json"
    "ato_chat/dbAto" // Ganti dengan nama paket Anda yang sesuai
)

func (s *Server) GetAllRolesHandler(w http.ResponseWriter, r *http.Request) {
    roles, err := dbAto.GetAllRoles(s.DB)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(roles)
}

// CreateRoleHandler handles HTTP requests for creating new roles.
func (s *Server) CreateRoleHandler(w http.ResponseWriter, r *http.Request) {
    var roleData struct {
        RoleName string `json:"role_name"`
    }
    if err := json.NewDecoder(r.Body).Decode(&roleData); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    roleID, err := dbAto.CreateRole(s.DB, roleData.RoleName)
    if err != nil {
        http.Error(w, "Failed to create role", http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "message": "Role created successfully",
        "role_id": roleID,
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}