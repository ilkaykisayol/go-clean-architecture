## Go-Clean-Architecture

### Docker setup
#### Build Image
Here is an example:
```bash
$ docker build --build-arg GITHUB_USERNAME=dummy_username --build-arg GITHUB_TOKEN=ghp_DuMmyToKen -t go-clean-architecture:latest .
```

#### Run Container
Here is an example:  
```bash
$ docker run --rm -d -p 8080:8080 --name go-clean-architecture-app go-clean-architecture:latest
```

#### Docker Compose
You may also use docker compose. Create your own docker-compose.yaml file from docker-compose-example.yaml and put your
sensitive & personalized data based on your machine environment.

`$ docker-compose up`  
`$ docker-compose up --build` force build the image before running.

## Testing
### Generating Mocks
[GoMock](https://github.com/golang/mock) is used for mock generation. Please see repository page for installation and detailed instructions.

After generating mocks, you can use them to set specific behaviours based on predefined conditions. Please see [here](https://github.com/golang/mock#building-mocks) for instructions and examples.

[GoMock GitHub repository](https://github.com/golang/mock)

### Rules
Create a new file with a postfix like `_test.go` to test a go code inside any go file. Test files are siblings of actual go files. They should be inside the same folder near them.
#### Example
`sample_service.go` actual go file.  
`sample_service_test.go` test file.

### Run Tests via CLI
Execute this command on the root directory of your project. Command will run tests only inside the `internal` directory and its subdirectories.

`$ go test -v ./internal/...` runs tests and logs verbose messages  
`$ go test -v -cover ./internal/...` runs tests, shows coverage and logs verbose messages.

## Packages
### Gin Web Framework
Gin is used for REST communication, please refer to Gin Web Framework docs for further information.

[Gin Github repository](https://github.com/gin-gonic/gin)  
[Gin Framework web page](https://gin-gonic.com/)

### Swagger
[Swaggo GitHub repository](https://github.com/swaggo/swag)  
[Swaggo Gin Middlewware GitHub repository](https://github.com/swaggo/gin-swagger)

Swaggo generates swagger files into ./docs folder which is ignored to not track changes. If examples are not enough for
any given case, please refer
to [declarative comments format](https://swaggo.github.io/swaggo.io/declarative_comments_format/) documentation when
required.

`$ swag fmt` formats the swagger related comments to make it more readable.  
`$ swag init` generates files into ./docs folder.

### Db Migration
golang-migrate is used for migrate database of Go-Clean-Architecture. First of all golang-migrate cli must be installed. Follow [these steps](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) for install cli. After that use below command to start migrate;

```bash
$ migrate -source file://migration -database "postgres://username:password@localhost:5433/go-clean-architecture?sslmode=disable" up
```

### Input Validation
Go playground validator is used for input validation. Also Gin uses this package for validations.

[Validator GitHub repository](https://github.com/go-playground/validator)  
[Validator godoc page](https://pkg.go.dev/github.com/go-playground/validator/v10)

### Redis
go-redis package is used for distributed cache/redis communication.

[go-redis GitHub repository](https://github.com/go-redis/redis)  
[go-redis godoc page](https://pkg.go.dev/github.com/go-redis/redis/v8)  
[go-redis web page](https://redis.uptrace.dev/)

### PostgreSQL Driver
[Driver GitHub repository](https://github.com/lib/pq)  
[Driver godoc page](https://pkg.go.dev/github.com/lib/pq)

### Environment files
GoDotEnv is used for environment files.

[GoDotEnv GitHub repository](https://github.com/joho/godotenv)  
[GoDotEnv godoc page](https://pkg.go.dev/github.com/joho/godotenv?utm_source=godoc)

.env file holds sample keys for required environment variables. You need to create your own local .env.{Environment}
files to override with sensitive data.

Example local .env file names:  
`.env.Development`  
`.env.Staging`  
`.env.Production`

These files are ignored to not track changes.

### Structured Logging
Uber's zap package is used.

[zap GitHub repository](https://github.com/uber-go/zap)  
[zap godoc page](https://pkg.go.dev/go.uber.org/zap)