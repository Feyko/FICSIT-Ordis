package domain

type Response struct {
	Text       string
	MediaLinks []string
}

type Command struct {
	Name     string
	Aliases  []string
	Response Response
}

func (elem Command) ID() string {
	return elem.Name
}

func (elem Command) SearchFields() []string {
	return []string{"Name", "Response"}
}
