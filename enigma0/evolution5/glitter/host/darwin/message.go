package darwin

import (
	"fmt"

	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution5/glitter/host/darwin/event"
)

type Message struct {
	event.Message
	Data1 uint
	Data2 uint
}

func (m Message) String() string {
	switch m.Message {
	case event.Move:
		return fmt.Sprintf("Moved to (%d, %d)", m.Data1, m.Data2)
	case event.Resize:
		return fmt.Sprintf("Resized to (%d, %d)", m.Data1, m.Data2)
	case event.Closed, event.CloseReq:
		return fmt.Sprintf("Closed")
	case event.FocusGain:
		return fmt.Sprintf("Focus Gained")
	case event.FocusLose:
		return fmt.Sprintf("Focus Lost")
	case event.Restore:
		return fmt.Sprintf("Restored")
	case event.Minimize:
		return fmt.Sprintf("Minimized")
	default:
		return fmt.Sprintf("%v", m)
	}
}
