package shepard

type Clone[T any] interface {
	Clone() T
}
