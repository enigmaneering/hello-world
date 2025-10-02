package std

type signalMaker byte

var Signal signalMaker

type spark []Neural

func (signalMaker) Spark(neurons ...Neural) spark {
	return neurons
}
