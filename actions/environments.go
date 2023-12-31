package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/x/responder"

	"projectcollection/models"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Environment)
// DB Table: Plural (environments)
// Resource: Plural (Environments)
// Path: Plural (/environments)
// View Template Folder: Plural (/templates/environments/)

// EnvironmentsResource is the resource for the Environment model
type EnvironmentsResource struct {
	buffalo.Resource
}

// List gets all Environments. This function is mapped to the path
// GET /environments
func (v EnvironmentsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	// check the permissions of the user. If it hasn't got permission for the "environment.view" resource, return an error
	if !c.Value("current_user").(*models.User).HasPermissionFor("environment.view") {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("You don't have permission to view environments"))
	}

	environments := &models.Environments{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Environments from the DB
	if err := q.All(environments); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// Add the paginator to the context so it can be used in the template.
		c.Set("pagination", q.Paginator)

		c.Set("environments", environments)
		return c.Render(http.StatusOK, r.HTML("environments/index.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(environments))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(environments))
	}).Respond(c)
}

// Show gets the data for one Environment. This function is mapped to
// the path GET /environments/{environment_id}
func (v EnvironmentsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	// check the permissions of the user. If it hasn't got permission for the "environment.view" resource, return an error
	if !c.Value("current_user").(*models.User).HasPermissionFor("environment.view") {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("You don't have permission to view environments"))
	}

	// Allocate an empty Environment
	environment := &models.Environment{}

	// To find the Environment the parameter environment_id is used.
	if err := tx.Find(environment, c.Param("environment_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		c.Set("environment", environment)

		return c.Render(http.StatusOK, r.HTML("environments/show.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(environment))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(environment))
	}).Respond(c)
}

// New renders the form for creating a new Environment.
// This function is mapped to the path GET /environments/new
func (v EnvironmentsResource) New(c buffalo.Context) error {
	// check the permissions of the user. If it hasn't got permission for the "environment.create" resource, return an error
	if !c.Value("current_user").(*models.User).HasPermissionFor("environment.create") {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("You don't have permission to create environments"))
	}
	c.Set("environment", &models.Environment{})

	return c.Render(http.StatusOK, r.HTML("environments/new.plush.html"))
}

// Create adds a Environment to the DB. This function is mapped to the
// path POST /environments
func (v EnvironmentsResource) Create(c buffalo.Context) error {
	// Allocate an empty Environment
	environment := &models.Environment{}

	// Bind environment to the html form elements
	if err := c.Bind(environment); err != nil {
		return err
	}
	// check the permissions of the user. If it hasn't got permission for the "environment.create" resource, return an error
	if !c.Value("current_user").(*models.User).HasPermissionFor("environment.create") {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("You don't have permission to create environments"))
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(environment)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the new.html template that the user can
			// correct the input.
			c.Set("environment", environment)

			return c.Render(http.StatusUnprocessableEntity, r.HTML("environments/new.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "environment.created.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/environments/%v", environment.ID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.JSON(environment))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.XML(environment))
	}).Respond(c)
}

// Edit renders a edit form for a Environment. This function is
// mapped to the path GET /environments/{environment_id}/edit
func (v EnvironmentsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	// check the permissions of the user. If it hasn't got permission for the "environment.edit" resource, return an error
	if !c.Value("current_user").(*models.User).HasPermissionFor("environment.edit") {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("You don't have permission to edit environments"))
	}

	// Allocate an empty Environment
	environment := &models.Environment{}

	if err := tx.Find(environment, c.Param("environment_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("environment", environment)
	return c.Render(http.StatusOK, r.HTML("environments/edit.plush.html"))
}

// Update changes a Environment in the DB. This function is mapped to
// the path PUT /environments/{environment_id}
func (v EnvironmentsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	// check the permissions of the user. If it hasn't got permission for the "environment.edit" resource, return an error
	if !c.Value("current_user").(*models.User).HasPermissionFor("environment.edit") {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("You don't have permission to edit environments"))
	}

	// Allocate an empty Environment
	environment := &models.Environment{}

	if err := tx.Find(environment, c.Param("environment_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind Environment to the html form elements
	if err := c.Bind(environment); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(environment)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the edit.html template that the user can
			// correct the input.
			c.Set("environment", environment)

			return c.Render(http.StatusUnprocessableEntity, r.HTML("environments/edit.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "environment.updated.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/environments/%v", environment.ID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(environment))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(environment))
	}).Respond(c)
}

// Destroy deletes a Environment from the DB. This function is mapped
// to the path DELETE /environments/{environment_id}
func (v EnvironmentsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	// check the permissions of the user. If it hasn't got permission for the "environment.delete" resource, return an error
	if !c.Value("current_user").(*models.User).HasPermissionFor("environment.delete") {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("You don't have permission to delete environments"))
	}

	// Allocate an empty Environment
	environment := &models.Environment{}

	// To find the Environment the parameter environment_id is used.
	if err := tx.Find(environment, c.Param("environment_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(environment); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a flash message
		c.Flash().Add("success", T.Translate(c, "environment.destroyed.success"))

		// Redirect to the index page
		return c.Redirect(http.StatusSeeOther, "/environments")
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(environment))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(environment))
	}).Respond(c)
}
