# Steps to run the backend locally

- Run a mysql server
    docker run --name db -d -p3306:3306 -e MYSQL_ROOT_PASSWORD=password123 mysql:5.7
- Install SequelPro
    https://www.sequelpro.com/
- Create a db 
- Clone the repo - `git clone git@github.com:dootoday/dootoday-go.git apidootoday`
- copy the `config/localConfig_sample.yaml` to `config/localConfig.yaml`
- Run with `go run main.go`