package id

type Wrapper[T any] struct {
	ID_     string
	Wrapped T
}

func (w Wrapper[T]) ID() string {
	return w.ID_
}

type WrapperUpdate[T any] struct {
	Wrapped *T
}
