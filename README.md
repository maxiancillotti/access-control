# access-control

A REST API that authorizes access to your services via JWT, with a simple set up process without roles for now, just users.

## Installation

```bash
git clone https://github.com/maxiancillotti/access-control
```

## Run

### Containerized API & Database

Use [Docker Compose](https://docs.docker.com/compose/) to build and run the API with its own database in dedicated containers for each of them.

```bash
docker compose -f "docker-compose.yaml" up -d --build
```

Only by executing this command on the CLI with the Compose file just as it was provided you can take a look at how the API works.

*Voilà*, now **access-control** is up and running.

#### Consider:

- Config: Aside from testing you will have to set up some parameters in the Compose file. More on this below.

- Storage: The SQL Server image is not so lightweight. You will need 1.65GB of free space. But don't fret, you can connect an already existent instance so you can save up your drive space.

### Containerized API linked to an existent SQL Server instance

Use [Docker Compose](https://docs.docker.com/compose/) to build and run the API.

```bash
docker compose -f "docker-compose_API_ONLY.yaml" up -d --build
```

But before attempting this, you need to set up some parameters in the Compose file.

To create the database, scripts are provided in ./db/sqlserver/db_create_scripts/ to create schema and insert basic data.

## Config

Find all the config variables in the Compose file of your preference.

#### API Container

- Ports: Map a port in the host to a port in the container to run the API's HTTP server and access it from the host. Change the env var HTTPSERVER_HOST_PORT accordingly.

- JWT Secret Keys: 256 bits strings that are used to sign and encrypt the tokens. Do not use the examples on production.

- Timeouts / Duration strings: Time expressions composed by a decimal number and a string suffix that indicates a unit, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

- Database: connection parameters. Instance is optional. When running DB Container, only DB_PASSWORD env var needs to be updated replicating this action in said container config (see Password below).

#### Database Container

- Ports: Map a port in the host to the 1433 port in the container to access the database from the host.

- Password: set in build arg MSSQL_SA_PASSWORD and in the healthcheck test command after the -P flag between the double quotation marks.

- SQL Server Edition: Change build arg MSSQL_SA_PASSWORD if needed. Default: Developer.

- Await time (in seconds) until server bootup: Change build arg SECS_AWAIT_SVR_BOOTUP if needed. Default: 100. The SQL Server needs to be prepared before any connection attemp or it will fail. The database will be created at build time. Set a time that you consider enough in your environment to wait for the server to be initialized before attempting to create the database.

# Usage

### First steps

Set up some data in the system:

- Change the default password for your Admin.
- Create Users to which Permissions will be granted.
- Create Resources to which you want to control access.
- Create Permissions so that your Users can access those Resources.

## Domain Models & Use Cases

### Admins
Are the system administrator's users that will be necessary to authenticate everytime a read or write action is performed on any of the models.

A default admin is created from the outset, which username and password are, respectively: "admin" and "APIUserPassword".
Change this password or create a new admin and disable the former.

#### X-Admin Authorization
Authenticates an admin vía http header "Authorization". Basically equal to Basic Authorization but with a custom "X-Admin" type instead of "Basic" as a prefix in the value of the header. Credentials formating are the same.

#### CREATE

**Request**

```http
POST /api/admins HTTP/1.1
Content-Type: application/json
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk
 
{
  "username": "admindevops"
}
```
Body fields:

- username (string)

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "id":4,
  "username":"admindevops",
  "password":"_eI4PH0 ^4AwZ+_U/ F*;4w\u003cysAf'~Ucs${gS{^4XguO2QE_|Q;:WI%\u003eKB(t}xWq"
}
```
Body fields:
- id (int)
- password (string)
- username (string)

Password is randomly generated.


#### RETRIEVE by username

**Request**

```http
GET /api/admins HTTP/1.1
Content-Type: application/json
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk
 
{
  "username": "admindevops"
}
```
Body fields:

- username (string)

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "id":4,
  "username":"admindevops",
  "password":"_eI4PH0 ^4AwZ+_U/ F*;4w\u003cysAf'~Ucs${gS{^4XguO2QE_|Q;:WI%\u003eKB(t}xWq"
}
```
Body fields:
- id (int)
- password (string)
- username (string)


#### UPDATE Password

**Request**

```http
PATCH /api/admins/{id}/password HTTP/1.1
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk
```
URL parameters:
- id (int)

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "password":"_eI4PH0 ^4AwZ+_U/ F*;4w\u003cysAf'~Ucs${gS{^4XguO2QE_|Q;:WI%\u003eKB(t}xWq"
}
```
Body fields:
- password (string)

Password is randomly generated.


#### UPDATE Enabled State

**Request**

```http
PATCH /api/admins/{id}/enabled-state HTTP/1.1
Content-Type: application/json
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk
 
