package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go_game_api.com/internal/data"
	validator "go_game_api.com/internal/validators"
	"go_game_api.com/pkg/utils"
)

// GET ALL USERS
func (app *application) getUsers(w http.ResponseWriter, r *http.Request) {
	validator := validator.New()
	qs := r.URL.Query()

	if len(qs) > 0 {
		// TODO
		queryStrings := &data.PaginatedData{}

		queryStrings.Page = utils.ReadQueryInt(qs, "page", validator)
		queryStrings.PerPage = utils.ReadQueryInt(qs, "perPage", validator)
		queryStrings.Sort = utils.ReadQueryString(qs, "sort", validator)

		if !validator.Valid() {
			app.failedValidationResponse(w, r, validator.Errors)
			return
		}
	}

	dbUsers, err := app.models.Users.GetAll()

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var users []data.User

	for _, dbUser := range dbUsers {

		user := data.User{
			Id:        dbUser.Id,
			Username:  dbUser.Username,
			Email:     dbUser.Email,
			Passhash:  dbUser.Passhash,
			Role:      dbUser.Role,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
		}

		users = append(users, user)
	}

	err = utils.WriteJSON(w, http.StatusOK, users, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// GET USER BY ID
func (app *application) getUserById(w http.ResponseWriter, r *http.Request) {
	userIdStr := chi.URLParam(r, "id")

	userId, err := strconv.Atoi(userIdStr)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetById(userId)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, user, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// POST USER
func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var userInput data.UserInput

	err := utils.ReadJSON(w, r, &userInput)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Printf("User %v", userInput)
	v := validator.New()

	if data.ValidateUserInput(v, &userInput); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user := &data.User{
		Username: userInput.Username,
		Email:    userInput.Email,
		// TODO add proper hashing
		Passhash: userInput.Password,
		Role:     data.UserRole,
	}

	err = app.models.Users.Insert(user)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/users/%d", user.Id))

	err = utils.WriteJSON(w, http.StatusCreated, user, headers)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// GET ALL Scores
func (app *application) getAllScores(w http.ResponseWriter, r *http.Request) {
	dbScores, err := app.models.Scores.GetAll()

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var scores []data.Score

	for _, dbScore := range dbScores {

		user := data.Score{
			Id:        dbScore.Id,
			Level:     dbScore.Level,
			Moves:     dbScore.Moves,
			Time:      dbScore.Time,
			UserId:    dbScore.UserId,
			CreatedAt: dbScore.CreatedAt,
			UpdatedAt: dbScore.UpdatedAt,
		}

		scores = append(scores, user)
	}

	fmt.Printf("values %v", scores)
	err = utils.WriteJSON(w, http.StatusOK, scores, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// POST SCORE
func (app *application) createScore(w http.ResponseWriter, r *http.Request) {
	var scoreInput data.ScoreInput

	err := utils.ReadJSON(w, r, &scoreInput)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateScoreInput(v, &scoreInput); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	score := &data.Score{
		Level:  scoreInput.Level,
		Moves:  scoreInput.Moves,
		Time:   scoreInput.Time,
		UserId: 2,
	}

	err = app.models.Scores.Insert(score)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/scores/%d", score.Id))

	err = utils.WriteJSON(w, http.StatusCreated, score, headers)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
