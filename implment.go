package xmlDeserializer

type AfterUnmarshaler interface {
	AfterUnmarshal()
}

type BeforeUnmarshaler interface {
	BeforeUnmarshal()
}
