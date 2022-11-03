package domain

type Crash struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Regexes     []string         `json:"regexes"`
	Response    *Response        `json:"response"`
}

func (c Crash) SearchFields() []string {
	return []string{"Name", "Description", "Regexes", "Response"}
}

func (c Crash) ID() string {
	return c.Name
}

type CrashUpdate struct {
	Name        *string        `json:"name"`
	Description *string        `json:"description"`
	Regexes     []*string      `json:"regexes"`
	Response    *Response      `json:"response"`
}