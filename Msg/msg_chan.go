package Msg

type Student struct {
	StudentName string
	StudentNum  string
}

type Message struct {
	Stu Student
	Msg string
}

const (
	MaxMsg = 1024
)

var MessageChan = make(chan Message, MaxMsg)
