package dctycoon

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/nzin/dctycoon/global"

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

const (
	DIFFICULTY_EASY   = 0
	DIFFICULTY_MEDIUM = 1
	DIFFICULTY_HARD   = 2
)

type GameTimerSubscriber interface {
	ChangeSpeed(speed int)
	NewDay(timer *timer.GameTimer)
}

//
// Game contains the global picture of the game, i.e.
// - npactors (nonplayer) + player
// - marketplace items (DemandTemplate, ServerBundle)
// - load/save game functions
type Game struct {
	timer            *timer.GameTimer
	npactors         []*NPDatacenter
	player           *Player
	demandtemplates  []*supplier.DemandTemplate
	cronevent        *timer.GameCronEvent
	trends           *supplier.Trend
	dcmap            *DatacenterMap
	gameui           *GameUI
	timerevent       *sws.TimerEvent
	timersubscribers []GameTimerSubscriber
	currentSpeed     int
	gamestats        *GameStats
	debug            bool
}

//
// NewGame create a common place where we generate customer demand
// and assign them to the best provider (aka Datacenter)
func NewGame(quit *bool, root *sws.RootWidget, debug bool) *Game {
	log.Debug("NewGame()")

	g := &Game{
		timer:            nil,
		player:           nil,
		npactors:         make([]*NPDatacenter, 0, 0),
		demandtemplates:  make([]*supplier.DemandTemplate, 0, 0),
		cronevent:        nil,
		trends:           supplier.NewTrend(),
		dcmap:            NewDatacenterMap(),
		gameui:           nil,
		timerevent:       nil,
		timersubscribers: make([]GameTimerSubscriber, 0, 0),
		currentSpeed:     SPEED_STOP,
		gamestats:        NewGameStats(),
		debug:            debug,
	}
	g.gameui = NewGameUI(quit, root, g)

	// load demand templates
	if templatesarray, err := global.AssetDir("assets/demandtemplates"); err != nil {
		log.Error("Unable to load demand templates!")
	} else {
		for _, assetname := range templatesarray {
			template := supplier.DemandTemplateAssetLoad(assetname)
			if template == nil {
				log.Error("Error loading demand templates %s", assetname)
			} else {
				g.demandtemplates = append(g.demandtemplates, template)
			}
		}
	}

	g.dcmap.AddRackStatusSubscriber(g)

	return g
}

// RackStatusChange comes from interface DatacenterMap::RackStatusSubscriber
func (self *Game) RackStatusChange(x, y int32, rackstate int32) {
	log.Debug("Game::RackStatusChange(", x, ",", y, ",", rackstate, ")")
	if rackstate == RACK_MELTING {
		tileelement := self.dcmap.GetTile(x, y).TileElement()
		if tileelement.ElementType() == supplier.PRODUCT_RACK {
			rackelement := tileelement.(*RackElement)
			for _, item := range rackelement.GetRackServers() {
				if rand.Float32() > 0.5 {
					self.player.GetInventory().ScrapItem(item)
				}
			}
		}
	}
}

// GeneralOutage comes from interface DatacenterMap::RackStatusSubscriber
func (self *Game) GeneralOutage(bool) {

}

func (self *Game) ShowOpening() {
	self.gameui.ShowOpening()
}

func (self *Game) InitGame(locationid string, difficulty int32) {
	log.Debug("Game::InitGame(", locationid, ",", difficulty, ")")

	var initialcapital float64
	var nbopponents int32
	switch difficulty {
	case DIFFICULTY_MEDIUM:
		initialcapital = 100000.0
		nbopponents = 5
	case DIFFICULTY_HARD:
		initialcapital = 50000.0
		nbopponents = 7
	default:
		initialcapital = 200000.0
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

	self.gamestats.InitGame(self.player.GetInventory())

	// loading list of names
	availablenames := make([]NameList, 0, 0)
	if nameList, err := global.Asset("assets/namelist.json"); err != nil {
		availablenames = append(availablenames, NameList{Name: "John Doe", Male: true})
		nbopponents = 1
	} else {
		if err := json.Unmarshal(nameList, &availablenames); err != nil {
			log.Error("Game::InitGame: Unable to parse correctly assets/namelist.json")
		}
	}

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

		indexname := rand.Int() % len(availablenames)
		opponent.Init(self.timer, 200000, locationid, self.trends, profile, availablenames[indexname].Name, availablenames[indexname].Male)
		opponent.NewYearOperations()
		self.npactors = append(self.npactors, opponent)

		availablenames = append(availablenames[:indexname], availablenames[indexname+1:]...)
	}
	for _, s := range self.timersubscribers {
		s.NewDay(self.timer)
	}
	self.dcmap.SetGame(self.player.GetInventory(), self.player.GetLocation(), self.timer.CurrentTime)
	self.dcmap.InitMap("24_24_standard.json")
	//	self.dcmap.InitMap("3_4_room.json")
	self.gameui.SetGame(self.timer, self.player.GetInventory(), self.player.GetLedger(), self.trends, self.player.GetLocation(), self.dcmap)
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

	self.gamestats.LoadGame(self.player.GetInventory(), v["stats"].(map[string]interface{}))

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
	self.dcmap.SetGame(self.player.GetInventory(), self.player.GetLocation(), self.timer.CurrentTime)
	self.dcmap.LoadMap(v["map"].(map[string]interface{}))
	self.gameui.SetGame(self.timer, self.player.GetInventory(), self.player.GetLedger(), self.trends, self.player.GetLocation(), self.dcmap)
	self.gameui.ShowDC()
}

