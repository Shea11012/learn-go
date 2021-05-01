package main

import "fmt"

type iGun interface {
	setName(name string)
	setPower(power int)
	getName() string
	getPower() int
}

type gun struct {
	name string
	power int
}

func (g gun) setName(name string)  {
	g.name = name
}

func (g gun) setPower(power int) {
	g.power = power
}

func (g gun) getName() string {
	return g.name
}

func (g gun) getPower() int {
	return g.power
}

type ak47 struct {
	gun
}

func newAk47() iGun {
	return &ak47{
		gun{
			name:  "Ak47",
			power: 4,
		},
	}
}

type musket struct {
	gun
}

func newMusket() iGun {
	return &musket{
		gun{
			name: "Musket",
			power: 1,
		},
	}
}

func getGun(gunType string) (iGun,error) {
	switch gunType {
	case "ak47":
		return newAk47(),nil
	case "musket":
		return newMusket(),nil
	default:
		return nil,fmt.Errorf("wrong gun type passed")
	}
}

func main() {
	ak47,_ := getGun("ak47")
	musket,_ := getGun("musket")
	printDetails(ak47)
	printDetails(musket)
}

func printDetails(g iGun) {
	fmt.Printf("Gun: %s\nPower: %d\n",g.getName(),g.getPower())
}
