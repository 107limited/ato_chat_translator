package models

type UserRegistration struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Company  int64  `json:"company"`
}

type LoginReqData struct {
	UserID     int    `json:"user_id"`
	Com        string `json:"company"`
	ComId      int64  `json:"companyId"`
	Department string `json:"department"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Auth       int    `json:"auth"`
}

type UserLogin struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	UserID     int    `json:"user_id"`
	Com        string `json:"company"`
	ComId      int64  `json:"companyId"`
	Department string `json:"department"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Auth       int    `json:"auth"`
}

// User mendefinisikan struktur data untuk user.
type User struct {
    ID          int    `json:"id"`
    Email       string `json:"email"`
    Password    string `json:"password,omitempty"` // Sertakan omitempty untuk tidak mengirimkan password dalam response
    CompanyID   int    `json:"company_id,omitempty"`
    RoleID      int64  `json:"role_id,omitempty"`
    Name        string `json:"name,omitempty"`
    CompanyName string `json:"company_name,omitempty"` // Field baru untuk nama perusahaan
    RoleName    string `json:"role_name,omitempty"`    // Field baru untuk nama role
}

// Struct untuk menyimpan data validasi sementara
type RegistrationValidation struct {
	UserId    int    `json:"user_id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Company   string `json:"company"`
	CompanyID int    `json:"company_id"`
	// Anda dapat menambahkan field lain sesuai kebutuhan, seperti pesan kesalahan, dll.
}

type PersonalData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     int64  `json:"role"`
	RoleID   int64  `json:"role_id"`
}
