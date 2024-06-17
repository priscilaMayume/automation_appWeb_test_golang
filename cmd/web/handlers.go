package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
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

	if app.Session.Exists(r.Context(), "user") {
		td.User = app.Session.Get(r.Context(), "user").(data.User)
	}

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

func (app *application) UploadProfilePic(w http.ResponseWriter, r *http.Request) {
	
	// chama uma função que extrai um arquivo de um upload (requisição)
	files, err := app.UploadFiles(r, "./static/img")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// obtém o usuário da sessão
	user := app.Session.Get(r.Context(), "user").(data.User)

	// cria uma variável do tipo data.UserImage
	var i = data.UserImage{
		UserID: user.ID,
		FileName: files[0].OriginalFileName,
	}

	// insere a imagem do usuário em user_images
	_, err = app.DB.InsertUserImage(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// atualiza a variável de sessão "user"
	updatedUser, err := app.DB.GetUser(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.Session.Put(r.Context(), "user", updatedUser)

	// redireciona de volta para a página de perfil
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)


}

type UploadedFile struct {
	OriginalFileName string
	FileSize         int64
}

func (app *application) UploadFiles(r *http.Request, uploadDir string) ([]*UploadedFile, error) {
	var uploadedFiles []*UploadedFile

	// faz o parse do formulário multipart com um tamanho máximo de 5 MB
	err := r.ParseMultipartForm(int64(1024 * 1024 * 5))
	if err != nil {
		return nil, fmt.Errorf("the uploaded file is too big, and must be less than %d bytes", 1024*1024*5)
	}

	// itera sobre os arquivos no formulário multipart
	for _, fHeaders := range r.MultipartForm.File {
		for _, hdr := range fHeaders {
			uploadedFiles, err = func(uploadedFiles []*UploadedFile) ([]*UploadedFile, error) {
				var uploadedFile UploadedFile
				infile, err := hdr.Open()
				if err != nil {
					return nil, err
				}
				defer infile.Close()

				uploadedFile.OriginalFileName = hdr.Filename

				var outfile *os.File
				defer outfile.Close()

				// cria o arquivo no diretório de upload e copia o conteúdo do arquivo enviado para ele
				if outfile, err = os.Create(filepath.Join(uploadDir, uploadedFile.OriginalFileName)); nil != err {
					return nil, err
				} else {
					fileSize, err := io.Copy(outfile, infile)
					if err != nil {
						return nil, err
					}
					uploadedFile.FileSize = fileSize
				}

				uploadedFiles = append(uploadedFiles, &uploadedFile)

				return uploadedFiles, nil
			}(uploadedFiles)
			if err != nil {
				return uploadedFiles, err
			}
		}
	}

	return uploadedFiles, nil
}

