package interfaces

type AIService interface {
	GenerateText(prompt string) (string, error)
	GenerateImage(description string) ([]byte, error)
}
