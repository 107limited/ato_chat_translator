package translation

// Translator adalah antarmuka untuk menerjemahkan pesan.
type Translator interface {
	TranslateMessage(message string) (string, error)
}
