package scene

import (
	"fmt"

	"github.com/pankona/gomo-simra/simra"
	"github.com/pankona/gomo-simra/simra/image"
	"github.com/pankona/gomo-simra/simra/simlog"
)

// unit base implementation

type uniter interface {
	Initialize()
	GetID() string
	SetID(id string)
	SetPosition(float32, float32)
	GetPosition() (float32, float32)
	SetGame(g *game)
	IsSpawned() bool
	Dead()
	DoAction()
	GetUnitType() string
	SetUnitType(unittype string)
	GetHP() int
	GetMoveSpeed() float32
	GetCost() int
	GetTarget() uniter
	GetSprite() simra.Spriter
	GetAction() *action
	IsAlly() bool
	simra.Subscriber
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
	actionRecall
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
	simra               simra.Simraer
	sprite              simra.Spriter
	id                  string
	unittype            string
	action              *action
	game                *game
	moveSpeed           float32
	hp                  int
	attackinfo          *attackInfo
	target              uniter
	isSpawned           bool
	delayTimeToSummon   int64
	elapsedTimeToSummon int64
	isAlly              bool
	cost                int
}

func (u *unitBase) Initialize() {}

func (u *unitBase) GetID() string {
	return u.id
}

func (u *unitBase) SetID(id string) {
	u.id = id
}

func (u *unitBase) SetPosition(x, y float32) {
	// FIXME: don't cast to int
	u.sprite.SetPosition(x, y)
}

func (u *unitBase) GetPosition() (float32, float32) {
	// FIXME: dont' cast to flaot32
	return (float32)(u.sprite.GetPosition().X), (float32)(u.sprite.GetPosition().Y)
}

func (u *unitBase) SetGame(g *game) {
	u.game = g
}

func (u *unitBase) IsSpawned() bool {
	return u.isSpawned
}

// note that this is NOT delegate method.
// used by an object that composites unitBase.
func (u *unitBase) onEvent(c *command) {
	switch c.commandtype {
	case commandSpawn:
		d, ok := c.data.(uniter)
		if !ok {
			// unhandled event. ignore.
			return
		}

		if u.id == d.GetID() {
			// my spawn.
			u.action = newAction(actionSpawn, d)
		}

	case commandSpawned:
		d, ok := c.data.(uniter)
		if !ok {
			// unhandled event. ignore.
			return
		}

		if u.id != d.GetID() {
			// this spawn event is not for me.
			_, ok := d.(*sampleUnit)
			if ok {
				if u.isSpawned {
					simlog.Debugf("enemy's spawn %s is detected! kill them all!", d.GetID())

					// enemy's spawn. move to defeat.
					u.action = newAction(actionMoveToNearestTarget, nil)
				}
			}
			return
		}

	case commandAttack:
		d := c.data.(uniter)
		if u.GetID() != d.GetID() {
			// this is not for me. ignore
			break
		}
		// TODO: load in advance. don't do every time.
		texName := fmt.Sprintf("%s_atk.png", u.GetUnitType())
		tex := u.simra.NewImageTexture(texName, image.Rect(0, 0, 384, 384))
		u.sprite.ReplaceTexture(tex)

		u.action = newAction(actionAttack, u.target)

	case commandAttackEnd:
		d := c.data.(uniter)
		if u.GetID() != d.GetID() {
			// this is not for me. ignore
			break
		}

		// TODO: load in advance. don't do every time.
		texName := fmt.Sprintf("%s.png", u.GetUnitType())
		tex := u.simra.NewImageTexture(texName, image.Rect(0, 0, 384, 384))
		u.sprite.ReplaceTexture(tex)

	case commandDamage:
		d, ok := c.data.(*damage)
		if !ok {
			return
		}
		if u.id != d.unit.GetID() {
			return
		}

		u.hp -= d.damage
		simlog.Debugf("[DAMAGE] i'm [%s], HP = %d", u.GetID(), u.hp)
		if u.hp <= 0 {
			simlog.Debugf("[DEAD] i'm %s", u.GetID())
			u.game.eventqueue <- newCommand(commandDead, u)
		}

	case commandDead:
		d := c.data.(uniter)
		if u.GetID() == d.GetID() {
			u.action = newAction(actionDead, nil)
		}
		if u.target == d.GetTarget() {
			u.game.eventqueue <- newCommand(commandAttackEnd, u)
		}
		if len(u.game.uniters) == 0 {
			u.action = nil
			break
		}

	case commandRecall:
		d := c.data.(uniter)
		if u.GetID() != d.GetID() {
			// this is not for me. ignore
			break
		}

		u.action = newAction(actionRecall, nil)

	default:
		// nop
	}
}

func (u *unitBase) DoAction() {}

func (u *unitBase) GetUnitType() string {
	return u.unittype
}

func (u *unitBase) SetUnitType(unittype string) {
	u.unittype = unittype
}

func (u *unitBase) GetTarget() uniter {
	return u.target
}

func (u *unitBase) GetSprite() simra.Spriter {
	return u.sprite
}

func (u *unitBase) GetAction() *action {
	return u.action
}

func (u *unitBase) GetHP() int {
	return u.hp
}

func (u *unitBase) GetMoveSpeed() float32 {
	return u.moveSpeed
}

func (u *unitBase) GetCost() int {
	return u.cost
}

func (u *unitBase) IsAlly() bool {
	return u.isAlly
}

