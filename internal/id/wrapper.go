package id

func Wrap[T any](v T, id string) Wrapper[T] {
	return Wrapper[T]{
		Wrapped: v,
		ID_:     id,
	}
}

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
