## Save Conversation
POST http://localhost:8080/api/conversations

Request Body :
{
  "user_id": "user1",
  "speaker": "Ato",
  "company_id": 107,
  "chat_room_id": "room13",
  "original_message": "Hello, what's time is it??"
}

Response Body :
{
    "conversations": [
        {
            "speaker": "Ato",
            "original_message": "Hello, what's time is it??",
            "translated_message": "こんにちは、今何時ですか？"
        }
    ]
}

## Get All Conversation
GET http://localhost:8080/api/conversations

Response :
    {
        "id": 4,
        "japanese_text": "こんにちは、今何時ですか？",
        "english_text": "Hello, what's time is it??",
        "speaker": "",
        "user_id": "user1",
        "company_id": 107,
        "chat_room_id": "room13",
        "original_message": "",
        "translated_message": "",
        "created_at": "0001-01-01T00:00:00Z"
    },

## Translate Message
POST http://localhost:8080/api/translate

Request Body :
{
    "speaker": "107",
    "original_message": "あなたはどこに住んでいますか？"
}

Response Body :
{
    "conversations": [
        {
            "speaker": "107",
            "original_message": "あなたはどこに住んでいますか？",
            "translated_message": "Where do you live?"
        }
    ]
}