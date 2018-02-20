# Application demonstrates my go skills #

## Installation ##

- Install [go](https://golang.org/doc/install) (I used 1.9.2)

- Install [mongoDB](https://docs.mongodb.com/manual/installation/) (I used 3.6.1)

- Clone repo in you own directory in $GOPATH/src/venues (directory's name is important) for example:

    `cd $GOPATH/src && git clone REPO_URL venues`

- Install dependencies6, I use dep, so I ran:

    `go get -u github.com/golang/dep/cmd/dep && dep ensure`

- Configure application via '.env' file, you have template as .env.template, for example:

    `touch .env && echo 'PORT=8000\nMONGO_ADDRESS=127.0.0.1\nMONGO_DB_NAME=venues\nMONGO_DB_NAME_TEST=venues\n' > .env`

- Run mongod, for example:

    `mongod --fork --syslog`

- Run tests:

    `make test`

- Run app!:

    `go run main.go`

## Usage ##

- Create new restaurant:

    `curl -X POST -H "Content-Type: application/json" -d '{"name": "Dish1"}' 'localhost:8000/restaurants/5a8ad983591b381c73797521/dish'`

- Get all restaurants:

    `curl -X GET -H "Content-Type: application/json" 'localhost:8000/restaurants'`

    * for ordering add `ordering` query parameter as `-rating` or `rating`

    * for filtering by city add `city` query params

    * for select page by page add `page` param
