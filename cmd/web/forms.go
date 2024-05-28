package main

import (
	"net/url"
	"strings"
)

// errors é um tipo de conveniência, para que possamos ter uma função associada ao nosso mapa.
type errors map[string][]string

func (e errors) Get(field string) string {
	errorSlice := e[field]
	if len(errorSlice) == 0 {
		return ""
	}

	return errorSlice[0]
}

// Add adiciona uma mensagem de erro para um determinado campo do formulário.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Form é o tipo usado para instanciar a validação de formulário
type Form struct {
	Data   url.Values
	Errors errors
}

// NewForm inicializa uma estrutura de formulário
func NewForm(data url.Values) *Form {
	return &Form{
		Data:   data,
		Errors: map[string][]string{},
	}
}

// Has verifica se o formulário possui um determinado campo
func (f *Form) Has(field string) bool {
	x := f.Data.Get(field)
	if x == "" {
		return false
	}
	return true
}

// Required verifica os campos obrigatórios
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Data.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "Este campo não pode estar em branco")
		}
	}
}

// Check é uma verificação genérica de validação. Podemos passar qualquer expressão
// que avalie como um booleano como o primeiro parâmetro.
func (f *Form) Check(ok bool, key, message string) {
	if !ok {
		f.Errors.Add(key, message)
	}
}

// Valido retorna true se não houver erros, caso contrário false
func (f *Form) Valido() bool {
	return len(f.Errors) == 0
}
