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
	"time"

	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var userData models.RegistrationValidation
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if userData.Email == "" || userData.Password == "" || userData.CompanyID == 0 {
		http.Error(w, "Email, password, and company_id are required fields", http.StatusBadRequest)
		return
	}

	if !dbAto.IsValidEmail(userData.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Asumsikan fungsi CreateToken sudah ada dan mengembalikan token JWT
	token, err := jwtforreg.CreateTokenOrSession(userData.Email, userData.CompanyID)
	if err != nil {
		http.Error(w, "Failed to create token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Menyiapkan respons dengan informasi akun
    response := map[string]interface{}{
        "message": "Registration initiated successfully. Please complete your personal data.",
        "token":   token,
        "account": map[string]interface{}{
            "email":     userData.Email,
            "company_id": userData.CompanyID, // Asumsikan userData memiliki field CompanyID
        },
    }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) PersonalDataHandler(w http.ResponseWriter, r *http.Request) {
    var personalData models.PersonalData
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&personalData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    authToken := r.Header.Get("Authorization")
    email, companyID, err := jwt.ValidateTokenOrSession(authToken)
    if err != nil {
        http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
        return
    }

    if personalData.Name == "" || personalData.RoleID == 0 {
        http.Error(w, "Name and RoleID are required fields", http.StatusBadRequest)
        return
    }

    hashedPassword, err := dbAto.HashPassword("temporary_password") // Ensure you handle passwords correctly
    if err != nil {
        http.Error(w, "Failed to hash password: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Assuming you've included email, companyID in the personalData struct or obtained them from the token
    _, err = s.DB.Exec("INSERT INTO users (email, password, company_id, name, role_id) VALUES (?, ?, ?, ?, ?)",
        email, hashedPassword, companyID, personalData.Name, personalData.RoleID)
    if err != nil {
        http.Error(w, "Failed to save user to database: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Here, replace the token generation with dbAto.GenerateTokenAndLogin as previously done
    token, err := dbAto.GenerateTokenAndLogin(s.DB, email, companyID)
    if err != nil {
        http.Error(w, "Failed to generate token and login: "+err.Error(), http.StatusInternalServerError)
        return
    }
	
	// Dapatkan ID pengguna berdasarkan email dari token
	userId, err := dbAto.GetUserIDByEmail(s.DB, email)
	if err != nil {
		http.Error(w, "Failed to get user ID: "+err.Error(), http.StatusInternalServerError)
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

// handler Get User By Company Name
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

func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Set cookie untuk 'auth_token' dengan nilai kosong dan tanggal kadaluarsa di masa lalu
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,                    // Pastikan ini diaktifkan jika Anda selalu menggunakan HTTPS
		SameSite: http.SameSiteStrictMode, // Menambahkan SameSite untuk keamanan tambahan
	})

	// Kirim respons sukses logout
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\": \"Logout berhasil\"}"))
}
