# Knit

Knit is a code generation toolkit that simplifies the process of adding and maintaining custom code generators in any project.

## Installation
### Homebrew
```shell
brew install knitcodegen/tap/knit
```

### Go Install
```
go install github.com/knitcodegen/knit@latest
```

## Example
Here's an example of `knit`'s usage with OpenAPI and template literals.

schema.yml
```yml
openapi: "3.0.0"
info:
  version: 1.0.0
  title: Example Codegen
paths:
  /pets:
    post:
      operationId: CreatePet
  /cars:
    post:
      operationId: CreateCar
```

example.go:
```go
package example

/*
  Generates void functions for every endpoint

  @knit input ./schema.yml
  @knit loader openapi
  @knit template tmpl`
    {{ range $k, $v := .Paths }} 
      func {{ .Post.OperationID }}() error {
        return nil
      }
    {{end}}  
  `
*/
// @+knit
// @!knit
```

Now running `knit example.go` will update the file:
```go
package example

/*
  Generates void functions for every endpoint

  @knit input ./schema.yml
  @knit loader openapi
  @knit template tmpl`
    {{ range $k, $v := .Paths }} 
      func {{ .Post.OperationID }}() error {
        return nil
      }
    {{end}}  
  `
*/
// @+knit

func CreatePet() error {
  return nil
}

func CreateCar() error {
  return nil
}

// @!knit
```

## Annotations
Annotations in `knit` are used to identify code generator options and the output location of the generated code. 

It takes a combination of both option annotations and codegen annotations to successfully use knit to insert generated code into a file.

### Option Annotations
Options are defined on the opening annotations in the following format:
```
@knit <option> <value>
```

See below for a list of available options.

### Codegen Annotations
The location in which the generated code is inserted into a code file is dictated by the open/close knit annotations:

```
// @+knit
  < code is generated here >
// @!knit
```

### Annotation Parsing Algorithm
It may be helpful to understand how the parsing algorithm works when knitting multiple generators in one file.  

The parser algorithm first splits the entire code file into `blocks` by the ending annotation: `@!knit`

It then loops over each block and uses regex to match and parse any generator options. 

```go
package test

func foo() {}
func bar() {}

// @knit input schema.yml
// @knit loader openapi
// @knit template ./file.tmpl
// @+knit

func Generated() { foo() }

// @!knit
//       ^ new "block" was split here

/// ... more code

/*
  Another generator definition entirely.
  This won't use the previously defined options.

  @knit input schema.yml
  @knit loader openapi
  @knit template tmpl`
    func Generated2 { bar() }
  `
*/
// @+knit

func Generated2() { bar() }

// @!knit
```

## Options
Options allow you to define the behavior of your code generator in a series of key-value pairs.

```
@knit <option> <value>
```

All options support environment variable expansion. All `$env` variables starting with `$` will be expanded before the option line is parsed. 
### `input`
The `input` option specifies the input file. 

Relative paths to files are resolved using the directory in which `knit` has been executed.

Currently remote input loading is not available but is planned. Please see issue #5 for more details.

### `loader`
The `loader` option specifies the loader used to load the input file. 

These are the available loader types:
- `yml`
- `json`
- `openapi`

### `template`
The `template` option specifies either a template file or template literal. 

Currently the only supported template engine is the `text/template` package native to Go. 

#### File
```
/*
  @knit input ./schema.yml
  @knit loader openapi
  @knit template ./relative/path/to/file.tmpl
*/
```
Relative paths to files are resolved using the directory in which `knit` has been executed.

Currently remote template loading is not available but is planned. Please see issue #5 for more details.
#### Template Literal
```
/*
  @knit input ./schema.yml
  @knit loader openapi
  @knit template tmpl`
    type Generated struct {
    {{ range $k, $v := .Paths }} 
        {{ .Post.OperationID }} string
    {{end}}
    }
  `
*/
```
Template literals are defined within backticks and must be prefixed with the extension of the file the template would otherwise be defined in.

As a rule of thumb, if a template literal exceeds 10 lines, it should probably be promoted to its own template file.