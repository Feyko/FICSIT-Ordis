package domain

type Command struct {
	Name     string
	Aliases  []string
	Response Response
}

type CommandUpdate struct {
	Update

	Name     *string
	Aliases  []string
	Response *ResponseUpdate
}

func (elem Command) ID() string {
	return elem.Name
}

func (elem Command) SearchFields() []string {
	return []string{"Name", "Response"}
}