func (u *unitBase) doAction(a *action) {
	switch a.actiontype {
	case actionSpawn:
		u.elapsedTimeToSummon++
		if u.elapsedTimeToSummon <= u.delayTimeToSummon {
			// still summoning...
			break
		}
		u.elapsedTimeToSummon = 0

		d := a.data.(uniter)
		u.sprite.SetScale(64, 64)
		u.SetPosition(d.GetPosition())
		u.isSpawned = true
		u.GetSprite().AddTouchListener(&unitTouchListener{
			sprite: u.GetSprite(),
			uniter: u,
			game:   u.game,
		})

		// start moving to target
		u.game.eventqueue <- newCommand(commandSpawned, u)
		u.action = newAction(actionMoveToNearestTarget, nil)

		if u.IsAlly() {
			u.game.ongoingSummon--
		}

	case actionAttack:
		target := a.data.(uniter)
		if !canAttackToTarget(u, target) {
			u.game.eventqueue <- newCommand(commandAttackEnd, u)
			u.action = newAction(actionMoveToNearestTarget, nil)
			break
		}

		if u.game.currentFrame-u.attackinfo.lastAttackTime >=
			(int64)(u.attackinfo.cooltime*framePerSec) {
			simlog.Debugf("[ATTACK] i'm %s", u.GetID())
			u.attackinfo.lastAttackTime = u.game.currentFrame
			u.game.eventqueue <- newCommand(commandAttacking, (uniter)(u))
			u.game.eventqueue <- newCommand(commandDamage, &damage{target, u.attackinfo.power})
		}

	default:
		// nop
	}
}

func (u *unitBase) Dead() {
	u.sprite.RemoveAllTouchListener()
	// FIXME: ???
	u.sprite.SetScale(1, 1)
	u.SetPosition(-1, -1)
	u.action = nil
	u.target = nil
	u.isSpawned = false
}

func getUnitByUnitType(sim simra.Simraer, unittype string) *unitBase {
	switch unittype {
	case "player1":
		return &unitBase{
			simra:     sim,
			moveSpeed: 1.5,
			hp:        50,
			attackinfo: &attackInfo{
				attackRange: 50,
				power:       15,
				cooltime:    2,
			},
			delayTimeToSummon: 3 * framePerSec,
			isAlly:            true,
			cost:              20,
		}

	case "player2":
		return &unitBase{
			simra:     sim,
			moveSpeed: 1.0,
			hp:        75,
			attackinfo: &attackInfo{
				attackRange: 50,
				power:       20,
				cooltime:    3,
			},
			delayTimeToSummon: 4 * framePerSec,
			isAlly:            true,
			cost:              30,
		}

	case "player3":
		return &unitBase{
			simra:     sim,
			moveSpeed: 0.5,
			hp:        30,
			attackinfo: &attackInfo{
				attackRange: 200,
				power:       20,
				cooltime:    3,
			},
			delayTimeToSummon: 5 * framePerSec,
			isAlly:            true,
			cost:              40,
		}

	case "enemy1":
		return &unitBase{
			simra:     sim,
			moveSpeed: 0.5,
			hp:        30,
			attackinfo: &attackInfo{
				attackRange: 50,
				power:       10,
				cooltime:    3,
			},
			isAlly: false,
			cost:   20,
		}

	case "enemy2":
		return &unitBase{
			simra:     sim,
			moveSpeed: 1.0,
			hp:        45,
			attackinfo: &attackInfo{
				attackRange: 200,
				power:       15,
				cooltime:    4,
			},
			isAlly: false,
			cost:   40,
		}
	}

	return nil
}

type unitTouchListener struct {
	sprite simra.Spriter
	uniter uniter
	game   *game
}

func (u *unitTouchListener) OnTouchBegin(x, y float32) {
	// nop
}

func (u *unitTouchListener) OnTouchMove(x, y float32) {
	// nop
}

func (u *unitTouchListener) OnTouchEnd(x, y float32) {
	u.game.eventqueue <- newCommand(commandUpdateSelection, u)
}

func newUnit(id, unittype string, game *game) uniter {
	// TODO: sample unit implemenation
	// unit type should be specified and switch here
	var u uniter
	switch unittype {
	case "player1":
		fallthrough
	case "player2":
		fallthrough
	case "player3":
		u = &player{
			unitBase:          getUnitByUnitType(game.simra, unittype),
			delayTimeToRecall: 3 * framePerSec,
		}
	case "enemy1":
		fallthrough
	case "enemy2":
		u = &sampleUnit{unitBase: getUnitByUnitType(game.simra, unittype)}
	default:
		panic("unknown unittype!")
	}
	u.SetID(id)
	u.SetGame(game)
	u.SetUnitType(unittype)

	// call each unit's initialize function
	u.Initialize()

	return u
}

type commandtype int

const (
	commandGameStarted commandtype = iota
	commandSpawn
	commandSpawned
	commandAttack
	commandAttacking
	commandAttackEnd
	commandDamage
	commandDead
	commandRecall
	commandRecalled
	commandGoToInitialState
	commandGoToRunningState
	commandUpdateSelection
	commandUnsetSelection
	commandShowMessage
	commandHideMessage
	commandWin
	commandLose
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

func killUnit(u uniter, umap map[string]uniter) {
	u.Dead()
	delete(umap, u.GetID())
}
