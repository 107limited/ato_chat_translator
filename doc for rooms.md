# Create Room
Endpoint ini digunakan untuk membuat chat room baru antara dua pengguna.

### URL: /api/chatrooms
### Method: POST
### Auth Required: Yes
### Permissions Required: Tidak ada

## Data yang Diperlukan dalam Request Body:

```json
{
  "user1_id": int,
  "user2_id": int
}
```

### Response Sukses:

Code: 201 CREATED
Content:

```json
{
  "chat_room_id": int,
  "message": "Chat room successfully created."
}

```

### Response Error:

Code: 400 BAD REQUEST
Content:
```json
{
  "message": "Invalid request body"
}
```

# Get Room
Endpoint ini digunakan untuk mendapatkan informasi chat room yang sudah ada antara dua pengguna.

### URL: /api/chatrooms/{user1_id}/{user2_id}
### Method: GET
### Auth Required: Yes
### Permissions Required: Tidak ada

### Parameter URL:

`user1_id`: ID pengguna pertama
`user2_id`: ID pengguna kedua

### Response Sukses:

Code: 200 OK
Content:
```json
{
  "chat_room_id": int,
  "user1_id": int,
  "user2_id": int,
  "message": "Chat room details retrieved successfully."
}
```
### Response Error:

Code: 404 NOT FOUND
Content:
```json
{
  "message": "Chat room not found."
}
```