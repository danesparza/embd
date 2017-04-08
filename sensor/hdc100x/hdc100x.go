/*
	This is a library for the HDC1000 Humidity & Temp Sensor
	Designed specifically to work with the HDC1000 sensor from Adafruit
	----> https://www.adafruit.com/products/2635

	These sensors use I2C to communicate, 2 pins are required to
	interface

	Written by Dan Esparza.

	Taken in large part from the original Arduino based
	source provided by Limor Fried / Ladyada here:
	https://github.com/adafruit/Adafruit_HDC1000_Library/
*/

package hdc100x

import (
	"encoding/binary"
	"sync"
	"time"

	"log"

	"github.com/kidoman/embd"
)

const (
	// I2CAddress is the address of this device on the i2c bus
	I2CAddress = 0x40

	hdc1000Temp     = 0x00
	hdc1000Humidity = 0x01

	configLocation = 0x02
	configRST      = (1 << 15)
	configHeat     = (1 << 13)
	configMode     = (1 << 12)
	configBatt     = (1 << 11)
	configTres14   = 0
	configTres11   = (1 << 10)
	configHres14   = 0
	configHres11   = (1 << 8)
	configHres8    = (1 << 9)

	serial1  = 0xFB
	serial2  = 0xFC
	serial3  = 0xFD
	manufID  = 0xFE
	deviceID = 0xFF

	pollDelay = 250
)

// HDC100x represents a HDC1000 series temperature and humidity sensor.
type HDC100x struct {
	Bus  embd.I2CBus
	Poll int

	//	The stuff below this can change
	oss uint

	ac1, ac2, ac3      int16
	ac4, ac5, ac6      uint16
	b1, b2, mb, mc, md int16
	b5                 int32
	calibrated         bool
	cmu                sync.RWMutex

	//	I like the idea of using channels to communicate values
	temps     chan uint16
	pressures chan int32
	altitudes chan float64
	quit      chan struct{}
}

// New returns a handle to a HDC100x sensor.
func New(bus embd.I2CBus) *HDC100x {
	return &HDC100x{Bus: bus, Poll: pollDelay}
}

// Temperature returns the current temperature reading.
func (d *HDC100x) Temperature() (float64, error) {
	select {
	case t := <-d.temps:
		temp := float64(t) / 10
		return temp, nil
	default:
		log.Println("hdc100x: no temps available... measuring")
		/*
			t, err := d.measureTemp()
			if err != nil {
				return 0, err
			}
			temp := float64(t) / 10
		*/
		temp := float64(0)
		return temp, nil
	}
}

func (d *HDC100x) measureTemp() (uint32, error) {
	/*
		if err := d.calibrate(); err != nil {
			return 0, err
		}
	*/

	utemp, err := d.readUncompensatedTemp()
	if err != nil {
		return 0, err
	}

	//	The Adafruit example goes through some machinations
	//	to find the actual temperature.  I figure that would go right here
	//	temp := d.calcTemp(utemp)

	//	For now, just return the uncalculated temp
	return utemp, nil
}

func (d *HDC100x) readUncompensatedTemp() (uint32, error) {

	if err := d.Bus.WriteByte(I2CAddress, hdc1000Temp); err != nil {
		return 0, err
	}
	time.Sleep(50)

	//	I think we need to read multiple bytes based on
	//	this line:
	//	https://github.com/adafruit/Adafruit_HDC1000_Library/blob/master/Adafruit_HDC1000.cpp#L120
	tempBytes, err := d.Bus.ReadBytes(I2CAddress, 4)
	if err != nil {
		return 0, err
	}

	//	Because we now have a byte slice, we need to
	//	convert to an unsigned int:
	temp := binary.BigEndian.Uint32(tempBytes)

	return temp, nil
}

// Run starts the sensor data acquisition loop.
func (d *HDC100x) Run() {

	//	Spawn a new goroutine
	go func() {
		d.quit = make(chan struct{})
		timer := time.Tick(time.Duration(d.Poll) * time.Millisecond)

		var temp uint16
		var pressure int32
		var altitude float64

		for {
			select {
			case <-timer:
				//	Measure temp
				/*
					t, err := d.measureTemp()
					if err == nil {
						temp = t
					}
					if err == nil && d.temps == nil {
						d.temps = make(chan uint16)
					}
				*/

				//	Measure humidity
				/*
					p, a, err := d.measurePressureAndAltitude()
					if err == nil {
						pressure = p
						altitude = a
					}
					if err == nil && d.pressures == nil && d.altitudes == nil {
						d.pressures = make(chan int32)
						d.altitudes = make(chan float64)
					}
				*/
			case d.temps <- temp:
			case d.pressures <- pressure:
			case d.altitudes <- altitude:
			case <-d.quit:
				d.temps = nil
				d.pressures = nil
				d.altitudes = nil
				return
			}
		}
	}()

	return
}

// Close the connection to the sensor
func (d *HDC100x) Close() {
	if d.quit != nil {
		d.quit <- struct{}{}
	}
}
