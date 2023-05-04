# GameVN Ultimate Bot

GVN multi-purpose bot. Inspired from: [go_api_boilerplate](https://github.com/yhagio/go_api_boilerplate)

Schema: [Link](https://drawsql.app/teams/clip-sub/diagrams/gvn-ultimate-bot)

# Requirements

- Golang (should be latest stable version, 1.18 as of now)
- A running Postgres DB instance. Some recommendations:
  - A local database using [Postgres.app](https://postgresapp.com/) if you are using Mac OS
  - A managed Heroku Postgres. It has 1GB storage and good enough for single-developer: https://www.heroku.com/postgres

## Up and running

- On Windows, if you are prompted by Windows Firewall, be sure to add `/out` folder into Firewall's whitelist.
- To generate Graphql:

```
go run github.com/99designs/gqlgen
```

- To run the app

```
go run main.go
```

## Role implicit types

0 = Regular role
1 = Nickname color role
2 = "Mobile" role
5 = "Covid" role
