package std

import "sync"

// Synchro represents a way to synchronize execution across threads.
//
// To send execution using a synchro, first create one using make.  Then, Engage the synchro
// (non-blocking) from the thread you wish to execute on.  The calling thread can then Send
// actions (blocking) to the other thread for intermittent execution.
//
//		 global -
//	   var synchro = make(std.Synchro)
//
//		 main loop -
//		  for ... {
//	    ...
//		   synchro.Engage()
//		   ...
//		  }
//
//		 sender -
//	   synchro.Send(func() { ... })
type Synchro chan *syncAction

// syncAction represents a "waitable" action.
type syncAction struct {
	sync.WaitGroup
	Action func()
}

// Send sends the provided action over the synchro channel and waits for it to be executed.
func (s Synchro) Send(action func()) {
	syn := &syncAction{Action: action}
	syn.Add(1)
	s <- syn
	syn.Wait()
}

// Engage asynchronously handles -all- of the currently incoming actions on the Synchro channel before returning control.
func (s Synchro) Engage() {
	for {
		select {
		case syn := <-s:
			syn.Action()
			syn.Done()
		default:
			return
		}
	}
}

// EngageOnce synchronously reads a single action on the Synchro channel before returning control.
func (s Synchro) EngageOnce() {
	syn := <-s
	syn.Action()
	syn.Done()
}
