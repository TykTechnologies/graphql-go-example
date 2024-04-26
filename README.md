# graphql-go-example

A sample GraphQL service that implements all CRUD operations.

This service uses [https://github.com/graphql-go/graphql](https://github.com/graphql-go/graphql) as GraphQL implementation.

## Installing

Clone the repository and run the following command:

```shell
go run main.go
```

### Docker

First, you need to build the container: 

```shell
docker build -t graphql-go-example .
```

Then, run the container:
```shell
docker run -p 8080:8080 graphql-go-example:latest
```

## Usage

```
./graphql-go-example -h
Usage of ./graphql-go-example:
  -address string
        address to listen (default "localhost")
  -help
        usage
  -port int
        port to bind (default 8080)
```

## Schema

```graphql
type Query {
  """Get product list"""
  list: [Product]

  """Get product by id"""
  product(id: Int): Product
}

type Mutation {
  """Create new product"""
  create(price: Float!, name: String!, info: String): Product

  """Delete product by id"""
  delete(id: Int!): Product

  """Update product by id"""
  update(info: String, price: Float, id: Int!, name: String): Product
}

type Product {
  id: Int
  info: String
  name: String
  price: Float
}
```

### Get product by id

Query:
```graphql
query ProductById($id: Int){
    product(id: $id) {
        name
        info
        price
    }
}
```

Variables:
```json
{
    "id": 1
}
```

Result:
```json
{
  "data": {
    "product": {
      "info": "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
      "name": "Chicha Morada",
      "price": 7.99
    }
  }
}
```


### Listing all products

Query:

```graphql
query {
    list {
        id
        name
        info
        price
    }
}
```

Result:

```json
{
    "data": {
        "list": [
            {
                "id": 1,
                "info": "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
                "name": "Chicha Morada",
                "price": 7.99
            },
            {
                "id": 2,
                "info": "Chicha de jora is a corn beer chicha prepared by germinating maize, extracting the malt sugars, boiling the wort, and fermenting it in large vessels (traditionally huge earthenware vats) for several days (wiki)",
                "name": "Chicha de jora",
                "price": 5.95
            },
            {
                "id": 3,
                "info": "Pisco is a colorless or yellowish-to-amber colored brandy produced in winemaking regions of Peru and Chile (wiki)",
                "name": "Pisco",
                "price": 9.95
            }
        ]
    }
}
```

### Create a new product

Query: 
```graphql
mutation CreateProduct($name: String!, $info: String, $price: Float!) {
    create(name:$name, info:$info, price:$price){
        id,
        name,
        info,
        price
    }
}
```

Variables:
```json
{
    "name": "Inca Kola",
    "info": "Inca Kola is a soft drink that was created in Peru in 1935 by British immigrant Joseph Robinson Lindley using lemon verbena (wiki)",
    "price": 1.99
}
```

Result:

```json
{
    "data": {
        "create": {
            "id": 73797,
            "info": "Inca Kola is a soft drink that was created in Peru in 1935 by British immigrant Joseph Robinson Lindley using lemon verbena (wiki)",
            "name": "Inca Kola",
            "price": 1.99
        }
    }
}
```

### Updating an existing product with ID

Query:
```graphql
mutation UpdateProduct($id: Int!, $price: Float!) {
    update(id:$id, price:$price){
        id,
        name,
        info,
        price
    }
}
```

Variables:
```json
{
    "id": 1,
    "price": 6.99
}
```

```json
{
    "data": {
        "update": {
            "id": 1,
            "info": "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
            "name": "Chicha Morada",
            "price": 6.99
        }
    }
}
```

### Deleting a product with ID

Query:
```graphql
mutation DeleteProduct($id: Int!) {
    delete(id: $id) {
        id,
        name,
        info,
        price
    }
}
```

Variables:
```json
{
    "id": 1
}
```

Result:
```json
{
    "data": {
        "delete": {
            "id": 1,
            "info": "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
            "name": "Chicha Morada",
            "price": 6.99
        }
    }
}
```

## Debug endpoint

This GraphQL service provides a debug endpoint for inspecting queries. It only accepts `GET` method. 

Make a `GET` request to `/debug/requests`.

Response:

```json
[
    {
        "headers": {
            "Accept": "*/*",
            "Accept-Encoding": "gzip, deflate, br",
            "Connection": "keep-alive",
            "Content-Length": "96",
            "Content-Type": "application/json",
            "Cookie": "csrf_token=mu3IHceOjlJi9BstUN7Wj8zClvljXzWD2DgiP66bfjc=",
            "Postman-Token": "650e8a7d-6c37-4eac-865e-142f2ba4d924",
            "User-Agent": "PostmanRuntime/7.37.3"
        },
        "body": "{\"query\":\"query {\\n    list {\\n        id\\n        name\\n        info\\n        price\\n    }\\n}\"}",
        "date": "Friday, 26-Apr-24 15:45:32 +03"
    }
]
```

## Introspection

This GraphQL service also supports introspection. 

```graphql
{
  __schema {
    types {
      name
    }
  }
}
```

This will return all defined types in the schema.