package std

// Synapse represents the basic synapse type.  As synapses are rife with nuance and subtlety, to create a synapse,
// please see the std.Signal factory
type Synapse chan<- any

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
