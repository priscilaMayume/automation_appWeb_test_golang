package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/priscilaMayume/automation_appWeb_test_golang/pkg/data"
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

// Handler para a página Profile
func (app *application) Profile(w http.ResponseWriter, r *http.Request) {
	// Renderiza o template
	_ = app.render(w, r, "profile.page.gohtml", &TemplateData{})
}


// Estrutura de dados para passar para o template
type TemplateData struct {
	IP   string
	Data map[string]any
	Error string
	Flash string
	User data.User
}

// Função para renderizar o template
func (app *application) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) error {
	// Faz o parse do template do disco
	parsedTemplate, err := template.ParseFiles(path.Join(pathToTemplates, t), path.Join(pathToTemplates, "base.layout.gohtml"))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	// Obtém o IP do contexto
	td.IP = app.ipFromContext(r.Context())
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Flash = app.Session.PopString(r.Context(), "flash")

	// Executa o template, passando os dados
	err = parsedTemplate.Execute(w, td)
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
		// redirect para o a pg de login com msg de erro
		app.Session.Put(r.Context(), "error", "Invalid login Credential")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Obtém os dados do formulário
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := app.DB.GetUserByEmail(email)
	if err != nil {
		// redirect para o a pg de login com msg de erro
		app.Session.Put(r.Context(), "error", "Invalid login!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if !app.authenticate(r, user, password) {
		app.Session.Put(r.Context(), "error", "Invalid Login!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// autenticação do user

	// Redirect de autenticação errada

	// Previnir fixation attack
	_ = app.Session.RenewToken(r.Context())

	// armazena menssagem de sucesso na sessao

	//redirect para outras paginas
	app.Session.Put(r.Context(), "flash", "Successfully Logged In!")
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

func (app *application) authenticate(r *http.Request, user *data.User, password string) bool {
	// Verifica se a senha fornecida corresponde à senha armazenada para o usuário
	if valid, err := user.PasswordMatches(password); err != nil || !valid {
		// Se houve um erro ou a senha não for válida, retorna falso
		return false
	}

	// Se a senha for válida, armazena a estrutura do usuário na sessão
	app.Session.Put(r.Context(), "user", user)
	// Retorna verdadeiro para indicar que a autenticação foi bem-sucedida
	return true
}
