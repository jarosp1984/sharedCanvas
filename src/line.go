package main

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	canvasWidth  = 800
	canvasHeight = 600
	minLineWidth = 1
	maxLineWidth = 50
)

var colorPattern = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

type Line struct {
	ID    int    `json:"id"`
	X1    int    `json:"x1"`
	Y1    int    `json:"y1"`
	X2    int    `json:"x2"`
	Y2    int    `json:"y2"`
	Color string `json:"color"`
	Width int    `json:"width"`
}

func (line Line) Validate() error {
	if err := validateCoordinate("x1", line.X1, 0, canvasWidth); err != nil {
		return err
	}
	if err := validateCoordinate("y1", line.Y1, 0, canvasHeight); err != nil {
		return err
	}
	if err := validateCoordinate("x2", line.X2, 0, canvasWidth); err != nil {
		return err
	}
	if err := validateCoordinate("y2", line.Y2, 0, canvasHeight); err != nil {
		return err
	}
	if !colorPattern.MatchString(line.Color) {
		return errors.New("color must be a hex value like #000000")
	}
	if line.Width < minLineWidth || line.Width > maxLineWidth {
		return fmt.Errorf("width must be between %d and %d", minLineWidth, maxLineWidth)
	}

	return nil
}

func validateCoordinate(name string, value int, min int, max int) error {
	if value < min || value > max {
		return fmt.Errorf("%s must be between %d and %d", name, min, max)
	}

	return nil
}
