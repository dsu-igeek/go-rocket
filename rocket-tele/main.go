package main

import (
	"github.com/igeekinc/go-rocket/pkg/core"
	"github.com/igeekinc/go-rocket/pkg/rocket"
	"log"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

func main() {
	// Load i2c drivers
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	ri := &core.RocketInfo{}

	bmpReader := core.NewBMPReader(bus, ri, 0x77)
	go bmpLoop(bmpReader)

	gpsReader, err := core.InitGPSReader(ri, "/dev/ttyS0", 9600, 8, 1)
	if err != nil {
		log.Fatal(err)
	}
	go gpsLoop(gpsReader)

	gpsReporter, err := rocket.InitRocketReporter(ri, "/dev/ttyUSB0", 57600, 8, 1)
	reporterLoop(gpsReporter)
}

func gpsLoop(gr core.GPSReader) {
	err := gr.UpdateFromGPSLoop()
	if err != nil {
		log.Fatal(err)
	}
}

func bmpLoop(bmp *core.BMPReader) {
	err := bmp.UpdateFromBMPLoop()
	if err != nil {
		log.Fatal(err)
	}
}
func reporterLoop(gr rocket.RocketReporter) {
	err := gr.RocketReporterLoop()
	if err != nil {
		log.Fatal(err)
	}
}
