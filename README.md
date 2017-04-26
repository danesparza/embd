# embd
Repository for [embd](https://github.com/kidoman/embd) based sensors

## Sensors included
* [Texas instruments HDC1008](https://learn.adafruit.com/adafruit-hdc1008-temperature-and-humidity-sensor-breakout/overview) - Temperature and Humidity sensor.  
* [Pimoroni Enviro-phat](http://docs.pimoroni.com/envirophat/) - Sensor package for Raspberry pi.  Support for the lsm303d accelerometer.

## Examples
### HDC100x
Don't forget to include the following import statements:

```Go
import (
  "github.com/danesparza/embd/sensor/hdc100x" // The sensor
  "github.com/kidoman/embd"
  _ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
)
```

Example code:
```Go
// Init I2C
if err := embd.InitI2C(); err != nil {
		panic(err)
}
defer embd.CloseI2C()

// Init the I2C bus and create a sensor object
bus := embd.NewI2CBus(1)
sensor := hdc100x.New(bus)

//	Get temperature from the sensor
fTemp, err := sensor.Temperature()
if err != nil {
  log.Fatal(err)
}
log.Printf("Temp from embd: %.2f°", fTemp)

//	Get humidity from the sensor
humidity, err := sensor.Humidity()
if err != nil {
  log.Fatal(err)
}
log.Printf("Humidity from embd: %.1f%%", humidity)
```
