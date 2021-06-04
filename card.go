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
	// Parses command line arguments
	args := os.Args[1:]

	if len(args) < 17 {
		fmt.Print(`Invalid arguments!`)
		return
	}

	// The idol's name
	name := args[0]

	// Whether or not the group image exists (if true, it's a soloist)
	hasGroup, err := strconv.ParseBool(args[1])

	if err != nil {
		log.Fatal(err)
		return
	}

	idolImage := args[2]  // Path to the image of the idol
	groupImage := args[3] // Path to the image of the group logo
	frameImage := args[4] // Path to the frame image
	maskImage := args[5]  // Path to the dye image

	if err != nil {
		log.Fatal(err)
		return
	}

	// CMY color for the dye
	colorC, err := strconv.ParseFloat(args[6], 64)

	if err != nil {
		log.Fatal(err)
		return
	}

	colorM, err := strconv.ParseFloat(args[7], 64)

	if err != nil {
		log.Fatal(err)
		return
	}

	colorY, err := strconv.ParseFloat(args[8], 64)

	if err != nil {
		log.Fatal(err)
		return
	}

	// Large or small card (large = 770x1100, small = 350x500)
	large, err := strconv.ParseBool(args[9])

	if err != nil {
		log.Fatal(err)
		return
	}

	// Path to the font file to use for the text
	font := args[10]

	// RGB color for the text
	textColorR, err := strconv.Atoi(args[11])

	if err != nil {
		log.Fatal(err)
		return
	}

	textColorG, err := strconv.Atoi(args[12])

	if err != nil {
		log.Fatal(err)
		return
	}

	textColorB, err := strconv.Atoi(args[13])

	if err != nil {
		log.Fatal(err)
		return
	}

	// Whether or not to overlay the dye image onto the frame or vice versa (true = frame over dye)
	overlay, err := strconv.ParseBool(args[14])

	if err != nil {
		log.Fatal(err)
		return
	}

	// The path where the card image will be saved to
	outPath := args[15]

	// Draw the actual image
	Draw(name, hasGroup, idolImage, groupImage, frameImage, maskImage, colorC, colorM, colorY, large, font, textColorR, textColorG, textColorB, outPath, overlay)

}

func Draw(name string, hasGroup bool, idolImage string, groupImage string, frameImage string, maskImage string, colorC float64, colorM float64, colorY float64, large bool, font string, textColorR int, textColorG int, textColorB int, outPath string, overlay bool) error {

	// Size of the card image, including blank space
	var sizeX int
	var sizeY int

	if large {
		sizeX = 770
		sizeY = 1100
	} else {
		sizeX = 350
		sizeY = 500
	}

	// Creates a blank image
	dc := gg.NewContext(sizeX, sizeY)

	// Text size
	nameSize := 30.0

	if large {
		nameSize *= 2.2
	}

	// Load font from the path
	err := dc.LoadFontFace(font, nameSize)

	if err != nil {
		log.Fatal(err)
		return err
	}

	// Load idol image from the path
	idol, err := gg.LoadImage(idolImage)

	if err != nil {
		log.Fatal(err)
		return err
	}

	// Load frame image from the path
	frame, err := gg.LoadImage(frameImage)

	if err != nil {
		log.Fatal(err)
		return err
	}

	// Execute an ImageMagick conversion on the dye mask colorizing it to
	// the appropriate color in CMY. The output buffer is put through stdout,
	// which we then read from.
	buff, err := exec.Command("convert", maskImage, "-colorize", fmt.Sprintf("%f", colorC)+","+fmt.Sprintf("%f", colorM)+","+fmt.Sprintf("%f", colorY), "png:-").Output()

	if err != nil {
		log.Fatal(err)
		return err
	}

	// Decode the colored dye image from a buffer to an image
	mask, _, err := image.Decode(bytes.NewReader(buff))

	if err != nil {
		log.Fatal(err)
		return err
	}

	// Distance from the top left edge to the idol image
	idolOffsetX := 48.0
	idolOffsetY := 54.5

	if large {
		idolOffsetX *= 2.2
		idolOffsetY = 120
	}

	// Draw the idol
	dc.DrawImage(idol, int(math.Floor(idolOffsetX)), int(math.Floor(idolOffsetY)))

	// If overlay is true, the dye image goes on top of the frame
	if overlay {
		dc.DrawImage(mask, 0, -1)
		dc.DrawImage(frame, 0, -1)
	} else {
		dc.DrawImage(frame, 0, 0)
		dc.DrawImage(mask, 0, -1)
	}

	// Set the text color
	dc.SetRGB(float64(textColorR/255), float64(textColorG/255), float64(textColorB/255))

	// x and y coordinates to place the name text
	textX := 50.0
	nameY := 436.0

	if large {
		textX *= 2.2
		nameY *= 2.2
	}

	// Draw the idol name
	dc.DrawString(name, textX, nameY)

	if hasGroup {
		// Colorize the group logo, this is again put read through stdout
		groupImage, err := exec.Command("convert", groupImage, "-channel", "RGB", "-negate", "-fill", "rgb("+strconv.Itoa(textColorR)+","+strconv.Itoa(textColorG)+","+strconv.Itoa(textColorB)+")", "-colorize", "100", "png:-").Output()

		if err != nil {
			log.Fatal(err)
			return err
		}

		// Decode the colorized group logo from a buffer to an image
		groupBytes, _, err := image.Decode(bytes.NewReader(groupImage))

		dc.DrawImage(groupBytes, 0, 0)

	}

	// This is going to store the bytes of our image
	writer := new(bytes.Buffer)

	err = png.Encode(writer, dc.Image())

	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Convert the buffer to Base64 so that we can read it from node
	data := base64.StdEncoding.EncodeToString(writer.Bytes())

	// Save it to the output
	dc.SavePNG(outPath)

	// Put the Base64 string through stdout so node can read it
	fmt.Printf("%v", data)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return nil
}
