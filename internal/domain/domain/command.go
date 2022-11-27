package domain

type Command struct {
	Name     string   `repos:"search"`
	Aliases  []string `repos:"search"`
	Response Response `repos:"search"`
}

func (elem Command) ID() string {
	return elem.Name
}

type CommandUpdate struct {
	Name     *string
	Aliases  []string
	Response *ResponseUpdate
}