func (self *Game) SaveGame(filename string) {
	log.Debug("Game::SaveGame(", filename, ")")
	gamefile, err := os.Create(filename)
	if err != nil {
		log.Error("Not able to create savegame: ", err.Error())
		return
	}
	gamefile.WriteString("{")
	gamefile.WriteString(fmt.Sprintf(`"map": %s,`, self.dcmap.SaveMap()) + "\n")
	gamefile.WriteString(fmt.Sprintf(`"trends": %s,`, self.trends.Save()) + "\n")
	gamefile.WriteString(fmt.Sprintf(`"clock": %s,`, self.timer.Save()) + "\n")
	gamefile.WriteString(fmt.Sprintf(`"stats": %s,`, self.gamestats.Save()) + "\n")
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
	actors := make([]supplier.Actor, 0, 0)
	for _, a := range self.npactors {
		actors = append(actors, a)
	}
	actors = append(actors, self.player)

	// check if a bundle has to be paid (and renewed)
	for _, a := range actors {
		for _, sb := range a.GetInventory().GetServerBundles() {
			//
			if sb.Date.Day() == self.timer.CurrentTime.Day() {
				// pay monthly fee
				sb.PayMontlyFee(a.GetLedger(), self.timer.CurrentTime)
			}

			if sb.Date.Month() == self.timer.CurrentTime.Month() &&
				sb.Date.Day() == self.timer.CurrentTime.Day() {
				// if we don't renew, we drop it
				if rand.Float64() >= sb.Renewalrate {
					for _, c := range sb.Contracts {
						c.Item.Pool.Release(c.Item, c.Nbcores, c.Ramsize, c.Disksize)
					}
					a.GetInventory().RemoveServerBundle(sb)
				}
			}
		}
	}

	// check if we are the 1st of January
	if self.timer.CurrentTime.Day() == 1 && self.timer.CurrentTime.Month() == 1 {
		for _, np := range self.npactors {
			np.NewYearOperations()
		}
	}

	// generate new demand
	for _, d := range self.demandtemplates {
		if !d.Beginningdate.After(self.timer.CurrentTime) { // before or equals
			if rand.Float64() >= (365.0-float64(d.Howoften))/365.0 {
				// here we will generate a new server demand (and see if it is fulfill)
				demand := d.InstanciateDemand()
				log.Debug("Game::GenerateDemandAndFee(): A new demand has been created: ", demand.ToString())
				serverbundle, actor := demand.FindOffer(actors, self.timer.CurrentTime)
				if serverbundle != nil {
					actor.GetInventory().AddServerBundle(serverbundle)
					serverbundle.PayMontlyFee(actor.GetLedger(), self.timer.CurrentTime)
				}
				self.gamestats.TriggerDemandStat(self.timer.CurrentTime, demand, actor, serverbundle)
			}
		}
	}
	// generate electricity outage
	for _, a := range actors {
		// TBD, if outage = true, image of the provider decrease
		a.GetInventory().GeneratePowerlineOutage(a.GetLocation().Electricityfailrate)
	}

	// check for over heat
	self.dcmap.ComputeOverLimits()

	// pay other montly fee:
	// - power lines
	// - location renting
	if self.timer.CurrentTime.Day() == 1 {
		for _, a := range actors {
			consumption, _, _ := a.GetInventory().GetGlobalPower()
			a.GetLedger().PayUtility(a.GetInventory().GetMonthlyPowerlinesPrice()+consumption*24*30*a.GetLocation().Electricitycost/1000, self.timer.CurrentTime)
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

// ChangeGameSpeed allows to pause, or resume game speed. 3 values allowed: SPEED_STOP, SPEED_FORWARD, SPEED_FASTFORWARD
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

// GetCurrentSpeed returns the current speed, i.e SPEED_STOP, SPEED_FORWARD, SPEED_FASTFORWARD
func (self *Game) GetCurrentSpeed() int {
	return self.currentSpeed
}

// GetNPActors returns the list of opponents in this play
func (self *Game) GetNPActors() []*NPDatacenter {
	return self.npactors
}

// GetPlayer returns the player info, nothing fancy here
func (self *Game) GetPlayer() *Player {
	return self.player
}

// GetPlayer returns the stats central repo
func (self *Game) GetGameStats() *GameStats {
	return self.gamestats
}

func (self *Game) GetDebug() bool {
	return self.debug
}

type NameList struct {
	Name string
	Male bool
}
