# Taxonomy Service
![Coverage](https://img.shields.io/badge/Coverage-40.9%25-yellow)
Taxonomy is a self-hosted, lightweight and simple (yet functional) full microservice for terms management.
You could use it either independent microservice or embed library.

_**Termservice is still in development, but it already has stable API and used in production. See TODO section below to learn about planned improvements.**_

### Purpose
> Terms everywhere! I really exhausted to implement terms management services in every project. The core idea of `termservice`
is wide usage universal solution for every project where terms needed.
> 
> `Termservice` could be embeded, has GraphQL and GRPC APIs, also provided powerfull command line interface. 
Of course, it's database agnostic and made with clean architecture principles. 
If you want to modify every layer you should implement its interface.

## Features
- Command line interface
- Database agnostic, uses [ent](https://entgo.io/) inside
- GraphQL server

## Overview
Connect every object with terms. Each object relate with term via _namespace_ and _entity_id_.
- Namespace 

### Restrictions
- Vocabulary name and parent id should be unique pair
- vocabulary required for every term
- name and vocabulary id should be unique pair

> Maybe you think how to create term without vocabulary.
I'll try to explain why I made decision to make categories required for every term.
What is a term without vocabulary? It's a record in a database where vocabulary_id equals NULL.
And when all terms without vocabulary will be requested, query with vocabulary_id equeals NULL will executed. 
In this case NULL is id of vocabulary. And it's strange, because NULL isn't convinient value for identifier.


## Cli interface
```shell
Usage:
  internal [command]

Available Commands:
  vocabulary    Operations with categories
  help        Help about any command
  init        Initiate service, create tables in a database
  namespace   CRUD operations with namespaces
  rel         Work with references
  serve       Run API server
  term         CRUD operations with terms

Flags:
      --dsn string   Data source name (connection information) (default "sqlite://./termservice.db?cache=shared&_fk=1")
  -h, --help         help for internal
  -v, --verbose      Make some output more verbose.

Use "termservice [command] --help" for more information about a command.
```

## Run GraphQL API in Docker
Make Dockerfile
```dockerfile
FROM dmalykh/termservice:latest
EXPOSE 8080
CMD ["termservice", "serve", "graphql", "-p", "8080"]
```
And run it!
```shell
docker build -t internal .
docker run -it --rm -p 8081:8080 internal   
```
Open http://127.0.0.1:8081/ to get acquainted with GraphiQL!


## TODO
- [ ] Getting started
- [ ] GraphQL API tests
- [ ] GRPC API
- [ ] Documentation
- [ ] Publish API specification
- [ ] Nested namespaces
- [ ] Make default namespace on init
- [ ] Make default vocabulary on init
- [ ] docker-compose
- [ ] Add telemetry and metrics
- [ ] https://github.com/rivo/tview
