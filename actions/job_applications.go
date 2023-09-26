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
// Model: Singular (JobApplication)
// DB Table: Plural (job_applications)
// Resource: Plural (JobApplications)
// Path: Plural (/job_applications)
// View Template Folder: Plural (/templates/job_applications/)

// JobApplicationsResource is the resource for the JobApplication model
type JobApplicationsResource struct {
	buffalo.Resource
}

// List gets all JobApplications. This function is mapped to the path
// GET /job_applications
func (v JobApplicationsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	jobApplications := &models.JobApplications{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all JobApplications from the DB
	if err := q.All(jobApplications); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// Add the paginator to the context so it can be used in the template.
		c.Set("pagination", q.Paginator)

		c.Set("jobApplications", jobApplications)
		return c.Render(http.StatusOK, r.HTML("job_applications/index.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(jobApplications))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(jobApplications))
	}).Respond(c)
}

// Show gets the data for one JobApplication. This function is mapped to
// the path GET /job_applications/{job_application_id}
func (v JobApplicationsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty JobApplication
	jobApplication := &models.JobApplication{}

	// To find the JobApplication the parameter job_application_id is used.
	if err := tx.Find(jobApplication, c.Param("job_application_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		c.Set("jobApplication", jobApplication)

		return c.Render(http.StatusOK, r.HTML("job_applications/show.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(jobApplication))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(jobApplication))
	}).Respond(c)
}

// Destroy deletes a JobApplication from the DB. This function is mapped
// to the path DELETE /job_applications/{job_application_id}
func (v JobApplicationsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty JobApplication
	jobApplication := &models.JobApplication{}

	// To find the JobApplication the parameter job_application_id is used.
	if err := tx.Find(jobApplication, c.Param("job_application_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(jobApplication); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a flash message
		c.Flash().Add("success", T.Translate(c, "jobApplication.destroyed.success"))

		// Redirect to the index page
		return c.Redirect(http.StatusSeeOther, "/job_applications")
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(jobApplication))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(jobApplication))
	}).Respond(c)
}
