package xmlDeserializer

import (
	"encoding/xml"
	"testing"
)

//define IParser interface
type IParser interface {
	Parse(string) error
}

//define NotifyParser struct
type NotifyParser struct {
	XMLName   xml.Name       `xml:"NotifyParser"`
	Name	string	`xml:"Name"`
}

func (this *NotifyParser)Parse(string) error{
	return nil
}

//define NotifyParser struct
type CallParser struct {
	XMLName   xml.Name       `xml:"CallParser"`
	Index	int	`xml:"Index"`
}

func (this *CallParser)Parse(string) error{
	return nil
}

//define IEqualRuler interface
type IEqualRuler interface {
	Equal(indata string) bool
}

//define EqualRulerA struct and EqualRulerB struct
type EqualRulerA struct {
	XMLName   xml.Name       `xml:"EqualRuleA"`
	AName	string	`xml:"Name"`
	IsMatchAnyString	bool	`xml:"MatchAny"`
	MatchString	string	`xml:"Match"`
}

func (this *EqualRulerA)Equal(indata string) bool{
	return false
}

type EqualRulerB struct {
	XMLName   xml.Name       `xml:"EqualRuleB"`
	BName	string	`xml:"Name"`
	IsMatchAnyString	bool	`xml:"MatchAny"`
	MatchString	string	`xml:"Match"`
}

func (this *EqualRulerB)Equal(indata string) bool{
	return false
}

type RootModel struct {
	XMLName   xml.Name           `xml:"Root"`

	InstancePtrArray	[]*EqualRulerB `xml:"Rules>EqualRuleB"`

	SingleInterface	IEqualRuler   `xml:"Factory.EqualRuler"`
	InterfaceArray	[]IParser `xml:"Factory.Parser"`
	InterfacePtrDeepChildren	[]IEqualRuler `xml:"EqualRulers>Factory.EqualRuler"`
}

//define instance factory
var instanceMap map[string]map[string]interface{}

func disposeXmlInstanceFactory() {
	instanceMap = nil
}

func initXmlInstanceFactory() {
	instanceMap = make(map[string]map[string]interface{})

	equalRuleMap := make(map[string]interface{})
	instanceMap["EqualRuler"] = equalRuleMap
	equalRuleMap["EqualRuleA"] = EqualRulerA{}
	equalRuleMap["EqualRuleB"] = EqualRulerB{}

	parserMap := make(map[string]interface{})
	instanceMap["Parser"] = parserMap
	parserMap["NotifyParser"] = NotifyParser{}
	parserMap["CallParser"] = CallParser{}
}

//test for Xml Deserialize
func TestXmlDeserializer(t *testing.T){
	testXml :=
		`<Root>
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
	</Root>`

	initXmlInstanceFactory()
	defer disposeXmlInstanceFactory()
	xmlDeserializer := NewDeserializer("Factory", instanceMap)
	instance := &RootModel{}
	err := xmlDeserializer.Deserialize([]byte(testXml),instance)
	if err != nil{
		t.Fail()
	}

	if instance.InstancePtrArray == nil || len(instance.InstancePtrArray) != 3{
		t.Fail()
	}
	if instance.InstancePtrArray[0].BName != "arr1" ||instance.InstancePtrArray[0].MatchString != "callin1"{
		t.Fail()
	}
	if instance.InstancePtrArray[1].BName != "arr2" ||instance.InstancePtrArray[1].MatchString != "callin2"{
		t.Fail()
	}
	if instance.InstancePtrArray[2].BName != "arr3" ||instance.InstancePtrArray[2].MatchString != "callin3"{
		t.Fail()
	}

	var ok bool
	ruleA,ok := instance.SingleInterface.(*EqualRulerA)
	if !ok{
		t.Fail()
	}
	if ruleA.AName != "singleInstance" || ruleA.MatchString != "callin"{
		t.Fail()
	}

	parser,ok := instance.InterfaceArray[0].(*NotifyParser)
	if !ok{
		t.Fail()
	}
	if parser.Name != "pppp"{
		t.Fail()
	}
	parser,ok = instance.InterfaceArray[1].(*NotifyParser)
	if !ok{
		t.Fail()
	}
	if parser.Name != "mmmm"{
		t.Fail()
	}
	parser2,ok := instance.InterfaceArray[2].(*CallParser)
	if !ok{
		t.Fail()
	}
	if parser2.Index != 111{
		t.Fail()
	}
	parser2,ok = instance.InterfaceArray[3].(*CallParser)
	if !ok{
		t.Fail()
	}
	if parser2.Index != 222{
		t.Fail()
	}

	ruleA2,ok := instance.InterfacePtrDeepChildren[0].(*EqualRulerA)
	if !ok{
		t.Fail()
	}
	if ruleA2.AName != "sub1" || ruleA2.MatchString != "a1"{
		t.Fail()
	}
	ruleB2,ok := instance.InterfacePtrDeepChildren[1].(*EqualRulerB)
	if !ok{
		t.Fail()
	}
	if ruleB2.BName != "sub2" || ruleB2.MatchString != "b2"{
		t.Fail()
	}
	ruleA3,ok := instance.InterfacePtrDeepChildren[2].(*EqualRulerA)
	if !ok{
		t.Fail()
	}
	if ruleA3.AName != "sub3" || ruleA3.MatchString != "a3"{
		t.Fail()
	}
}


