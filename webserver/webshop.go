package webserver

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/dctycoon/supplier"

	"html/template"

	"github.com/nzin/dctycoon"
	log "github.com/sirupsen/logrus"
)

var emptyAnswer = `<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx on Debian!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx on Debian!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working on Debian. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a></p>

<p>
      Please use the <tt>reportbug</tt> tool to report bugs in the
      nginx package with Debian. However, check <a
      href="http://bugs.debian.org/cgi-bin/pkgreport.cgi?ordering=normal;archive=0;src=nginx;repeatmerged=0">existing
      bug reports</a> before reporting a new bug.
</p>

<p><em>Thank you for using debian and nginx.</em></p>


</body>
</html>`

var game *dctycoon.Game

type IndexPage struct {
	game *dctycoon.Game
}

type DedicatedOffer struct {
	Name      string
	Cpu       string
	Ram       string
	Disk      string
	Vt        string
	Price     string
	Remaining string
}

type VpsOffer struct {
	Name      string
	Cpu       string
	Ram       string
	Disk      string
	Price     string
	Remaining string
}

type IndexTemplate struct {
	Dedicated          []*DedicatedOffer
	Vps                []*VpsOffer
	Companyname        string
	ElectricalNetworks string
	Location           string
}

func (self *IndexPage) fillDedicated() []*DedicatedOffer {
	dedicatedOffers := make([]*DedicatedOffer, 0, 0)
	pool := game.GetPlayer().GetInventory().GetDefaultPhysicalPool()
	for _, offer := range game.GetPlayer().GetOffers() {
		if offer.Vps == false {
			vtstring := "no"
			if offer.Vt == true {
				vtstring = "yes"
			}
			dedicated := &DedicatedOffer{
				Name:      offer.Name,
				Cpu:       fmt.Sprintf("%dx Altium", offer.Nbcores),
				Ram:       global.AdjustMega(offer.Ramsize),
				Disk:      global.AdjustMega(offer.Disksize),
				Vt:        vtstring,
				Price:     fmt.Sprintf("%.0f", offer.Price),
				Remaining: fmt.Sprintf("%d", pool.HowManyFit(offer.Nbcores, offer.Ramsize, offer.Disksize, offer.Vt)),
			}
			dedicatedOffers = append(dedicatedOffers, dedicated)
		}
	}
	return dedicatedOffers
}

func (self *IndexPage) fillVps() []*VpsOffer {
	vpsOffers := make([]*VpsOffer, 0, 0)
	pool := game.GetPlayer().GetInventory().GetDefaultVpsPool()
	for _, offer := range game.GetPlayer().GetOffers() {
		if offer.Vps == true {
			vps := &VpsOffer{
				Name:      offer.Name,
				Cpu:       fmt.Sprintf("%dx Altium", offer.Nbcores),
				Ram:       global.AdjustMega(offer.Ramsize),
				Disk:      global.AdjustMega(offer.Disksize),
				Price:     fmt.Sprintf("%.0f", offer.Price),
				Remaining: fmt.Sprintf("%d", pool.HowManyFit(offer.Nbcores, offer.Ramsize, offer.Disksize, offer.Vt)),
			}
			vpsOffers = append(vpsOffers, vps)
		}
	}
	return vpsOffers
}

func (self *IndexPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//	fmt.Println("IndexPage")

	t := template.New("index page")
	if indexAsset, err := global.Asset("assets/webshop/index.html"); err == nil {
		t, _ := t.Parse(string(indexAsset))

		nbpowerlines := 0
		for i := 0; i < 3; i++ {
			if self.game.GetPlayer().GetInventory().GetPowerlines()[i] != supplier.POWERLINE_NONE {
				nbpowerlines++
			}
		}

		variables := IndexTemplate{
			Dedicated:          self.fillDedicated(),
			Vps:                self.fillVps(),
			Companyname:        self.game.GetPlayer().GetCompanyName(),
			Location:           self.game.GetPlayer().GetLocation().Name,
			ElectricalNetworks: fmt.Sprintf("%d", nbpowerlines),
		}
		t.Execute(w, variables)
	} else {
		fmt.Fprintf(w, emptyAnswer)
	}
}

type StaticPage struct{}

func (self *StaticPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	mimetype := "text/plain"
	switch filepath.Ext(path) {
	case ".css":
		mimetype = "text/css"
	case ".js":
		mimetype = "text/javascript"
	case ".png":
		mimetype = "image/png"
	}
	w.Header().Set("Content-Type", mimetype)
	//	fmt.Println("StaticPage", path)
	if asset, err := global.Asset("assets/webshop" + path); err == nil {
		w.Write(asset)
	} else {
		fmt.Fprintf(w, emptyAnswer)
	}
}

func CheckGameRunning(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if the game is not loaded we just serve static content
		if game.GetPlayer() == nil && strings.HasPrefix(r.URL.Path, "/static/") == false {
			if asset, err := global.Asset("assets/webshop/shopclosed.html"); err == nil {
				w.Write(asset)
			} else {
				fmt.Fprintf(w, emptyAnswer)
			}
			return
		}
		h.ServeHTTP(w, r)
	})
}

func Webshop(gameObject *dctycoon.Game) {
	game = gameObject
	indexPage := &IndexPage{game: gameObject}
	staticContent := &StaticPage{}

	http.Handle("/", CheckGameRunning(indexPage))
	http.Handle("/static/", CheckGameRunning(staticContent))
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Error("ListenAndServe: ", err)
	}
}
