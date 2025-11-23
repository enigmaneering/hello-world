package event

type Message int

// Event types
const (
	Resize    Message = 1
	Move      Message = 2
	FocusGain Message = 3
	FocusLose Message = 4
	CloseReq  Message = 5
	Closed    Message = 6
	Minimize  Message = 7
	Restore   Message = 8
)
