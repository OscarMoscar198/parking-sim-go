package models

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"sync"
	"time"
)

var (
	Gray           = color.RGBA{R: 30, G: 30, B: 30, A: 255}
	Black          = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	mutexExitEnter sync.Mutex
)

const (
	lambda         = 2.0
	MaxWait    int = 10
	MaxParking int = 20
)

type Parking struct {
	waitCars            []*Car
	parking             [MaxParking]*Car
	entrace             *Car
	exit                *Car
	out                 *Car
	semQuit             chan bool
	semRenderNewCarWait chan bool
}

func NewParking(sENCW chan bool, sQ chan bool) *Parking {
	parking := &Parking{
		semRenderNewCarWait: sENCW,
		semQuit:             sQ,
	}
	return parking
}

func (p *Parking) MakeParking() {
	for i := range p.parking {
		car := NewSpaceCar()
		p.parking[i] = car
	}
}

func (p *Parking) MakeOutStation() *Car {
	p.out = NewSpaceCar()
	return p.out
}

func (p *Parking) MakeExitStation() *Car {
	p.exit = NewSpaceCar()
	return p.exit
}

func (p *Parking) MakeEntraceStation() *Car {
	p.entrace = NewSpaceCar()
	return p.entrace
}

func (p *Parking) GenerateCars() {
	i := 20
	for {
		select {
		case <-p.semQuit:
			fmt.Printf("GenerateCars Close")
			return
		default:
			interarrivalTime := -math.Log(1-rand.Float64()) / lambda //Poisson.
			time.Sleep(time.Duration(interarrivalTime * float64(time.Second)))
			if len(p.waitCars) < MaxWait {
				car := NewCar(i, p.semQuit)
				i++
				p.waitCars = append(p.waitCars, car)
				p.semRenderNewCarWait <- true
			}
		}
	}
}

func (p *Parking) CheckParking() {
	for {
		select {
		case <-p.semQuit:
			fmt.Printf("CheckParking Close")
			return
		default:
			index := p.SearchSpace()
			if index != -1 && !p.WaitCarsIsEmpty() {
				mutexExitEnter.Lock()
				p.MoveToEntrace()
				p.MoveToPark(index)
				mutexExitEnter.Unlock()
			}
		}
	}
}

func (p *Parking) MoveToEntrace() {
	car := p.PopWaitCars()
	p.entrace.ReplaceData(car)
	time.Sleep(1 * time.Second)
}

func (p *Parking) MoveToPark(index int) {
	p.parking[index].ReplaceData(p.entrace)
	p.parking[index].text.Show()

	p.entrace.ReplaceData(NewSpaceCar())
	go p.parking[index].StartCount(index)
	time.Sleep(1 * time.Second)

}

func (p *Parking) OutCarToExit() {
	for {
		select {
		case <-p.semQuit:
			fmt.Printf("CarExit Close")
			return
		default:
			if !WaitExitCarsIsEmpty() {
				mutexExitEnter.Lock()
				car := PopExitWaitCars()

				p.MoveToExit(car.ID)
				p.MoveToOut()
				mutexExitEnter.Unlock()

				time.Sleep(1 * time.Second)
				p.out.ReplaceData(NewSpaceCar())
			}
		}
	}
}

func (p *Parking) MoveToExit(index int) {
	p.exit.ReplaceData(p.parking[index])
	p.parking[index].text.Hide()
	p.parking[index].ReplaceData(NewSpaceCar())
	time.Sleep(1 * time.Second)
}

func (p *Parking) MoveToOut() {
	p.out.ReplaceData(p.exit)
	p.exit.ReplaceData(NewSpaceCar())
	time.Sleep(1 * time.Second)
}

func (p *Parking) SearchSpace() int {
	for s := range p.parking {
		if p.parking[s].GetID() == -1 {
			return s
		}
	}
	return -1
}

func (p *Parking) PopWaitCars() *Car {
	car := p.waitCars[0]
	if !p.WaitCarsIsEmpty() {
		p.waitCars = p.waitCars[1:]
	}
	return car
}

func (p *Parking) WaitCarsIsEmpty() bool {
	return len(p.waitCars) == 0
}

func (p *Parking) GetWaitCars() []*Car {
	return p.waitCars
}

func (p *Parking) GetEntraceCar() *Car {
	return p.entrace
}

func (p *Parking) GetExitCar() *Car {
	return p.exit
}

func (p *Parking) GetParking() [MaxParking]*Car {
	return p.parking
}

func (p *Parking) ClearParking() {
	for i := range p.parking {
		p.parking[i] = nil
	}
}
