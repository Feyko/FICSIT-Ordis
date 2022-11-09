package domain

type CrashUpdate struct {
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	Regexes     []string  `json:"regexes"`
	Response    *Response `json:"response"`
}

type Crash struct {
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Regexes     []string `json:"regexes"`
	Response    Response `json:"response"`
}

func (c Crash) SearchFields() []string {
	return []string{"Name", "Description", "Regexes", "Response"}
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
