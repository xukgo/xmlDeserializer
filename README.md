# xmlDeserializer
fill struct filed from XML easily in go,It is very important that it support interface field.

The xmlDeserializer package is a lightweight, pure go package that deserializer XML to an object. Its design was inspired by my previous job about java xml.
Because golang's reflection function is not as powerful as java,I found a tool that can deserialize the interface and couldn't find it for a long time.
Now I write one myself with registered factory mode.

Some of the package's capabilities and features:

	Dependent the library:beevik/etree.
	Compatible with official deserialization.
	Need to register the interface factory first.
	Support for custom interface prefix flags.
	Support deep search prefix.
	Support for deep recursive deserialization.
	
Deserializer an XML document
The following example deserializer an XML document to an object with interfaces.

	//code here
	type RootModel struct {
		XMLName xml.Name `xml:"Root"`

		InstancePtrArray []*EqualRulerB `xml:"Rules>EqualRuleB"`
		SingleInterface	IEqualRuler `xml:"Factory.EqualRuler"`
		InterfaceArray	[]IParser `xml:"Factory.Parser"`
		InterfacePtrDeepChildren []IEqualRuler `xml:"EqualRulers>Factory.EqualRuler"`
	}

	var instanceMap map[string]map[string]interface{}
	
	initXmlInstanceFactory()
	defer disposeXmlInstanceFactory()
	xmlDeserializer := xmlUtil.NewDeserializer("Factory", instanceMap)
	instance := &RootModel{}
	err := xmlDeserializer.Deserialize([]byte(testXml),instance)
	if err != nil{
		t.Fail()
	}

source xml string below:

	<?xml version="1.0" encoding="UTF-8"?>
	<Root>
		<Rules>
			<EqualRuleB>
				<Name>arr1</Name>
				<Match>callin1</Match>
			</EqualRuleB>
			<EqualRuleB>
				<Name>arr2</Name>
				<Match>callin2</Match>
			</EqualRuleB>
			<EqualRuleB>
				<Name>arr3</Name>
				<Match>callin3</Match>
			</EqualRuleB>
		</Rules>
		<EqualRuleA>
			<Name>singleInstance</Name>
			<Match>callin</Match>
		</EqualRuleA>
		<NotifyParser>
			<Name>pppp</Name>
		</NotifyParser>
		<NotifyParser>
			<Name>mmmm</Name>
		</NotifyParser>
		<CallParser>
			<Index>111</Index>
		</CallParser>
		<CallParser>
			<Index>222</Index>
		</CallParser>
		<EqualRulers>
			<EqualRuleA>
				<Name>sub1</Name>
				<Match>a1</Match>
			</EqualRuleA>
			<EqualRuleB>
				<Name>sub2</Name>
				<Match>b2</Match>
			</EqualRuleB>
			<EqualRuleA>
				<Name>sub3</Name>
				<Match>a3</Match>
			</EqualRuleA>
		</EqualRulers>
	</Root>
	
Contributing:
This project accepts contributions. Just fork the repo and submit a pull request!

Please forgive me for my poor English.	author-xukuan