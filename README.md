# Knit

`knit` is a code generation toolkit that simplifies the process of adding and maintaining custom code generators in any project. The tool is written in Go, but its usage is not limited to just Go projects. Generate code into any file using just an input and a template. Inputs can be in a variety of formats including plain `json` or `yaml`, or use custom loaders to parse [openapi](https://swagger.io/specification/) and [graphql](https://graphql.org/) files. 

- [Usage Docs](https://github.com/knitcodegen/knit/blob/develop/docs/usage.md)
- [Go Template Engine](https://pkg.go.dev/text/template)
- [Sprig Template Functions](https://github.com/Masterminds/sprig)


## Installation
### Homebrew
Users with [homebrew](https://brew.sh/) can simply install via tap:
```shell
brew install knitcodegen/tap/knit
```

### Go Install
If you're working in a Go environment
```
go install github.com/knitcodegen/knit/cmd/knit@latest
```

## Examples
The following examples will use `knit` to generate code from an openapi 3.0 specification.

openapi.yml
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
template.tmpl
```tmpl
type Server interface {
{{ range $k, $v := .Paths -}} 
  {{ .Post.OperationID }}(w http.ResponseWriter, r *http.Request)
{{end}}
}
```
example.go:
```go
package example

/*
  @knit input ./openapi.yml
  @knit loader openapi3
  @knit template ./template.tmpl 
*/
// @+knit
// @!knit
```

Now running `knit example.go` will update the file:
```go
package example

/*
  @knit input ./openapi.yml
  @knit loader openapi3
  @knit template ./template.tmpl 
*/
// @+knit

type Server interface {
  CreatePet(w http.ResponseWriter, r *http.Request)
  CreateCar(w http.ResponseWriter, r *http.Request)
}

// @!knit
```

The same code can also be generated from the command line. Instead of inserting the result into a file, it's written to stdout:
```sh
knit generate \
  --input="./openapi.yml" \
  --loader="openapi3" \
  --template="./template.tmpl" > codegen.go
```
codegen.go:
```go
type Server interface {
  CreatePet(w http.ResponseWriter, r *http.Request)
  CreateCar(w http.ResponseWriter, r *http.Request)
}
```