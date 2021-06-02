package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"

	"github.com/fogleman/gg"
)

func main() {
	args := os.Args[1:]

	if len(args) < 17 {
		fmt.Print(`Invalid arguments!`)
		return
	}

	name := args[0]
	serial, err := strconv.Atoi(args[1])

	if err != nil {
		log.Fatal(err)
		return
	}

	hasGroup, err := strconv.ParseBool(args[2])

	if err != nil {
		log.Fatal(err)
		return
	}

	idolImage := args[3]
	groupImage := args[4]
	frameImage := args[5]
	maskImage := args[6]
	cardId, err := strconv.Atoi(args[7])

	if err != nil {
		log.Fatal(err)
		return
	}

	colorC, err := strconv.ParseFloat(args[8], 64)

	if err != nil {
		log.Fatal(err)
		return
	}

	colorM, err := strconv.ParseFloat(args[9], 64)

	if err != nil {
		log.Fatal(err)
		return
	}

	colorY, err := strconv.ParseFloat(args[10], 64)

	if err != nil {
		log.Fatal(err)
		return
	}

	large, err := strconv.ParseBool(args[11])

	if err != nil {
		log.Fatal(err)
		return
	}

	font := args[12]

	textColorR, err := strconv.Atoi(args[13])

	if err != nil {
		log.Fatal(err)
		return
	}

	textColorG, err := strconv.Atoi(args[14])

	if err != nil {
		log.Fatal(err)
		return
	}

	textColorB, err := strconv.Atoi(args[15])

	if err != nil {
		log.Fatal(err)
		return
	}

	overlay, err := strconv.ParseBool(args[16])

	if err != nil {
		log.Fatal(err)
		return
	}

	outPath := args[17]

	Draw(name, serial, hasGroup, idolImage, groupImage, frameImage, maskImage, cardId, colorC, colorM, colorY, large, font, textColorR, textColorG, textColorB, outPath, overlay)

}

func Draw(name string, serial int, hasGroup bool, idolImage string, groupImage string, frameImage string, maskImage string, cardId int, colorC float64, colorM float64, colorY float64, large bool, font string, textColorR int, textColorG int, textColorB int, outPath string, overlay bool) error {

	var sizeX int
	var sizeY int

	if large {
		sizeX = 770
		sizeY = 1100
	} else {
		sizeX = 350
		sizeY = 500
	}

	dc := gg.NewContext(sizeX, sizeY)

	nameSize := 30.0
	// serialSize := 20.0

	if large {
		nameSize *= 2.2
		// serialSize *= 2.2
	}

	err := dc.LoadFontFace(font, nameSize)

	if err != nil {
		log.Fatal(err)
		return err
	}

	idol, err := gg.LoadImage(idolImage)

	if err != nil {
		log.Fatal(err)
		return err
	}

	frame, err := gg.LoadImage(frameImage)

	if err != nil {
		log.Fatal(err)
		return err
	}

	buff, err := exec.Command("convert", maskImage, "-colorize", fmt.Sprintf("%f", colorC)+","+fmt.Sprintf("%f", colorM)+","+fmt.Sprintf("%f", colorY), "png:-").Output()

	if err != nil {
		log.Fatal(err)
		return err
	}

	mask, _, err := image.Decode(bytes.NewReader(buff))

	if err != nil {
		log.Fatal(err)
		return err
	}

	idolOffsetX := 48.0
	idolOffsetY := 54.5

	if large {
		idolOffsetX *= 2.2
		idolOffsetY = 120
	}

	dc.DrawImage(idol, int(math.Floor(idolOffsetX)), int(math.Floor(idolOffsetY)))

	if overlay {
		dc.DrawImage(mask, 0, -1)
		dc.DrawImage(frame, 0, -1)
	} else {
		dc.DrawImage(frame, 0, 0)
		dc.DrawImage(mask, 0, -1)
	}

	dc.SetRGB(float64(textColorR/255), float64(textColorG/255), float64(textColorB/255))

	textX := 50.0
	nameY := 436.0
	serialY := 421.0

	if large {
		textX *= 2.2
		nameY *= 2.2
		serialY *= 2.2
	}

	dc.DrawString(name, textX, nameY)

	if hasGroup {
		groupImage, err := exec.Command("convert", groupImage, "-channel", "RGB", "-negate", "-fill", "rgb("+strconv.Itoa(textColorR)+","+strconv.Itoa(textColorG)+","+strconv.Itoa(textColorB)+")", "-colorize", "100", "png:-").Output()

		if err != nil {
			log.Fatal(err)
			return err
		}

		groupBytes, _, err := image.Decode(bytes.NewReader(groupImage))

		dc.DrawImage(groupBytes, 0, 0)

	}

	writer := new(bytes.Buffer)

	err = png.Encode(writer, dc.Image())

	if err != nil {
		log.Fatal(err)
		return nil
	}

	data := base64.StdEncoding.EncodeToString(writer.Bytes())

	dc.SavePNG(outPath)
	fmt.Printf("%v", data)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return nil
}
