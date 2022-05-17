package domain

type Update struct {
	AddInsteadOfSet bool
}

type Response struct {
	Text       string
	MediaLinks []string
}

type ResponseUpdate struct {
	Update

	Text       *string
	MediaLinks []string
}
