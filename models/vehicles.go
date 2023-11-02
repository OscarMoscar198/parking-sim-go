package models

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Vehicle struct {
	ID         int
	rectangule *canvas.Rectangle
	time       int
	semQuit    chan bool
}

const (
	minDuration = 5
	maxDuration = 7
)

var (
	exitCars []*Vehicle
)

func NewSpaceVehicle() *Vehicle {
	rectangule := canvas.NewRectangle(White)
	rectangule.SetMinSize(fyne.NewSquareSize(float32(30)))

	car := &Vehicle{
		ID:         -1,
		rectangule: rectangule,
		time:       0,
	}

	return car
}

func NewVehicle(id int, sQ chan bool) *Vehicle {
	purpleColor := color.RGBA{R: 128, G: 0, B: 128, A: 255}
	time := rand.Intn(maxDuration-minDuration) + minDuration

	rectangule := canvas.NewRectangle(purpleColor)
	rectangule.SetMinSize(fyne.NewSquareSize(float32(30)))

	car := &Vehicle{
		ID:         id,
		rectangule: rectangule,
		time:       time,
		semQuit:    sQ,
	}

	return car
}

func (c *Vehicle) StartCount(id int) {
	for {
		select {
		case <-c.semQuit:
			fmt.Printf("StartCount Close")
			return
		default:
			if c.time <= 0 {
				c.ID = id
				exitCars = append(exitCars, c)
				return
			}
			c.time--
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (c *Vehicle) GetRectangle() *canvas.Rectangle {
	return c.rectangule
}

func (c *Vehicle) ReplaceData(car *Vehicle) {
	c.ID = car.ID
	c.time = car.time
	c.rectangule.FillColor = car.rectangule.FillColor
}

func (c *Vehicle) GetTime() int {
	return c.time
}

func (c *Vehicle) GetID() int {
	return c.ID
}

func GetWaitVehicles() []*Vehicle {
	return exitCars
}

func PopExitWaitVehicles() *Vehicle {
	car := exitCars[0]
	if !WaitExitVehiclesIsEmpty() {
		exitCars = exitCars[1:]
	}
	return car
}

func WaitExitVehiclesIsEmpty() bool {
	return len(exitCars) == 0
}
