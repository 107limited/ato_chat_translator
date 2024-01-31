package web

import (
	"ato_chat/dbAto"
	"ato_chat/jwt"
	"ato_chat/models"
	"encoding/json"
	"log"

	"net/http"
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

	// Membuat token JWT setelah validasi dan hashing password
	token, err := jwt.CreateTokenOrSession(userData.Email, hashedPassword, userData.CompanyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// validationData := models.RegistrationValidation{
	// 	Email:    userData.Email,
	// 	Password: userData.Password,
	// 	Company:  userData.Company,
	// }

	// Kirim token sebagai respons
	response := map[string]interface{}{"message": "User registered successfully", "token": token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	// Anda dapat menyimpan data validasiData di sini atau mengirimkannya ke langkah berikutnya
	// dalam proses pendaftaran.
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
	token := r.Header.Get("Authorization") // Sesuaikan dengan cara Anda mengirim token
	email, companyId, err := jwt.ValidateTokenOrSession(token)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	

	// Validasi data personal data
	if personalData.Name == "" || personalData.RoleID == 0 {
		http.Error(w, "Name and RoleID are required fields", http.StatusBadRequest)
		return
	}

	hashedPassword, err := dbAto.HashPassword(personalData.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Membuat struct User dengan data yang diperlukan
	user := models.User{
		Email:     email,
		Password:  hashedPassword, // Menggunakan hashedPassword dari token
		CompanyID: companyId,      // Pastikan ini sesuai dengan tipe data di database
		Name:      personalData.Name,
		RoleID:    personalData.RoleID, // Pastikan menggunakan RoleID
	}

	// Menciptakan pengguna dan login
    token, err = dbAto.CreateUserAndLogin(s.DB, user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Kirim token sebagai bagian dari respons jika tidak ada error
    response := map[string]interface{}{
        "message": "User registered and logged in successfully",
        "token":   token, // Token dikirimkan kembali ke klien
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated) // Status 201 menandakan 'Created'
    json.NewEncoder(w).Encode(response) // Mengirimkan token dalam format JSON
}

func (s *Server) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse data login dari permintaan
	var loginData models.UserLogin
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Lakukan autentikasi berdasarkan username dan password
	user, err := dbAto.AuthenticateUser(loginData.Username, loginData.Password)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Load RSA private key from database
	privateKey, err := dbAto.LoadKey()
	if err != nil {
		http.Error(w, "Failed to load RSA private key", http.StatusInternalServerError)
		return
	}

	// Generate JWT token using your custom implementation
	token, err := jwt.GenerateToken(user, privateKey)
	if err != nil {
		http.Error(w, "Failed to generate JWT token", http.StatusInternalServerError)
		return
	}

	// Berikan respons dengan token JWT
	response := map[string]interface{}{"token": token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
