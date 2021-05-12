package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"time"

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

	colorR, err := strconv.Atoi(args[8])

	if err != nil {
		log.Fatal(err)
		return
	}

	colorG, err := strconv.Atoi(args[9])

	if err != nil {
		log.Fatal(err)
		return
	}

	colorB, err := strconv.Atoi(args[10])

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

	outPath := args[16]

	Draw(name, serial, hasGroup, idolImage, groupImage, frameImage, maskImage, cardId, colorR, colorG, colorB, large, font, textColorR, textColorG, textColorB, outPath)

}

func Draw(name string, serial int, hasGroup bool, idolImage string, groupImage string, frameImage string, maskImage string, cardId int, colorR int, colorG int, colorB int, large bool, font string, textColorR int, textColorG int, textColorB int, outPath string) error {
	start := time.Now()

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
	serialSize := 20.0

	if large {
		nameSize *= 2.2
		serialSize *= 2.2
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

	var mask image.Image
	_cardId := strconv.Itoa(cardId)
	dyeUrl := "./cache/masks/" + _cardId

	if !large {
		dyeUrl += "_small"
	}

	if _, err := os.Stat(dyeUrl); err == nil {
		mask, err = gg.LoadImage(dyeUrl)

		if err != nil {
			log.Fatal(err)
			return err
		}

	} else if os.IsNotExist(err) {
		err = exec.Command("convert", maskImage, "-fill", "rgb("+strconv.Itoa(colorR)+","+strconv.Itoa(colorG)+","+strconv.Itoa(colorB)+")", "-colorize", "100", dyeUrl).Run()

		if err != nil {
			log.Fatal(err)
			return err
		}

		mask, err = gg.LoadImage(dyeUrl)

		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	idolOffsetX := 47.0
	idolOffsetY := 54.0

	if large {
		idolOffsetX *= 2.2
		idolOffsetY *= 2.2
	}

	dc.DrawImage(idol, int(math.Floor(idolOffsetX)), int(math.Floor(idolOffsetY)))
	dc.DrawImage(frame, 0, 0)
	dc.DrawImage(mask, 0, -1)

	dc.SetRGB(float64(textColorR/255), float64(textColorG), float64(textColorB))

	textX := 50.0
	nameY := 445.0
	serialY := 421.0

	if large {
		textX *= 2.2
		nameY *= 2.2
		serialY *= 2.2
	}

	dc.DrawString(name, textX, nameY)

	err = dc.LoadFontFace(font, serialSize)

	if err != nil {
		log.Fatal(err)
		return err
	}

	dc.DrawString("#"+strconv.Itoa(serial), textX, serialY)

	if hasGroup {
		groupCacheLoc := "./cache/groups/" + _cardId
		var group image.Image

		if !large {
			groupCacheLoc += "_small"
		}

		if _, err := os.Stat(groupCacheLoc); err == nil {
			group, err = gg.LoadImage(groupCacheLoc)

			if err != nil {
				log.Fatal(err)
				return err
			}

		} else if os.IsNotExist(err) {
			err = exec.Command("convert", groupImage, "-fill", "rgb("+strconv.Itoa(textColorR)+","+strconv.Itoa(textColorG)+","+strconv.Itoa(textColorB)+")", "-colorize", "100", groupCacheLoc).Run()

			if err != nil {
				log.Fatal(err)
				return err
			}

			group, err = gg.LoadImage(groupCacheLoc)

			if err != nil {
				log.Fatal(err)
				return err
			}
		}

		dc.DrawImage(group, 0, 0)

	}

	dc.SavePNG(outPath)

	duration := time.Since(start)
	fmt.Print(duration)
	return nil
}
