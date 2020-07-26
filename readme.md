# Steps to run the backend locally

- Run a mysql server
    docker run --name db -d -p3306:3306 -e MYSQL_ROOT_PASSWORD=password123 mysql:5.7
- Install SequelPro
    https://www.sequelpro.com/
- Create a db 
- Clone the repo - `git clone git@github.com:dootoday/dootoday-go.git apidootoday`
- copy the `config/localConfig_sample.yaml` to `config/localConfig.yaml`
- Run with `go run main.go`


# Run server
```
ENVIRONMENT=production \
DB_DRIVER=mysql \
DB_NAME=dootoday \
DB_PORT=3306 \
DB_HOSTNAME=dootoday.ci5opzkdzsli.ap-south-1.rds.amazonaws.com \
DB_USERNAME=root \
DB_PASSWORD=hTHLkRG5Due9Fms6KBF3 \
SERVER_PORT=9060 \
SERVER_HOSTNAME=api.doo.today \
DEBUG=true \
ACCESS_TOKEN_SECRET=access-this-is-a-very-secret-string-1 \
REFRESH_TOKEN_SECRET=refresh-this-is-a-very-secret-string-1 \
RP_API_KEY='rzp_test_oW4N8eXjSQAzY8' \
RP_API_SECRET='ds1XkOEFtZQ8BsdAkPN4Nh5n' \
FRONT_END_BASE='https://doo.today' \
BACK_END_BASE='https://api.doo.today' \
DOO_TODAY_LOGO='https://dootoday-assets.s3.ap-south-1.amazonaws.com/logo-200x200.png' \
DOO_TODAY_NAME='DooToday', \
DOO_TODAY_DESC='Daily task simplified' \
REDIS_HOST='ec2-15-207-85-133.ap-south-1.compute.amazonaws.com' \
REDIS_PORT=6379 \
REDIS_PASSWORD=mFP4NRm6ks5pG99KhGbeUCrM4AgH9S9mj8U9vH69XkKG7MGLZu6jFj82h3WB \
./main
```