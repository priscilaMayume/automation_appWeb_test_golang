package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	    // Configura o caminho relativo para os templates
		relativePath := "./../../templates/"
		absPath, err := filepath.Abs(relativePath)
		if err != nil {
			t.Fatalf("Error getting absolute path: %v", err)
		}
		pathToTemplates = absPath

	// Definição dos testes com nome, URL e código de status esperado
	var theTests = []struct {
		name               string
		url                string
		expectedStatusCode int
	}{
		{"home", "/", http.StatusOK},
		{"404", "/fish", http.StatusNotFound},
	}

	routes := app.routes()

	// Cria um servidor de teste
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	// Percorre os dados de teste
	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		// Verifica o código de status da resposta
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s: expected status %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestAppHome(t *testing.T) {
	// Configura o caminho relativo para os templates
	relativePath := "./../../templates/"
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		t.Fatalf("Error getting absolute path: %v", err)
	}
	pathToTemplates = absPath

	// Definição dos testes com nome, valor a ser colocado na sessão e HTML esperado
	var tests = []struct {
		name         string
		putInSession string
		expectedHTML string
	}{
		{"first visit", "", "<small>From Session:"},
		{"second visit", "hello world!", "<small>From Session: hello world!"},
	}

	// Percorre os dados de teste
	for _, e := range tests {
		// Cria uma requisição HTTP GET para a raiz "/"
		req, _ := http.NewRequest("GET", "/", nil)
		req = addContextAndSessionToRequest(req, app)
		_ = app.Session.Destroy(req.Context())

		// Coloca um valor na sessão, se necessário
		if e.putInSession != "" {
			app.Session.Put(req.Context(), "test", e.putInSession)
		}

		// Cria um gravador de resposta HTTP
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.Home)
		handler.ServeHTTP(rr, req)

		// Verifica o código de status da resposta
		if rr.Code != http.StatusOK {
			t.Errorf("TestAppHome returned wrong status code; expected 200 but got %d", rr.Code)
		}

		// Lê o corpo da resposta
		body, _ := io.ReadAll(rr.Body)
		if !strings.Contains(string(body), e.expectedHTML) {
			t.Errorf("%s: did not find %s in response body", e.name, e.expectedHTML)
		}
	}
}

func TestApp_renderWithBadTemplate(t *testing.T) {
    // Configura o caminho relativo para os templates
    relativePath := "./../../templates/"
    absPath, err := filepath.Abs(relativePath)
    if err != nil {
        t.Fatalf("Error getting absolute path: %v", err)
    }
    pathToTemplates = absPath

    // Define o caminho para o template com erro
    pathToTemplates = "./testdata/"

    // Cria uma requisição HTTP GET para a raiz "/"
    req, _ := http.NewRequest("GET", "/", nil)
    req = addContextAndSessionToRequest(req, app)

    // Cria um gravador de resposta HTTP
    rr := httptest.NewRecorder()

    // Tenta renderizar o template "bad.page.gohtml"
    err = app.render(rr, req, "bad.page.gohtml", &TemplateData{})

    // Verifica se ocorreu um erro ao renderizar o template
    if err == nil {
        t.Error("esperado erro ao renderizar template ruim, mas não ocorreu")
    }
}

// getCtx retorna um contexto com um valor adicionado
func getCtx(req *http.Request) context.Context {
	ctx := context.WithValue(req.Context(), contextUserKey, "unknown")
	return ctx
}

// addContextAndSessionToRequest adiciona contexto e sessão à requisição
func addContextAndSessionToRequest(req *http.Request, app application) *http.Request {
	req = req.WithContext(getCtx(req))
	ctx, _ := app.Session.Load(req.Context(), req.Header.Get("X-Session"))
	return req.WithContext(ctx)
}