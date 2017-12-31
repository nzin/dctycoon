package dctycoon

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/nzin/dctycoon/accounting"
	"github.com/nzin/dctycoon/supplier"
	"github.com/nzin/dctycoon/timer"
	"github.com/nzin/sws"
	log "github.com/sirupsen/logrus"
)

// Actor -> Player or NPDatacenter
type Actor interface {
	GetInventory() *supplier.Inventory
	GetLedger() *accounting.Ledger
	//	Save()
	//	Load(gametimer)
}

//
// the global picture of the game, i.e.
// - npactors (nonplayer) + player
// - marketplace items (DemandTemplate, ServerBundle)
// - load/save game functions
type Game struct {
	timer           *timer.GameTimer
	npactors        []Actor
	player          *Player
	demandtemplates []*supplier.DemandTemplate
	serverbundles   []*supplier.ServerBundle
	cronevent       *timer.GameCronEvent
	trends          *supplier.Trend
	gameui          *GameUI
}

//
// NewGame create a common place where we generate customer demand
// and assign them to the best provider (aka Datacenter)
func NewGame(quit *bool, root *sws.RootWidget) *Game {
	log.Debug("NewGame()")

	g := &Game{
		timer:           nil,
		player:          nil,
		npactors:        make([]Actor, 0, 0),
		demandtemplates: make([]*supplier.DemandTemplate, 0, 0),
		serverbundles:   make([]*supplier.ServerBundle, 0, 0),
		cronevent:       nil,
		trends:          supplier.NewTrend(),
		gameui:          NewGameUI(quit, root),
	}

	return g
}

func (self *Game) InitGame(initialcapital float64, location string) {
	log.Debug("Game::InitGame(", initialcapital, ")")
	if self.cronevent != nil {
		self.timer.RemoveCron(self.cronevent)
	}
	self.timer = timer.NewGameTimer()
	self.cronevent = self.timer.AddCron(-1, -1, -1, func() {
		self.GenerateDemandAndFee()
	})
	self.trends.Init(self.gameui.eventpublisher, self.timer)
	self.player = NewPlayer()
	self.player.Init(self.timer, initialcapital, location)
	self.gameui.InitGame(self.timer, self.player.GetInventory(), self.player.GetLedger(), self.trends)
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

	self.timer = timer.NewGameTimer()
	self.cronevent = self.timer.AddCron(-1, -1, -1, func() {
		self.GenerateDemandAndFee()
	})
	self.timer.Load(v["clock"].(map[string]interface{}))
	self.trends.Load(v["trends"].(map[string]interface{}), self.gameui.eventpublisher, self.timer)
	self.player = NewPlayer()
	self.player.LoadGame(self.timer, v["player"].(map[string]interface{}))
	self.gameui.LoadGame(v, self.timer, self.player.GetInventory(), self.player.GetLedger(), self.trends)
	self.gameui.ShowDC()
}

func (self *Game) SaveGame(filename string) {
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
	gamefile.WriteString(fmt.Sprintf(`"player": { %s }`, self.player.Save()) + "\n")
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
