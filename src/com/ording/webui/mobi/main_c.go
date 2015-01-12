package mobi

import (
	"com/ording/entity"
	"net/http"
	"github.com/atnet/gof/app"
)

type mainC struct {
	app.Context
}

func (this *mainC) Login(w http.ResponseWriter, r *http.Request) {
	this.Context.Template().Render(w, "views/ucenter/login.html", nil)
}

func (this *mainC) Index(w http.ResponseWriter, r *http.Request, p *entity.Partner) {
	w.Write([]byte(p.Name))
}
