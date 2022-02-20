# Usage
## Generator
Whether you are using `knit` from the command line or embedded in a file, the behaviour of a code generator is defined using a series of key-value option pairs. 

### `input`
The `input` option specifies an input file or literal. This option is _required_ for all generators.

#### File
An `input` can be defined as a path to a file. The two default input file types are `json` and `yaml`. 

Relative paths to files are resolved using the directory in which `knit` has been executed.

Currently remote file loading is not available but is planned for a future release. Please follow [#5](https://github.com/knitcodegen/knit/issues/5) for details and updates.

#### Literal
An `input` can also be defined as a literal. A literal is a multiline string surrounded by backticks prefixed by the extension of the file the text would otherwise reside in.
```text
yaml`
schema:
  paths:
    /dogs:
      post:
        summary: "Creates a \`dog\`"
        operationId: CreateDog
`
```
Literals will be fed to their respective loader "literally", so be weary of spacing and formatting for input types like `yaml`.

Backtick characters in literals can be escaped using a prefixed backslash.

### `loader`
The `loader` option specifies the loader used to load the input file.

These are the available loader types:
- `yml` / `yaml`
- `json`
- `openapi3`

These are the planned loader types:
- `protobuf` (follow [#1](https://github.com/knitcodegen/knit/issues/1))
- `graphql` (follow [#2](https://github.com/knitcodegen/knit/issues/2))

Sometimes the loader type can be inferred from the input file type but it is important to understand when to be explicit.

As an example, technically an openapi specification could be loaded by the `yaml` loader, however, more specialized loaders like `openapi3` know how to resolve file references and validate schema adherence.

### `template`
The `template` option specifies a template file or literal. This option is _required_ for all generators.

Currently the only supported template engines are the ones native to Go:
- `text/template` 
- `html/template`

#### File

#### Literal
A `template` can also be defined as a literal. A literal is a multiline string surrounded by backticks prefixed by the extension of the file the text would otherwise reside in.
```
tmpl`
  type Generated struct {
    {{ range $k, $v := .Map }} 
      {{ $k }} string \`json:"{{ $v }}"\`
    {{end}}
  }
`
```

Backtick characters in literals can be escaped using a prefixed backslash.

## CLI
`knit` has a command line interface that allows you to load inputs and execute templates. All generated code is sent directly to stdout so it can be appended to a file or piped to another tool.

```sh
knit generate \
  --input="./openapi.yml" \
  --loader="openapi3" \
  --template="./template.tmpl" > codegen.go
```

## Annotations
Annotations allow `knit` to embed generated code into a file. Annotations are used to identify code generator options and the output location of the generated code. 

It takes a combination of both option annotations and codegen annotations to successfully use knit to insert generated code into a file.

`knit` is capable of inserting generated code into any file, as long as that file supports comments where the options can be defined.

example.ts:
```ts
/*
  @knit input ./sizes.yml
  @knit template tmpl`
    {{ range $k, $v := .Sizes }} 
      export const {{ $k }} = "$v";
    {{ end }}
  ` 
*/
// @+knit
// @!knit
```
Running `knit example.ts` produces the following:
```ts
/*
  @knit input ./sizes.yml
  @knit template tmpl`
    {{ range $k, $v := .Sizes }} 
      export const {{ $k }} = "$v";
    {{ end }}
  ` 
*/
// @+knit
export const SIZE_SM = "small"
export const SIZE_MD = "medium"
export const SIZE_LG = "large"
// @!knit
```


### Option Annotations
Options are defined on the opening annotations in the following format:
```
@knit <option> <value>
```

All annotated options support environment variable expansion. All `$env` variables starting with `$` will be expanded before the option line is parsed. 
```
@knit input $OPENAPI_SPEC
```


### Codegen Annotations
The location in which the generated code is inserted into a code file is dictated by the open/close knit annotations:

```
// @+knit
  < code is generated here >
// @!knit
```