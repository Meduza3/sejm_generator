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
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

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
	Value    int
	Currency Currency
}

type Currency string

const (
	Trust   Currency = "trust"
	Cash    Currency = "cash"
	Scandal Currency = "scandal"
)

// Values 0-10
type LegislationCard struct {
	ArtPath  string
	Title    string
	Opinions [10]Opinion
	Effects  [7]int
	Cost     Cost
}

type Symbol string

const (
	NoSymbol  Symbol = ""
	Reflect   Symbol = "reflect"
	Table     Symbol = "table"
	Paperclip Symbol = "paperclip"
)

type ActionCard struct {
	ArtPath     string
	Title       string
	Description string
	Symbol      Symbol
	Cost        Cost
	RedText     string
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

func ParseLegislationInput(input string) LegislationCard {
	inputParts := strings.Split(input, ">")

	// Basic validation of input parts length
	if len(inputParts) < 5 {
		log.Fatalf("Invalid input: expected at least 5 parts, got %d in input: %s", len(inputParts), input)
	}

	artPath := inputParts[0]
	title := inputParts[1]
	opinions := stringToOpinions(inputParts[2])
	effects := stringToEffects(inputParts[3])
	cost, _ := strconv.Atoi(inputParts[4])

	return NewLegislationCard(artPath, title, opinions, effects, cost)
}

func ParseActionInput(input string) ActionCard {
	inputParts := strings.Split(input, ">")

	// Basic validation of input parts length
	if len(inputParts) < 6 {
		log.Fatalf("Invalid input: expected at least 6 parts, got %d in input: %s", len(inputParts), input)
	}

	artPath := inputParts[0]
	title := inputParts[1]
	description := inputParts[2]
	var symbol Symbol
	switch inputParts[3] {
	case "reflect":
		symbol = Reflect
	case "table":
		symbol = Table
	case "paperclip":
		symbol = Paperclip
	default:
		symbol = NoSymbol
	}
	var currency Currency
	switch inputParts[4] {
	case "trust":
		currency = Trust
	case "cash":
		currency = Cash
	case "scandal":
		currency = Scandal
	}
	cost, _ := strconv.Atoi(inputParts[5])
	var redtext string
	if len(inputParts) == 7 {
		redtext = inputParts[6]
	}

	return NewActionCard(artPath, title, description, symbol, cost, currency, redtext)
}
func NewActionCard(artPath, title, description string, symbol Symbol, cost int, currency Currency, redtext string) ActionCard {
	return ActionCard{
		ArtPath:     artPath + ".png",
		Title:       title,
		Description: description,
		Symbol:      symbol,
		Cost: Cost{
			Value:    cost,
			Currency: currency,
		},
		RedText: redtext,
	}
}

func NewLegislationCard(artPath string, title string, opinions [10]Opinion, effects [7]int, cost int) LegislationCard {

	return LegislationCard{
		ArtPath:  artPath + ".png",
		Title:    title,
		Opinions: opinions,
		Effects:  effects,
		Cost: Cost{
			Value:    cost,
			Currency: Cash,
		},
	}
}

// Średnica kółka: 300px
// 1cm = 300px

const DirName = "generated"

func main() {
	fmt.Println("------------------SEJM GENERATOR--------------------")
	fmt.Println("1 - Generate legislation cards")
	fmt.Println("2 - Generate action cards")

	var selection int
	fmt.Printf("Please enter your selection: ")
	_, err := fmt.Scanf("%d", &selection)
	if err != nil {
		fmt.Println("Failed to read input:", err)
		time.Sleep(500 * time.Millisecond)

		return
	}

	switch selection {
	case 1:
		clearConsole()
		legislationCardsLoop()
	case 2:
		clearConsole()
		actionCardsLoop()
	default:
		fmt.Printf("%d is not an option. exiting", selection)
		time.Sleep(500 * time.Millisecond)
		os.Exit(1)
	}
}

func clearConsole() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	case "linux", "darwin":
		cmd = exec.Command("clear")
	default:
		fmt.Println("Unsupported platform")
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

func actionCardsLoop() {
	fmt.Println("------------------SEJM GENERATOR--------------------")
	fmt.Println("card code: filename>card title>description>symbol>costtype>cost>[optional red description]")
	fmt.Println("filename: without .png")
	fmt.Println("description: long description of the action")
	fmt.Println("symbol: reflect, table or paperclip")
	fmt.Println("Costtype: trust, cash or scandal")
	fmt.Println("cost: in [-10,10]")
	if _, err := os.Stat(DirName); os.IsNotExist(err) {
		err := os.Mkdir(DirName, 0755)
		if err != nil {
			fmt.Printf("Critical error - create directory %s yourself\n", DirName)
			os.Exit(1)
		}
	}

	for {
		fmt.Println("----------------------------------------------------")
		fmt.Println("")
		fmt.Println("Input card code:")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		// If the input is empty, prompt the user again
		if input == "" {
			fmt.Println("Input cannot be empty. Please enter the card code or type 'exit' to quit.")
			continue
		}

		if input == "exit" {
			fmt.Println("Exiting...")
			os.Exit(0)
		}
		card := ParseActionInput(input)

		err := drawActionCard(card, DirName+"/"+card.ArtPath)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\n")
		fmt.Printf("Generated %s -> wygenerowane/%s\n", card.ArtPath, card.ArtPath)
		fmt.Printf("\n")
	}

}

func legislationCardsLoop() {
	fmt.Println("------------------SEJM GENERATOR--------------------")
	fmt.Println("card code: filename>card title>opinions>effects>cost")
	fmt.Println("filename: without .png")
	fmt.Println("opinions: (1,2,2,-2,-2,0,0,0,0,-1)")
	fmt.Println("effects: (0,0,0,1,-2,0,1)")
	fmt.Println("cost: in [-10,10]")

	if _, err := os.Stat(DirName); os.IsNotExist(err) {
		err := os.Mkdir(DirName, 0755)
		if err != nil {
			fmt.Printf("Critical error - create directory %s yourself\n", DirName)
			os.Exit(1)
		}
	}

	for {
		fmt.Println("----------------------------------------------------")
		fmt.Println("")
		fmt.Println("Input card code:")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		// If the input is empty, prompt the user again
		if input == "" {
			fmt.Println("Input cannot be empty. Please enter the card code or type 'exit' to quit.")
			continue
		}

		if input == "exit" {
			fmt.Println("Exiting...")
			os.Exit(0)
		}
		card := ParseLegislationInput(input)

		err := drawLegislationCard(card, DirName+"/"+card.ArtPath)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\n")
		fmt.Printf("Generated %s -> generated/%s\n", card.ArtPath, card.ArtPath)
		fmt.Printf("\n")
	}
}

func drawActionCard(card ActionCard, filename string) error {
	backgroundFile, err := os.Open("assets/action_printcard.png")
	if err != nil {
		return err
	}
	defer backgroundFile.Close()

	backgroundImage, err := png.Decode(backgroundFile)
	if err != nil {
		return err
	}
	backgroundImage = resize.Resize(1680, 2580, backgroundImage, resize.Lanczos3)

	backgroundImage, err = addArt(backgroundImage, card.ArtPath)
	if err != nil {
		return fmt.Errorf("in drawCard(): %v", err)
	}

	backgroundImage = addTitle(backgroundImage, card.Title)

	backgroundImage, err = addDescription(backgroundImage, card.Description)
	if err != nil {
		return fmt.Errorf("in drawCard(): %v", err)
	}

	if card.RedText != "" {
		backgroundImage, err = addRibbon(backgroundImage)
		if err != nil {
			return fmt.Errorf("in drawCard(): %v", err)
		}
		backgroundImage, err = addRedText(backgroundImage, card.RedText)
		if err != nil {
			return fmt.Errorf("in drawCard(): %v", err)
		}
	}

	backgroundImage, err = addSymbol(backgroundImage, card.Symbol)
	if err != nil {
		return fmt.Errorf("in drawCard(): %v", err)
	}

	if card.Cost.Value != 0 {
		backgroundImage, err = addCost(backgroundImage, card.Cost)
		if err != nil {
			return fmt.Errorf("in drawCard(): %v", err)
		}
	}

	result, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err := png.Encode(result, backgroundImage); err != nil {
		fmt.Println("Error encoding output image:", err)
		return err
	}
	return nil
}

func drawLegislationCard(card LegislationCard, filename string) error {
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
	backgroundImage = resize.Resize(1680, 2580, backgroundImage, resize.Lanczos3)

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

func addSymbol(backgroundImage image.Image, symbol Symbol) (*image.RGBA, error) {
	resultImage := image.NewRGBA(backgroundImage.Bounds())
	draw.Draw(resultImage, backgroundImage.Bounds(), backgroundImage, image.Point{}, draw.Src)

	symbolPath := "assets/symbol"
	switch symbol {
	case Reflect:
		symbolPath += "/reflect.png"
	case Table:
		symbolPath += "/table.png"
	case Paperclip:
		symbolPath += "/paperclip.png"
	default:
		symbolPath += string(symbol)
	}

	var err error
	symbolFile, err := os.Open(symbolPath)
	if err != nil {
		return nil, fmt.Errorf("in drawSymbol(): Failed to open file: %v", err)
	}

	symbolImage, err := png.Decode(symbolFile)
	if err != nil {
		return nil, fmt.Errorf("in drawSymbol(): Failed to decode file: %v", err)
	}

	symbolImage = resize.Resize(300, 300, symbolImage, resize.Lanczos3)
	symbolPosition := image.Point{X: 90, Y: 2200}
	draw.Draw(resultImage, symbolImage.Bounds().Add(symbolPosition), symbolImage, image.Point{}, draw.Over)
	return resultImage, nil
}

func addRedText(backgroundImage image.Image, redtext string) (*image.RGBA, error) {
	resultImage := image.NewRGBA(backgroundImage.Bounds())
	draw.Draw(resultImage, backgroundImage.Bounds(), backgroundImage, image.Point{}, draw.Src)

	// Load the font
	fontBytes, err := os.ReadFile("assets/sylfaen.ttf")
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	// Create a new context for drawing text
	context := freetype.NewContext()
	context.SetFont(font)
	fontSize := 135 - 50
	context.SetFontSize(float64(fontSize))
	context.SetClip(resultImage.Bounds())
	context.SetDst(resultImage)
	context.SetSrc(image.NewUniform(color.White))

	face := truetype.NewFace(font, &truetype.Options{
		Size: float64(165 - 50),
	})

	maxWidth := resultImage.Bounds().Dx() - 40 // Set a maximum width for text lines, with padding

	// Split the title into lines based on available width
	lines := splitTextIntoLines(face, redtext, maxWidth)

	// Draw each line of text
	y := 2200 + 90
	for _, line := range lines {
		lineWidth := textWidth(face, line)
		x := ((resultImage.Bounds().Dx() - lineWidth) / 2) + 200 // Center text horizontally
		x += 150
		pt := freetype.Pt(x, y)
		_, err = context.DrawString(line, pt)
		if err != nil {
			return nil, err
		}
		y += int(context.PointToFixed(float64(fontSize)) >> 6)
	}

	return resultImage, nil
}

func addRibbon(backgroundImage image.Image) (*image.RGBA, error) {
	resultImage := image.NewRGBA(backgroundImage.Bounds())
	draw.Draw(resultImage, backgroundImage.Bounds(), backgroundImage, image.Point{}, draw.Src)
	ribbonFile, err := os.Open("assets/ribbon.png")
	if err != nil {
		return nil, fmt.Errorf("in addRibbon(): %v", err)
	}

	ribbonImage, err := png.Decode(ribbonFile)
	if err != nil {
		return nil, fmt.Errorf("in addRibbon(): %v", err)
	}

	ribbonImage = resize.Resize(1680, 2580, ribbonImage, resize.Lanczos3)
	ribbonPosition := image.Point{X: 0, Y: 2200}
	draw.Draw(resultImage, ribbonImage.Bounds().Add(ribbonPosition), ribbonImage, image.Point{}, draw.Over)
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

	switchOnValue := func(filepath string) string {
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

		return filepath
	}

	filepath := "assets/ceny/"
	if cost.Currency == Cash {
		filepath += "cash/"
		filepath = switchOnValue(filepath)
	} else if cost.Currency == Trust {
		filepath += "trust/"
		filepath = switchOnValue(filepath)
	} else if cost.Currency == Scandal {
		filepath += "scandal/"
		filepath = switchOnValue(filepath)
	}
	filepath += ".png"
	return filepath
}

func addDescription(backgroundImage image.Image, title string) (*image.RGBA, error) {
	resultImage := image.NewRGBA(backgroundImage.Bounds())
	draw.Draw(resultImage, backgroundImage.Bounds(), backgroundImage, image.Point{}, draw.Src)

	// Load the font
	fontBytes, err := os.ReadFile("assets/sylfaen.ttf")
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	// Create a new context for drawing text
	context := freetype.NewContext()
	context.SetFont(font)
	fontSize := 135 - 50
	context.SetFontSize(float64(fontSize))
	context.SetClip(resultImage.Bounds())
	context.SetDst(resultImage)
	context.SetSrc(image.NewUniform(color.Black))

	face := truetype.NewFace(font, &truetype.Options{
		Size: float64(165 - 50),
	})

	maxWidth := resultImage.Bounds().Dx() - 40 // Set a maximum width for text lines, with padding

	// Split the title into lines based on available width
	lines := splitTextIntoLines(face, title, maxWidth)

	// Draw each line of text
	y := 1320 + 90
	for _, line := range lines {
		lineWidth := textWidth(face, line)
		x := (resultImage.Bounds().Dx() - lineWidth) / 2 // Center text horizontally
		x += 150
		pt := freetype.Pt(x, y)
		_, err = context.DrawString(line, pt)
		if err != nil {
			return nil, err
		}
		y += int(context.PointToFixed(float64(fontSize)) >> 6)
	}

	return resultImage, nil
}

// splitTextIntoLines splits the input text into multiple lines such that each line fits within the maxWidth.
func splitTextIntoLines(face font.Face, text string, maxWidth int) []string {
	words := strings.Fields(text)
	var lines []string
	var currentLine string

	for _, word := range words {
		testLine := currentLine + " " + word
		if textWidth(face, strings.TrimSpace(testLine)) > maxWidth {
			if currentLine == "" {
				// If the current line is empty, add the word to the line anyway (to prevent infinite loop)
				currentLine = word
			}
			lines = append(lines, strings.TrimSpace(currentLine))
			currentLine = word
		} else {
			currentLine = testLine
		}
	}
	if currentLine != "" {
		lines = append(lines, strings.TrimSpace(currentLine))
	}

	return lines
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
