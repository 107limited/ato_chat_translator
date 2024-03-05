package web

import (
	jwtforreg "ato_chat/JwtForReg"
	"ato_chat/dbAto"
	"sync"

	//"ato_chat/jwt"
	"ato_chat/models"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"

	"net/http"

	"github.com/gorilla/mux"
)

// generateUniqueKey generates a unique key using UUID
func generateUniqueKey() string {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("Failed to generate unique key: %v", err)
	}
	return newUUID.String()
}

var temporaryStorage = make(map[string]string)

// saveToTemporaryStorage saves a value with a unique key to temporary storage
func saveToTemporaryStorage(key, value string) error {
	// Simpan value dengan key ke dalam map
	// Dalam implementasi nyata, Anda mungkin ingin mempertimbangkan untuk menangani kasus di mana key sudah ada
	temporaryStorage[key] = value
	return nil
}

// getFromTemporaryStorage retrieves a value by its key from temporary storage
func getFromTemporaryStorage(key string) (string, error) {
	// Dapatkan value berdasarkan key dari map
	value, exists := temporaryStorage[key]
	if !exists {
		return "", fmt.Errorf("key not found in temporary storage")
	}
	return value, nil
}

func (s *Server) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {

	var logEntry = logrus.WithFields(logrus.Fields{
		"handler": "RegisterUserHandler",
	})

	var userData models.RegistrationValidation
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logEntry.WithError(err).Error("Failed to decode user data")
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

	if !dbAto.IsValidEmail(userData.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Periksa apakah email sudah diambil
	emailTaken, err := dbAto.IsEmailTaken(s.DB, userData.Email)
	if err != nil {
		// Gagal melakukan query ke database
		http.Error(w, "Failed to check email availability: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if emailTaken {
		// Email sudah diambil
		http.Error(w, "Email address is already registered", http.StatusConflict)
		return
	}

	// Hash password pengguna
	hashedPassword, err := dbAto.HashPassword(userData.Password) // Gunakan password asli dari pengguna
	if err != nil {
		http.Error(w, "Failed to hash password: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Simpan hashed password ke tempat penyimpanan sementara
	hashedPasswordKey := generateUniqueKey()                        // Fungsi untuk menghasilkan kunci unik
	err = saveToTemporaryStorage(hashedPasswordKey, hashedPassword) // Pseudocode
	if err != nil {
		http.Error(w, "Failed to save hashed password: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Asumsikan fungsi CreateTokenOrSession dimodifikasi untuk menerima hashedPasswordKey
	token, err := jwtforreg.CreateTokenOrSessionWithHashedPasswordKey(userData.Email, userData.CompanyID, hashedPasswordKey)
	if err != nil {
		http.Error(w, "Failed to create token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// After successful token creation and before returning success response
	logEntry.WithFields(logrus.Fields{
		"email":     userData.Email,
		"companyID": userData.CompanyID,
		"token":     token, // Assuming token variable is available here
	}).Info("Registration successful. Token generated for updating personal data.")

	// Menyiapkan respons dengan informasi akun
	response := map[string]interface{}{
		"message": "Registration initiated successfully. Please complete your personal data.",
		"token":   token,
		"account": map[string]interface{}{
			"email":      userData.Email,
			"company_id": userData.CompanyID, // Asumsikan userData memiliki field CompanyID
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Personal Data Handler

func (s *Server) PersonalDataHandler(w http.ResponseWriter, r *http.Request) {
	var personalData models.PersonalData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&personalData); err != nil {
		log.Printf("Error decoding personal data: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authToken := r.Header.Get("Authorization")
	email, companyID, hashedPasswordKey, err := jwtforreg.ValidateTokenAndGetHashedPasswordKey(authToken)
	if err != nil {
		log.Printf("Invalid token: %v", err)
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Gunakan hashedPasswordKey untuk mengambil hashed password dari tempat penyimpanan sementara
	hashedPassword, err := getFromTemporaryStorage(hashedPasswordKey)
	if err != nil {
		log.Printf("Failed to retrieve hashed password: %v", err)
		http.Error(w, "Failed to retrieve hashed password: "+err.Error(), http.StatusInternalServerError)
		return
	}

	roleID, err := dbAto.CheckRoleExistsOrCreate(s.DB, personalData.RoleName)
	if err != nil {
		log.Printf("Failed to process role: %v", err)
		http.Error(w, "Failed to process role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = s.DB.Exec("INSERT INTO users (email, password, company_id, name, role_id) VALUES (?, ?, ?, ?, ?)",
		email, hashedPassword, companyID, personalData.Name, roleID)
	if err != nil {
		log.Printf("Failed to save user to database: %v", err)
		http.Error(w, "Failed to save user to database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := dbAto.GenerateTokenAndLogin(s.DB, email, companyID)
	if err != nil {
		log.Printf("Failed to generate token and login: %v", err)
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

	// Assuming success at this point
	log.Printf("User %v updated personal data and logged in successfully", email)

	// Kirim token serta detail akun yang berhasil diupdate sebagai bagian dari respons
	response := map[string]interface{}{
		"message": "User logged in and updated Personal Data successfully",
		"token":   token,
		"account": map[string]interface{}{
			"id":           userId,                // ID pengguna yang berhasil diupdate
			"name":         personalData.Name,     // Nama dari data personal
			"role_name":    personalData.RoleName, // Gunakan freetext untuk nama role
			"company_name": companyName,           // Nama perusahaan dari company_id
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// Misalkan Anda memiliki ConnectionManager yang sudah diperluas seperti ini:
type ConnectionManager struct {
	Connections     map[string]map[*websocket.Conn]struct{}
	UserConnections map[int]*websocket.Conn // Mapping dari userID ke WebSocket connection
	mu              sync.Mutex
}

// Fungsi untuk menandai pengguna sebagai online dan menyimpan koneksi WebSocket mereka
func (manager *ConnectionManager) SetUserOnline(userID int, conn *websocket.Conn) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.UserConnections[userID] = conn
}

// Fungsi untuk menandai pengguna sebagai offline
func (manager *ConnectionManager) SetUserOffline(userID int) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	delete(manager.UserConnections, userID)
}

func (s *Server) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse data login dari permintaan
	var loginData models.UserLogin
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log sebelum memanggil fungsi AuthenticateUser
	log.Printf("Attempting to authenticate user: Email: %s", loginData.Email)

	// Lakukan autentikasi berdasarkan email dan password
	user, err := dbAto.AuthenticateUser(loginData.Email, loginData.Password)
	if err != nil {
		// Log jika terjadi error selama autentikasi
		log.Printf("Authentication failed for user: Email: %s, Error: %v", loginData.Email, err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Log sukses autentikasi
	log.Printf("User authenticated successfully: Email: %s", user.Email)

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
	// Misalnya, Anda memiliki mapping yang lebih kompleks atau query database untuk mendapatkan data yang sesuai.
	// Contoh sederhana dengan mapping hard-coded:
	mapping := map[string]string{
		"107": "ATO",
		"ATO": "107",
	}

	// Cek apakah identifier ada dalam mapping.
	if translated, ok := mapping[identifier]; ok {
		return translated, nil
	}

	// Jika identifier tidak dikenali, kembalikan error yang lebih informatif.
	return "", fmt.Errorf("invalid identifier: %s is not recognized as a valid company name or ID", identifier)
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
		log.Printf("Error getting users by company identifier: %v", err) // Tambahkan ini
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
