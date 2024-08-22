package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
)

type Opinion int

const (
	ExtraAgainst Opinion = iota - 2
	Against
	Indifferent
	For
	ExtraFor
)

type Cost struct {
	Value       int
	Currency    Currency
	AddToBudget bool
}

type Currency string

const (
	Trust Currency = "Trust"
	Cash  Currency = "Cash"
)

// Values 0-10
type LegislationCard struct {
	ArtPath  string
	Title    string
	Opinions [10]Opinion
	Effects  [7]int
	Cost     Cost
}

var (
	grupyImagePaths = []string{
		"assets/grupy/kat.png",
		"assets/grupy/prg.png",
		"assets/grupy/soc.png",
		"assets/grupy/pzc.png",
		"assets/grupy/rob.png",
		"assets/grupy/nar.png",
		"assets/grupy/glo.png",
		"assets/grupy/eko.png",
		"assets/grupy/sam.png",
		"assets/grupy/cen.png",
	}
	wskaznikiImagePaths = []string{
		"assets/wsk/DochodMinus.png",         //0
		"assets/wsk/DochodPlus.png",          //1
		"assets/wsk/ZatrudnienieMinus.png",   //2
		"assets/wsk/ZatrudnieniePlus.png",    //3
		"assets/wsk/InfrastrukturaMinus.png", //4
		"assets/wsk/InfrastrukturaPlus.png",  //5
		"assets/wsk/WolnoscMinus.png",        //6
		"assets/wsk/WolnoscPlus.png",
		"assets/wsk/BezpieczenstwoMinus.png",
		"assets/wsk/BezpieczenstwoPlus.png",
		"assets/wsk/ZdrowieMinus.png",
		"assets/wsk/ZdrowiePlus.png",
		"assets/wsk/InflacjaMinus.png",
		"assets/wsk/InflacjaPlus.png",
	}
)

const cm = 300

func ParseInput(input string) LegislationCard {
	inputParts := strings.Split(input, ".")
	artPath := inputParts[0]
	title := inputParts[1]
	opinions := stringToOpinions(inputParts[2])
	effects := stringToEffects(inputParts[3])
	cost, _ := strconv.Atoi(inputParts[4])

	return NewLegislationCard(artPath, title, opinions, effects, cost)
}
func NewLegislationCard(artPath string, title string, opinions [10]Opinion, effects [7]int, cost int) LegislationCard {

	return LegislationCard{
		ArtPath:  artPath + ".png",
		Title:    title,
		Opinions: opinions,
		Effects:  effects,
		Cost: Cost{
			Value:       cost,
			Currency:    Cash,
			AddToBudget: false,
		},
	}
}

// Średnica kółka: 300px
// 1cm = 300px

