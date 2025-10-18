package std

type signalMaker[T any] byte

func Signal[T any]() signalMaker[T] {
	return signalMaker[T](0)
}

type impulse[T any] []*Thought[T]

func (signalMaker[T]) Impulse(thoughts ...*Thought[T]) any {
	return thoughts
}

type spark[T any] []*Thought[T]

func (signalMaker[T]) Spark(thoughts ...*Thought[T]) any {
	return thoughts
}
