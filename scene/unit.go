package scene

import (
	"github.com/pankona/gomo-simra/simra"
)

// unit base implementation

type uniter interface {
	Initialize()
	GetID() string
	simra.Subscriber
	SetPosition(float32, float32)
	GetPosition() (float32, float32)
	DoAction()
}

type attackInfo struct {
	attackRange    int
	power          int
	cooltime       float32 // second
	lastAttackTime int64   // frame
}

type position struct {
	x float32
	y float32
}

type actiontype int

const (
	// SPAWN spawns an unit
	actionSpawn actiontype = iota
	actionDead
	actionMoveToNearestTarget
	actionAttack
)

type action struct {
	actiontype actiontype
	data       interface{}
}

func newAction(a actiontype, d interface{}) *action {
	return &action{
		actiontype: a,
		data:       d,
	}
}

type unitBase struct {
	simra.Subscriber
	sprite    simra.Sprite
	id        string
	action    *action
	game      *game
	moveSpeed float32
}

func (u *unitBase) GetID() string {
	return u.id
}

func (u *unitBase) SetPosition(x, y float32) {
	u.sprite.X = x
	u.sprite.Y = y
}

func (u *unitBase) GetPosition() (float32, float32) {
	return u.sprite.X, u.sprite.Y
}

func newUnit(id, unittype string, game *game) uniter {
	// TODO: sample unit implemenation
	// unit type should be specified and switch here
	var u uniter
	switch unittype {
	case "player":
		u = &player{
			unitBase:   &unitBase{id: id, game: game, moveSpeed: 0.5},
			attackinfo: &attackInfo{attackRange: 50, power: 15, cooltime: 2},
		}
	default:
		// TODO: remove later
		u = &sampleUnit{
			unitBase:   &unitBase{id: id, game: game, moveSpeed: 0.5},
			attackinfo: &attackInfo{attackRange: 50, power: 15, cooltime: 2},
		}
	}

	// call each unit's initialize function
	u.Initialize()
	return u
}

type commandtype int

const (
	// SPAWN spawns an unit
	commandSpawn commandtype = iota
	// Damage to specified unit
	commandDamage
	commandDead
	commandGoToInitialState
	commandGoToRunningState
	commandUpdateSelection
	commandUnsetSelection
)

type command struct {
	commandtype commandtype
	data        interface{}
}

type damage struct {
	unit   uniter
	damage int
}

func newCommand(c commandtype, d interface{}) *command {
	return &command{commandtype: c, data: d}
}
