package grifts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"projectcollection/models"
	"time"

	. "github.com/gobuffalo/grift/grift"
	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/ssh"
)

type application struct {
	ID            uuid.UUID `json:"id" db:"id"`
	ProjectID     uuid.UUID `json:"project_id" db:"project_id"`
	ClientID      uuid.UUID `json:"client_id" db:"client_id"`
	RuntimeID     uuid.UUID `json:"runtime_id" db:"runtime_id"`
	DatabaseID    uuid.UUID `json:"database_id" db:"database_id"`
	EnvironmentID uuid.UUID `json:"environment_id" db:"environment_id"`
	Repository    string    `json:"repository" db:"repository"`
	Branch        string    `json:"branch" db:"branch"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

var _ = Namespace("processor", func() {

	Desc("application", "Application job processor")
	Add("application", func(c *Context) error {
		// Read the job_application table for new entries (where the processed_at is null)
		nextJobs := models.JobApplications{}
		models.DB.Where("processed_at is null").All(&nextJobs)
		for _, job := range nextJobs {
			// setup the execution time
			job.ProcessedAt = nulls.NewTime(time.Now())
			// execute the job.
			newApp := &application{}
			err := json.Unmarshal([]byte(job.NewParams.String), &newApp)
			if err != nil {
				// TODO: log error
				job.Response = nulls.NewString("Failed to unmarshal the job parameters. " + err.Error() + " " + job.NewParams.String)
				models.DB.Update(&job)
				continue
			}
			// execute the create project job.
			project := models.Project{}
			models.DB.Where("id = ?", newApp.ProjectID).First(&project)
			client := models.Client{}
			models.DB.Where("id = ?", newApp.ClientID).First(&client)
			dbtype := models.Dbtype{}
			models.DB.Where("id = ?", newApp.DatabaseID).First(&dbtype)
			runtime := models.Runtime{}
			models.DB.Where("id = ?", newApp.RuntimeID).First(&runtime)
			environment := models.Environment{}
			models.DB.Where("id = ?", newApp.EnvironmentID).First(&environment)

			hosts := models.Hosts{}
			models.DB.Where("environment_id = ?", environment.ID).All(&hosts)
			// create an application struct to pass to the server.
			app := &models.Application{
				Project:     &project,
				Client:      &client,
				Database:    &dbtype,
				Runtime:     &runtime,
				Environment: &environment,
				Repository:  newApp.Repository,
				Branch:      newApp.Branch,
			}

			response := ""
			for _, host := range hosts {
				response += executeServerCommand(&host, app, job.Type)
			}
			job.Response = nulls.NewString(response)
			models.DB.Update(&job)
		}
		return nil
	})

})

func executeServerCommand(host *models.Host, data *models.Application, command string) string {
	key, err := ioutil.ReadFile(host.SSHKey)
	if err != nil {
		return fmt.Sprintf("unable to read private key: %v - %v", err, host)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Sprintf("unable to parse private key: %v", err)
	}
	config := &ssh.ClientConfig{
		User: host.SSHUser,
		Auth: []ssh.AuthMethod{
			// Add in password check here for moar security.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.IP, host.SSHPort), config)
	if err != nil {
		return fmt.Sprintf("unable to connect to %s:%d with %s: %v", host.IP, host.SSHPort, host.SSHUser, err)
	}
	defer client.Close()
	ss, err := client.NewSession()
	if err != nil {
		return fmt.Sprintf("unable to create SSH session: %v", err)
	}
	defer ss.Close()
	// Creating the buffer which will hold the remotly executed command's output.
	var stdoutBuf bytes.Buffer
	var commandString string
	ss.Stdout = &stdoutBuf
	switch command {
	case "create":
		commandString = fmt.Sprintf("/usr/local/bin/setup-project.sh %s %s %s %s \"%s\"",
			data.Client.Name,
			data.Project.Name,
			data.Runtime.Name,
			data.Database.Name,
			data.Repository+" / "+data.Branch)
	case "destroy":
		commandString = fmt.Sprintf("/usr/local/bin/remove-project.sh %s %s",
			data.Client.Name,
			data.Project.Name)
	default:
		return "Unknown command"
	}

	ss.Run(commandString)
	return stdoutBuf.String()
}
