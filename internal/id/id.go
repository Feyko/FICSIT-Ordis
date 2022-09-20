package id

type IDer interface {
	ID() string
}

type Searchable interface {
	IDer
	SearchFields() []string
}
