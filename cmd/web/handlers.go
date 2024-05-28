package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
)

var pathToTemplates = "./templates/"

// Função para lidar com a rota Home
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	_ = app.render(w, r, "home.page.gohtml", &TemplateData{})

}

type TemplateData struct {
	IP   string
	Data map[string]any
}

func (app *application) render(w http.ResponseWriter, r *http.Request, t string, data *TemplateData) error {
		// analisar o modelo a partir do disco
	parsedTemplate, err := template.ParseFiles(path.Join(pathToTemplates, t))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	data.IP = app.ipFromContext(r.Context())

	// executa o modelo, passando-lhe dados, se existirem
	err = parsedTemplate.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	log.Println(email, password)
	fmt.Fprint(w, email)

	
}