package repo

type CollectionNotFoundError struct{}

func (c CollectionNotFoundError) Error() string {
	return "collection not found"
}

var ErrCollectionNotFound CollectionNotFoundError
