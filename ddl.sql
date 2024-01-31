CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    company_id INT,
    role_id INT,
    name VARCHAR(255)
    FOREIGN KEY (company_id) REFERENCES companies(id),
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

CREATE TABLE companies (
    id INT PRIMARY KEY AUTO_INCREMENT,
    company_name VARCHAR(255) NOT NULL
);

CREATE TABLE roles (
    id INT PRIMARY KEY AUTO_INCREMENT,
    role_name VARCHAR(255) NOT NULL
);

CREATE TABLE chat_room (
    id INT PRIMARY KEY AUTO_INCREMENT,
    chat_room_name VARCHAR(255) NOT NULL
);


CREATE TABLE conversations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    japanese_text TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
    english_text TEXT,
    user_id VARCHAR(255),
    company_id INT,
    chat_room_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (chat_room_id) REFERENCES chat_room(id),
    FOREIGN KEY (company_id) REFERENCES companies(id)
);

CREATE TABLE `sec_m` (
  `private_key` blob DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;