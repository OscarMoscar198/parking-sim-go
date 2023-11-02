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

type View struct {
	window               fyne.Window
	waitRectangleStation [models.MaxWait]*canvas.Rectangle
}

var Gray = color.RGBA{R: 30, G: 30, B: 30, A: 255}

var parking *models.Parking

func NewView(window fyne.Window) *View {
	View := &View{window: window}

	semQuit = make(chan bool)
	semRenderNewCarWait = make(chan bool)

	parking = models.NewParking(semRenderNewCarWait, semQuit)
	View.MakeScene()
	View.StartSimulation()

	return View
}

func (p *View) MakeScene() {

	containerView := container.New(layout.NewVBoxLayout())
	containerParkingOut := container.New(layout.NewHBoxLayout())

	containerParkingOut.Add(p.MakeWaitStation())
	containerParkingOut.Add(layout.NewSpacer())
	containerParkingOut.Add(p.MakeExitStation())
	containerParkingOut.Add(layout.NewSpacer())

	containerView.Add(containerParkingOut)
	containerView.Add(layout.NewSpacer())
	containerView.Add(p.MakeParkingLotEntrance())
	containerView.Add(layout.NewSpacer())
	containerView.Add(p.MakeEnterAndExitStation())
	containerView.Add(layout.NewSpacer())
	containerView.Add(p.MakeParking())
	containerView.Add(layout.NewSpacer())

	p.window.SetContent(containerView)
	p.window.Resize(fyne.NewSize(500, 500))
	p.window.CenterOnScreen()
}

func (p *View) MakeParking() *fyne.Container {
	parkingContainer := container.New(layout.NewGridLayout(5))
	parking.MakeParking()

	parkingArray := parking.GetParking()
	for i := 0; i < len(parkingArray); i++ {
		if i == 10 {
			addSpace(parkingContainer)
		}
		parkingContainer.Add(container.NewCenter(parkingArray[i].GetRectangle()))
	}
	return container.NewCenter(parkingContainer)
}

func (p *View) MakeWaitStation() *fyne.Container {
	parkingContainer := container.New(layout.NewGridLayout(5))
	for i := len(p.waitRectangleStation) - 1; i >= 0; i-- {
		car := models.NewSpaceVehicle().GetRectangle()
		p.waitRectangleStation[i] = car
		p.waitRectangleStation[i].Hide()
		parkingContainer.Add(p.waitRectangleStation[i])
	}
	return parkingContainer
}

func (p *View) MakeExitStation() *fyne.Container {
	out := parking.MakeOutStation()
	return container.NewCenter(out.GetRectangle())
}

func (p *View) MakeEnterAndExitStation() *fyne.Container {
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

func (p *View) MakeParkingLotEntrance() *fyne.Container {
	EntraceContainer := container.New(layout.NewGridLayout(3))
	EntraceContainer.Add(layout.NewSpacer())
	return EntraceContainer
}

func (p *View) RenderNewCarWaitStation() {
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

func (p *View) RenderUpdate() {
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

func (p *View) StartSimulation() {
	go parking.GenerateCars()
	go parking.OutCarToExit()
	go parking.CheckParking()
	go p.RenderNewCarWaitStation()
	go p.RenderUpdate()

}

func (p *View) RestartSimulation() {
	close(semQuit)
	NewView(p.window)
}

func addSpace(parkingContainer *fyne.Container) {
	for j := 0; j < 5; j++ {
		parkingContainer.Add(layout.NewSpacer())
	}
}

