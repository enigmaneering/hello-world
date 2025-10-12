package std

type Neuron[T any] struct {
	name      string
	action    func(Impulse[T])
	potential func(Impulse[T]) bool
	cleanup   func(Impulse[T])
}

func (n Neuron[T]) Named() string {
	return n.name
}

func (n Neuron[T]) Action(imp Impulse[T]) {
	if n.action != nil {
		n.action(imp)
	}
}

func (n Neuron[T]) Potential(imp Impulse[T]) bool {
	if n.potential == nil {
		return true
	}
	return n.potential(imp)
}

func (n Neuron[T]) Cleanup(imp Impulse[T]) {
	if n.cleanup != nil {
		n.cleanup(imp)
	}
}

func NewNeuron[T any](named string, action func(Impulse[T]), potential func(Impulse[T]) bool, cleanup ...func(Impulse[T])) Neuron[T]
{

}
