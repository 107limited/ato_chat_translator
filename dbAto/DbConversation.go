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

	query := `
	SELECT 
        cr.id AS chat_room_id, 
        CASE 
            WHEN u.id = cr.user1_id THEN cr.user1_id
            WHEN u.id = cr.user2_id THEN cr.user2_id
        END AS user_id, 
        u.name AS user_name, 
        c.company_name, 
        cr.created_at
    FROM 
        chat_room cr
    JOIN 
        users u ON u.id = cr.user1_id OR u.id = cr.user2_id
    JOIN 
        companies c ON u.company_id = c.id
    WHERE 
        cr.user1_id = ? OR cr.user2_id = ?
    ORDER BY 
        cr.created_at DESC;
	`

	rows, err := db.Query(query, userID, userID) // Memasukkan userID dua kali untuk kedua placeholder

	if err != nil {
		return nil, fmt.Errorf("error querying chat rooms by user ID: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var room models.ChatRoomDetail
		var createdAtBytes []byte // Gunakan ini untuk menerima data timestamp sebagai []byte

		// Sesuaikan urutan Scan berdasarkan urutan kolom dalam query Anda.
		if err := rows.Scan(&room.ChatRoomID, &room.UserID, &room.Name, &room.CompanyName, &createdAtBytes); err != nil {
			return nil, fmt.Errorf("error scanning chat room row: %v", err)
		}

		// Misalkan createdAtBytes adalah []byte yang dihasilkan dari scanning database
		createdAtString := string(createdAtBytes)

		// Gunakan format yang sesuai dengan string waktu Anda, yaitu "2006-01-02 15:04:05" untuk "YYYY-MM-DD HH:MM:SS"
		room.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtString)
		if err != nil {
			return nil, fmt.Errorf("error parsing created_at: %v", err)
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over chat rooms: %v", err)
	}

	return rooms, nil
}
