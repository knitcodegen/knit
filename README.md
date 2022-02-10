# Knit

Knit is an inline code generation tool that combines the power of Go's text/template package with automatic spec file loading.

## Example


```yml
openapi: "3.0.0"
info:
  version: 1.0.0
  title: Example Codegen
paths:
  /pets:
    post:
      operationId: GetPet
```

example.go:
```go
package example

// @knit input $SCHEMA_FILE
// @knit loader openapi
// @knit template ./templates/openapi/test.tmpl

// @!knit
```

Running `knit` against this file will first load the `schema` file using the openapi `loader`, then execute the `template` using the data from the `loader`. The resulting text will then be inserted between the `@knit` annotations. The annotations are maintained so the code can be generated again when the schema changes.

Here is the resulting code:
```go
package example

// @knit input $SCHEMA_FILE
// @knit loader openapi
// @knit template ./templates/openapi/test.tmpl

type Generated struct {
    Pet string
}

// @!knit
```

## Annotations
Define generators you'd like to knit into your codebase using the `@knit` annotations.

The parser algorithm first splits the target code file by the ending annotation: `@!knit` -- It then works backwards on each split block to identify the options for the generator.

Options are defined on the opening annotations in the following format:
```
@knit <option> <value>
```

Annotations also support environment variables. All `$env` variables will be expanded before the option line is parsed. 

## Options
### Input
The `input` option specifies the input file. This must be a file on your system. 

Relative paths are resolved using the directory in which you've run `knit`

### Loader
The `loader` option specifies the loader used to load the input file. 

Currently the only loader type is `openapi`

### Template
The `template` option specifies the template file. 

Relative paths are resolved using the directory in which you've run `knit`

## Demo
To demo `knit`, please clone the repository and run the following:

```sh
go run cmd/main.go -glob "./test/*.go"
```