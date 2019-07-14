/*
Package mlx90615 implements reading of the MLX90615 sensor using Go.

*/
package mlx90615

import (
	"fmt"
	"math"

	"github.com/d2r2/go-i2c"
)

// Appropriate registers for the MLX90615
const (
	RegisterAmbient    = 0x26
	RegisterObject     = 0x27
	RegisterEmissivity = 0x13
)

// MLX90615 is the sensor itself
type MLX90615 struct {
	// fieldless!
}

// NewMLX90615 returns new sensor instance
func NewMLX90615() *MLX90615 {
	v := &MLX90615{}
	return v
}

func readingToTemperature(raw int16) (temp float64) {
	return float64(raw)*0.02 - 273.15
}

// ReadAmbientTemperature returns the ambient temperature and an error value
func (mlxobj *MLX90615) ReadAmbientTemperature(i2cObj *i2c.I2C) (temp float64, err error) {
	reading, err := i2cObj.ReadRegS16LE(RegisterAmbient)
	if err != nil {
		return 0, err
	}
	temp = readingToTemperature(reading)
	return
}

// ReadObjectTemperature returns the object temperature and an error value
func (mlxobj *MLX90615) ReadObjectTemperature(i2cObj *i2c.I2C) (temp float64, err error) {
	reading, err := i2cObj.ReadRegS16LE(RegisterObject)
	if err != nil {
		return 0, err
	}
	temp = readingToTemperature(reading)
	return
}

// ReadEmissivity reads the emissivity of the MLX90615. Emissivity is a positive
// float no more than 1.00.
func (mlxobj *MLX90615) ReadEmissivity(i2cObj *i2c.I2C) (emissivity float64, err error) {
	reading, err := i2cObj.ReadRegS16LE(RegisterEmissivity)
	if err != nil {
		return 0, err
	}
	return float64(reading) / 16384, nil
}

// WriteEmissivity sets the emissivity of the MLX90615. Emissivity must be
// a positive float no more than 1.00, otherwise will return an error
func (mlxobj *MLX90615) WriteEmissivity(i2cObj *i2c.I2C, emissivity float64) error {
	if emissivity <= 0 || emissivity > 1 {
		return fmt.Errorf("Emissivity must be positive and less than 1. Got %v instead", emissivity)
	}
	var toWrite = int16(math.Round(emissivity * 16384))
	if err := i2cObj.WriteRegS16LE(RegisterEmissivity, toWrite); err != nil {
		return err
	}
	return nil
}

// Reset "resets" the MLX90615 by writing 1 as emissivity
func (mlxobj *MLX90615) Reset(i2cObj *i2c.I2C) error {
	return mlxobj.WriteEmissivity(i2cObj, 1)
}
