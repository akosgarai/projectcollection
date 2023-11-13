package grifts

import (
	"projectcollection/models"

	"github.com/gobuffalo/grift/grift"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		// insert the dbtypes: no, mysql, postgres
		dbTypeNames := []string{"no", "mysql", "postgres"}
		for _, name := range dbTypeNames {
			dbType := models.Dbtype{}
			// if the dbtype doesn't exist, create it
			models.DB.Where("name = ?", name).First(&dbType)
			if dbType.Name == "" {
				dbType = models.Dbtype{
					Name: name,
					ID:   uuid.Must(uuid.NewV4()),
				}
				models.DB.Create(&dbType)
			}
		}
		// insert the runtime: noPHP, PHP71FPM, PHP74FPM, PHP81FPM
		runtimeNames := []string{"noPHP", "PHP71FPM", "PHP74FPM", "PHP81FPM"}
		for _, name := range runtimeNames {
			runtime := models.Runtime{}
			// if the runtime doesn't exist, create it
			models.DB.Where("name = ?", name).First(&runtime)
			if runtime.Name == "" {
				runtime = models.Runtime{
					Name: name,
					ID:   uuid.Must(uuid.NewV4()),
				}
				models.DB.Create(&runtime)
			}
		}
		// insert the environments
		stagingEnv := models.Environment{}
		productionEnv := models.Environment{}
		models.DB.Where("name = ?", "staging").First(&stagingEnv)
		models.DB.Where("name = ?", "production").First(&productionEnv)
		if stagingEnv.Name == "" {
			stagingEnv = models.Environment{
				Name: "staging",
				ID:   uuid.Must(uuid.NewV4()),
			}
			models.DB.Create(&stagingEnv)
		}
		if productionEnv.Name == "" {
			productionEnv = models.Environment{
				Name: "production",
				ID:   uuid.Must(uuid.NewV4()),
			}
			models.DB.Create(&productionEnv)
		}
		// insert the hosts (ip: staging or production)
		stagingHost := models.Host{}
		productionHost := models.Host{}
		models.DB.Where("ip = ?", "staging").First(&stagingHost)
		models.DB.Where("ip = ?", "production").First(&productionHost)
		if stagingHost.IP == "" {
			stagingHost = models.Host{
				IP:            "staging",
				EnvironmentID: stagingEnv.ID,
				ID:            uuid.Must(uuid.NewV4()),
				Name:          "server-staging",
				SSHUser:       "scriptexecutor",
				SSHPort:       2222,
				SSHKey:        "/root/.ssh/id_rsa_shared",
			}
			models.DB.Create(&stagingHost)
		}
		if productionHost.IP == "" {
			productionHost = models.Host{
				IP:            "production",
				EnvironmentID: productionEnv.ID,
				ID:            uuid.Must(uuid.NewV4()),
				Name:          "server-production",
				SSHUser:       "scriptexecutor",
				SSHPort:       2222,
				SSHKey:        "/root/.ssh/id_rsa_shared",
			}
			models.DB.Create(&productionHost)
		}
		// create one test client
		testClient := models.Client{}
		models.DB.Where("name = ?", "testclient").First(&testClient)
		if testClient.Name == "" {
			testClient = models.Client{
				Name: "testclient",
				ID:   uuid.Must(uuid.NewV4()),
			}
			models.DB.Create(&testClient)
		}
		// create one test project
		testProject := models.Project{}
		models.DB.Where("name = ?", "testproject").First(&testProject)
		if testProject.Name == "" {
			testProject = models.Project{
				Name: "testproject",
				ID:   uuid.Must(uuid.NewV4()),
			}
			models.DB.Create(&testProject)
		}

		// create the admin user
		adminUser := models.User{}
		models.DB.Where("email = ?", "admin@email.com").First(&adminUser)
		ph, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if adminUser.Email == "" {
			adminUser = models.User{
				Email:        "admin@email.com",
				ID:           uuid.Must(uuid.NewV4()),
				PasswordHash: string(ph),
			}
			models.DB.Create(&adminUser)
		}
		// create the developer user
		developerUser := models.User{}
		models.DB.Where("email = ?", "developer@email.com").First(&developerUser)
		ph, _ = bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if developerUser.Email == "" {
			developerUser = models.User{
				Email:        "developer@email.com",
				ID:           uuid.Must(uuid.NewV4()),
				PasswordHash: string(ph),
			}
			models.DB.Create(&developerUser)
		}
		// create the user user
		userUser := models.User{}
		models.DB.Where("email = ?", "user@email.com").First(&userUser)
		ph, _ = bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if userUser.Email == "" {
			userUser = models.User{
				Email:        "user@email.com",
				ID:           uuid.Must(uuid.NewV4()),
				PasswordHash: string(ph),
			}
			models.DB.Create(&userUser)
		}
		// create the roles
		roleNames := map[string]string{"sysadmin": "allow everything", "developer": "allow to see the list / view pages", "user": "allow to see the application list / view pages"}
		roles := map[string]models.Role{}
		for name, description := range roleNames {
			role := models.Role{}
			models.DB.Where("name = ?", name).First(&role)
			if role.Name == "" {
				role = models.Role{
					Name:        name,
					Description: description,
					ID:          uuid.Must(uuid.NewV4()),
				}
				models.DB.Create(&role)
			}
			roles[name] = role
		}
		// create the resources
		resourceNames := []string{"dbtype", "runtime", "environment", "host", "client", "project", "application", "alias"}
		actionNames := []string{"view", "create", "edit", "delete"}
		for _, resourceName := range resourceNames {
			for _, actionName := range actionNames {
				resource := models.Resource{}
				models.DB.Where("name = ?", resourceName+"."+actionName).First(&resource)
				if resource.Name == "" {
					resource = models.Resource{
						Name: resourceName + "." + actionName,
						ID:   uuid.Must(uuid.NewV4()),
					}
					models.DB.Create(&resource)
				}
			}
		}
		// assign the resources to the roles
		roleResources := map[string][]string{
			"sysadmin":  {"dbtype.view", "dbtype.create", "dbtype.edit", "dbtype.delete", "runtime.view", "runtime.create", "runtime.edit", "runtime.delete", "environment.view", "environment.create", "environment.edit", "environment.delete", "host.view", "host.create", "host.edit", "host.delete", "client.view", "client.create", "client.edit", "client.delete", "project.view", "project.create", "project.edit", "project.delete", "application.view", "application.create", "application.edit", "application.delete", "alias.view", "alias.create", "alias.edit", "alias.delete"},
			"developer": {"dbtype.view", "runtime.view", "environment.view", "host.view", "client.view", "project.view", "application.view", "alias.view"},
			"user":      {"application.view", "alias.view"},
		}
		for roleName, resourceNames := range roleResources {
			role := models.Role{}
			models.DB.Where("name = ?", roleName).First(&role)
			for _, resourceName := range resourceNames {
				resource := models.Resource{}
				models.DB.Where("name = ?", resourceName).First(&resource)
				roleToResource := models.RoleToResource{}
				models.DB.Where("role_id = ? AND resource_id = ?", role.ID, resource.ID).First(&roleToResource)
				if roleToResource.RoleID == uuid.Nil {
					roleToResource = models.RoleToResource{
						RoleID:     role.ID,
						ResourceID: resource.ID,
						ID:         uuid.Must(uuid.NewV4()),
					}
					models.DB.Create(&roleToResource)
				}
			}
		}
		// assign the roles to the users
		userRoleMap := map[uuid.UUID]string{adminUser.ID: "sysadmin", developerUser.ID: "developer", userUser.ID: "user"}
		for userID, roleName := range userRoleMap {
			userToRole := models.UserToRole{}
			models.DB.Where("user_id = ?", userID).First(&userToRole)
			if userToRole.UserID == uuid.Nil {
				userToRole = models.UserToRole{
					UserID: userID,
					RoleID: roles[roleName].ID,
					ID:     uuid.Must(uuid.NewV4()),
				}
				models.DB.Create(&userToRole)
			}
		}
		return nil
	})

})
