package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type Image struct {
}

type Opinion int

const (
	Indifferent Opinion = iota
	For
	Against
)

type Cost struct {
	Value     int
	Recurring bool
}

// Values 0-10
type Card struct {
	Art     Image
	Title   string
	For     []Opinion
	Against []Opinion
	Effects []int
	Cost    Cost
}

var (
	row1Count  int = 0
	row2Count  int = 0
	card       Card
	imagePaths = []string{
		"assets/grupy/Socjalni.png",
		"assets/grupy/Ekolodzy.png",
		"assets/grupy/Centrysci.png",
		"assets/grupy/Globalisci.png",
		"assets/grupy/Katolicy.png",
		"assets/grupy/Narodowcy.png",
		"assets/grupy/Progresywni.png",
		"assets/grupy/Przedsiebiorcy.png",
		"assets/grupy/Robotnicy.png",
		"assets/grupy/Samorzadowcy.png",
		"assets/grupy/Socjalni.png",
		"assets/grupy/Ekolodzy.png",
		"assets/grupy/Centrysci.png",
		"assets/grupy/Globalisci.png",
		"assets/grupy/Katolicy.png",
		"assets/grupy/Narodowcy.png",
		"assets/grupy/Progresywni.png",
		"assets/grupy/Przedsiebiorcy.png",
		"assets/grupy/Robotnicy.png",
		"assets/grupy/Samorzadowcy.png",
	}
	wskaznikiImagePaths = []string{
		"assets/wsk/InflacjaMinus.png",
		"assets/wsk/InflacjaPlus.png",
		"assets/wsk/DochodMinus.png",
		"assets/wsk/DochodPlus.png",
		"assets/wsk/ZatrudnienieMinus.png",
		"assets/wsk/ZatrudnieniePlus.png",
		"assets/wsk/BezpieczenstwoMinus.png",
		"assets/wsk/BezpieczenstwoPlus.png",
		"assets/wsk/WolnoscMinus.png",
		"assets/wsk/WolnoscPlus.png",
		"assets/wsk/InfrastrukturaMinus.png",
		"assets/wsk/InfrastrukturaPlus.png",
		"assets/wsk/ZdrowieMinus.png",
		"assets/wsk/ZdrowiePlus.png",
	}
	sliderValues = make([]int, 7)
)

func main() {
	card = Card{
		za:      make([]bool, 10),
		przeciw: make([]bool, 10),
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
	images := make([]*canvas.Image, 20)
	checkboxes := make([]fyne.CanvasObject, 20)
	for i := 0; i < 20; i++ {
		checkbox := widget.NewCheck("", func(checked bool) {
			if checked {
				if i < 10 {
					card.za[i] = true
					row1Count++
				} else {
					card.przeciw[i-10] = true
					row2Count++
				}
			} else {
				if i < 10 {
					card.za[i] = false
					row1Count--
				} else {
					card.przeciw[i-10] = false
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
	row1 := container.NewHBox(checkboxes[0:10]...)  // First row with 10 checkboxes
	row2 := container.NewHBox(checkboxes[10:20]...) // Second row with the remaining 10 checkboxes

	text1 := widget.NewLabel("Grupy Interesów ZA")
	text2 := widget.NewLabel("Grupy Interesów PRZECIW")

	sliderRows := make([]fyne.CanvasObject, 0, 6)

	for i := 0; i < 7; i++ {
		leftImage := canvas.NewImageFromFile(wskaznikiImagePaths[i*2])
		leftImage.FillMode = canvas.ImageFillContain
		rightImage := canvas.NewImageFromFile(wskaznikiImagePaths[i*2+1])
		rightImage.FillMode = canvas.ImageFillContain

		leftImage.SetMinSize(fyne.NewSize(50, 50))
		rightImage.SetMinSize(fyne.NewSize(50, 50))
		slider := widget.NewSlider(-3, 3)
		slider.Value = 0
		slider.Step = 1

		slider.OnChanged = func(value float64) {
			sliderValues[i] = int(value)
			updateImage(image, titleEntry.Text)
		}

		row := container.New(layout.NewGridLayoutWithColumns(3),
			leftImage,
			container.New(layout.NewVBoxLayout(),
				layout.NewSpacer(),
				slider,
				layout.NewSpacer(),
			),
			rightImage,
		)
		sliderRows = append(sliderRows, row)
	}
	sliders := container.NewVBox(sliderRows...)
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
	checkColumn := container.NewVBox(fileButton, titleEntry, text1, row1, text2, row2, sliders, loreEntry)

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
	zaPositions := generateRandomPointsZA(row1Count)
	przeciwPositions := generateRandomPointsPRZECIW(row2Count)
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

	index := 0
	for i, checked := range card.za {
		if checked {
			overlayFile, err := os.Open(imagePaths[i])
			if err != nil {
				panic(err)
			}
			defer overlayFile.Close()

			overlayImage, err := png.Decode(overlayFile)
			if err != nil {
				panic(err)
			}

			overlayPosition := zaPositions[index]
			overlayPoint := image.Point{X: overlayPosition.X, Y: overlayPosition.Y}
			draw.Draw(resultImage, overlayImage.Bounds().Add(overlayPoint), overlayImage, image.Point{}, draw.Over)
			index++
		}
	}

	index = 0
	for i, checked := range card.przeciw {
		if checked {
			overlayFile, err := os.Open(imagePaths[i])
			if err != nil {
				panic(err)
			}
			defer overlayFile.Close()

			overlayImage, err := png.Decode(overlayFile)
			if err != nil {
				panic(err)
			}

			overlayPosition := przeciwPositions[index]
			overlayPoint := image.Point{X: overlayPosition.X, Y: overlayPosition.Y}
			draw.Draw(resultImage, overlayImage.Bounds().Add(overlayPoint), overlayImage, image.Point{}, draw.Over)
			index++
		}
	}

	for i, value := range sliderValues {
		if value != 0 {
			if value > 0 {
				overlayFile, _ := os.Open(wskaznikiImagePaths[i*2+1])
				overlayImage, _ := png.Decode(overlayFile)
				overlayPoint := image.Point{X: 113, Y: 2000}
				draw.Draw(resultImage, overlayImage.Bounds().Add(overlayPoint), overlayImage, image.Point{}, draw.Over)
			} else {

				overlayFile, _ := os.Open(wskaznikiImagePaths[i*2])
				overlayImage, _ := png.Decode(overlayFile)
				overlayPoint := image.Point{X: 800, Y: 2000}
				draw.Draw(resultImage, overlayImage.Bounds().Add(overlayPoint), overlayImage, image.Point{}, draw.Over)
			}
		}
	}

	addLabel(resultImage, title, 470, 1100, 200)

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

func generateRandomPointsZA(n int) []Point {
	ymin := 1200
	ymax := 1750
	xmin := 20
	xmax := 500

	points := make([]Point, n)
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	for i := 0; i < n; i++ {
		points[i] = Point{
			X: rand.Intn(xmax-xmin+1) + xmin,
			Y: rand.Intn(ymax-ymin+1) + ymin,
		}
	}
	return points
}

func generateRandomPointsPRZECIW(n int) []Point {
	ymin := 1200
	ymax := 1750
	xmin := 750
	xmax := 1250

	points := make([]Point, n)
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	for i := 0; i < n; i++ {
		points[i] = Point{
			X: rand.Intn(xmax-xmin+1) + xmin,
			Y: rand.Intn(ymax-ymin+1) + ymin,
		}
	}
	return points
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