{
  "enabled_state": true
}
```
URL parameters:

- id (int)

Body fields:

- enabled_state (bool)


**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "success_message": "admin enabled state updated OK"
}
```
Body fields:
- success_message (string)


#### DELETE

**Request**

```http
DELETE /api/admins/{id} HTTP/1.1
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk
```
URL parameters:

- id (int)

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "success_message": "admin deleted OK"
}
```
Body fields:
- success_message (string)


### Users
Are the clients of the APIs you want to control access to.


#### CREATE

**Request**

```http
POST /api/users HTTP/1.1
Content-Type: application/json
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk
 
{
  "username": "apicustomers"
}
```
Body fields:

- username (string)

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "id":4,
  "username":"apicustomers",
  "password":"_eI4PH0 ^4AwZ+_U/ F*;4w\u003cysAf'~Ucs${gS{^4XguO2QE_|Q;:WI%\u003eKB(t}xWq"
}
```
Body fields:
- id (int)
- password (string)
- username (string)

Password is randomly generated.


#### RETRIEVE by username

**Request**

```http
GET /api/users HTTP/1.1
Content-Type: application/json
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk
 
{
  "username": "apicustomers"
}
```
Body fields:

- username (string)

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "id":4,
  "username":"apicustomers",
  "password":"_eI4PH0 ^4AwZ+_U/ F*;4w\u003cysAf'~Ucs${gS{^4XguO2QE_|Q;:WI%\u003eKB(t}xWq"
}
```
Body fields:
- id (int)
- password (string)
- username (string)


#### UPDATE Password

**Request**

```http
PATCH /api/users/{id}/password HTTP/1.1
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk
```
URL parameters:
- id (int)

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "password":"_eI4PH0 ^4AwZ+_U/ F*;4w\u003cysAf'~Ucs${gS{^4XguO2QE_|Q;:WI%\u003eKB(t}xWq"
}
```
Body fields:
- password (string)

Password is randomly generated.


#### UPDATE Enabled State

**Request**

```http
PATCH /api/users/{id}/enabled-state HTTP/1.1
Content-Type: application/json
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk
 
{
  "enabled_state": true
}
```
URL parameters:

- id (int)

Body fields:

- enabled_state (bool)


**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "success_message": "user enabled state updated OK"
}
```
Body fields:
- success_message (string)


#### DELETE

**Request**

```http
DELETE /api/users/{id} HTTP/1.1
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk
```
URL parameters:

- id (int)

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "success_message": "user deleted OK"
}
```
Body fields:
- success_message (string)


#### Resources
The resources to which you need to control access.


#### CREATE

**Request**

```http
POST /api/resources HTTP/1.1
Content-Type: application/json
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk
 
{
  "path": "/customers"
}
```
Body fields:

- path (string)

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "id":2,
  "path": "/customers"
}
```
Body fields:
- id (int)
- path (string)


#### DELETE

**Request**

```http
DELETE /api/resources/{id} HTTP/1.1
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk
```
URL parameters:

- id (int)

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "success_message": "resource deleted OK"
}
```
Body fields:
- success_message (string)


