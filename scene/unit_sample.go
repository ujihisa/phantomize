package scene

import (
	"image"

	"github.com/pankona/gomo-simra/simra"
)

type sampleUnit struct {
	*unitBase
}

func (u *sampleUnit) Initialize() {
	assetName := u.game.assetNameByUnitType(u.unittype)
	simra.GetInstance().AddSprite(&u.sprite)
	tex := simra.NewImageTexture(assetName,
		image.Rect(0, 0, 384, 384))
	u.sprite.ReplaceTexture(tex)

}

func (u *sampleUnit) OnEvent(i interface{}) {
	c, ok := i.(*command)
	if !ok {
		panic("unexpected command received. fatal.")
	}

	switch c.commandtype {
	default:
		u.unitBase.onEvent(c)
	}
}

func (u *sampleUnit) DoAction() {
	a := u.action
	if a == nil {
		// idle
		return
	}

	switch a.actiontype {
	case actionMoveToNearestTarget:
		u.target = nearestUnit(u.unitBase, u.game.players)
		if u.target == nil {
			break
		}
		moveToTarget(u.unitBase, u.target)

		if canAttackToTarget(u.unitBase, u.target) {
			u.game.eventqueue <- newCommand(commandAttack, u)
		}

	case actionDead:
		// i'm dead!
		killUnit(u, u.game.uniters)

	default:
		u.unitBase.doAction(a)
	}
}
