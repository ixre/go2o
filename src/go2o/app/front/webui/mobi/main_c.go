package mobi

import (
	"github.com/atnet/gof/app"
	"go2o/core/domain/interface/partner"
	"net/http"
)

type mainC struct {
	app.Context
}

func (this *mainC) Login(w http.ResponseWriter, r *http.Request) {
	this.Context.Template().Render(w, "views/ucenter/login.html", nil)
}

func (this *mainC) Index(w http.ResponseWriter, r *http.Request, p *partner.ValuePartner) {
	w.Write([]byte(p.Name))
}
