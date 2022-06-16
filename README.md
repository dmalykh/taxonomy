# UNDER DEVELOPMENT

# TagService
TagService is full microservice for tags management. 
You could use it either independent microservice or embed library.

## Overview
Connect every object with tags. Each object relate with tag via _namespace_ and _entity_id_.
- Namespace 

## Features
- Command line interface
- Database agnostic, uses [ent](https://entgo.io/) inside
- GraphQl server
- GRPC server and client

## Cli interface
```
serve run microservice
tag list
category list
        create name title
namespace list
        create name
```

Category 
 - name and parent id should be unique pair

Tag
 - category required for tag. 
 - name and category id should be unique pair

Maybe you think "How could I create tag without category". 
I'll try to explain why I made decision to make categories required for every tag.
What is a tag without category? It's a record in a database where category_id equals NULL.
And when all tags without category will be requested, query with category_id equeals NULL will executed. In this case NULL is id of category. And it's stange

RENAME server dir to tagservice

## TODO
- [ ] Getting started
- [ ] Documentation
- [ ] Publish API specification
- [ ] Nested namespaces
- [ ] Default namespace 
- [ ] Main category