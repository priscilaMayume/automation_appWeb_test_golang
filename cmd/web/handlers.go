package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
)

// Define o caminho para os templates
var pathToTemplates = "./templates/"

// Handler para a página inicial
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	// Mapa para armazenar os dados do template
	var td = make(map[string]interface{})

	// Verifica se existe uma chave na sessão
	if app.Session.Exists(r.Context(), "test") {
		msg := app.Session.GetString(r.Context(), "test")
		td["test"] = msg
	} else {
		// Se não existir, adiciona uma chave com a data e hora atual
		app.Session.Put(r.Context(), "test", "Hit this page at "+time.Now().UTC().String())
	}

	// Renderiza o template
	_ = app.render(w, r, "home.page.gohtml", &TemplateData{Data: td})
}

// Estrutura de dados para passar para o template
type TemplateData struct {
	IP   string
	Data map[string]interface{}
}

// Função para renderizar o template
func (app *application) render(w http.ResponseWriter, r *http.Request, t string, data *TemplateData) error {
	// Faz o parse do template do disco
	parsedTemplate, err := template.ParseFiles(path.Join(pathToTemplates, t), path.Join(pathToTemplates, "base.layout.gohtml"))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	// Obtém o IP do contexto
	data.IP = app.ipFromContext(r.Context())

	// Executa o template, passando os dados
	err = parsedTemplate.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

// Handler para a página de login
func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	// Faz o parse dos dados do formulário
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Valida os dados do formulário
	form := NewForm(r.PostForm)
	form.Required("email", "password")

	if !form.Valid() {
		fmt.Fprint(w, "failed validation")
		return
	}

	// Obtém os dados do formulário
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	log.Println(email, password)

	fmt.Fprint(w, email)
}
