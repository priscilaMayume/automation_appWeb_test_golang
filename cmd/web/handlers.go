package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
)

var pathToTemplates string // Referência externa à variável pathToTemplates que está no setup.go

// Handler para a página inicial
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	log.Println("Home handler called") // Debug
	var td = make(map[string]any)

	// Verifica se a sessão "test" existe
	if app.Session.Exists(r.Context(), "test") {
		msg := app.Session.GetString(r.Context(), "test")
		td["test"] = msg
		log.Println("Session exists with value:", msg) // Debug
	} else {
		// Adiciona um valor à sessão "test" com a data e hora atual
		newValue := "Hit this page at " + time.Now().UTC().String()
		app.Session.Put(r.Context(), "test", newValue)
		log.Println("Session created with value:", newValue) // Debug
	}
	err := app.render(w, r, "home.page.gohtml", &TemplateData{Data: td})
	if err != nil {
		log.Println("Error rendering template:", err) // Debug
	}
}

// Estrutura para os dados do template
type TemplateData struct {
	IP   string
	Data map[string]any
}

// Método para renderizar templates
func (app *application) render(w http.ResponseWriter, r *http.Request, t string, data *TemplateData) error {
	// Faz o parse do template a partir do disco.
	parsedTemplate, err := template.ParseFiles(path.Join(pathToTemplates, t), path.Join(pathToTemplates, "base.layout.gohtml"))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Println("Error parsing templates:", err) // Debug
		return err
	}

	data.IP = app.ipFromContext(r.Context())

	// Executa o template, passando os dados, se houver
	err = parsedTemplate.Execute(w, data)
	if err != nil {
		log.Println("Error executing template:", err) // Debug
		return err
	}

	return nil
}

// Handler para login
func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err) // Debug
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Valida os dados
	form := NewForm(r.PostForm)
	form.Required("email", "password")

	// Verifica se o formulário é válido
	if !form.Valid() {
		fmt.Fprint(w, "failed validation")
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	log.Println(email, password)

	fmt.Fprint(w, email)
}
