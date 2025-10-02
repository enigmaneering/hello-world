package std

type signalMaker byte

var Signal signalMaker

type spark []Neural

func (signalMaker) Spark(neurons ...Neural) spark {
	return neurons
}

type decay int

func (signalMaker) Decay(impulseDelay ...uint) decay {
	d := uint(0)
	if len(impulseDelay) > 0 {
		d = uint(impulseDelay[0])
	}
	return decay(d)
}
func (signalMaker) DecayClose(impulseDelay ...uint) decay {
	d := 0
	if len(impulseDelay) > 0 {
		d = -int(impulseDelay[0])
	}
	return decay(d)
}

type closer byte

// Close will tell a decayed synapse it's safe to close its channel and fully exit.
//
// NOTE: If sent to a non-decayed channel, this will be discarded.
func (signalMaker) Close() closer {
	return closer(0)
}
