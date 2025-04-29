# GraphQL to Java (gql2j)

A simple (and dummy) Go application that converts GraphQL schemas to Java classes.

## Usage

```bash
go run main.go -schema=/path/to/schema.graphql -output=./output -package=com.example.model
```

### Parameters

- `-schema`: (Required) Path to the GraphQL schema file
- `-output`: (Optional) Output directory for generated Java classes (default: "output")
- `-package`: (Optional) Java package name for generated classes (default: "com.example.model")

## Example

For a GraphQL schema:

```graphql
type User {
  id: ID!
  name: String
  email: String!
  posts: [Post]
}

type Post {
  id: ID!
  title: String!
  content: String
  author: User!
}

enum UserRole {
  ADMIN
  USER
  GUEST
}
```

The tool will generate corresponding Java classes in the specified output directory.