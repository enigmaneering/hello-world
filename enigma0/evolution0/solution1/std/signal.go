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

type decay int

// Decay waits the specified number of activations (or 0 if omitted) before decaying the synapse.
//
// NOTE: If you'd like to also close the synaptic channel, please use DecayClose
func (signalMaker) Decay(impulseDelay ...uint) decay {
	d := uint(0)
	if len(impulseDelay) > 0 {
		d = uint(impulseDelay[0])
	}
	return decay(d)
}

// DecayClose waits the specified number of activations (or 0 if omitted) before decaying the synapse
// and immediately closing its synaptic channel.
//
// NOTE: Closing the channel is often a source of panics, please prefer Decay if possible.
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

type mute byte

func (signalMaker) Mute() mute {
	return mute(0)
}

type unmute byte

func (signalMaker) Unmute() unmute {
	return unmute(0)
}

type shutdown bool

func (signalMaker) Shutdown(close ...bool) shutdown {
	return shutdown(len(close) == 0 && close[0])
}
