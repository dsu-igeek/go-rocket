package ui

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/adrianmo/go-nmea"
	"github.com/igeekinc/go-rocket/pkg/core"
	"github.com/igeekinc/go-rocket/pkg/ground"
	"time"
)

func RunRocketTrackerUI(rocketReceiver *ground.RocketReceiver) {
	a := app.New()
	w := a.NewWindow("Rocket Tracker")
	//w.SetFullScreen(true)
	//w.SetContent(widget.NewLabel("Hello World!"))
	kRocketPosLabel := widget.NewLabel("Rocket Pos: ")
	rocketPosLabel := widget.NewLabel("")
	kOurPosLabel := widget.NewLabel("Our Pos: ")
	ourPosLabel := widget.NewLabel("")
	kDistanceLabel := widget.NewLabel("Distance: ")
	distanceLabel := widget.NewLabel("")
	kLaunchModeLabel := widget.NewLabel("Launch mode: ")
	launchModeLabel := widget.NewLabel("false")
	rtu := rocketTrackerUI{
		rocketReceiver:       rocketReceiver,
		gpsPosLabel:          rocketPosLabel,
		ourPosLabel:          ourPosLabel,
		distanceLabel:        distanceLabel,
		launchModeStateLabel: launchModeLabel,
	}
	launchButton := widget.NewButton("Launch", func() {
		fmt.Println("Launch pressed")
		rocketReceiver.SendLaunchMode()
	})
	infoCanvas := container.New(layout.NewGridLayout(2),
		kRocketPosLabel,
		rocketPosLabel,
		kOurPosLabel,
		ourPosLabel,
		kDistanceLabel,
		distanceLabel,
		kLaunchModeLabel,
		launchModeLabel,
		launchButton)
	w.SetContent(infoCanvas)

	go rtu.updateLoop()
	w.ShowAndRun()
}

type rocketTrackerUI struct {
	rocketReceiver       *ground.RocketReceiver
	gpsPosLabel          *widget.Label
	ourPosLabel          *widget.Label
	distanceLabel        *widget.Label
	launchModeStateLabel *widget.Label
}

func (recv *rocketTrackerUI) updateLoop() {
	for {
		updateLabelWithGPS(recv.gpsPosLabel, recv.rocketReceiver.RocketInfo.GPS, recv.rocketReceiver.RocketInfo.Altitude)
		updateLabelWithGPS(recv.ourPosLabel, recv.rocketReceiver.GPS, 0)
		distance := core.Distance(recv.rocketReceiver.GPS.Latitude, recv.rocketReceiver.GPS.Longitude,
			recv.rocketReceiver.RocketInfo.GPS.Latitude, recv.rocketReceiver.RocketInfo.GPS.Longitude)
		distanceStr := fmt.Sprintf("%0f M", distance)
		fmt.Println(distanceStr)
		recv.distanceLabel.SetText(distanceStr)
		if recv.rocketReceiver.RocketInfo.Recording {
			recv.launchModeStateLabel.SetText("true")
		} else {
			recv.launchModeStateLabel.SetText("false")
		}
		time.Sleep(time.Second) // Would be nicer if we waited on a channel from RocketInfo
	}
}

func updateLabelWithGPS(label *widget.Label, gps nmea.GGA, altitude float64) {
	labelStr := ""
	if gps.FixQuality == nmea.Invalid {
		labelStr = "Invalid GPS"
	} else {
		labelStr = fmt.Sprintf("Lat: %09.5f Lon:%09.5f Alt: %.2fM", gps.Latitude, gps.Longitude, altitude)
	}
	fmt.Println(labelStr)
	label.SetText(labelStr)
}
