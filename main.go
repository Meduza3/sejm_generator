package main

import (
	"fmt"
	"os"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/dialog"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type Card struct {
	title string
	za map[int]bool
	przeciw map[int]bool
}

var (
	row1Count int = 0
	row2Count int = 0
)

func main() {
	card := Card{
		za:       make(map[int]bool),
		przeciw:  make(map[int]bool),
	}
	for i := 0; i < 10; i++ {
		card.za[i] = false
		card.przeciw[i] = false
	}
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("SEJM Generator")


	fileButton := widget.NewButton("Select Image", func() {
		dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if file == nil {
				return
			}
			defer file.Close()


		}, w).Show()
	})



	image := canvas.NewImageFromFile("result.png")
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(300, 300))

	updateImage := func(img *canvas.Image, title string) {
		drawImage(title)
		img.File = "result.png"
		img.Refresh()
		fmt.Println("Image updated!")
	}

	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Tytul ustawy...")
	titleEntry.OnSubmitted = func(content string) {
		updateImage(image, content)
	}
	// Create checkboxes
	imagePaths := make([]string, 20)
	for i := 0; i < 20; i++ {
		imagePaths[i] = "assets/grupy/ekolodzy.png"
	}
	images := make([]*canvas.Image, 20)
	checkboxes := make([]fyne.CanvasObject, 20)
	for i := 0; i < 20; i++ {
		checkbox := widget.NewCheck("", func(checked bool) {
			if checked {
				if i < 10 {
					card.za[i] = true
					row1Count++
				} else {
					card.przeciw[i] = true
					row2Count++
				}
			} else {
				if i < 10 {
					card.za[i - 10] = false
					row1Count--
				} else {
					card.przeciw[i - 10] = false
					row2Count--
				}
			}
			updateImage(image, titleEntry.Text)
		})
		img := canvas.NewImageFromFile(imagePaths[i])
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(50, 50))
		images[i] = img
		checkboxes[i] = container.NewVBox(img, checkbox)
	}

	// Arrange checkboxes in two rows
	row1 := container.NewHBox(checkboxes[0:10]...) // First row with 10 checkboxes
	row2 := container.NewHBox(checkboxes[10:20]...) // Second row with the remaining 10 checkboxes

	text1 := widget.NewLabel("Grupy Interesów ZA")
	text2 := widget.NewLabel("Grupy Interesów PRZECIW")


	// Save interface
	filenameEntry := widget.NewEntry()
	filenameEntry.SetPlaceHolder("Nazwa...")

	saveButton := widget.NewButton("Zapisz", func() {
		filename := filenameEntry.Text
		fmt.Printf("%s.png saved\n", filename)
	})

	loreEntry := widget.NewEntry()
	loreEntry.SetPlaceHolder("Ciekawostki...")
	// Combine the two rows into a single column
	checkColumn := container.NewVBox(fileButton, titleEntry, text1, row1, text2, row2, loreEntry)

	// Create an additional column for other content if needed
	otherColumn := container.NewVBox(
		widget.NewLabel("Result:"),
		image,
		widget.NewLabel("Nazwa Pliku:"),
		filenameEntry,
		saveButton,
	)

	// Use a grid layout with 2 columns to layout the two main columns
	grid := container.New(layout.NewGridLayout(2), checkColumn, otherColumn)

	w.SetContent(grid)
	w.ShowAndRun()
}

func drawImage(title string) {
	za_positions := GrupyCountToPositions(row1Count)
	baseFile, err := os.Open("assets/card.png")
	if err != nil {
		panic(err)
	}
	defer baseFile.Close()

	baseImage, err := png.Decode(baseFile)
	if err != nil {
		panic(err)
	}

	resultImage := image.NewRGBA(baseImage.Bounds())
	draw.Draw(resultImage, baseImage.Bounds(), baseImage, image.Point{}, draw.Over)

	for i := 0; i < row1Count; i++ {
		overlayFile, err := os.Open("assets/grupy/ekolodzy.png")
		if err != nil {
			panic(err)
		}
		defer overlayFile.Close()

		overlayImage, err := png.Decode(overlayFile)
		if err != nil {
			panic(err)
		}

		overlayPosition := za_positions[i]
		overlayPoint := image.Point{X: overlayPosition.X, Y: overlayPosition.Y}
		draw.Draw(resultImage, overlayImage.Bounds().Add(overlayPoint), overlayImage, image.Point{}, draw.Over)
	}

	addLabel(resultImage, title, 670, 1000, 200)

	// Save the result image to "result.png"
	outFile, err := os.Create("result.png")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, resultImage); err != nil {
		panic(err)
	}
}

type Point struct {
	X int
	Y int
}
func GrupyCountToPositions(count int) []Point {
	switch count {
		case 0:
			return nil
		case 1:
			return []Point{{151, 1592}}
		default:
			return []Point{{151, 1592}, {500, 1592}}
	}
}

func addLabel(img *image.RGBA, label string, x, y, fontSize int) {
	fontBytes, err := ioutil.ReadFile("assets/font.ttf")
	if err != nil {
		fmt.Println(err)
		return
	}
	ttf, err := opentype.Parse(fontBytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	face, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    float64(fontSize),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	col := color.RGBA{255, 255, 255, 255} // white color
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}
	d.DrawString(label)
}