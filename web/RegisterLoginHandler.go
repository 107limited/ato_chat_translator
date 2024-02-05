package web

import (
	jwtforreg "ato_chat/JwtForReg"
	"ato_chat/dbAto"
	"ato_chat/jwt"
	"ato_chat/models"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse data pendaftaran dari permintaan
	var userData models.RegistrationValidation
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validasi data pendaftaran
	if userData.Email == "" || userData.Password == "" {
		http.Error(w, "Email and password are required fields", http.StatusBadRequest)
		return
	}

	// Validasi format email
	if !dbAto.IsValidEmail(userData.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}
	err := s.DB.Ping()
	if err != nil {
		log.Println("Database ping failed:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	emailTaken, err := dbAto.IsEmailTaken(s.DB, userData.Email)
	// Tangani err dan emailTaken sesuai kebutuhan

	if err != nil {
		// Handle error, misalnya dengan mengembalikan HTTP 500
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if emailTaken {
		http.Error(w, "Email address is already registered", http.StatusConflict)
		return
	}

	// Hash password sebelum menyimpan ke database
	hashedPassword, err := dbAto.HashPassword(userData.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Simpan email, password yang di-hash, dan company_id ke database
	_, err = s.DB.Exec("INSERT INTO users (email, password, company_id) VALUES (?, ?, ?)",
		userData.Email, hashedPassword, userData.CompanyID)
	if err != nil {
		http.Error(w, "Failed to save user to database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Buat token JWT (misalnya menggunakan userID yang baru dibuat)
	token, err := jwtforreg.CreateTokenOrSession(userData.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Kirim token serta detail akun yang berhasil didaftarkan sebagai bagian dari respons
	response := map[string]interface{}{
		"message": "User registered successfully, complete personal data!",
		"token":   token,
		"account": map[string]interface{}{
			"email":     userData.Email,     // email dari data pendaftaran
			"companyId": userData.CompanyID, // companyId dari data pendaftaran
			// Tambahkan informasi lain jika diperlukan
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)   // Status 201 menandakan 'Created'
	json.NewEncoder(w).Encode(response) // Mengirimkan respons dalam format JSON

}

func (s *Server) PersonalDataHandler(w http.ResponseWriter, r *http.Request) {
	// Parse data personal data dari permintaan
	var personalData models.PersonalData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&personalData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Mendapatkan dan memvalidasi token dari header
	authToken := r.Header.Get("Authorization") // Sesuaikan dengan cara Anda mengirim token
	email, _, err := jwt.ValidateTokenOrSession(authToken)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Dapatkan ID pengguna berdasarkan email dari token
	userId, err := dbAto.GetUserIDByEmail(s.DB, email)
	if err != nil {
		http.Error(w, "Failed to get user ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Gunakan email dan companyId sesuai kebutuhan

	// Validasi data personal data
	if personalData.Name == "" || personalData.RoleID == 0 {
		http.Error(w, "Name and RoleID are required fields", http.StatusBadRequest)
		return
	}

	// Membuat struct User dengan data yang diperlukan
	user := models.User{
		Name:   personalData.Name,
		RoleID: personalData.RoleID,
	}

	// Update data pengguna dengan name dan role yang baru
	_, err = s.DB.Exec("UPDATE users SET name = ?, role_id = ? WHERE email = ?",
		personalData.Name, personalData.RoleID, email) // 'email' didapat dari token
	if err != nil {
		http.Error(w, "Failed to update user data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Di PersonalDataHandler atau RegisterUserHandler setelah sukses menyimpan data pengguna
	token, err := dbAto.GenerateTokenAndLogin(s.DB, user.Email, user.CompanyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Mendapatkan companyId dari database
	companyId, err := dbAto.GetCompanyIDByEmail(s.DB, email)
	if err != nil {
		http.Error(w, "Failed to get company ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Dapatkan company_name berdasarkan companyId
	companyName, err := dbAto.GetCompanyNameByID(s.DB, companyId)
	if err != nil {
		http.Error(w, "Failed to get company name: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Kirim token serta detail akun yang berhasil diupdate sebagai bagian dari respons
	response := map[string]interface{}{
		"message": "User logged in and updated Personal Data successfully",
		"token":   token,
		"account": map[string]interface{}{
			"id":           userId,              // ID pengguna yang berhasil diupdate
			"name":         personalData.Name,   // Nama dari data personal
			"role_id":      personalData.RoleID, // RoleID dari data personal
			"company_name": companyName,         // Nama perusahaan dari company_id
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (s *Server) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse data login dari permintaan
	var loginData models.UserLogin
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Lakukan autentikasi berdasarkan email dan password
	user, err := dbAto.AuthenticateUser(loginData.Email, loginData.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token untuk pengguna yang berhasil login
	token, err := jwtforreg.CreateToken(user.Email, user.ID)
	if err != nil {
		http.Error(w, "Failed to generate JWT token", http.StatusInternalServerError)
		return
	}
	// Dapatkan informasi tambahan dari database
	name, companyName, roleName, err := dbAto.GetAdditionalUserInfo(s.DB, user.Email)
	if err != nil {
		http.Error(w, "Failed to get additional user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Dapatkan ID pengguna berdasarkan email dari token
	userId, err := dbAto.GetUserIDByEmail(s.DB, user.Email)
	if err != nil {
		http.Error(w, "Failed to get user ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Berikan respons dengan token JWT dan informasi akun
	response := map[string]interface{}{
		"token": token,
		"account": map[string]interface{}{
			"id":           userId,      // ID dari akun yang berhasil login
			"email":        user.Email,  // Email dari akun yang berhasil login
			"name":         name,        // Nama dari akun yang berhasil login
			"company_name": companyName, // Nama perusahaan dari akun yang berhasil login
			"role_name":    roleName,    // Nama role dari akun yang berhasil login
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (s *Server) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := dbAto.GetAllUsers(s.DB)
	if err != nil {
		log.Printf("Failed to get users: %v", err)
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("Failed to encode users to JSON: %v", err)
		http.Error(w, "Failed to process data", http.StatusInternalServerError)
	}
}

func (s *Server) GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "user id is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, err := dbAto.GetUserById(s.DB, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Dalam package web atau package yang sesuai untuk HTTP handler

func (s *Server) GetUsersByCompanyIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyIdStr, ok := vars["companyId"]
	if !ok {
		http.Error(w, "Company ID is required", http.StatusBadRequest)
		return
	}

	companyId, err := strconv.Atoi(companyIdStr)
	if err != nil {
		http.Error(w, "Invalid Company ID", http.StatusBadRequest)
		return
	}

	users, err := dbAto.GetUsersByCompanyId(s.DB, companyId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Fungsi helper untuk menerjemahkan antara companyName dan companyId.
func TranslateCompanyIdentifier(identifier string) (string, error) {
	// Misalnya, Anda memiliki mapping sederhana ini:
	if identifier == "107" {
		return "ATO", nil // Misalkan "ATO" adalah nama untuk company_id 107
	} else if identifier == "ATO" {
		return "107", nil // Misalkan Anda ingin menerjemahkan "ATO" menjadi ID 107
	}
	return "", fmt.Errorf("invalid identifier")
}

//handler Get User By Company Name 
func (s *Server) GetUsersByCompanyIdentifierHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    identifier := vars["companyIdentifier"] // Ini bisa berupa ID atau nama

    // Terjemahkan identifier ke bentuk yang diinginkan.
    translatedIdentifier, err := TranslateCompanyIdentifier(identifier)
    if err != nil {
        http.Error(w, "Invalid company identifier", http.StatusBadRequest)
        return
    }

    var users []models.User
    if translatedIdentifier == "ATO" {
        // Lakukan query berdasarkan company_name jika hasil terjemahannya adalah "ATO"
        users, err = dbAto.GetUsersByCompanyName(s.DB, translatedIdentifier)
    } else {
        // Asumsikan hasil terjemahan adalah "107", lakukan query berdasarkan company_id
        users, err = dbAto.GetUsersByCompanyName(s.DB, translatedIdentifier)
    }

    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}
