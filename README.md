# GraphQL to Java (gql2j)

A Go tool that generates Java classes from GraphQL schemas with support for Lombok, JSR-303 validation, custom type mappings, and directives.

## Installation

```bash
go install github.com/source-c/go-gql2j/cmd/gql2j@latest
```

Or build from source:

```bash
git clone https://github.com/source-c/go-gql2j.git
cd go-gql2j
go build -o gql2j ./cmd/gql2j
```

## Quick Start

```bash
# Basic usage
gql2j -schema schema.graphql -output ./generated -package com.example.model

# With config file
gql2j -config gql2j.yaml

# With Lombok and validation
gql2j -schema schema.graphql -output ./generated -package com.example.model -lombok -validation
```

## CLI Flags

| Flag | Description |
|------|-------------|
| `-config` | Path to YAML config file |
| `-schema` | GraphQL schema path (overrides config) |
| `-output` | Output directory (overrides config) |
| `-package` | Java package name (overrides config) |
| `-java-version` | Target Java version: 8, 11, 17, 21 |
| `-lombok` | Enable Lombok annotations |
| `-lombok-disable` | Disable Lombok annotations |
| `-validation` | Enable JSR-303 validation |
| `-validation-disable` | Disable JSR-303 validation |
| `-validation-package` | Validation package: `jakarta` or `javax` |
| `-clean` | Clean output directory before generating |
| `-verbose` | Enable verbose output |
| `-version` | Print version information |

## Configuration File

Create a `gql2j.yaml` file (see `gql2j.yaml.example` for full options):

```yaml
schema:
  path: "./schema.graphql"
  includes:
    - "./types/*.graphql"

output:
  directory: "./generated"
  package: "com.example.model"

java:
  version: 17
  fieldVisibility: "private"
  collectionType: "List"
  nullableHandling: "wrapper"
  naming:
    fieldCase: "camelCase"
    classSuffix: ""
    interfacePrefix: ""

typeMappings:
  scalars:
    DateTime:
      javaType: "java.time.LocalDateTime"
      imports: ["java.time.LocalDateTime"]
    UUID:
      javaType: "java.util.UUID"
      imports: ["java.util.UUID"]

features:
  lombok:
    enabled: true
    data: true
    builder: true
    noArgsConstructor: true
  validation:
    enabled: true
    package: "jakarta"
    notNullOnNonNull: true
```

## Supported Directives

| Directive | Target | Effect |
|-----------|--------|--------|
| `@skip` | Type, Field | Exclude from generation |
| `@javaName(name: "...")` | Type, Field, Enum Value | Override Java name |
| `@javaType(type: "...", imports: [...])` | Field | Custom Java type |
| `@deprecated(reason: "...")` | Field, Enum Value | Add `@Deprecated` |
| `@annotation(value: "...", imports: [...])` | Type, Field | Add custom annotation |
| `@constraint(...)` | Field | JSR-303 validation |
| `@lombok(exclude: [...], include: [...])` | Type | Per-type Lombok config |
| `@collection(type: "Set")` | Field | Override collection type |

### Directive Examples

```graphql
type User @lombok(exclude: ["builder"]) {
  id: ID! @javaType(type: "java.util.UUID", imports: ["java.util.UUID"])
  email: String! @constraint(pattern: "^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$")
  tags: [String] @collection(type: "Set")
  legacyField: String @deprecated(reason: "Use newField instead")
}

type Entity @annotation(value: "@Entity", imports: ["jakarta.persistence.Entity"]) {
  id: ID!
}

type Internal @skip {
  secret: String
}
```

### Constraint Directive Options

```graphql
field: String @constraint(
  minLength: 1,
  maxLength: 100,
  min: 0,
  max: 150,
  pattern: "^[a-z]+$",
  notNull: true,
  notBlank: true,
  email: true
)
```

## Example

For a GraphQL schema:

```graphql
type User {
  id: ID!
  name: String
  email: String!
  posts: [Post]
  role: UserRole
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

interface Node {
  id: ID!
}

input CreatePostInput {
  title: String!
  content: String!
}
```

Generated `User.java` (with Lombok and validation enabled):

```java
package com.example.model;

import jakarta.validation.constraints.NotNull;
import java.util.List;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
public class User {

    @NotNull
    private String id;

    private String name;

    @NotNull
    private String email;

    private List<Post> posts;

    private UserRole role;
}
```

## Library Usage

Use gql2j as a library in your Go code:

```go
package main

import (
    "fmt"
    "github.com/source-c/go-gql2j/pkg/api"
)

func main() {
    result, err := api.GenerateToDir(api.Options{
        SchemaPath:       "schema.graphql",
        OutputDir:        "./generated",
        Package:          "com.example.model",
        JavaVersion:      17,
        EnableLombok:     true,
        EnableValidation: true,
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Generated %d files\n", len(result.Files))
}
```

## Nullable Handling

Configure how nullable GraphQL fields are represented in Java:

| Mode | Description | Example |
|------|-------------|---------|
| `wrapper` | Use wrapper types (default) | `Integer`, `Boolean` |
| `optional` | Use `Optional<T>` | `Optional<Integer>` |
| `annotation` | Use `@Nullable` annotation | `@Nullable Integer` |

## License

MIT
