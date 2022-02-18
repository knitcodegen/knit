package test

/*
  @knit input $SCHEMA_FILE
  @knit loader openapi
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
