package scene

import (
	"image"
	"image/color"
	"time"

	"github.com/pankona/gomo-simra/simra"
	"github.com/pankona/phantomize/scene/config"
)

type message struct {
	game *game
}

func (m *message) OnEvent(i interface{}) {
	c, ok := i.(*command)
	if !ok {
		// should be a command. ignore.
		return
	}

	switch c.commandtype {
	case commandShowMessage:
		message := c.data.(string)

		sprite := simra.NewSprite()
		simra.GetInstance().AddSprite2(sprite)

		tex := simra.NewTextTexture(
			message,
			40, // fontsize
			color.RGBA{255, 255, 255, 255},
			image.Rect(0, 0, config.ScreenWidth, 80),
		)
		sprite.ReplaceTexture2(tex)
		sprite.X, sprite.Y = config.ScreenWidth/2, 300
		sprite.W, sprite.H = config.ScreenWidth, 80

		go func() {
			select {
			case <-time.After(2 * time.Second):
				m.game.eventqueue <- newCommand(commandHideMessage, sprite)
			}
		}()

	case commandHideMessage:
		s := c.data.(*simra.Sprite)
		simra.GetInstance().RemoveSprite(s)
	}
}