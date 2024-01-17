# Ato Chat App

Ato Chat App is a conversation application that leverages the GPT-4 model to translate messages between Japanese and English. The application is developed using the Go programming language and utilizes standard HTTP without relying on any specific web framework. The integration with the GPT-4 model enhances the application's ability to provide accurate translations and a seamless conversational experience

## Installation

1. Make sure you have Go installed on your system.
2. Clone this repository to your local directory:

   
```bash
   git clone https://github.com/ilhamfw/ato-chat
   cd ato-chat-app
   ```
   
1. Install the required dependencies:

```bash
    go mod init
```


2. Create a .env file and configure the necessary environment variables, such as the API Key for OpenAI.

## Usage
Run the application:

Copy code
```bash
go run main.go
```
Open the application in your browser at http://localhost:8080.

# API Endpoints

## Save Conversation

Save a conversation in the system.

### Endpoint

`POST http://localhost:8080/api/conversations`

### Request Body

```json
{
  "user_id": "user1",
  "speaker": "Ato",
  "company_id": 107,
  "chat_room_id": "room13",
  "original_message": "Hello, what's time is it??"
}
```

### Response Body
```json
{
  "conversations": [
    {
      "speaker": "Ato",
      "original_message": "Hello, what's time is it??",
      "translated_message": "こんにちは、今何時ですか？"
    }
  ]
}
```
### A. Screeshoot Save Conversation
![Save Conversation](image/1.%20Save%20Conversation.jpg)


## Translate Message (POST /api/translate)
Translate a message.
### Endpoint
`POST http://localhost:8080/api/translate`

### Request Body
```json 
{
  "user_id": "user123",
  "speaker": "Ato",
  "company_id": 107,
  "chat_room_id": "room123",
  "original_message": "Hello, how are you?"
}
```

### Response Body
```json
{
  "conversations": [
    {
      "speaker": "107",
      "original_message": "あなたはどこに住んでいますか？",
      "translated_message": "Where do you live?"
    }
  ]
}
```
### B. Screenshot Translate Japanese-English
![Translate Message](image/2.%20Translate%20Japanese-English.jpg)
### C. Screenshot Translate English-Japanese
![Translate Message](image/3.%20Translate%20English-Japanese.jpg)


## Get All Conversations (GET /api/conversations)
Retrieve all conversations.
### Endpoint
`GET http://localhost:8080/api/conversations`
### Response Body
```json
[
  {
    "ID": 1,
    "JapaneseText": "こんにちは",
    "EnglishText": "Hello",
    "UserID": "user123",
    "CompanyID": 107,
    "ChatRoomID": "room123",
    "OriginalMessage": "Hello",
    "CreatedAt": "2024-01-01T00:00:00Z"
  },
  // ... other conversations ...
]
```
### D. Screenshoot Get All Conversation
![Get All Conversation](image/4.%20Get%20All%20Conversation.jpg)

### E. Result at Database ato_chat Table conversations using MariaDB at DBeaver
![Result](image/5.%20Table%20Conversations%20using%20Maria%20Db%20with%20Dbeaver.jpg)



