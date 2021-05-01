package bridge

type printer interface {
	printFile()
}

type Computer interface {
	print()
	setPrinter(printer)
}
