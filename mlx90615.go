/*
Package mlx90615 implements reading of the MLX90615 sensor using Go.
*/
package mlx90615

import (
	"math"
	"strconv"

	"periph.io/x/periph/conn/i2c/i2creg"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/host"
)

// Appropriate registers for the MLX90615
const (
	RegisterAmbient    = 0x26
	RegisterObject     = 0x27
	RegisterEmissivity = 0x13
)

func readingToTemperature(raw []byte) (temp float64) {
	if len(raw) != 2 {
		return -1 // no way...
	}
	temp = float64((uint16(raw[1]&0x007F) << 8) | uint16(raw[0]))
	temp = temp*0.02 - 273.15
	return
}

func emissivityToBytes(emissivity float64) []byte {
	// i guess little endian?
	bufferInt := uint16(math.Round(16384 * emissivity))
	return []byte{byte(bufferInt & 0xFF), byte(bufferInt << 8)}
}

// MLX90615 is the sensor itself
type MLX90615 struct {
	// unexported fields ;)
	bus *i2c.BusCloser
	dev *i2c.Dev
}

// NewMLX90615 returns new sensor instance
func NewMLX90615(addr uint8, bus int) (*MLX90615, error) {
	if _, err := host.Init(); err != nil {
		return nil, err
	}
	busObj, err := i2creg.Open(strconv.Itoa(bus))
	if err != nil {
		return nil, err
	}
	devObj := i2c.Dev{Bus: busObj, Addr: uint16(addr)}
	return &MLX90615{&busObj, &devObj}, nil
}

func (mlx *MLX90615) readRegister(register byte, size int) ([]byte, error) {
	buffer := make([]byte, size)
	if err := mlx.dev.Tx([]byte{register}, buffer); err != nil {
		return nil, err
	}
	return buffer, nil
}

// ReadAmbientTemperature returns the ambient temperature and an error value
func (mlx *MLX90615) ReadAmbientTemperature() (temp float64, err error) {
	buffer, err := mlx.readRegister(RegisterAmbient, 2)
	if err != nil {
		return 0, err
	}
	return readingToTemperature(buffer), nil
}

// ReadObjectTemperature returns the object temperature and an error value
func (mlx *MLX90615) ReadObjectTemperature() (temp float64, err error) {
	buffer, err := mlx.readRegister(RegisterObject, 2)
	if err != nil {
		return 0, err
	}
	return readingToTemperature(buffer), nil
}

// ReadEmissivity reads the emissivity of the MLX90615.
// Emissivity is a positive float no more than 1.00.
func (mlx *MLX90615) ReadEmissivity() (emissivity float64, err error) {
	buffer, err := mlx.readRegister(RegisterAmbient, 2)
	if err != nil {
		return 0, err
	}
	return float64(uint16(buffer[1])<<8|uint16(buffer[0])) / 16384, nil
}
