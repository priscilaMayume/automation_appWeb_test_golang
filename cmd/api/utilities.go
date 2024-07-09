package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Função para escrever uma resposta JSON
func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, wrap ...string) error {
	// out irá conter a versão final do JSON para enviar ao cliente
	var out []byte

	// decide se o payload JSON será encapsulado em uma tag geral JSON
	if len(wrap) > 0 {
		// encapsulador
		wrapper := make(map[string]interface{})
		wrapper[wrap[0]] = data
		jsonBytes, err := json.Marshal(wrapper)
		if err != nil {
			return err
		}
		out = jsonBytes
	} else {
		// sem encapsulador
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return err
		}
		out = jsonBytes
	}

	// define o tipo de conteúdo e o status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// escreve o JSON na resposta
	_, err := w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

// Função para enviar uma resposta JSON de erro
func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	type jsonError struct {
		Message string `json:"message"`
	}

	theError := jsonError{
		Message: err.Error(),
	}

	_ = app.writeJSON(w, statusCode, theError, "error")
}

// Função para ler e decodificar um JSON da requisição
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1024 * 1024 // um megabyte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// tentativa de decodificar os dados
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	// certifique-se de que há apenas um valor JSON no payload
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
