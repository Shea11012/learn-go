package bridge

import "fmt"

type windows struct {
	printer
}

func (w *windows) print() {
	fmt.Println("print window")
	w.printer.printFile()
}

func (w *windows) setPrinter(p printer) {
	w.printer = p
}

