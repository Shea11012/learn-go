package main

import "fmt"

type client struct {}

func (c *client) insertLightningConnectorIntoComputer(com computer) {
	fmt.Println("client inserts lightning connector into computer")
	com.insertIntoLightningPort()
}

type computer interface {
	insertIntoLightningPort()
}

type mac struct {}

func (m *mac) insertIntoLightningPort() {
	fmt.Println("lightning connector is plugged into mac machine")
}

type windows struct {}

func (w *windows) insertIntoUSBPort() {
	fmt.Println("USB connector is plugged into window machine")
}

type windowAdapter struct {
	windowMachine *windows
}

func (w *windowAdapter) insertIntoLightningPort() {
	fmt.Println("adapter converts lightning signal to usb")
	w.windowMachine.insertIntoUSBPort()
}

func main() {
	client := &client{}
	mac := &mac{}

	client.insertLightningConnectorIntoComputer(mac)

	windowMachine := &windows{}
	windowMachineAdapter := &windowAdapter{
		windowMachine: windowMachine,
	}

	client.insertLightningConnectorIntoComputer(windowMachineAdapter)
}