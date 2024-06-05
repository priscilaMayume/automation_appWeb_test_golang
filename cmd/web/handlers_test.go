package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	var theTests = []struct {
		name               string
		url                string
		expectedStatusCode int
	}{
		{"home", "/", http.StatusOK},  // Teste para a página inicial com status OK
		{"404", "/fish", http.StatusNotFound},  // Teste para uma página não encontrada
	}

	routes := app.routes()

	// Cria um servidor de teste
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	// Itera pelos dados de teste
	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s: expected status %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestAppHome(t *testing.T) {
	var tests = []struct {
		name         string
		putInSession string
		expectedHTML string
	}{
		{"first visit", "", "<small>From Session:"},  // Teste para a primeira visita à página
		{"second visit", "hello, world!", "<small>From Session: hello, world!"},  // Teste para a segunda visita à página
	}

	for _, e := range tests {
		// Cria uma requisição HTTP GET para a raiz "/"
		req, _ := http.NewRequest("GET", "/", nil)

		req = addContextAndSessionToRequest(req, app)
		_ = app.Session.Destroy(req.Context())

		if e.putInSession != "" {
			app.Session.Put(req.Context(), "test", e.putInSession)
		}

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
	// Define o caminho para os templates com um template ruim
	pathToTemplates = "./testdata/"

	req, _ := http.NewRequest("GET", "/", nil)
	req = addContextAndSessionToRequest(req, app)
	rr := httptest.NewRecorder()

	err := app.render(rr, req, "bad.page.gohtml", &TemplateData{})
	if err == nil {
		t.Error("expected error from bad template, but did not get one")
	}
	pathToTemplates = "./../../templates/"
	
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

// Testa a função de login do aplicativo
func Test_app_Login(t *testing.T) {
	// Define os casos de teste com dados de entrada e resultados esperados
	var tests = []struct{
		name string // Nome do caso de teste
		postedData url.Values // Dados do formulário enviados na requisição
		expectedStatusCode int // Código de status HTTP esperado na resposta
		expectedLoc string // URL esperada para redirecionamento
	}{
		{
			name: "valid login", // Caso de teste para login válido
			postedData: url.Values{
				"email": {"admin@example.com"},
				"password": {"secret"},
			},
			expectedStatusCode: http.StatusSeeOther, // Espera redirecionamento 303
			expectedLoc: "/user/profile", // Espera redirecionar para perfil do usuário
		},
		{
			name: "missing form data", // Caso de teste para dados do formulário ausentes
			postedData: url.Values{
				"email": {""},
				"password": {""},
			},
			expectedStatusCode: http.StatusSeeOther, // Espera redirecionamento 303
			expectedLoc: "/", // Espera redirecionar para a página inicial
		},
		{
			name: "user not found", // Caso de teste para usuário não encontrado
			postedData: url.Values{
				"email": {"you@there.com"},
				"password": {"password"},
			},
			expectedStatusCode: http.StatusSeeOther, // Espera redirecionamento 303
			expectedLoc: "/", // Espera redirecionar para a página inicial
		},
		{
			name: "bad credentials", // Caso de teste para credenciais incorretas
			postedData: url.Values{
				"email": {"admin@example.com"},
				"password": {"password"},
			},
			expectedStatusCode: http.StatusSeeOther, // Espera redirecionamento 303
			expectedLoc: "/", // Espera redirecionar para a página inicial
		},
	}

	// Itera sobre cada caso de teste
	for _, e := range tests {
		// Cria uma nova requisição HTTP POST com os dados do formulário
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(e.postedData.Encode()))
		// Adiciona contexto e sessão à requisição
		req = addContextAndSessionToRequest(req, app)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		
		// Cria um ResponseRecorder para capturar a resposta
		rr := httptest.NewRecorder()
		// Define o handler da função de login
		handler := http.HandlerFunc(app.Login)
		// Executa a requisição
		handler.ServeHTTP(rr, req)

		// Verifica se o código de status retornado é o esperado
		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s: retornou o código de status errado; esperava %d, mas obteve %d", e.name, e.expectedStatusCode, rr.Code)
		}

		// Verifica se a URL de redirecionamento é a esperada
		actualLoc, err := rr.Result().Location()
		if err == nil {
			if actualLoc.String() != e.expectedLoc {
				t.Errorf("%s: esperava redirecionamento para %s, mas foi para %s", e.name, e.expectedLoc, actualLoc.String())
			}
		} else {
			t.Errorf("%s: cabeçalho de localização não definido", e.name)
		}
	}
}
