package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// TestForm_Has verifica se o formulário tem um campo específico.
func TestForm_Has(t *testing.T) {
	form := NewForm(nil)

	has := form.Has("whatever")
	if has {
		t.Error("o formulário mostra ter um campo quando não deveria")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = NewForm(postedData)

	has = form.Has("a")
	if !has {
		t.Error("mostra que o formulário não tem um campo quando deveria ter")
	}
}

// TestForm_Required verifica se os campos obrigatórios são validados corretamente.
func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := NewForm(r.PostForm)

	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("o formulário mostra ser válido quando campos obrigatórios estão faltando")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData

	form = NewForm(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("mostra que o POST não tem campos obrigatórios quando deveria ter")
	}
}

// TestForm_Check verifica se o método Check está funcionando corretamente.
func TestForm_Check(t *testing.T) {
	form := NewForm(nil)

	form.Check(false, "password", "password é obrigatório")
	if form.Valid() {
		t.Error("Valid() retorna falso, mas deveria ser verdadeiro ao chamar Check()")
	}
}

// TestForm_ErrorGet verifica se o método Get de Errors está funcionando corretamente.
func TestForm_ErrorGet(t *testing.T) {
	form := NewForm(nil)
	form.Check(false, "password", "password é obrigatório")
	s := form.Errors.Get("password")

	if len(s) == 0 {
		t.Error("deveria retornar um erro de Get, mas não retorna")
	}

	s = form.Errors.Get("whatever")
	if len(s) != 0 {
		t.Error("não deveria retornar um erro, mas retorna um")
	}
}
