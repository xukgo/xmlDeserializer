package xmlDeserializer

type IXmlUnmarshaler interface {
	AfterUnmarshal() error
}
