package scene

import (
	"image"
	"image/color"

	"github.com/pankona/gomo-simra/simra"
	"github.com/pankona/phantomize/scene/config"
)

type instruction struct {
	sprite *simra.Sprite
	game   *game
}

func (in *instruction) initialize() {
	in.sprite = simra.NewSprite()
	simra.GetInstance().AddSprite2(in.sprite)
}

func (in *instruction) OnEvent(i interface{}) {
	c := i.(*command)
	switch c.commandtype {
	case commandGameStarted:
		tex := simra.NewTextTexture(
			"choose first unit from here ↓",
			40, // fontsize
			color.RGBA{255, 255, 255, 255},
			image.Rect(0, 0, config.ScreenWidth, 80),
		)
		in.sprite.ReplaceTexture2(tex)
		in.sprite.X, in.sprite.Y = config.ScreenWidth/2+240, 260
		in.sprite.W, in.sprite.H = config.ScreenWidth, 80

	case commandUpdateSelection:
		tex := simra.NewTextTexture(
			"↑ tap field to summon the unit",
			40, // fontsize
			color.RGBA{255, 255, 255, 255},
			image.Rect(0, 0, config.ScreenWidth, 80),
		)
		in.sprite.ReplaceTexture2(tex)
		in.sprite.X, in.sprite.Y = config.ScreenWidth/2+275, 260
		in.sprite.W, in.sprite.H = config.ScreenWidth, 80

	case commandSpawn:
		simra.GetInstance().RemoveSprite(in.sprite)

	}
}