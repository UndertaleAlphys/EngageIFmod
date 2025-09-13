package utils

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func CreateBackground(imagePath string, translucency float64) *canvas.Image {
	background := canvas.NewImageFromFile(imagePath)
	background.FillMode = canvas.ImageFillStretch
	background.SetMinSize(fyne.NewSize(900, 700))
	background.Translucency = translucency
	return background
}
