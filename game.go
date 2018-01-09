package dctycoon

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/nzin/dctycoon/global"

	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/dctycoon/timer"
	"github.com/nzin/sws"
	log "github.com/sirupsen/logrus"
)

const (
	SPEED_STOP        = iota
	SPEED_FORWARD     = iota
	SPEED_FASTFORWARD = iota
)

// Actor -> Player or NPDatacenter
type Actor interface {
	GetInventory() *supplier.Inventory
	GetLedger() *accounting.Ledger
}

type GameTimerSubscriber interface {
	ChangeSpeed(speed int)
	NewDay(timer *timer.GameTimer)
}

//
// the global picture of the game, i.e.
// - npactors (nonplayer) + player
// - marketplace items (DemandTemplate, ServerBundle)
// - load/save game functions
type Game struct {
	timer            *timer.GameTimer
	npactors         []*NPDatacenter
	player           *Player
	demandtemplates  []*supplier.DemandTemplate
	serverbundles    []*supplier.ServerBundle
	cronevent        *timer.GameCronEvent
	trends           *supplier.Trend
	gameui           *GameUI
	timerevent       *sws.TimerEvent
	timersubscribers []GameTimerSubscriber
	currentSpeed     int
}

//
// NewGame create a common place where we generate customer demand
// and assign them to the best provider (aka Datacenter)
func NewGame(quit *bool, root *sws.RootWidget) *Game {
	log.Debug("NewGame()")

	g := &Game{
		timer:            nil,
		player:           nil,
		npactors:         make([]*NPDatacenter, 0, 0),
		demandtemplates:  make([]*supplier.DemandTemplate, 0, 0),
		serverbundles:    make([]*supplier.ServerBundle, 0, 0),
		cronevent:        nil,
		trends:           supplier.NewTrend(),
		gameui:           nil,
		timerevent:       nil,
		timersubscribers: make([]GameTimerSubscriber, 0, 0),
		currentSpeed:     SPEED_STOP,
	}
	g.gameui = NewGameUI(quit, root, g)

	return g
}

func (self *Game) ShowOpening() {
	self.gameui.ShowOpening()
}

func (self *Game) InitGame(locationid string, difficulty int32) {
	log.Debug("Game::InitGame(", locationid, ",", difficulty, ")")

	var initialcapital float64
	var nbopponents int32
	switch difficulty {
	case 1:
		initialcapital = 10000.0
		nbopponents = 5
	case 2:
		initialcapital = 5000.0
		nbopponents = 7
	default:
		initialcapital = 20000.0
		nbopponents = 3
	}
	if self.cronevent != nil {
		self.timer.RemoveCron(self.cronevent)
	}
	if self.timerevent != nil {
		self.timerevent.StopRepeat()
	}
	self.timerevent = nil

	self.timer = timer.NewGameTimer()
	self.cronevent = self.timer.AddCron(-1, -1, -1, func() {
		self.GenerateDemandAndFee()
	})
	self.trends.Init(self.gameui.eventpublisher, self.timer)
	self.player = NewPlayer()
	self.player.Init(self.timer, initialcapital, locationid)
	self.gameui.InitGame(self.timer, self.player.GetInventory(), self.player.GetLedger(), self.trends)

	// opponents
	self.npactors = make([]*NPDatacenter, 0, 0)
	for nb := int32(0); nb < nbopponents; nb++ {
		opponent := NewNPDatacenter()

		var locationarray []string
		for k, _ := range supplier.AvailableLocation {
			locationarray = append(locationarray, k)
		}
		locationid := locationarray[rand.Int()%len(locationarray)]

		var profile string
		if profilesarray, err := global.AssetDir("assets/npdatacenter"); err != nil {
			profile = "mono_r100_r200.json"
		} else {
			profile = profilesarray[rand.Int()%len(profilesarray)]
		}

		opponent.Init(self.timer, 20000, locationid, self.trends, profile)
		self.npactors = append(self.npactors, opponent)
	}
	for _, s := range self.timersubscribers {
		s.NewDay(self.timer)
	}
	self.gameui.ShowDC()
}

