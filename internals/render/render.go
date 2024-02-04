package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Vikram222726/bookings/internals/config"
	"github.com/Vikram222726/bookings/internals/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

// This fn will read files from disk, parse them and load them to browser
// everytime this endpoint is rendered...
func RenderTemplateWithoutCache(w http.ResponseWriter, temp string){
	parsedTemplate, _ := template.ParseFiles("./templates/" + temp, "./templates/base.layout.tmpl")
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("Got error while parsing template")
		return
	}
	fmt.Println("Successfully parsed template!")
}


// In order to stop reading templated from disk, parsing that and loading
// that file to browser, we can use cache in form on map DS..
var templateCache = make(map[string]*template.Template)

func RenderTemplateWithCache(w http.ResponseWriter, temp string){
	var tmplPtr *template.Template
	var err error

	// Check if temp template is already present in templateCache...
	_, isMap := templateCache[temp]
	if !isMap{
		// Add template to the templateCache
		fmt.Println("Reading template from disk and adding to cache...")
		err = AddTemplateToCache(temp)
		if err != nil {
			fmt.Println("Got error while adding template to cache", err)
		}
	}else{
		// Fetch template ptr from cache and return the response...
		fmt.Println("Reading template from cache...")
	}

	tmplPtr = templateCache[temp]
	err = tmplPtr.Execute(w, nil)
	if err != nil{
		fmt.Println("Got error while writing template to browser", err)
	}
}

func AddTemplateToCache(temp string) error {
	templateList := []string{
		fmt.Sprintf("./templates/%s", temp),
		"./templates/base.layout.tmpl",
	}

	parsedTemplate, err := template.ParseFiles(templateList...)
	if err != nil{
		return err
	}

	templateCache[temp] = parsedTemplate
	return nil
}

func NewTemplate(a *config.AppConfig){
	app = a
}

// Now we are going to render the template using third and more complex
// method, where:
// 1. We'll first create cache for all the templates the very first time
// 2. We'll fetch the required asked template from this cache
// 3. We'll render the template on the browser...

func AddDefaultData(tempData *models.TemplateData, r *http.Request) *models.TemplateData {
	// Add logic to render different template data based on different temp values
	tempData.CSRFToken = nosurf.Token(r)

	return tempData
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, temp string, tempData *models.TemplateData){
	// First Create cache..
	var tempCache map[string]*template.Template
	if app.UseCache{
		tempCache = app.TemplateCache
	}else{
		tempCache, _ = CreateTemplateCache()
	}

	// Second fetch the required template from the cache
	reqdTemplate, tempPresent := tempCache[temp]
	if !tempPresent{
		log.Fatal("template not present in temp cache")
	}

	// Third render the template to the browser
	// err = reqdTemplate.Execute(w, nil)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	buf := new(bytes.Buffer)

	tempData = AddDefaultData(tempData, r)

	err := reqdTemplate.Execute(buf, tempData)
	if err != nil{
		log.Fatal(err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateTemplateCache()(map[string]*template.Template, error){
	//tempCache := make(map[string]*template.Template)
	fmt.Println("Started creating new cache")
	tempCache := map[string]*template.Template{} // can also initialize map like this

	// Now get all of the files ending with *.page.tmpl
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return tempCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages{
		name := filepath.Base(page)
		// name will give us home.page.tmpl
		// Now we'll parse the page and store that inside new template with variable named as name
		tempSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return tempCache, nil
		}
		
		layoutFiles, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return tempCache, nil
		}

		if len(layoutFiles) > 0{
			tempSet, err = tempSet.ParseGlob("./templates/*.layout.tmpl")
			if err != nil{
				return tempCache, nil
			}
		}

		tempCache[name] = tempSet
	}

	return tempCache, nil
}