func main() {
	dirName := "generated"

	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.Mkdir(dirName, 0755)
		if err != nil {
			fmt.Printf("Critical error - create directory %s yourself\n", dirName)
			os.Exit(1)
		}
	}
	for {
		fmt.Println("Input card code:")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "exit" {
			os.Exit(0)
		}
		card := ParseInput(input)

		err := drawCard(card, dirName+"/"+card.ArtPath)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Generated %s -> wygenerowane/%s\n", card.ArtPath, card.ArtPath)
		fmt.Printf("\n")
	}
}
func drawCard(card LegislationCard, filename string) error {
	//Load base card png
	backgroundFile, err := os.Open("assets/print_card.png")
	if err != nil {
		return err
	}
	backgroundImage, err := png.Decode(backgroundFile)
	defer backgroundFile.Close()
	if err != nil {
		return err
	}

	backgroundImage, err = addArt(backgroundImage, card.ArtPath)
	if err != nil {
		return fmt.Errorf("in drawCard(): %v", err)
	}

	backgroundImage, err = addStamps(backgroundImage, card.Opinions[:])
	if err != nil {
		return fmt.Errorf("in drawCard(): %v", err)
	}

	backgroundImage, err = addEffects(backgroundImage, card.Effects[:])
	if err != nil {
		return fmt.Errorf("in drawCard(): %v", err)
	}

	if card.Cost.Value != 0 {
		backgroundImage, err = addCost(backgroundImage, card.Cost)
		if err != nil {
			return fmt.Errorf("in drawCard(): %v", err)
		}
	}

	backgroundImage = addTitle(backgroundImage, card.Title)

	result, err := os.Create(filename)
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

func addArt(backgroundImage image.Image, artPath string) (*image.RGBA, error) {

	//Add card art to backgroundImage
	artFile, err := os.Open(artPath)
	if err != nil {
		return nil, fmt.Errorf("in addArt(): Failed to open %s", artPath)
	}
	var art image.Image
	art, err = png.Decode(artFile)
	defer artFile.Close()
	if err != nil {
		art, err = jpeg.Decode(artFile)
		if err != nil {
			return nil, fmt.Errorf("in addArt(): Failed to decode png or jpg of %s", artPath)
		}
	}
	resultImage := image.NewRGBA(backgroundImage.Bounds())

	//First resize the art to fit the card.

	draw.Draw(resultImage, backgroundImage.Bounds(), backgroundImage, image.Point{}, draw.Over)

	artWidth := uint(56 * 30)
	artHeight := uint(33*30 + 1)

	art = resize.Resize(artWidth, artHeight, art, resize.Lanczos3)
	draw.Draw(resultImage, art.Bounds(), art, image.Point{}, draw.Over)

	artWidth = uint(5 * cm)
	artHeight = uint(3 * cm)

	art = resize.Resize(artWidth, artHeight, art, resize.Lanczos3)
	artPosition := image.Point{X: 90, Y: 90}
	draw.Draw(resultImage, art.Bounds().Add(artPosition), art, image.Point{}, draw.Over)
	return resultImage, nil
}

type group struct {
	id      int
	opinion Opinion
}

func addStamps(backgroundImage image.Image, opinions []Opinion) (*image.RGBA, error) {

	resultImage := image.NewRGBA(backgroundImage.Bounds())
	draw.Draw(resultImage, backgroundImage.Bounds(), backgroundImage, image.Point{}, draw.Src)

	fors := make([]group, 0)
	againsts := make([]group, 0)
	for idx, op := range opinions {
		if op == For {
			fors = append(fors, group{idx, For})
		} else if op == Against {
			againsts = append(againsts, group{idx, Against})
		} else if op == ExtraFor {
			fors = append(fors, group{idx, ExtraFor})
		} else if op == ExtraAgainst {
			againsts = append(againsts, group{idx, ExtraAgainst})
		}
	}

	if len(fors) > 4 {
		return nil, fmt.Errorf("in addStamps(): Too many For groups! Pick up to 4")
	}

	if len(againsts) > 4 {
		return nil, fmt.Errorf("in addStamps(): Too many Against groups! Pick up to 4")
	}

	image, err := drawStampFor(resultImage, fors)
	if err != nil {
		return nil, err
	}

	image, err = drawStampAgainst(image, againsts)
	if err != nil {
		return nil, err
	}
	return image, nil
	//drawStampAgainst(resultImage, againsts)
}

func drawStampFor(backgroundImage image.Image, fors []group) (*image.RGBA, error) {
	resultImage := image.NewRGBA(backgroundImage.Bounds())
	draw.Draw(resultImage, backgroundImage.Bounds(), backgroundImage, image.Point{}, draw.Src)

	for idx, group := range fors {
		path := grupyImagePaths[group.id]
		if group.opinion == ExtraFor {
			path = strings.Split(path, ".png")[0] + "U.png"
		}
		stampFile, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("in drawStampFor(): Failed to open file for stamp_id %d: %v", group.id, err)
		}
		defer stampFile.Close()

		stampImage, err := png.Decode(stampFile)
		if err != nil {
			return nil, fmt.Errorf("in drawStampFor(): Failed to decode png of stamp_id %d: %v", group.id, err)
		}

		stampImage = resize.Resize(300, 300, stampImage, resize.Lanczos3)

		var stampPosition image.Point
		switch idx {
		case 0:
			stampPosition = image.Point{X: 90 + 750 + 30, Y: 1200 + 30 + 90}
		case 1:
			stampPosition = image.Point{X: 390 + 60 + 750 + 60, Y: 1200 + 30 + 90}
		case 2:
			stampPosition = image.Point{X: 90 + 750 + 30, Y: 1200 + 30 + 300 + 90 + 90}
		case 3:
			stampPosition = image.Point{X: 390 + 60 + 750 + 60, Y: 1200 + 30 + 300 + 90 + 90}
		}
		draw.Draw(resultImage, stampImage.Bounds().Add(stampPosition), stampImage, image.Point{}, draw.Over)
	}
	return resultImage, nil
}

