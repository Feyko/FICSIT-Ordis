package repo

type Error string

func (e Error) Error() string {
	return string(e)
}

var (
	ErrCollectionNotFound Error = "collection not found"
	ErrElementNotFound    Error = "element not found"
)
