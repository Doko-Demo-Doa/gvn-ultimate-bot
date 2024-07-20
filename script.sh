up() {
  docker-compose up
}

down() {
  docker-compose down
}

gql() {
  go run github.com/99designs/gqlgen -v
}

swag() {
  swag init -g app/app.go
}