#### RETRIEVE by Path

**Request**

```http
GET /api/resources HTTP/1.1
Content-Type: application/json
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk

{
  "path": "/customers"
}
```

Body fields:
- path (string)

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "id":2,
  "path": "/customers"
}
```
Body fields:
- id (int)
- path (string)


#### RETRIEVE All

**Request**

```http
GET /api/resources HTTP/1.1
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk
```

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "id":2,
  "path": "/customers"
},
{
  "id":3,
  "path": "/products"
}
```
Body fields:
- id (int)
- path (string)

#### HTTP Methods
Find the IDs associated to each of the HTTP Methods so you can use it to grant permissions. 


#### RETRIEVE by Name

**Request**

```http
GET /api/http-methods HTTP/1.1
Content-Type: application/json
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk
 
{
  "name": "POST"
}
```
Body fields:

- name (string)

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "id": 1,
  "name": "POST"
}
```
Body fields:
- id (int)
- name (string)

#### RETRIEVE All

**Request**

```http
GET /api/http-methods HTTP/1.1
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk
```

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "id": 1,
  "name": "POST"
},
{
  "id": 2,
  "name": "GET"
},
...
```
Body fields:
- id (int)
- name (string)

### Users REST Permissions
Having already created Users, Resources and looking up the HTTP Methods IDs, you can use this data to grant permissions for your REST APIs.

#### CREATE

**Request**

```http
POST /api/users-rest-permissions HTTP/1.1
Content-Type: application/json
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk

{
  "user_id": 1,
  "permission": {
    "resource_id": 2,
    "method_id": 3
  }
}
```
Body fields:

- user_id (int)
- resource_id (int)
- method_id (int)

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "user_id": 1,
  "permission": {
    "resource_id": 2,
    "method_id": 3
  }
}
```
Body fields:
- user_id (int)
- permission (obj{resource_id, method_id})
- resource_id (int)
- method_id (int)

#### DELETE

**Request**

```http
DELETE /api/users-rest-permissions HTTP/1.1
Content-Type: application/json
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk

{
  "user_id": 1,
  "permission": {
    "resource_id": 2,
    "method_id": 3
  }
}
```
Body fields:
- user_id (int)
- permission (obj{resource_id, method_id})
- resource_id (int)
- method_id (int)

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "success_message": "permissions deleted OK"
}
```
Body fields:
- success_message (string)

#### RETRIEVE All by UserID

**Request**

```http
GET /api/users/{userID}/rest-permissions HTTP/1.1
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk
```
URL parameters:

- userID (int)

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "user_id": 1,
  "permissions_ids": {
    "resource_id": 2,
    "method_ids": {
      1,
      2,
      3,
    }
  }
}
```
Body fields:
- user_id (int)
- permissions_ids (obj{resource_id, method_ids})
- resource_id (int)
- method_ids ([]int)

#### RETRIEVE All by UserID with Descriptions

**Request**

```http
GET /api/users/{userID}/rest-permissions-with-descriptions HTTP/1.1
Authorization: X-Admin YWRtaW46QVBJVXNlclBhc3N3b3Jk
```
URL parameters:

- userID (int)

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "user_id": 1,
  "permissions": {
    {
      "resource": {
        "id": 1,
        "path": "/customers"
      },
      "methods": {
        {
          "id": 1,
          "name": "POST"
        },
        {
          "id": 2,
          "name": "GET"
        }
      }
    },
    {
      "resource": {
        "id": 2,
        "path": "/products"
      },
      "methods": {
        {
          "id": 1,
          "name": "POST"
        },
        {
          "id": 2,
          "name": "GET"
        }
      }
    }
  }
}
```
Body fields:
- user_id (int)
- permissions ([]obj{resource, []methods})
- resource (obj{id (int), path (string)})
- methods ([]obj{id (int), name (string)})

### Error responses

For any of the former cases.

```http
HTTP/1.1 409 Conflict
Content-Type: application/json
 
{
  "error_message":"error creating admin: username already exists"
}
```

