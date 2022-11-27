package domain

type Response struct {
	Text       *string  `json:"text" repos:"search"`
	MediaLinks []string `json:"mediaLinks"`
}

type ResponseUpdate struct {
	Text       *string
	MediaLinks []string
}