func (self *Game) LoadGame(filename string) {
	log.Debug("Game::LoadGame(", filename, ")")
	gamefile, err := os.Open(filename)
	if err != nil {
		log.Error("Game::LoadGame(): ", err.Error())
		return
	}
	var v map[string]interface{}
	jsonParser := json.NewDecoder(gamefile)
	if err = jsonParser.Decode(&v); err != nil {
		log.Error("Game::LoadGame(): parsing game file ", err.Error())
		return
	}
	gamefile.Close()

	if self.cronevent != nil {
		self.timer.RemoveCron(self.cronevent)
	}
	if self.timerevent != nil {
		self.timerevent.StopRepeat()
	}
	self.timerevent = nil

	self.timer = timer.NewGameTimer()
	self.cronevent = self.timer.AddCron(-1, -1, -1, func() {
		self.GenerateDemandAndFee()
	})
	self.timer.Load(v["clock"].(map[string]interface{}))
	self.trends.Load(v["trends"].(map[string]interface{}), self.gameui.eventpublisher, self.timer)
	self.player = NewPlayer()
	self.player.LoadGame(self.timer, v["player"].(map[string]interface{}))
	self.gameui.LoadGame(v, self.timer, self.player.GetInventory(), self.player.GetLedger(), self.trends)

	opponents := v["opponents"].([]interface{})
	// opponents
	self.npactors = make([]*NPDatacenter, 0, 0)
	for _, o := range opponents {
		opponent := NewNPDatacenter()
		opponent.LoadGame(self.timer, self.trends, o.(map[string]interface{}))
		self.npactors = append(self.npactors, opponent)
	}
	for _, s := range self.timersubscribers {
		s.NewDay(self.timer)
	}
	self.gameui.ShowDC()
}

func (self *Game) SaveGame(filename string) {
	log.Debug("Game::SaveGame(", filename, ")")
	gamefile, err := os.Create(filename)
	if err != nil {
		log.Error("Not able to create savegame: ", err.Error())
		return
	}
	data := self.gameui.SaveGame()
	gamefile.WriteString("{")
	gamefile.WriteString(fmt.Sprintf(`%s,`, data) + "\n")
	gamefile.WriteString(fmt.Sprintf(`"trends": %s,`, self.trends.Save()) + "\n")
	gamefile.WriteString(fmt.Sprintf(`"clock": %s,`, self.timer.Save()) + "\n")
	gamefile.WriteString(fmt.Sprintf(`"player": { %s },`, self.player.Save()) + "\n")
	gamefile.WriteString(`"opponents": [`)
	for i, o := range self.npactors {
		if i > 0 {
			gamefile.WriteString(",\n")
		}
		gamefile.WriteString(o.Save())
	}
	gamefile.WriteString("]")
	gamefile.WriteString("}\n")

	gamefile.Close()
}

//
// GenerateDemandAndFee generate randomly new demand (and check if a datacenter can handle it)
func (self *Game) GenerateDemandAndFee() {
	log.Debug("Game::GenerateDemandAndFee()")
	// check if a bundle has to be paid (and renewed)
	for _, sb := range self.serverbundles {
		//
		if sb.Date.Day() == self.timer.CurrentTime.Day() {
			// pay monthly fee

			// TBD
		}

		if sb.Date.Month() == self.timer.CurrentTime.Month() &&
			sb.Date.Day() == self.timer.CurrentTime.Day() {
			// if we don't renew, we drop it
			if rand.Float64() >= sb.Renewalrate {
				for _, c := range sb.Contracts {
					c.Item.Pool.Release(c.Item, c.Nbcores, c.Ramsize, c.Disksize)
				}
			}
		}
	}

	inventories := make([]*supplier.Inventory, 0, 0)
	for _, a := range self.npactors {
		inventories = append(inventories, a.GetInventory())
	}
	inventories = append(inventories, self.player.GetInventory())

	// generate new demand
	for _, d := range self.demandtemplates {
		if d.Beginningdate.After(self.timer.CurrentTime) {
			if rand.Float64() >= (365.0-float64(d.Howoften))/365.0 {
				// here we will generate a new server demand (and see if it is fulfill)
				demand := d.InstanciateDemand()
				serverbundle := demand.FindOffer(inventories, self.timer.CurrentTime)
				if serverbundle != nil {
					self.serverbundles = append(self.serverbundles, serverbundle)
				}
			}
		}
	}
}

func (self *Game) AddGameTimerSubscriber(subscriber GameTimerSubscriber) {
	for _, s := range self.timersubscribers {
		if s == subscriber {
			return
		}
	}
	self.timersubscribers = append(self.timersubscribers, subscriber)
}

func (self *Game) RemoveGameTimerSubscriber(subscriber GameTimerSubscriber) {
	for i, s := range self.timersubscribers {
		if s == subscriber {
			self.timersubscribers = append(self.timersubscribers[:i], self.timersubscribers[i+1:]...)
			break
		}
	}
}

func (self *Game) ChangeGameSpeed(speed int) {
	if self.timerevent != nil {
		self.timerevent.StopRepeat()
	}
	self.currentSpeed = speed
	for _, s := range self.timersubscribers {
		s.ChangeSpeed(speed)
	}

	if speed == SPEED_FORWARD || speed == SPEED_FASTFORWARD {
		tick := 2 * time.Second
		if speed == SPEED_FASTFORWARD {
			tick = time.Second / 2
		}
		self.timerevent = sws.TimerAddEvent(time.Now().Add(tick), tick, func(evt *sws.TimerEvent) {
			self.timer.TimerClock()
			for _, s := range self.timersubscribers {
				s.NewDay(self.timer)
			}
		})
	}
}

func (self *Game) GetCurrentSpeed() int {
	return self.currentSpeed
}
