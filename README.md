# Project register app powered by Buffalo

## Prerequisites

- Generate an ssh private and public key pair under the .ssh directory.

```bash
ssh-keygen -t rsa -b 4096 -C "email@address.com"
```

- Create the authorized_keys file under the .ssh directory and add the public key as content.

```bash
cat .ssh/id_rsa.pub > .ssh/authorized_keys
```

## Run the application

```bash
docker compose up -d --build
```

If you point your browser to [http://127.0.0.1:3000](http://127.0.0.1:3000) you should see a "Welcome to Buffalo!" page.

## Basic setup

### Runtimes

Runtime [list page](http://127.0.0.1:3000/runtimes/). Click to "Create New Runtime" button. Example values: PHP71FPM, PHP74FPM, PHP81FPM, noPHP

### Dbtypes

Dbytpe [list page](http://127.0.0.1:3000/dbtypes/). Click to "Create New Dbtype" button. Example values: no, mysql

### Environments

Environment [list page](http://127.0.0.1:3000/environments/). Click to "Create New Environment" button. Example values: staging, production

### Clients

Client [list page](http://127.0.0.1:3000/clients/). Click to "Create New Client" button. Example values: test-client

### Projects

Project [list page](http://127.0.0.1:3000/projects/). Click to "Create New Project" button. Example values: test-project

### Hosts

Host [list page](http://127.0.0.1:3000/hosts/). Click to "Create New Host" button.

Example values for a staging environment:
- Name: server-staging
- IP: staging (the name of the staging environment in the docker system)
- Environment: staging (selectable from the previously inserted environments)
- SSH User: scriptexecutor (the name of the user in the staging environment container)
- SSH Port: 2222 (ssh port on the docker environment containers)
- SSH Key: /root/.ssh/id_rsa (the name of the private key file. It is attached under the /root/.ssh directory in the application container.)

Example values for a production environment:
- Name: server-production
- IP: production (the name of the production environment in the docker system)
- Environment: production (selectable from the previously inserted environments)
- SSH User: scriptexecutor (the name of the user in the production environment container)
- SSH Port: 2222 (ssh port on the docker environment containers)
- SSH Key: /root/.ssh/id_rsa (the name of the private key file. It is attached under the /root/.ssh directory in the application container.)

### Application

Application [list page](http://127.0.0.1:3000/applications/). Click to "Create New Application" button.

Example values for a staging application:
- Owner email: email@address.com (It does not send email yet, this value is only stored in the applications table)
- Project: test-project (The first project option is selected by default)
- Client: test-client (The first client option is selected by default)
- Runtime: PHP71FPM (The first runtime option is selected by default)
- Database: no (The first dbtype option is selected by default)
- Environment: staging (The first environment option is selected by default)

## Application initialization

Once an application has been created an entry will be created to the job_application table. It is processed by a task that is executed with cron.

job_application table:
- new_params (json of the expected values of the application)
- orig_params (json of the original values of the application - on case of edit)
- processed_at (default null)
- response

## Generator commands

```console
buffalo generate resource host name ip environment_id:uuid ssh_user ssh_port:int ssh_key
```

## Users, Roles, Auth

The basics of the users, roles and authentication has been implemented based on the [documentation](https://gobuffalo.io/documentation/guides/auth/) and this [post](https://zoharpeled.wordpress.com/2020/11/07/role-based-security-database-design/).

