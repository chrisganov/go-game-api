package main

import (
	"net/http"

	"github.com/chrisganov/go-game-api/pkg/utils"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// TODO:
	responseBody := map[string]string{
		"status":      "available",
		"environment": "dev",
		"version":     "1",
	}

	err := utils.WriteJSON(w, http.StatusOK, responseBody, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
