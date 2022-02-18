package loader

type SchemaLoader interface {
	LoadFromData(data []byte) (interface{}, error)
	//LoadFromFile(location string) (interface{}, error)
}
