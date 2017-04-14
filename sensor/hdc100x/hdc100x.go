/*
	This is a library for the HDC1008 Humidity & Temp Sensor
	Designed specifically to work with the HDC100x sensor from Adafruit
	----> https://www.adafruit.com/products/2635

	These sensors use I2C to communicate, 2 pins are required to
	interface

	Written by Dan Esparza.

	For more information on the I2C operations for this sensor
	see the spec sheet located here (specifically pages 10-15):
	http://www.adafruit.com/datasheets/hdc1008.pdf
*/

package hdc100x

import (
	"time"

	"log"

	"github.com/kidoman/embd"
)

const (
	//	I2C Address
	hdc1000Address byte = 0x40

	//	Registers
	hdc1000TempRegister           byte = 0x00 // Temperature register
	hdc1000HumidityRegister       byte = 0x01 // Humidity register
	hdc1000ConfigRegister         byte = 0x02 // Config register
	hdc1000ManufacturerIDRegister byte = 0xFE // Manufacturer ID
	hdc1000DeviceIDRegister       byte = 0xFF // Device ID
	hdc1000SerialHighRegister     byte = 0xFB
	hdc1000SerialMidRegister      byte = 0xFC
	hdc1000SerialBottomRegister   byte = 0xFD

	//	Configuration register bits
	configReset          = 0x8000 // bit 15
	configHeaterEnable   = 0x2000 // bit 13
	configAquisitionMode = 0x1000 // bit 12
	configBatteryStatus  = 0x0800 // bit 11

	configTemperatureResolution      = 0x0400 // bit 10
	configTemperatureResolution14Bit = 0x0000 // bit 0
	configTemperatureResolution11Bit = 0x0400 // bit 10

	configHumidityResolutionHBit  = 0x0200 // bit 9
	configHumidityResolutionLBit  = 0x0100 // bit 8
	configHumidityResolution14Bit = 0x0000 // bit 0
	configHumidityResolution11Bit = 0x0100 // bit 8
	configHumidityResolution8Bit  = 0x0200 // bit 9

	serial1  = 0xFB
	serial2  = 0xFC
	serial3  = 0xFD
	manufID  = 0xFE
	deviceID = 0xFF
)

// HDC100x represents a HDC1000 series temperature and humidity sensor.
type HDC100x struct {
	Bus embd.I2CBus
}

// New returns a handle to a HDC100x sensor.
func New(bus embd.I2CBus) *HDC100x {
	//	Initialize by setting the config register:
	bus.WriteBytes(hdc1000Address, []byte{hdc1000ConfigRegister})

	//	Return a new object
	return &HDC100x{Bus: bus}
}

// Temperature returns the current temperature reading in Fahrenheit
func (d *HDC100x) Temperature() (float64, error) {
	//	Send the command to get the temp
	if err := d.Bus.WriteByte(hdc1000Address, hdc1000TempRegister); err != nil {
		return 0, err
	}

	time.Sleep(65 * time.Millisecond)

	//	Read 2 byte temperature data
	tempData, err := d.Bus.ReadBytes(hdc1000Address, 2)
	if err != nil {
		log.Fatal(err)
	}

	//	Combine 2 bytes into a single float64 & calculate temp
	w := float64(uint32(tempData[0])<<8 + uint32(tempData[1]))
	cTemp := (((w / 65536.0) * 165.0) - 40.0) //	Celsius
	fTemp := (cTemp * 1.8) + 32               //	Fahrenheit

	//	Return the calculated temp
	return fTemp, nil
}

// Humidity returns the current relative humidity percentage
func (d *HDC100x) Humidity() (float64, error) {
	//	Send the command to get the humidity
	if err := d.Bus.WriteByte(hdc1000Address, hdc1000HumidityRegister); err != nil {
		return 0, err
	}

	time.Sleep(65 * time.Millisecond)

	//	Read 2 byte humidity data
	humData, err := d.Bus.ReadBytes(hdc1000Address, 2)
	if err != nil {
		log.Fatal(err)
	}

	//	Combine 2 bytes into a single float64 & calculate humidity
	h := float64(uint32(humData[0])<<8 + uint32(humData[1]))
	humidity := (h / 65536.0) * 100

	//	Return the humidity
	return humidity, nil
}
