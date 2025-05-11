package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	//var input struct {
	//	Username string `json:"username"`
	//	Email    string `json:"email"`
	//	Password string `json:"password"`
	//}

	fmt.Fprintf(w, "Create user")

	//err := app.readJSON(w, r, &input)
	//if err != nil {
	//	app.badRequestResponse(w, r, err)
	//	return
	//}
	//
	//user := &models.User{
	//	Username: input.Username,
	//	Email:    input.Email,
	//	Password: input.Password,
	//}
	//
	//err = app.users.Insert(user)
	//if err != nil {
	//	app.serverErrorResponse(w, r, err)
	//	return
	//}
	//
	//headers := make(http.Header)
	//headers.Set("Location", fmt.Sprintf("/users/%d", user.ID))
	//
	//err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, headers)
	//if err != nil {
	//	app.serverErrorResponse(w, r, err)
	//}
}

func (app *application) showUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	user, err := app.users.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	user, err := app.users.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user.Username = input.Username
	user.Email = input.Email
	user.Password = input.Password

	err = app.users.Update(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) authenticateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.users.Authenticate(input.Email, input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/users/%d", user.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
