package xmlDeserializer

type IAfterUnmarshaler interface {
	AfterUnmarshal() error
}

type IUnmarshaler interface {
	DeserializeXML(xmlContent string, facoryMap map[string]map[string]interface{}) error
}
