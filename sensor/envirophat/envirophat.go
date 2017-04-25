/*
	This is a library for the Enviro-phat Sensor
	Designed specifically to work with the Enviro-phat sensor from Pimoroni
	----> https://shop.pimoroni.com/products/enviro-phat

	These sensors use I2C to communicate, 2 pins are required to
	interface

	Written by Dan Esparza.

	For more information on the I2C operations for the accelerometer sensor
	see the spec sheet located here (specifically pages 22-29):
	http://www.st.com/resource/en/datasheet/lsm303d.pdf
*/

package envirophat

import (
	"encoding/binary"
	"math"

	"github.com/kidoman/embd"
)

const (
	//	I2C Addresses
	lsm303dAddress byte = 0x1D //	Accelerometer

	//	Configuration register bits
	CTRL_REG0 = 0x1F
	CTRL_REG1 = 0x20
	CTRL_REG2 = 0x21
	CTRL_REG3 = 0x22
	CTRL_REG4 = 0x23
	CTRL_REG5 = 0x24
	CTRL_REG6 = 0x25
	CTRL_REG7 = 0x26

	//	Acceleromter outputs:
	OUT_X_L_A = 0x28
	OUT_X_H_A = 0x29
	OUT_Y_L_A = 0x2A
	OUT_Y_H_A = 0x2B
	OUT_Z_L_A = 0x2C
	OUT_Z_H_A = 0x2D

	//	Mag scales
	magScale2  = 0x00 // full-scale is +/- 2 Gauss
	magScale4  = 0x20 // +/- 4 Guass
	magScale8  = 0x40 // +/- 8 Guass
	magScale12 = 0x60 // +/- 12 Guass

	//	Accelerometer scale
	ACCEL_SCALE = 2 // +/- 2gs
)

// Envirophat represents a Pimoroni enviro-phat with multiple sensor types.
type Envirophat struct {
	Bus embd.I2CBus
}

// New returns a handle to a Envirophat object.
func New(bus embd.I2CBus) *Envirophat {
	//	Initialize accelerometer:
	bus.WriteByteToReg(lsm303dAddress, CTRL_REG1, byte(0x57))          // 0x57 = ODR=50hz, all accel axes on ## maybe 0x27 is Low Res?
	bus.WriteByteToReg(lsm303dAddress, CTRL_REG2, byte((3<<6)|(0<<3))) // set full scale +/- 2g
	bus.WriteByteToReg(lsm303dAddress, CTRL_REG3, byte(0x00))          // No interrupt
	bus.WriteByteToReg(lsm303dAddress, CTRL_REG4, byte(0x00))          // No interrupt
	bus.WriteByteToReg(lsm303dAddress, CTRL_REG5, byte((4 << 2)))      // 0x10 = mag 50Hz output rate
	bus.WriteByteToReg(lsm303dAddress, CTRL_REG6, byte(magScale2))     // Magnetic Scale +/1 1.3 Guass
	bus.WriteByteToReg(lsm303dAddress, CTRL_REG7, byte(0x00))          // 0x00 continuous conversion mode

	//	Return a new object
	return &Envirophat{Bus: bus}
}

// Accelerometer returns the current accelerometer vector
func (d *Envirophat) Accelerometer() (x float64, y float64, z float64, err error) {

	////////////////////
	// 	READ 'X' AXIS:
	////////////////////
	valx1, err := d.Bus.ReadByteFromReg(lsm303dAddress, byte(OUT_X_H_A))
	if err != nil {
		return 0, 0, 0, err
	}
	valx2, err := d.Bus.ReadByteFromReg(0x1D, byte(OUT_X_L_A))
	if err != nil {
		return 0, 0, 0, err
	}

	//	Two's complement and scale adjustment:
	valx2s := int16(binary.BigEndian.Uint16([]byte{valx1, valx2}))
	x = float64(valx2s) / math.Pow(2, 15) * ACCEL_SCALE

	////////////////////
	// 	READ 'Y' AXIS:
	////////////////////
	valy1, err := d.Bus.ReadByteFromReg(lsm303dAddress, byte(OUT_Y_H_A))
	if err != nil {
		return 0, 0, 0, err
	}
	valy2, err := d.Bus.ReadByteFromReg(0x1D, byte(OUT_Y_L_A))
	if err != nil {
		return 0, 0, 0, err
	}

	//	Two's complement and scale adjustment:
	valy2s := int16(binary.BigEndian.Uint16([]byte{valy1, valy2}))
	y = float64(valy2s) / math.Pow(2, 15) * ACCEL_SCALE

	////////////////////
	// 	READ 'Z' AXIS:
	////////////////////
	valz1, err := d.Bus.ReadByteFromReg(lsm303dAddress, byte(OUT_Z_H_A))
	if err != nil {
		return 0, 0, 0, err
	}
	valz2, err := d.Bus.ReadByteFromReg(0x1D, byte(OUT_Z_L_A))
	if err != nil {
		return 0, 0, 0, err
	}

	//	Two's complement and scale adjustment:
	valz2s := int16(binary.BigEndian.Uint16([]byte{valz1, valz2}))
	z = float64(valz2s) / math.Pow(2, 15) * ACCEL_SCALE

	return x, y, z, nil
}
