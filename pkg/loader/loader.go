package loader

type SchemaLoader interface {
	LoadFromFile(location string) (interface{}, error)
}
