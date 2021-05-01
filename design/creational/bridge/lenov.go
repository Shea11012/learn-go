package bridge

import "fmt"

type lenov struct {}

func (l *lenov) printFile() {
	fmt.Println("print lenov")
}

type hp struct {

}

func (h *hp) printFile() {
	fmt.Println("print hp")
}

