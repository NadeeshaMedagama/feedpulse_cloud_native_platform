package models

type AIAnalyzeRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type AISummaryRequest struct {
	Text string `json:"text"`
}
