package handlers

import (
	"fmt"
	"net/http"

	"github.com/Vikram222726/bookings/pkg/config"
	"github.com/Vikram222726/bookings/pkg/models"
	"github.com/Vikram222726/bookings/pkg/render"
)

var Repo *Repository

type Repository struct{
	AppConfig *config.AppConfig
}

func NewRepo(a *config.AppConfig) (* Repository){
	return &Repository{
		AppConfig: a,
	}
}

func SetRepo(r *Repository){
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request){
	remoteIP := r.RemoteAddr
	m.AppConfig.Session.Put(r.Context(), "remote_ip", remoteIP) 

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request){
	// perform some logic here
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello World"
	fmt.Println("Session Data:", m.AppConfig.Session)

	remoteIP := m.AppConfig.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// send the data to be rendered on browser via template..
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
