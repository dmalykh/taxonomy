# TagService
TagService is a self-hosted, lightweight and simple (yet functional) full microservice for tags management.
You could use it either independent microservice or embed library.

_**Tagservice is still in development, but it already has stable API and used in production. See TODO section below to learn about planned improvements.**_

### Purpose
> Tags everywhere! I really exhausted to implement tags management services in every project. The core idea of `tagservice`
is wide usage universal solution for every project where tags needed.
> 
> `Tagservice` could be embeded, has GraphQL and GRPC APIs, also provided powerfull command line interface. 
Of course, it's database agnostic and made with clean architecture principles. 
If you want to modify every layer you should implement its interface.

## Features
- Command line interface
- Database agnostic, uses [ent](https://entgo.io/) inside
- GraphQL server

## Overview
Connect every object with tags. Each object relate with tag via _namespace_ and _entity_id_.
- Namespace 

### Restrictions
- Category name and parent id should be unique pair
- category required for every tag
- name and category id should be unique pair

> Maybe you think how to create tag without category.
I'll try to explain why I made decision to make categories required for every tag.
What is a tag without category? It's a record in a database where category_id equals NULL.
And when all tags without category will be requested, query with category_id equeals NULL will executed. 
In this case NULL is id of category. And it's strange, because NULL isn't convinient value for identifier.


## Cli interface
```shell
Usage:
  tagservice [command]

Available Commands:
  category    Operations with categories
  help        Help about any command
  init        Initiate service, create tables in a database
  namespace   CRUD operations with namespaces
  rel         Work with relations
  serve       Run API server
  tag         CRUD operations with tags

Flags:
      --dsn string   Data source name (connection information) (default "sqlite://./tagservice.db?cache=shared&_fk=1")
  -h, --help         help for tagservice
  -v, --verbose      Make some output more verbose.

Use "tagservice [command] --help" for more information about a command.
```

## Run GraphQL API in Docker
Make Dockerfile
```dockerfile
FROM dmalykh/tagservice:latest
EXPOSE 8080
CMD ["tagservice", "serve", "graphql", "-p", "8080"]
```
And run it!
```shell
docker build -t tagservice .
docker run -it --rm -p 8081:8080 tagservice   
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
- [ ] Make default category on init
- [ ] docker-compose
