package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
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
	card            Card
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
		Title: "Depenalizacja piractwa cyfrowego",
	}

	err := drawCard(card, "testcard")
	if err != nil {
		log.Fatal(err)
	}
}

func drawCard(card Card, filename string) error {
	//Load base card png
	backgroundFile, err := os.Open("assets/card.png")
	backgroundImage, err := png.Decode(backgroundFile)
	defer backgroundFile.Close()
	if err != nil {
		return err
	}

	//Add card art to backgroundImage
	artFile, err := os.Open("assets/piractwo.png")
	art, err := png.Decode(artFile)
	defer artFile.Close()
	if err != nil {
		return nil
	}

	imageWithArt := addArt(backgroundImage, art)
	result, err := os.Create("result.png")
	defer result.Close()

	if err := png.Encode(result, imageWithArt); err != nil {
		fmt.Println("Error encoding output image:", err)
		return err
	}
	return nil
}

func addArt(backgroundImage image.Image, art image.Image) *image.RGBA {
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