func drawStampAgainst(backgroundImage image.Image, againsts []group) (*image.RGBA, error) {
	resultImage := image.NewRGBA(backgroundImage.Bounds())
	draw.Draw(resultImage, backgroundImage.Bounds(), backgroundImage, image.Point{}, draw.Src)

	for idx, group := range againsts {
		path := grupyImagePaths[group.id]
		if group.opinion == ExtraAgainst {
			path = strings.Split(path, ".png")[0] + "U.png"
		}
		stampFile, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("in drawStampAgainst(): Failed to draw stamp_id %d", group.id)
		}
		stampImage, err := png.Decode(stampFile)
		if err != nil {
			return nil, fmt.Errorf("in drawStampAgainst(): Failed to decode png of stamp_id %d", group.id)
		}

		stampImage = resize.Resize(300, 300, stampImage, resize.Lanczos3)

		var stampPosition image.Point

		switch idx {
		case 0:
			stampPosition = image.Point{X: 90 + 30, Y: 1200 + 90 + 30}
		case 1:
			stampPosition = image.Point{X: 390 + 60 + 60, Y: 1200 + 90 + 30}
		case 2:
			stampPosition = image.Point{X: 90 + 30, Y: 1200 + 30 + 300 + 90 + 90}
		case 3:
			stampPosition = image.Point{X: 390 + 60 + 60, Y: 1200 + 30 + 300 + 90 + 90}
		}
		draw.Draw(resultImage, stampImage.Bounds().Add(stampPosition), stampImage, image.Point{}, draw.Over)
	}

	return resultImage, nil
}

func addEffects(backgroundImage image.Image, effects []int) (*image.RGBA, error) {
	resultImage := image.NewRGBA(backgroundImage.Bounds())
	draw.Draw(resultImage, backgroundImage.Bounds(), backgroundImage, image.Point{}, draw.Src)

	var err error
	effCount := 0
	for idx, val := range effects {
		if val != 0 {
			resultImage, err = drawEffect(resultImage, idx, val, effCount)
			if err != nil {
				return nil, fmt.Errorf("in addEffects(): %v", err)
			}
			effCount++
		}
	}
	return resultImage, nil
}

func drawEffect(backgroundImage *image.RGBA, effect_id int, effect_val int, effCount int) (*image.RGBA, error) {
	var effectFile *os.File
	var err error
	if effect_val > 0 {
		effectFile, err = os.Open(wskaznikiImagePaths[effect_id*2+1])
	} else if effect_val < 0 {
		effectFile, err = os.Open(wskaznikiImagePaths[effect_id*2])
	}
	if err != nil {
		return nil, fmt.Errorf("in drawEffect(): Failed to open image path: %v", err)
	}
	defer effectFile.Close()

	effImage, err := png.Decode(effectFile)
	if err != nil {
		return nil, fmt.Errorf("in drawEffect(): failed to decode image path: %v", err)
	}
	//First resize the art to fit the card.
	wskWidth := uint(1 * cm)
	wskHeight := uint(1.14 * cm)

	effImage = resize.Resize(wskWidth, wskHeight, effImage, resize.Lanczos3)

	x := 90 + 30 + effCount*(65+int(wskWidth))
	y := 2460 - int(wskHeight) - 120
	abs_val := int(math.Abs(float64(effect_val)))
	for i := abs_val; i > 0; i-- {
		effPostion := image.Point{X: x, Y: y + i*50}
		draw.Draw(backgroundImage, effImage.Bounds().Add(effPostion), effImage, image.Point{}, draw.Over)
	}

	return backgroundImage, nil
}

