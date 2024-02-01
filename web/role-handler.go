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
