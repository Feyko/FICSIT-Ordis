package domain

type CrashUpdate struct {
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	Regexes     []string  `json:"regexes"`
	Response    *Response `json:"response"`
}

type Crash struct {
	Name        string   `json:"name" repos:"search"`
	Description *string  `json:"description" repos:"search"`
	Regexes     []string `json:"regexes" repos:"search"`
	Response    Response `json:"response" repos:"search"`
}

func (c Crash) ID() string {
	return c.Name
}

type CrashMatch struct {
	MatchedText string `json:"matchedText"`
	Crash       *Crash `json:"crash"`
	CharSpan    *Span  `json:"charSpan"`
}

type Span struct {
	Start int `json:"start"`
	End   int `json:"end"`
}
