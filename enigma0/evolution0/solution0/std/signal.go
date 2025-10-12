package std

type signalMaker byte

var Signal signalMaker

type spark *Impulse

func (signalMaker) Spark(source ...*Impulse) spark {
	if len(source) > 0 {
		return source[0]
	}
	return nil
}
