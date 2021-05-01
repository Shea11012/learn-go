package bridge

import "fmt"

type Mac struct {
	printer
}

func (m *Mac) print() {
	fmt.Println("print mac")
	m.printer.printFile()
}

func (m *Mac) setPrinter(p printer) {
	m.printer = p
}
