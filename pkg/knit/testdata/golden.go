package test

/*
  @knit input $SCHEMA_FILE
  @knit loader openapi3
  @knit template tmpl`
    // HelloWorld is a struct
    type HelloWorld struct {

    }
  `
*/
// @+knit

// HelloWorld is a struct
type HelloWorld struct {
}

// @!knit

/*
  @knit input $SCHEMA_FILE
  @knit loader openapi3
  @knit template tmpl`
    // HolaMundo is a struct
    type HolaMundo struct {

    }
  `
*/
// @+knit /!\ GENERATED CODE. DO NOT EDIT

// HolaMundo is a struct
type HolaMundo struct {
}

// @!knit /!\ GENERATED CODE. DO NOT EDIT!

/*
  @knit input ./schemas/yml/test.yml
  @knit loader yml
  @knit template tmpl`
    type Generated2 struct {
    {{ range $k, $v := .test }}
        {{ $k }} string \`yml:"{{ $v }}"\`
    {{end}}
    }
  `
*/
// @+knit

type Generated2 struct {
	A string `yml:"b"`

	C string `yml:"d"`

	E string `yml:"f"`
}

// @!knit

/*
  @knit input yml`
    test:
      A: "b"
      C: "d"
      E: "f"
  `
  @knit template tmpl`
    type Generated3 struct {
    {{ range $k, $v := .test }}
        {{ $k }} string \`yml:"{{ $v }}"\`
    {{end}}
    }
  `
*/
// @+knit

type Generated3 struct {
	A string `yml:"b"`

	C string `yml:"d"`

	E string `yml:"f"`
}

// @!knit

/*
  @knit input json`{
    "test": {
      "A": "b",
      "C": "d",
      "E": "f"
    }
  }`
  @knit template tmpl`
    type GeneratedFromJSON struct {
    {{ range $k, $v := .test }}
        {{ $k }} string \`json:"{{ $v }}"\`
    {{end}}
    }
  `
*/
// @+knit

type GeneratedFromJSON struct {
	A string `json:"b"`

	C string `json:"d"`

	E string `json:"f"`
}

// @!knit
