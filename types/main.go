package types

type MessageFmt struct {
	Id        int    `json:"id"`
	UserId    int    `json:"user_id"`
	ToId      int    `json:"to_id"`
	Message   string `json:"message"`
	Date      int    `json:"date"`
	RoomId    int    `json:"chat_room_id"`
	CompanyId int    `json:"company_id"`
	CreatedAt string `json:"created_at"`
	English   string `json:"english_text"`
	Japanese  string `json:"japanese_text"`
	Speaker   string `json:"speaker"`
}

type IMessage struct {
	Id      int    `json:"id"`
	UserId  int    `json:"user_id"`
	ToId    int    `json:"to_id"`
	Message string `json:"message"`
	Date    int    `json:"date"`
	Token   string `json:"token"`
}

type LastMessage struct {
	English  string `json:"english"`
	Japanese string `json:"japanese"`
	UserID   int    `json:"user_id"`
	Date     int64  `json:"date"`
}