func addCost(backgroundImage image.Image, cost Cost) (*image.RGBA, error) {
	resultImage := image.NewRGBA(backgroundImage.Bounds())
	draw.Draw(resultImage, backgroundImage.Bounds(), backgroundImage, image.Point{}, draw.Src)
	costFile, err := os.Open(costToFilepath(cost))
	if err != nil {
		return nil, fmt.Errorf("addCost oops")
	}
	costImage, err := png.Decode(costFile)
	if err != nil {
		return nil, fmt.Errorf("addCost oops")
	}
	x := 0 + 90 + 15
	y := 20 + 90 + 15
	costPosition := image.Point{X: x, Y: y}
	draw.Draw(resultImage, costImage.Bounds().Add(costPosition), costImage, image.Point{}, draw.Over)
	return resultImage, nil

}

func costToFilepath(cost Cost) string {
	filepath := "assets/ceny/"
	switch cost.Value {
	case -5:
		filepath += "minus5"
	case -4:
		filepath += "minus4"
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
	case 4:
		filepath += "plus4"
	case 5:
		filepath += "plus5"
	default:
		filepath += "minus1"
	}
	if !cost.AddToBudget {
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
	fontSize := 135
	context.SetFontSize(float64(fontSize))
	context.SetClip(resultImage.Bounds())
	context.SetDst(resultImage)
	context.SetSrc(image.NewUniform(color.White))

	face := truetype.NewFace(font, &truetype.Options{
		Size: 165,
	})

	y := 1020 + 90
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
	} else {
		y += 100
	}

	lines := strings.Split(title, "\n")

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

func replaceSubstringInSlice(slice []string, oldSubstr, newSubstr string) []string {
	for i, str := range slice {
		slice[i] = strings.ReplaceAll(str, oldSubstr, newSubstr)
	}
	return slice
}

func stringToEffects(input string) [7]int {
	// Step 1: Remove the parentheses
	input = strings.Trim(input, "()")

	// Step 2: Split the string by commas
	stringParts := strings.Split(input, ",")

	// Step 3: Initialize an array of type [7]int
	var result [7]int

	// Step 4: Convert the string values to integers and assign to the array
	for i, str := range stringParts {
		if i < 7 { // Ensure we don't go out of bounds
			value, err := strconv.Atoi(str)
			if err != nil {
				fmt.Println("Error converting string to int:", err)
				continue // Skip the error and continue processing the rest
			}
			result[i] = value
		}
	}

	// Step 5: Return the result
	return result
}

func stringToOpinions(input string) [10]Opinion {
	// Step 1: Remove the parentheses
	input = strings.Trim(input, "()")

	// Step 2: Split the string by commas
	stringParts := strings.Split(input, ",")

	// Step 3: Initialize an array of type [10]int
	var result [10]Opinion

	// Step 4: Convert the string values to integers and assign to the array
	for i, str := range stringParts {
		if i < 10 { // Ensure we don't go out of bounds
			value, err := strconv.Atoi(str)
			if err != nil {
				fmt.Println("Error converting string to int:", err)
				continue // Skip the error and continue processing the rest
			}
			result[i] = Opinion(value)
		}
	}

	// Step 5: Return the result
	return result
}
