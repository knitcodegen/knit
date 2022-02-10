package test

// @knit input $SCHEMA_FILE
// @knit loader openapi
// @knit template ./templates/openapi/test.tmpl
type Generated struct {
	Pet string
}

// @!knit

func foo() {

}

func bar() {

}
