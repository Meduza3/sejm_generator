package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
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
	ArtPath  string
	Title    string
	Opinions [10]Opinion
	Effects  [7]int
	Cost     Cost
}

var (
	grupyImagePaths = []string{
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
)

func main() {
	//Initialize a card
	card := Card{
		ArtPath:  "assets/piractwo.png",
		Title:    "Depenalizacja piractwa cyfrowego",
		Opinions: [10]Opinion{For, Against, For, Against, For, For, Against, For, Against, For},
		Effects:  [7]int{0, 1, 0, -1, 1, 0, -1},
	}

	err := drawCard(card, "testcard")
	if err != nil {
		log.Fatal(err)
	}
}

func drawCard(card Card, filename string) error {
	//Load base card png
	backgroundFile, err := os.Open("assets/card.png")
	if err != nil {
		return err
	}
	backgroundImage, err := png.Decode(backgroundFile)
	defer backgroundFile.Close()
	if err != nil {
		return err
	}

	backgroundImage = addArt(backgroundImage, card.ArtPath)

	backgroundImage = addStamps(backgroundImage, card.Opinions[:])

	backgroundImage = addEffects(backgroundImage, card.Effects[:])
	if card.Cost.Value > 0 {
		backgroundImage = addCost(backgroundImage, card.Cost)
	}

	backgroundImage = addTitle(backgroundImage, card.Title)

	result, err := os.Create(filename + ".png")
	if err != nil {
		return err
	}
	defer result.Close()
	if err := png.Encode(result, backgroundImage); err != nil {
		fmt.Println("Error encoding output image:", err)
		return err
	}
	return nil
}

func addArt(backgroundImage image.Image, artPath string) *image.RGBA {

	//Add card art to backgroundImage
	artFile, err := os.Open(artPath)
	if err != nil {
		fmt.Println("addArt oops")
	}
	art, err := png.Decode(artFile)
	defer artFile.Close()
	if err != nil {
		fmt.Println("addArt oops2")
	}
	//First resize the art to fit the card.
	artWidth := uint(300 * 5)
	artHeight := uint(300 * 3)
	art = resize.Resize(artWidth, artHeight, art, resize.Lanczos3)
	resultImage := image.NewRGBA(backgroundImage.Bounds())
	draw.Draw(resultImage, backgroundImage.Bounds(), backgroundImage, image.Point{}, draw.Src)
	artPosition := image.Point{X: 0, Y: 0}
	draw.Draw(resultImage, art.Bounds().Add(artPosition), art, image.Point{}, draw.Over)
	return resultImage
}

func addStamps(backgroundImage image.Image, opinions []Opinion) *image.RGBA {
	resultImage := image.NewRGBA(backgroundImage.Bounds())
	draw.Draw(resultImage, backgroundImage.Bounds(), backgroundImage, image.Point{}, draw.Src)
	for idx, op := range opinions {
		time.Sleep(1 * time.Millisecond)
		if op == For {
			resultImage = drawStampFor(resultImage, idx)
		} else if op == Against {
			resultImage = drawStampAgainst(resultImage, idx)
		}
	}
	return resultImage
}

func drawStampFor(backgroundImage image.Image, stamp_id int) *image.RGBA {
	resultImage := image.NewRGBA(backgroundImage.Bounds())
	draw.Draw(resultImage, backgroundImage.Bounds(), backgroundImage, image.Point{}, draw.Src)

	stampFile, err := os.Open(grupyImagePaths[stamp_id])
	if err != nil {
		fmt.Println("drawStampFor oops")
	}
	stampImage, err := png.Decode(stampFile)
	if err != nil {
		fmt.Println("drawStampFor oops2")
	}

	stampImage = resize.Resize(300, 300, stampImage, resize.Lanczos3)
	x := rand.Intn(345 + 1)
	y := rand.Intn(1725-1200+1) + 1200
	stampPosition := image.Point{X: x, Y: y}
	draw.Draw(resultImage, stampImage.Bounds().Add(stampPosition), stampImage, image.Point{}, draw.Over)
	return resultImage
}

func drawStampAgainst(backgroundImage image.Image, stamp_id int) *image.RGBA {
	resultImage := image.NewRGBA(backgroundImage.Bounds())
	draw.Draw(resultImage, backgroundImage.Bounds(), backgroundImage, image.Point{}, draw.Src)

	stampFile, err := os.Open(grupyImagePaths[stamp_id])
	if err != nil {
		fmt.Println("drawStampAgainst oops")
	}
	stampImage, err := png.Decode(stampFile)
	if err != nil {
		fmt.Println("drawStampAgainst oops2")
	}

	stampImage = resize.Resize(300, 300, stampImage, resize.Lanczos3)
	x := rand.Intn(1200-765+1) + 765
	y := rand.Intn(1725-1200+1) + 1200
	stampPosition := image.Point{X: x, Y: y}
	draw.Draw(resultImage, stampImage.Bounds().Add(stampPosition), stampImage, image.Point{}, draw.Over)
	return resultImage
}

func addEffects(backgroundImage image.Image, effects []int) *image.RGBA {
	resultImage := image.NewRGBA(backgroundImage.Bounds())
	draw.Draw(resultImage, backgroundImage.Bounds(), backgroundImage, image.Point{}, draw.Src)
	goodOffset := 0
	badOffset := 0
	for idx, val := range effects {
		if val != 0 {
			if val > 0 {
				resultImage = drawGoodEffect(resultImage, idx*2+1, goodOffset)
				goodOffset++
			} else { // val < 0
				resultImage = drawBadEffect(resultImage, idx*2, badOffset)
				badOffset++
			}
		}
	}
	return resultImage
}

func drawGoodEffect(backgroundImage image.Image, effect_id, offset int) *image.RGBA {
	resultImage := image.NewRGBA(backgroundImage.Bounds())
	draw.Draw(resultImage, backgroundImage.Bounds(), backgroundImage, image.Point{}, draw.Src)
	effectFile, err := os.Open(wskaznikiImagePaths[effect_id])
	if err != nil {
		fmt.Println("drawGoodEffect oops")
	}
	effectImage, err := png.Decode(effectFile)
	if err != nil {
		fmt.Println("drawGoodEffect oops2")
	}
	effectImage = resize.Resize(360, 300, effectImage, resize.Lanczos3)
	x := 765 + offset*150
	y := 2010
	effectPosition := image.Point{X: x, Y: y}
	draw.Draw(resultImage, effectImage.Bounds().Add(effectPosition), effectImage, image.Point{}, draw.Over)
	return resultImage
}

func drawBadEffect(backgroundImage image.Image, effect_id, offset int) *image.RGBA {
	resultImage := image.NewRGBA(backgroundImage.Bounds())
	draw.Draw(resultImage, backgroundImage.Bounds(), backgroundImage, image.Point{}, draw.Src)
	effectFile, err := os.Open(wskaznikiImagePaths[effect_id])
	if err != nil {
		fmt.Println("drawBadEffect oops")
	}
	effectImage, err := png.Decode(effectFile)
	if err != nil {
		fmt.Println("drawBadEffect oops2")
	}
	effectImage = resize.Resize(360, 300, effectImage, resize.Lanczos3)
	x := 10 + offset*150
	y := 2010
	effectPosition := image.Point{X: x, Y: y}
	draw.Draw(resultImage, effectImage.Bounds().Add(effectPosition), effectImage, image.Point{}, draw.Over)
	return resultImage
}

func addCost(backgroundImage image.Image, cost Cost) *image.RGBA {
	resultImage := image.NewRGBA(backgroundImage.Bounds())
	draw.Draw(resultImage, backgroundImage.Bounds(), backgroundImage, image.Point{}, draw.Src)
	costFile, err := os.Open(costToFilepath(cost))
	if err != nil {
		fmt.Println("addCost oops")
	}
	costImage, err := png.Decode(costFile)
	if err != nil {
		fmt.Println("addCost oops2")
	}
	x := 0
	y := 20
	costPosition := image.Point{X: x, Y: y}
	draw.Draw(resultImage, costImage.Bounds().Add(costPosition), costImage, image.Point{}, draw.Over)
	return resultImage
}

func costToFilepath(cost Cost) string {
	filepath := "assets/ceny/"
	switch cost.Value {
	case -3:
		filepath += "minus3"
	case -2:
		filepath += "minus2"
	case -1:
		filepath += "minus1"
	case 1:
		filepath += "plus1"
	case 2:
		filepath += "plus2"
	case 3:
		filepath += "plus3"
	default:
		filepath += "minus1"
	}
	if cost.Recurring {
		filepath += "raz"
	}
	filepath += ".png"
	return filepath
}

func addTitle(backgroundImage image.Image, title string) *image.RGBA {
	resultImage := image.NewRGBA(backgroundImage.Bounds())
	draw.Draw(resultImage, backgroundImage.Bounds(), backgroundImage, image.Point{}, draw.Src)
	fontBytes, err := os.ReadFile("assets/sylfaen.ttf")
	if err != nil {
		fmt.Println("addTitle oops")
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println("addTitle oops2")
	}
	context := freetype.NewContext()
	context.SetFont(font)
	fontSize := 150
	context.SetFontSize(float64(fontSize))
	context.SetClip(resultImage.Bounds())
	context.SetDst(resultImage)
	context.SetSrc(image.NewUniform(color.White))

	face := truetype.NewFace(font, &truetype.Options{
		Size: 150,
	})

	if len(title) > 11 {
		mid := len(title) / 2
		left := strings.LastIndex(title[:mid], " ")
		right := strings.Index(title[mid:], " ") + mid

		if left == -1 {
			left = 0
		}
		if right == -1 || right >= len(title) {
			right = len(title)
		}
		if mid-left < right-mid {
			mid = left
		} else {
			mid = right
		}
		title = title[:mid] + "\n" + title[mid+1:]
	}

	lines := strings.Split(title, "\n")

	y := 1020
	for _, line := range lines {
		lineWidth := textWidth(face, line)
		x := (resultImage.Bounds().Dx() - lineWidth) / 2 // Center text horizontally
		pt := freetype.Pt(x, y)
		_, err = context.DrawString(line, pt)
		if err != nil {
			fmt.Println("addTitle oop3")
		}
		y += int(context.PointToFixed(float64(fontSize)) >> 6)
	}

	return resultImage
}

func textWidth(face font.Face, text string) int {
	width := 0
	for _, char := range text {
		aw, ok := face.GlyphAdvance(rune(char))
		if !ok {
			continue
		}
		width += int(aw >> 6)
	}
	return width
}

func randomCard() Card {
	rand.Seed(time.Now().UnixNano())

	// Generate random opinions
	var opinions [10]Opinion
	for i := range opinions {
		opinions[i] = Opinion(rand.Intn(3)) // Random value between 0 (Indifferent), 1 (For), and 2 (Against)
	}

	// Generate random effects
	var effects [7]int
	nonZeroCount := rand.Intn(2) + 2 // 2 or 3 non-zero effects
	for i := 0; i < nonZeroCount; i++ {
		index := rand.Intn(7)
		for effects[index] != 0 { // Ensure unique non-zero indices
			index = rand.Intn(7)
		}
		effects[index] = rand.Intn(10) + 1 - 5 // Random non-zero value between 1 and 10
	}

	// Generate random cost
	cost := Cost{
		Value:     rand.Intn(6) - 3,  // Random cost value between -3 and 3
		Recurring: rand.Intn(2) == 1, // Random boolean for recurring cost
	}

	// Generate random title
	titles := []string{"Zakaz handlu w niedziele", "Prywatyzacja uniwersytetów", "Małżeństwa jednopłciowe z adopcją"} // Example titles
	title := titles[rand.Intn(len(titles))]

	return Card{
		Title:    title,
		Opinions: opinions,
		Effects:  effects,
		Cost:     cost,
	}
}
