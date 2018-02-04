package webshop

import (
	"fmt"
	"net/http"

	"github.com/nzin/dctycoon/global"

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

func (self *IndexPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//	fmt.Println("IndexPage")

	t := template.New("index page")
	if indexAsset, err := global.Asset("assets/webshop/index.tmpl"); err == nil {
		t, _ := t.Parse(string(indexAsset))
		t.Execute(w, nil)
	} else {
		fmt.Fprintf(w, emptyAnswer)
	}
}

type StaticPage struct{}

func (self *StaticPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	//	fmt.Println("StaticPage", path)
	if asset, err := global.Asset("assets/webshop" + path); err == nil {
		fmt.Fprintf(w, string(asset))
	} else {
		fmt.Fprintf(w, emptyAnswer)
	}
}

func CheckGameRunning(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if game.GetPlayer() == nil {
			fmt.Fprintf(w, emptyAnswer)
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
