CREATE TABLE conversations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    japanese_text TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
    english_text TEXT,
    user_id VARCHAR(255),
    company_id INT,
    chat_room_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
