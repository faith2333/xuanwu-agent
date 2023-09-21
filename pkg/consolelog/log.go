package consolelog

import "github.com/gookit/color"

type ConsoleLog struct {
}

func NewConsoleLog() *ConsoleLog {
	return &ConsoleLog{}
}

func (cl *ConsoleLog) Error(msg string) {
	color.Red.Println("ERROR: " + msg)
}

func (cl *ConsoleLog) Errorf(msg string, a ...any) {
	color.Red.Printf("ERROR: "+msg+"\n", a...)
}

func (cl *ConsoleLog) Waring(msg string) {
	color.Yellow.Println("WARING: " + msg)
}

func (cl *ConsoleLog) Waringf(msg string, a ...any) {
	color.Red.Printf("WARING: "+msg+"\n", a...)
}

func (cl *ConsoleLog) Info(msg string) {
	color.Green.Println("INFO: " + msg)
}

func (cl *ConsoleLog) Infof(msg string, a ...any) {
	color.Green.Printf("INFO: "+msg+"\n", a...)
}

func (cl *ConsoleLog) Printf(msg string, a ...any) {
	color.Green.Printf(msg, a...)
}