Body fields:
- error_message (string)


## Auth Workflow

### User Authentication

Creates a JSON Web Token (JWT) using Basic Authorization with the API User credentials.

Default duration: 30 min. Update AUTH_JWT_EXPIRATION_TIME_DURATION env var in Compose file if needed.

The User must request the token to **access-control** before attempting a request to a REST API Resource.

**Request**

```http
POST /api/users/token HTTP/1.1
Authorization: Basic dXNlcjpwYXNzd29yZA==
```

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
    "token": "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0..hcXjpXviA8_u8rUG5ZQlgA.sdqsDwz_mxzP7aTYGi3k5kZqbq1eiZ1K290-zbK0lDJVrOUJNdhdixAsKFx6GntEbc8IrilnbRhzol0QuNyPsXpJX14dWdoSHlGfA6MHSxbxvv3vuReSEse7yzFV6T8euDTjqrAveb2NgplA2B_c7mu2X-LWfUrWv1UdhJc8GlHig-SXQVgXsrAoR-D593NzcxdQMNFEqlu-8y_l7R6Lq4WQ6vJVIg6vxmqgNVZejpajHB7mnbt7-h3wyE8VQrqnCOJJI2h1jylq9ilMqyTHBYIy0CQA3058-H_1GhfENDM.IpA-BB2RHIDv4toaOtFWlw"
}
```

Body fields:
- token (string)

### User's Request to a Resource

User must make a request to a Resource using Authorization header, with the "Bearer" type prefix, so it can be authorized.

**Request**

```http
GET /api/customers HTTP/1.1
Authorization: Bearer eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0.._uIIQ_8G5-G3XGGaGGkyWQ.zuZDQB0q7XP4I99QytMcE06Rm9Ei5jZ988nc0E9LmsQgkftwPeFl1ucRplYePOA54k7B7wJQ6sM0qjfR3PUC_DHBRUCkeeoHN9PscN8UH25_P9qBV7LnP4ZXdoXkHamy98vzzXEJcNf3DSuRfP9cvBBy_qHWpO3q6wZ0udxDr6408gw2NfnFPQck1iXnGrEP60tp66c9krMY4Ls5f0Kw304ssvvRlQtCt_RVC7FntHq6szPlBIC6qS3rHbliq_R502aXiwLWLJYjuC-weZKXyrqHfpyJA6a_4zu476JktKU.Cvm8JmcuMKEAsQvjG2omcA
```

### Authorization

The requested API must make a request to **access-control** to authorize the User to access the Resource executing a given HTTP Method.

**Request**

```http
POST /api/users/token/authorize HTTP/1.1
Content-Type: application/json

{
    "token": "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0.._uIIQ_8G5-G3XGGaGGkyWQ.zuZDQB0q7XP4I99QytMcE06Rm9Ei5jZ988nc0E9LmsQgkftwPeFl1ucRplYePOA54k7B7wJQ6sM0qjfR3PUC_DHBRUCkeeoHN9PscN8UH25_P9qBV7LnP4ZXdoXkHamy98vzzXEJcNf3DSuRfP9cvBBy_qHWpO3q6wZ0udxDr6408gw2NfnFPQck1iXnGrEP60tp66c9krMY4Ls5f0Kw304ssvvRlQtCt_RVC7FntHq6szPlBIC6qS3rHbliq_R502aXiwLWLJYjuC-weZKXyrqHfpyJA6a_4zu476JktKU.Cvm8JmcuMKEAsQvjG2omcA",
    "resource_requested": "/customers",
    "method_requested": "GET"
}
```

Body fields:
- token (string)
- resource_requested (string)
- method_requested (string)

**Successful response**

```http
HTTP/1.1 200 OK
Content-Type: application/json
 
{
  "success_message": "Token authorization OK"
}
```
Body fields:
- success_message (string)


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[Apache License 2.0](https://choosealicense.com/licenses/apache-2.0/)