package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_application_handlers(t *testing.T) {
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

	pathToTemplates = "./../../templates/"

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
	// Cria uma requisição
	req, _ := http.NewRequest("GET", "/", nil)

	req = addContextAndSessionToRequest(req, app)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(app.Home)

	handler.ServeHTTP(rr, req)

	// Verifica o código de status
	if rr.Code != http.StatusOK {
		t.Errorf("TestAppHome returned wrong status code; expected 200 but got %d", rr.Code)
	}

	// Lê o corpo da resposta
	body, _ := io.ReadAll(rr.Body)
	if !strings.Contains(string(body), `<small>From Session:`) {
		t.Error("did not find correct text in html")
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
