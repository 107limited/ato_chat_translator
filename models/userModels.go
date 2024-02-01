package models

type UserRegistration struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Company  int64  `json:"company"`
}

type LoginReqData struct {
	UserID     string `json:"userId"`
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
	UserID     string `json:"userId"`
	Com        string `json:"company"`
	ComId      int64  `json:"companyId"`
	Department string `json:"department"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Auth       int    `json:"auth"`
}

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"` // Sertakan omitempty untuk tidak mengirimkan password dalam response
	CompanyID int    `json:"company_id,omitempty"`
	RoleID    int64  `json:"role_id,omitempty"`
	Name      string `json:"name,omitempty"`
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
