package domain

type Response struct {
	Text       *string  `json:"text"`
	MediaLinks []string `json:"mediaLinks"`
}

type ResponseUpdate struct {
	Text       *string
	MediaLinks []string
}
