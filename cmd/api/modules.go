package main

import (
	"assignment1/internal/data"
	"assignment1/internal/validator"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Add a createModuleInfoHandler for the "POST /v1/modules" endpoint
func (app *application) createModuleInfoHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		ModuleName     string   `json:"module_name"`
		ModuleDuration int32    `json:"module_duration"`
		ExamType       []string `json:"exam_type"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	module := &data.Module{
		ModuleName:     input.ModuleName,
		ModuleDuration: input.ModuleDuration,
		ExamType:       input.ExamType,
	}

	v := validator.New()
	// Call the ValidateModule() function
	if data.ValidateModule(v, module); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Insert() method on our movies model
	err = app.models.Modules.Insert(module)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/modules/%d", module.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"module": module}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// getModuleInfoHandler for the "GET /v1/modules/:id" endpoint
func (app *application) getModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {

		app.notFoundResponse(w, r)
		return
	}

	module, err := app.models.Modules.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"module": module}, nil)
	if err != nil {

		app.serverErrorResponse(w, r, err)
	}

}

// editModuleInfoHandler for the "PUT /v1/modules/:id" endpoint
func (app *application) editModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	module, err := app.models.Modules.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		ModuleName     string   `json:"module_name"`
		ModuleDuration int32    `json:"module_duration"`
		ExamType       []string `json:"exam_type"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	module.UpdatedAt = time.Now()
	module.ModuleName = input.ModuleName
	module.ModuleDuration = input.ModuleDuration
	module.ExamType = input.ExamType

	v := validator.New()

	if data.ValidateModule(v, module); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Modules.Update(module)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"module": module}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// deleteModuleInfoHandler for the "DELETE /v1/modules/:id" endpoint
func (app *application) deleteModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Modules.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "module successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
