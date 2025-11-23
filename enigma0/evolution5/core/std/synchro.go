package std

import "sync"

// Synchro represents a way to synchronize execution across threads.
//
// To send execution using a synchro, first create one using make - then, Engage the synchro
// from the thread you wish to execute on.  The calling thread can then Send blocking actions
// which the synchronizable thread should "intermittently" execute.
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

// Engage asynchronously handles the currently incoming actions on the Synchro channel before returning control.
//
// NOTE: If you'd like to process ALL available messages in a single engagement, rather than one, please provide 'true' to 'processAll'
func (s Synchro) Engage(processAll ...bool) {
	// This defaults to 'false' as SDL2 windows were intermittently not getting created if multiple were asked to be spawned in the same
	// loop cycle.  This was solved by putting a natural 'delay' between each message in the form of waiting a single loop cycle.
	// 		- Alex
	all := len(processAll) > 0 && processAll[0]
	for {
		select {
		case syn := <-s:
			syn.Action()
			syn.Done()
			if !all {
				return
			}
		default:
			return
		}
	}
}

// EngageBlocking synchronously handles the currently incoming actions on the Synchro channel before returning control.
//
// NOTE: If you'd like to process ALL available messages in a single engagement, rather than one, please provide 'true' to 'processAll'
func (s Synchro) EngageBlocking(processAll ...bool) {
	// This defaults to 'false' as SDL2 windows were intermittently not getting created if multiple were asked to be spawned in the same
	// loop cycle.  This was solved by putting a natural 'delay' between each message in the form of waiting a single loop cycle.
	// 		- Alex
	all := len(processAll) > 0 && processAll[0]
	for {
		select {
		case syn := <-s:
			syn.Action()
			syn.Done()
			if !all {
				return
			}
		}
	}
}
