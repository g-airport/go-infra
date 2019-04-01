package main

import (
	"flag"
	"fmt"
	"image/color"

	"github.com/g-airport/go-infra/avatar"
	"github.com/g-airport/go-infra/avatar/calc"
)

var colors = []uint32{
	0xff6200, 0x42c58e, 0x5a8de1, 0x785fe0,
}

func main() {
	flag.Parse()

	// init avatar builder, you need to tell builder ttf file and how to alignment text
	ab := avatarbuilder.NewAvatarBuilder("./SourceHanSansSC-Medium.ttf", &calc.SourceHansSansSCMedium{})
	ab.SetBackGroundColorHex(colors[3])
	ab.SetFrontGroundColor(color.White)
	ab.SetFontSize(60)
	ab.SetAvatarSize(200, 200)
	if err := ab.GenerateImageAndSave("哈哈哈", "./outCn.png"); err != nil {
		fmt.Println(err)
		return
	}
}
