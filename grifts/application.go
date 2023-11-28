package grifts

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"projectcollection/models"
	"time"

	. "github.com/gobuffalo/grift/grift"
	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

const (
	wsAddress = "localhost:3000"
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
		// create a web socket connection to the client
		u := url.URL{Scheme: "ws", Host: wsAddress, Path: "/wsbc"}
		conn, _, errWebsocket := websocket.DefaultDialer.Dial(u.String(), nil)
		if errWebsocket != nil {
			// Log error
			fmt.Printf("application job - websocket connection error: %s\n", errWebsocket.Error())
		}
		defer conn.Close()

		for _, job := range nextJobs {
			// setup the execution time
			job.ProcessedAt = nulls.NewTime(time.Now())
			var app *models.Application
			var appErr error
			switch job.Type {
			case "create":
				app, appErr = job.NewParamApplication()
			case "destroy":
				app, appErr = job.OrigParamApplication()
			default:
				appErr = errors.New("unknown job type")
			}
			if appErr != nil {
				// TODO: log error
				job.Response = nulls.NewString(appErr.Error())
				models.DB.Update(&job)
				if errWebsocket == nil {
					err := conn.WriteMessage(websocket.TextMessage, []byte(job.String()))
					if err != nil {
						fmt.Printf("application job - websocket write error: %s\n", err.Error())
					}
				}
				continue
			}
			hosts := models.Hosts{}
			models.DB.Where("environment_id = ?", app.Environment.ID).All(&hosts)

			response := ""
			for _, host := range hosts {
				response += executeServerCommand(&host, app, job.Type)
			}
			job.Response = nulls.NewString(response)
			models.DB.Update(&job)
			jobMsg := job.String()
			if errWebsocket == nil {
				err := conn.WriteMessage(websocket.TextMessage, []byte(jobMsg))
				if err != nil {
					fmt.Printf("application job - websocket write error: %s\n", err.Error())
				}
			}
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
		commandString = fmt.Sprintf("/usr/local/bin/destroy-project.sh %s %s",
			data.Client.Name,
			data.Project.Name)
	default:
		return "Unknown command"
	}

	ss.Run(commandString)
	return stdoutBuf.String()
}
