package dbAto

import (
	"ato_chat/models" // Adjust the import path based on your project structure.
	"database/sql"
	"fmt"
	"log"
	"time"
)

// GetAllChatRooms retrieves all chat rooms from the database.
func GetAllChatRooms(db *sql.DB) ([]models.ChatRoom, error) {
	var chatRooms []models.ChatRoom

	query := `SELECT id, user1_id, user2_id, created_at FROM chat_room ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying chat rooms: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cr models.ChatRoom
		var createdAtStr string // Use this to scan the timestamp
		if err := rows.Scan(&cr.ID, &cr.User1ID, &cr.User2ID, &createdAtStr); err != nil {
			log.Printf("Error scanning chat room: %v", err)
			return nil, err
		}

		// Convert the string timestamp to time.Time
		parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Error parsing created_at timestamp: %v", err)
			return nil, err
		}
		cr.CreatedAt = parsedTime // Assign the parsed time to the struct

		chatRooms = append(chatRooms, cr)
	}

	return chatRooms, nil
}

// GetChatRoomById retrieves details of a chat room by its ID.
func GetChatRoomById(db *sql.DB, chatRoomId int) (*models.ChatRoom, error) {
	var chatRoom models.ChatRoom
	var createdAtString string // Temporarily store the timestamp as a string

	query := `SELECT id, user1_id, user2_id, created_at FROM chat_room WHERE id = ?`
	err := db.QueryRow(query, chatRoomId).Scan(&chatRoom.ID, &chatRoom.User1ID, &chatRoom.User2ID, &createdAtString)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("chat room with ID %d not found", chatRoomId)
		}
		return nil, fmt.Errorf("error querying chat room by ID: %v", err)
	}

	// Parse the timestamp string into time.Time using the correct format.
	chatRoom.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtString)
	if err != nil {
		return nil, fmt.Errorf("error parsing created_at timestamp: %v", err)
	}

	return &chatRoom, nil
}

func GetChatRoomsByUserID(db *sql.DB, userID int) ([]models.ChatRoomDetail, error) {
	var rooms []models.ChatRoomDetail

	query := `SELECT 
    cr.id AS chat_room_id, 
    CASE 
        WHEN u.id = cr.user1_id THEN cr.user2_id 
        ELSE cr.user1_id 
    END AS partner_user_id,
    CASE 
        WHEN u.id = cr.user1_id THEN u2.name 
        ELSE u1.name 
    END AS partner_name,
    CASE 
        WHEN u.id = cr.user1_id THEN c2.company_name 
        ELSE c1.company_name 
    END AS company_name,
    cr.created_at,
    lm.english_text AS last_message_english,
    lm.japanese_text AS last_message_japanese,
    lm.user_id AS last_message_user_id,
    lm.date AS last_message_date
FROM 
    chat_room cr
JOIN 
    users u ON u.id = cr.user1_id OR u.id = cr.user2_id
JOIN 
    users u1 ON u1.id = cr.user1_id
JOIN 
    users u2 ON u2.id = cr.user2_id
JOIN 
    companies c1 ON u1.company_id = c1.id
JOIN 
    companies c2 ON u2.company_id = c2.id
LEFT JOIN (
    SELECT 
        chat_room_id, 
        english_text, 
        japanese_text,
        user_id,
        date,
        ROW_NUMBER() OVER(PARTITION BY chat_room_id ORDER BY created_at DESC) AS rn
    FROM 
        conversations
) lm ON cr.id = lm.chat_room_id AND lm.rn = 1
WHERE 
    u.id = ?
ORDER BY 
    lm.date DESC, cr.created_at DESC;

`

	rows, err := db.Query(query, userID) // Memasukkan userID tiga kali untuk ketiga placeholder

	if err != nil {
		return nil, fmt.Errorf("error querying chat rooms by user ID: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var room models.ChatRoomDetail
		var createdAtStr string
		var lastMessageEnglish, lastMessageJapanese sql.NullString
		var lastMessageUser sql.NullInt64 // Untuk menangani user_id yang bisa NULL
		var lastMessageDate sql.NullInt64

		if err := rows.Scan(&room.ChatRoomID, &room.UserID, &room.Name, &room.CompanyName, &createdAtStr, &lastMessageEnglish, &lastMessageJapanese, &lastMessageUser, &lastMessageDate); err != nil {
			return nil, fmt.Errorf("error scanning chat room row: %v", err)
		}

		// Konversi createdAtStr ke time.Time
		room.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing created_at timestamp: %v", err)
		}
		// Menetapkan nilai default menjadi 0 jika last_message_date adalah NULL
		var messageDate int64 = 0 // Menetapkan default value
		if lastMessageDate.Valid {
			messageDate = lastMessageDate.Int64 // Jika tidak NULL, gunakan nilai dari database
		}

		if lastMessageUser.Valid {
			// lastMessageUser.Int64 memiliki nilai user_id yang valid
			fmt.Println("Last message user ID:", lastMessageUser.Int64)
		} else {
			// Tidak ada user_id untuk pesan terakhir (mungkin karena tidak ada pesan)
			fmt.Println("No last message user ID")
		}

		// Contoh memasukkan ke dalam struktur response
		if lastMessageUser.Valid {
			room.LastMessageUser = lastMessageUser.Int64
		} else {
			room.LastMessageUser = 0 // Atau pilih untuk tidak menetapkan / menggunakan 'omitempty' di tag JSON
		}
		var userID int64 // Siapkan variabel untuk menampung user_id

		// Periksa apakah lastMessageUser memiliki nilai valid
		if lastMessageUser.Valid {
			userID = lastMessageUser.Int64 // Gunakan nilai int64 jika valid
		} else {
			userID = 0 // Atau nilai default yang diinginkan ketika user_id adalah NULL
		}

		room.LastMessage.Date = messageDate

		// Sekarang, gunakan userID yang sudah diolah saat membangun LastMessage
		room.LastMessage = models.LastMessage{
			English:  lastMessageEnglish.String,
			Japanese: lastMessageJapanese.String,
			UserID:   userID, // Gunakan userID yang sudah diolah
			Date:     messageDate,
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over chat rooms: %v", err)
	}

	return rooms, nil
}

// // GetLastMessageByChatRoomID retrieves the last message from the specified chat room.
// func  GetLastMessageByChatRoomID(chatRoomID int) (*models.Conversation, error) {
//     var lastMessage models.Conversation
//     query := `SELECT * FROM conversations WHERE chat_room_id = ? ORDER BY created_at DESC LIMIT 1`
//     row := db.QueryRow(query, chatRoomID)
//     if err := row.Scan(&lastMessage.ID, &lastMessage.JapaneseText, &lastMessage.EnglishText, /* ... other fields ... */); err != nil {
//         return nil, err
//     }
//     return &lastMessage, nil
// }
