package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/x/responder"

	"projectcollection/models"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Application)
// DB Table: Plural (applications)
// Resource: Plural (Applications)
// Path: Plural (/applications)
// View Template Folder: Plural (/templates/applications/)

// ApplicationsResource is the resource for the Application model
type ApplicationsResource struct {
	buffalo.Resource
}

// List gets all Applications. This function is mapped to the path
// GET /applications
func (v ApplicationsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	// check the permissions of the user. If it hasn't got permission for the "application.view" resource, return an error
	if !c.Value("current_user").(*models.User).HasPermissionFor("application.view") {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("You don't have permission to view applications"))
	}

	applications := &models.Applications{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Applications from the DB
	if err := q.Eager("Project").Eager("Client").Eager("Runtime").Eager("Database").Eager("Environment").Eager("Aliases.Alias").All(applications); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// Add the paginator to the context so it can be used in the template.
		c.Set("pagination", q.Paginator)

		c.Set("applications", applications)
		return c.Render(http.StatusOK, r.HTML("applications/index.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(applications))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(applications))
	}).Respond(c)
}

// Show gets the data for one Application. This function is mapped to
// the path GET /applications/{application_id}
func (v ApplicationsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	// check the permissions of the user. If it hasn't got permission for the "application.view" resource, return an error
	if !c.Value("current_user").(*models.User).HasPermissionFor("application.view") {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("You don't have permission to view applications"))
	}

	// Allocate an empty Application
	application := &models.Application{}

	// To find the Application the parameter application_id is used.
	if err := tx.Eager("Project").Eager("Client").Eager("Runtime").Eager("Database").Eager("Environment").Eager("Aliases.Alias").Find(application, c.Param("application_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		c.Set("application", application)

		return c.Render(http.StatusOK, r.HTML("applications/show.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(application))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(application))
	}).Respond(c)
}

// New renders the form for creating a new Application.
// This function is mapped to the path GET /applications/new
func (v ApplicationsResource) New(c buffalo.Context) error {
	// Get the DB connection from the context
	c.Set("application", &models.Application{})
	if err := v.setSelectLists(c); err != nil {
		return err
	}
	// check the permissions of the user. If it hasn't got permission for the "application.create" resource, return an error
	if !c.Value("current_user").(*models.User).HasPermissionFor("application.create") {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("You don't have permission to create applications"))
	}

	return c.Render(http.StatusOK, r.HTML("applications/new.plush.html"))
}

// Create adds a Application to the DB. This function is mapped to the
// path POST /applications
func (v ApplicationsResource) Create(c buffalo.Context) error {
	// check the permissions of the user. If it hasn't got permission for the "application.create" resource, return an error
	if !c.Value("current_user").(*models.User).HasPermissionFor("application.create") {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("You don't have permission to create applications"))
	}
	// Allocate an empty Application
	application := &models.Application{}

	// Bind application to the html form elements
	if err := c.Bind(application); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(application)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the new.html template that the user can
			// correct the input.
			c.Set("application", application)
			if err := v.setSelectLists(c); err != nil {
				return err
			}

			return c.Render(http.StatusUnprocessableEntity, r.HTML("applications/new.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}
	// Create an entry to the job_application table
	jobApplication := &models.JobApplication{
		Type:      "create",
		NewParams: nulls.NewString(application.String()),
	}
	tx.Create(jobApplication)

	// add the aliases to the application insert to the application_to_aliases table
	for _, alias := range application.NewAliases {
		applicationToAlias := &models.ApplicationToAlias{
			ApplicationID: application.ID,
			AliasID:       alias,
		}
		tx.Create(applicationToAlias)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "application.created.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/applications/%v", application.ID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.JSON(application))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.XML(application))
	}).Respond(c)
}

