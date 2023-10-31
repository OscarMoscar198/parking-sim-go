package view

import (
	"fmt"
	"image/color"
	"simulation/models"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

var semRenderNewCarWait chan bool
var semQuit chan bool

type ParkingView struct {
	window               fyne.Window
	waitRectangleStation [models.MaxWait]*canvas.Rectangle
}

var Gray = color.RGBA{R: 30, G: 30, B: 30, A: 255}

var parking *models.Parking

func NewParkingView(window fyne.Window) *ParkingView {
	parkingView := &ParkingView{window: window}

	semQuit = make(chan bool)
	semRenderNewCarWait = make(chan bool)

	parking = models.NewParking(semRenderNewCarWait, semQuit)
	parkingView.MakeScene()
	parkingView.StartSimulation()

	return parkingView
}

func (p *ParkingView) MakeScene() {

	containerParkingView := container.New(layout.NewVBoxLayout())
	containerParkingOut := container.New(layout.NewHBoxLayout())

	containerParkingOut.Add(p.MakeWaitStation())
	containerParkingOut.Add(layout.NewSpacer())
	containerParkingOut.Add(p.MakeExitStation())
	containerParkingOut.Add(layout.NewSpacer())

	containerParkingView.Add(containerParkingOut)
	containerParkingView.Add(layout.NewSpacer())
	containerParkingView.Add(p.MakeParkingLotEntrance())
	containerParkingView.Add(layout.NewSpacer())
	containerParkingView.Add(p.MakeEnterAndExitStation())
	containerParkingView.Add(layout.NewSpacer())
	containerParkingView.Add(p.MakeParking())
	containerParkingView.Add(layout.NewSpacer())

	p.window.SetContent(containerParkingView)
	p.window.Resize(fyne.NewSize(500, 500))
	p.window.CenterOnScreen()
}

func (p *ParkingView) MakeParking() *fyne.Container {
	parkingContainer := container.New(layout.NewGridLayout(5))
	parking.MakeParking()

	parkingArray := parking.GetParking()
	for i := 0; i < len(parkingArray); i++ {
		if i == 10 {
			addSpace(parkingContainer)
		}
		parkingContainer.Add(container.NewCenter(parkingArray[i].GetRectangle(), parkingArray[i].GetText()))
	}
	return container.NewCenter(parkingContainer)
}

func (p *ParkingView) MakeWaitStation() *fyne.Container {
	parkingContainer := container.New(layout.NewGridLayout(5))
	for i := len(p.waitRectangleStation) - 1; i >= 0; i-- {
		car := models.NewSpaceCar().GetRectangle()
		p.waitRectangleStation[i] = car
		p.waitRectangleStation[i].Hide()
		parkingContainer.Add(p.waitRectangleStation[i])
	}
	return parkingContainer
}

func (p *ParkingView) MakeExitStation() *fyne.Container {
	out := parking.MakeOutStation()
	return container.NewCenter(out.GetRectangle())
}

func (p *ParkingView) MakeEnterAndExitStation() *fyne.Container {
	parkingContainer := container.New(layout.NewGridLayout(5))
	parkingContainer.Add(layout.NewSpacer())
	entrace := parking.MakeEntraceStation()
	parkingContainer.Add(entrace.GetRectangle())
	parkingContainer.Add(layout.NewSpacer())
	exit := parking.MakeExitStation()
	parkingContainer.Add(exit.GetRectangle())
	parkingContainer.Add(layout.NewSpacer())
	return container.NewCenter(parkingContainer)
}

func (p *ParkingView) MakeParkingLotEntrance() *fyne.Container {
	EntraceContainer := container.New(layout.NewGridLayout(3))
	EntraceContainer.Add(makeBorder())
	EntraceContainer.Add(layout.NewSpacer())
	EntraceContainer.Add(makeBorder())
	return EntraceContainer
}

func (p *ParkingView) RenderNewCarWaitStation() {
	for {
		select {
		case <-semQuit:
			fmt.Printf("RenderNewCarWaitStation Close")
			return
		case <-semRenderNewCarWait:
			waitCars := parking.GetWaitCars()
			for i := len(waitCars) - 1; i >= 0; i-- {
				if waitCars[i].ID != -1 {
					p.waitRectangleStation[i].Show()
					p.waitRectangleStation[i].FillColor = waitCars[i].GetRectangle().FillColor
				}
			}
			p.window.Content().Refresh()
		}
	}
}

func (p *ParkingView) RenderUpdate() {
	for {
		select {
		case <-semQuit:
			fmt.Printf("RenderUpdate Close")
			return
		default:
			p.window.Content().Refresh()
			time.Sleep(1 * time.Second)
		}
	}
}

func (p *ParkingView) StartSimulation() {
	go parking.GenerateCars()
	go parking.OutCarToExit()
	go parking.CheckParking()
	go p.RenderNewCarWaitStation()
	go p.RenderUpdate()

}

func (p *ParkingView) RestartSimulation() {
	close(semQuit)
	NewParkingView(p.window)
}

func addSpace(parkingContainer *fyne.Container) {
	for j := 0; j < 5; j++ {
		parkingContainer.Add(layout.NewSpacer())
	}
}

func makeBorder() *canvas.Rectangle {
	square := canvas.NewRectangle(color.RGBA{R: 255, G: 255, B: 255, A: 0})
	square.SetMinSize(fyne.NewSquareSize(float32(30)))
	square.StrokeColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	square.StrokeWidth = float32(1)
	return square
}
