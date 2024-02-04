package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Vikram222726/bookings/internals/config"
	"github.com/Vikram222726/bookings/internals/models"
	"github.com/Vikram222726/bookings/internals/render"
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

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request){
	// perform some logic here
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello World"
	fmt.Println("Session Data:", m.AppConfig.Session)

	remoteIP := m.AppConfig.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// send the data to be rendered on browser via template..
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation allows the users to submit form and make reservation
func (m *Repository) Reservations(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{})
}

// General renders the general rooms service
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the major rooms service
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Search for rooms availability
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

type jsonResponse struct {
	Result string `json:"result"`
	Status int `json:"status"`
	Message string `json:"message"`
}

// This checks whether the hotel rooms are available for booking and returns a JSON response
func (m *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request){
	response := jsonResponse{
		Result: "success",
		Status: 200,
		Message: "Room Available!",
	}

	output, err := json.MarshalIndent(response, "", "    ")
	if err != nil{
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}
 
// Post data for rooms availability
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request){
	startDate := r.Form.Get("start")
	endDate := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("Rooms are available from start date %s to end date %s", startDate, endDate)))
}

// Provides contact information
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// Provides contact information
func (m *Repository) MakeReservations(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{})
}