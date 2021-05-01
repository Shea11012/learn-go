package bridge

import (
	"fmt"
	"testing"
)

func TestBridge(t *testing.T) {
	hpPrinter := &hp{}
	lenvoPrinter := &lenov{}

	macCmp := &Mac{}
	macCmp.setPrinter(hpPrinter)
	macCmp.print()
	fmt.Println()

	macCmp.setPrinter(lenvoPrinter)
	macCmp.print()
	fmt.Println()

	winCmp := &windows{}
	winCmp.setPrinter(hpPrinter)
	winCmp.print()
	fmt.Println()

	winCmp.setPrinter(lenvoPrinter)
	winCmp.print()
	fmt.Println()
}
