package xmlDeserializer

import (
	"encoding/xml"
	"fmt"
	"github.com/beevik/etree"
	"reflect"
	"strings"
)

type Deserializer struct{
	prefix string
	factory map[string]map[string]interface{}
}

func NewDeserializer(prefix string, factory map[string]map[string]interface{}) *Deserializer{
	model := Deserializer{
		prefix:prefix,
		factory:factory,
	}
	return &model
}

func (this *Deserializer) Deserialize(xmlContent []byte, instance interface{})error{
	err := xml.Unmarshal(xmlContent, instance)
	if err != nil {
		return err
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlContent); err != nil {
		return err
	}

	xmlNameString,err := GetXmlName(instance)
	if err != nil {
		return err
	}

	err = this.parseByElement(doc.SelectElement(xmlNameString), instance)
	if err != nil {
		return err
	}
	return nil
}

//codeXmlTag:a>b>Factory.typeName  like this
func (this *Deserializer) getMatchTagNodes(root *etree.Element, codeXmlTag string)[]*etree.Element{
	parents,xmlTag := parsePrefixXmlTag(codeXmlTag)
	mapTypeName := getMapTypeNameFromXmlTag(xmlTag, this.prefix)
	if mapTypeName == ""{
		return nil
	}
	instanceNames := resolveInstanceNames(this.factory, mapTypeName)
	if instanceNames == nil || len(instanceNames) == 0{
		return nil
	}

	selectNode := root
	for _,tagName := range parents{
		selectNode = getChildByTagName(selectNode, tagName)
		if selectNode == nil{
			return nil
		}
	}

	for idx := range instanceNames{
		instanceNames[idx] = strings.ToLower(instanceNames[idx])
	}

	var resArr []*etree.Element
	children := selectNode.ChildElements()
	for _,node := range children{
		if checkStringSliceContains(instanceNames, strings.ToLower(node.Tag)){
			resArr = append(resArr, node)
		}
	}

	return resArr
}

func (this *Deserializer) parseStrcutPtr(root *etree.Element, instance interface{}, tagName string) error{
	node := getChildByTagName(root, tagName)
	err := this.parseByElement(node, instance)
	if err != nil {
		return err
	}
	return nil
}

func (this *Deserializer) parsePrefixXmlInterfaceField(root *etree.Element, field reflect.Value, xmlTagName string) error {
	nodes := this.getMatchTagNodes(root, xmlTagName)
	if nodes == nil{
		return nil
	}
	if len(nodes) != 1{
		return fmt.Errorf("%s expect one interface node,acutal count is %d", xmlTagName, len(nodes))
	}

	selectNode := nodes[0]
	mapTypeName := getMapTypeNameFromXmlTag(xmlTagName, this.prefix)
	newInstance := resolveInstance(this.factory, mapTypeName, selectNode.Tag)
	if newInstance == nil{
		return fmt.Errorf("resolveInstance by %s return nil", xmlTagName)
	}
	err := unmarshalByElement(selectNode ,newInstance)
	if err != nil{
		return err
	}

	err = this.parseByElement(selectNode, newInstance)
	if err != nil{
		return err
	}

	field.Set(reflect.ValueOf(newInstance))
	return nil
}

func (this *Deserializer) parsePrefixXmlInterfaceSliceField(root *etree.Element, field reflect.Value, xmlTagName string) error {
	nodes := this.getMatchTagNodes(root, xmlTagName)
	if nodes == nil || len(nodes) == 0{
		return nil
	}

	mapTypeName := getMapTypeNameFromXmlTag(xmlTagName, this.prefix)
	sliceCarrier := make([]reflect.Value, 0)
	for _,node := range nodes{
		newInstance := resolveInstance(this.factory, mapTypeName, node.Tag)
		if newInstance == nil{
			return fmt.Errorf("resolveInstance by %s return nil", xmlTagName)
		}
		err := unmarshalByElement(node ,newInstance)
		if err != nil{
			return err
		}

		err = this.parseByElement(node, newInstance)
		if err != nil{
			return err
		}

		sliceCarrier = append(sliceCarrier, reflect.ValueOf(newInstance))
	}
	arrBind := reflect.Append(field, sliceCarrier...)
	field.Set(arrBind)
	return nil
}

func (this *Deserializer) parsePrefixXmlField(root *etree.Element, field reflect.Value, xmlTagName string) error {
	fieldKind := field.Kind()
	if fieldKind != reflect.Interface && fieldKind != reflect.Slice{
		return nil
	}

	if fieldKind == reflect.Interface{
		err := this.parsePrefixXmlInterfaceField(root, field, xmlTagName)
		if err != nil{
			return err
		}
	}

	if fieldKind == reflect.Slice {
		//这里判断必须为nil
		if !field.IsNil(){
			return nil
		}
		err := this.parsePrefixXmlInterfaceSliceField(root, field, xmlTagName)
		if err != nil{
			return err
		}
	}

	return nil
}

func (this *Deserializer) parseNotPrefixXmlField(root *etree.Element, field reflect.Value, xmlTagName string) error {
	fieldKind := field.Kind()
	if fieldKind == reflect.Struct {
		h := field.Addr().Interface()
		err := this.parseStrcutPtr(root, h, xmlTagName)
		if err != nil {
			return err
		}
	}
	if fieldKind == reflect.Ptr {
		//这里判断为nil，跳过空指针
		if field.IsNil(){
			return nil
		}

		h := field.Interface()

		err := this.parseStrcutPtr(root, h, xmlTagName)
		if err != nil {
			return err
		}
	}
	if fieldKind == reflect.Slice {
		//这里判断为nil，跳过空指针
		if field.IsNil(){
			return nil
		}
		nodes := getChildrenByTagName(root, xmlTagName)
		for i:=0;i< field.Len();i++{
			iv := field.Index(i).Interface()
			sliceElemKind := reflect.TypeOf(iv).Kind()

			if sliceElemKind == reflect.Struct {
				h := field.Index(i).Addr().Interface()
				err := this.parseByElement(nodes[i], h)
				if err != nil {
					return err
				}
			}

			if sliceElemKind == reflect.Ptr {
				h := field.Index(i).Interface()
				err := this.parseByElement(nodes[i], h)
				if err != nil {
					return err
				}
			}
			//fmt.Println(iv)
		}
	}

	return nil
}

func (this *Deserializer) parseByElement(root *etree.Element, instance interface{})error{
	var err error
	v := reflect.ValueOf(instance).Elem() // the struct variable

	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		theField := v.FieldByName(fieldInfo.Name)
		fieldKind := theField.Kind()
		//fmt.Println(fieldKind)
		if fieldKind != reflect.Ptr  && fieldKind != reflect.Struct && fieldKind != reflect.Interface && fieldKind != reflect.Slice{
			continue
		}
		fieldType := theField.Type()
		//fmt.Println(fieldType)
		if fieldType == xmlNameType{
			continue
		}

		tag := fieldInfo.Tag           // a reflect.StructTag
		xmlTagContent := tag.Get("xml")

		//去掉逗号后面内容 如 `xml:"nodeName,omitempty"`
		arr := strings.Split(xmlTagContent, ",")
		xmlTagName := arr[0]

		//不是factory前缀
		if !checkIsPrefixXmlTag(xmlTagName, this.prefix){
			err  = this.parseNotPrefixXmlField(root, theField, xmlTagName)
			if err != nil{
				return err
			}
			continue
		}

		//下面处理factory前缀的，目前只允许interface和它对应的数组
		err  = this.parsePrefixXmlField(root, theField, xmlTagName)
		if err != nil{
			return err
		}
	}

	return nil
}
