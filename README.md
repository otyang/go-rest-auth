# Auth-projects

---

## Powered by:

- `Fiber` - is an Express inspired web framework built on top of Fasthttp https://github.com/gofiber/fiber
- `Bun` - Simple and performant DB client for PostgreSQL, MySQL, and SQLite https://github.com/uptrace/bun

---

## Repository Structure:

A top-level directory layout:

```
.
├── cmd
├── config              # Configuration files
├── src                 # Source files
└── README.md
```

Based on Bob Martin’s clean architecture, we have the typical directory layout of `src` folder:

```
├── ...
├── src                 
│   ├── domain          # Entities 
│   ├── infrastructure  # Frameworks and Drivers (External Interfaces)
│   ├── interface       # Controllers Presenters (Interface Adapters)
│   ├── usecase         # Use Cases
│   └── ...             # etc.
└── ...
```

#### domain layer — Entities:

Entities is a domain model (layer) that has wide enterprise business rules and can be a set of data structures and
functions.

#### usecase layer — Use Cases:

In Use cases layer we have three directories: repository, presenter and interactor. Interactor is in charge of Input
Port and presenter is in charge of Output Port. Interactor has a set of methods of specific application business rules
depending on repository and presenter interface.

#### interface layer — Controllers Presenters:

In Interface Adapter layer, there are controller, presenter and repository folders. Controllers is in charge of the C of
MVC model and handles API requests that come from the outer layer. Repository is a specific implementation of repository
in Use Cases and stores any database handler as Gateway.

#### infrasctructure layer — External Interfaces:

We have datastore and router in infrastructure. Datastore is used as creating a database instance, in this case, we use
bun with PostgreSQL. Router is defined as a routing request using Fiber.

---

## Local deployment:

Clone the project to the directory. Create a folder rsa_keys and add private_key.pem and public_key.pem there. Create a config.yml file, in the conf folder, using config.yml.example.

#### Step by step creation of Postgres database inside Docker container:

Pull the official image of the Postgres database:

```
docker pull postgres
```

Run the Docker container:

```
docker run --name auth_project -e POSTGRES_PASSWORD=12345 -p 5436:5432 -d postgres
```

Get a bash shell in the container:

```
docker exec -it auth_project bash
```

Run psql — PostgreSQL interactive terminal:

```
psql -U postgres
```

Create a database:

```
CREATE DATABASE auth;
```

Check if it exists:

```
\l
```

Exit psql:

```
\q
```

Exit bash shell:

```
exit
```

Detailed:

https://hub.docker.com/_/postgres

https://medium.com/better-programming/connect-from-local-machine-to-postgresql-docker-container-f785f00461a7

#### Connecting to local Postgres database inside Docker container:

`POSTGRES_USER = postgres`

`POSTGRES_PASSWORD = 12345`

`POSTGRES_HOST = localhost`

`POSTGRES_PORT = 5436`

`POSTGRES_DATABASE = auth`

`POSTGRES_URL = postgres://postgres:12345@localhost:5436/auth_project?sslmode=disable`

#### Install `golang-migrate` library - database migrations and CLI:

```
brew install golang-migrate
```

Up migration 

```
migrate -database "postgres://postgres:12345@localhost:5436/auth?sslmode=disable"  -path src/infrastructure/storage/postgres/migrations up
```