package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.alexedwards.net/internal/data"
)
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Title   string   `json:"title"`
        Year    int32    `json:"year"`
        Runtime int32    `json:"runtime"`
        Genres  []string `json:"genres"`
    }
    // Use the new readJSON() helper to decode the request body into the input struct.
    // If this returns an error we send the client the error message along with a 400
    // Bad Request status code, just like before.
    err := app.readJSON(w, r, &input)
    if err != nil {
        app.errorResponse(w, r, http.StatusBadRequest, err.Error())
        return
    }
    fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
			// Use the new notFoundResponse() helper.
			app.notFoundResponse(w, r)
			return
	}
	movie := data.Movie{
			ID:        id,
			CreatedAt: time.Now(),
			Title:     "Casablanca",
			Runtime:   102,
			Genres:    []string{"drama", "romance", "war"},
			Version:   1,
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
			// Use the new serverErrorResponse() helper.
			app.serverErrorResponse(w, r, err)
	}
}