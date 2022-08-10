# access-control

A REST API that authorizes access to your services via JWT, with a configuration process that is not an overkill.

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

*Voilà*, now access-control is up and running.

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


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[Apache License 2.0](https://choosealicense.com/licenses/apache-2.0/)