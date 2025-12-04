package dtos

type AnswerQuestionResult struct {
	Status string
	Text   string
}

type CreateChatParams struct {
	Title      string
	DocumentID string
}

type UpdateChatParams struct {
	Title *string
}
