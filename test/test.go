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

func foo() {

}

func bar() {

}
