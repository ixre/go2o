package apiserv

import (
	"com/service/goclient"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/atnet/gof/app"
	"github.com/atnet/gof/net/jsv"
)

type websocketC struct {
	app.Context
}

func (this *websocketC) Login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func (this *websocketC) Test(w http.ResponseWriter, r *http.Request) {
	b, t, msg := goclient.Member.Login("newmin", "123000")
	if b {
		w.Write([]byte("[Login]:Sucessfull." + t))
	} else {
		w.Write([]byte("[Login]:Failed." + msg))
	}
}

func (this *websocketC) Partner(w http.ResponseWriter, r *http.Request) {
	buffer := goclient.Redirect.Post([]byte(fmt.Sprintf(
		`{"partner_id":"%s","secret":"%s"}>>Partner.GetPartner`,
		r.FormValue("partner_id"), r.FormValue("secret"))), 512)
	w.Write(buffer)
}

func (this *websocketC) Category(w http.ResponseWriter, r *http.Request) {
	buffer := goclient.Redirect.Post([]byte(fmt.Sprintf(
		`{"partner_id":"%s","secret":"%s"}>>Partner.Category`,
		r.FormValue("partner_id"), r.FormValue("secret"))), 2048)

	var v jsv.Result
	jsv.JsonCodec.Unmarshal(buffer, &v)
	b, _ := json.Marshal(v.Data)
	w.Write(b)
}