// Edit renders a edit form for a Application. This function is
// mapped to the path GET /applications/{application_id}/edit
func (v ApplicationsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	// check the permissions of the user. If it hasn't got permission for the "application.edit" resource, return an error
	if !c.Value("current_user").(*models.User).HasPermissionFor("application.edit") {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("You don't have permission to edit applications"))
	}

	// Allocate an empty Application
	application := &models.Application{}

	if err := tx.Eager().Find(application, c.Param("application_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("application", application)
	if err := v.setSelectLists(c); err != nil {
		return err
	}
	return c.Render(http.StatusOK, r.HTML("applications/edit.plush.html"))
}

// Update changes a Application in the DB. This function is mapped to
// the path PUT /applications/{application_id}
func (v ApplicationsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	// check the permissions of the user. If it hasn't got permission for the "application.edit" resource, return an error
	if !c.Value("current_user").(*models.User).HasPermissionFor("application.edit") {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("You don't have permission to edit applications"))
	}

	// Allocate an empty Application
	application := &models.Application{}

	if err := tx.Find(application, c.Param("application_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind Application to the html form elements
	if err := c.Bind(application); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(application)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the edit.html template that the user can
			// correct the input.
			c.Set("application", application)
			if err := v.setSelectLists(c); err != nil {
				return err
			}

			return c.Render(http.StatusUnprocessableEntity, r.HTML("applications/edit.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	// delete all the aliases from the application_to_aliases table
	applicationToAliases := &models.ApplicationToAliases{}
	if err := tx.Where("application_id = ?", application.ID).All(applicationToAliases); err != nil {
		return err
	}
	for _, applicationToAlias := range *applicationToAliases {
		if err := tx.Destroy(&applicationToAlias); err != nil {
			return err
		}
	}
	// add the aliases to the application insert to the application_to_aliases table
	for _, alias := range application.NewAliases {
		applicationToAlias := &models.ApplicationToAlias{
			ApplicationID: application.ID,
			AliasID:       alias,
		}
		tx.Create(applicationToAlias)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "application.updated.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/applications/%v", application.ID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(application))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(application))
	}).Respond(c)
}

// Destroy deletes a Application from the DB. This function is mapped
// to the path DELETE /applications/{application_id}
func (v ApplicationsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	// check the permissions of the user. If it hasn't got permission for the "application.delete" resource, return an error
	if !c.Value("current_user").(*models.User).HasPermissionFor("application.delete") {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("You don't have permission to delete applications"))
	}

	// Allocate an empty Application
	application := &models.Application{}

	// To find the Application the parameter application_id is used.
	if err := tx.Find(application, c.Param("application_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Create an entry to the job_application table
	jobApplication := &models.JobApplication{
		Type:       "destroy",
		OrigParams: nulls.NewString(application.String()),
	}
	tx.Create(jobApplication)

	if err := tx.Destroy(application); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a flash message
		c.Flash().Add("success", T.Translate(c, "application.destroyed.success"))

		// Redirect to the index page
		return c.Redirect(http.StatusSeeOther, "/applications")
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(application))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(application))
	}).Respond(c)
}

// setSelectLists sets the values for the selectors in the application forms.
func (v ApplicationsResource) setSelectLists(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	environments := models.Environments{}
	tx.All(&environments)
	c.Set("environments", environments)

	runtimes := models.Runtimes{}
	tx.All(&runtimes)
	c.Set("runtimes", runtimes)

	databases := models.Dbtypes{}
	tx.All(&databases)
	c.Set("databases", databases)

	clients := models.Clients{}
	tx.All(&clients)
	c.Set("clients", clients)

	projects := models.Projects{}
	tx.All(&projects)
	c.Set("projects", projects)

	aliases := models.Aliases{}
	aliasQuery := tx.Q()
	aliasQuery.LeftJoin("application_to_aliases", "aliases.id = application_to_aliases.alias_id")
	if c.Param("application_id") != "" {
		aliasQuery.Where("application_to_aliases.application_id IS NULL OR application_to_aliases.application_id = ?", c.Param("application_id"))
	} else {
		aliasQuery.Where("application_to_aliases.application_id IS NULL")
	}

	aliasQuery.All(&aliases)
	c.Set("aliases", aliases)

	return nil
}
