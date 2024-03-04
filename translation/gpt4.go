package translation

import (
	chat "ato_chat/chat"
	"ato_chat/models"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GPT4Client adalah klien untuk berinteraksi dengan API GPT-4
type GPT4Client struct {
	APIKey string
	Model  string
	DB     chat.ConversationRepository
}

func (c *GPT4Client) GetAllConversations() ([]*models.Conversation, error) {
	// Retrieve all conversations from the database
	conversations, err := c.DB.GetAllConversations()
	if err != nil {
		return nil, fmt.Errorf("error getting all conversations: %v", err)
	}

	return conversations, nil
}


// SaveConversation implements chat.ConversationRepository.
func (c *GPT4Client) SaveConversation(conv *models.Conversation) (sql.Result, error) {
	// Buat objek untuk menyimpan percakapan ke dalam database atau penyimpanan lainnya

	// Lakukan penyimpanan ke dalam database atau penyimpanan lainnya
	// Contoh: Simpan ke dalam database MySQL\
	res, err := c.DB.SaveConversation(conv)
	if err != nil {
		return nil, fmt.Errorf("error saving conversation: %v", err)
	}
	fmt.Printf("Saving conversation: %+v\n", conv)

	// Logika penyimpanan berhasil
	fmt.Printf("Conversation saved: %+v\n", conv)
	return res, nil
}

// NewGPT4Client creates a new instance of GPT4Client
func NewGPT4Client(apiKey, model string, db chat.ConversationRepository) *GPT4Client {
	return &GPT4Client{
		APIKey: apiKey,
		Model:  model,
		DB:     db,
	}
}

type APIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// TranslateMessage menerjemahkan pesan menggunakan API GPT-4
func (c *GPT4Client) TranslateMessage(prompt string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	// Membuat request body
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "system", "content": "translate following sentence to japanese without the romanji if the sentence in english, and do the opposite if sentence in japanese. Show only the translation:" + prompt},
		},
	})
	if err != nil {
		return "", err
	}

	// Membuat HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	// // Mengatur kunci API ke header
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// Mengirim request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Penanganan status code error
    if resp.StatusCode != http.StatusOK {
        if resp.StatusCode == http.StatusUnauthorized {
            return "", fmt.Errorf("API key expired or invalid")
        }
        return "", fmt.Errorf("received non-ok status code: %d", resp.StatusCode)
    }

	// Membaca respons
	var bodyBuffer bytes.Buffer
	_, err = io.Copy(&bodyBuffer, resp.Body)
	if err != nil {
		return "", err
	}

	// // Menguraikan JSON ke dalam struktur data
	var responseData map[string]interface{}
	err = json.Unmarshal(bodyBuffer.Bytes(), &responseData)
	if err != nil {
		return "", err
	}

	// // Mengambil nilai yang diinginkan dari struktur data
	fmt.Printf("%v", responseData)
	translatedMessage, ok := responseData["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	if !ok {
		return "", fmt.Errorf("failed to extract translated message from response")
	}

	return translatedMessage, nil
}
