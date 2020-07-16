package xmlDeserializer

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/beevik/etree"
	"reflect"
	"strings"
)

var xmlNameType = reflect.TypeOf(xml.Name{})

func checkStringSliceContains(slice []string, dest string) bool {
	for idx := range slice {
		if slice[idx] == dest {
			return true
		}
	}

	return false
}

func GetXmlName(instance interface{}) (string, error) {
	v := reflect.ValueOf(instance).Elem()

	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		currentField := v.FieldByName(fieldInfo.Name)
		fieldKind := currentField.Kind()
		if fieldKind != reflect.Struct {
			continue
		}
		fieldType := currentField.Type()
		if fieldType != xmlNameType {
			continue
		}

		tag := fieldInfo.Tag // a reflect.StructTag
		xmlTag := tag.Get("xml")
		arr := strings.Split(xmlTag, ",")
		tagHeader := arr[0]

		return tagHeader, nil
	}

	return "", fmt.Errorf("instance have no xmlName struct field")
}

func CloneElement(newRoot *etree.Element, oldRoot *etree.Element) {
	newRoot.Tag = oldRoot.Tag
	newRoot.SetText(oldRoot.Text())

	for _, item := range oldRoot.Attr {
		newRoot.CreateAttr(item.Key, item.Value)
	}
	for _, item := range oldRoot.ChildElements() {
		sub := newRoot.CreateElement(item.Tag)
		sub.SetText(item.Text())
		CloneElement(sub, item)
	}
}

func GetElementXml(elem *etree.Element, addInst bool) string {
	doc := etree.NewDocument()
	if addInst {
		doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	}
	doc.Indent(4)

	root := doc.Element.CreateElement("")
	CloneElement(root, elem)

	str, _ := doc.WriteToString()
	return str
}

func resolveInstanceNames(factory map[string]map[string]interface{}, typ string) []string {
	typDict, ok := factory[typ]
	if !ok {
		return nil
	}

	var arr []string
	for key := range typDict {
		arr = append(arr, key)
	}

	return arr
}

func resolveInstance(factory map[string]map[string]interface{}, typ string, name string) interface{} {
	typDict, ok := factory[typ]
	if !ok {
		return nil
	}

	resIns, ok := typDict[name]
	if !ok {
		return nil
	}

	tp := reflect.ValueOf(resIns).Type()
	h := reflect.New(tp).Interface()
	return h
}

func WrapNodeName(xmlstr string, wrapName string) string {
	endIndex := strings.Index(xmlstr, "?>")
	if endIndex > 0 {
		xmlstr = xmlstr[endIndex+2:]
	}
	bf := new(bytes.Buffer)
	if endIndex > 0 {
		bf.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
		bf.WriteString("\r\n")
	}
	bf.WriteString(fmt.Sprintf("<%s>", wrapName))
	bf.WriteString(xmlstr)
	bf.WriteString(fmt.Sprintf("\r\n</%s>", wrapName))
	return bf.String()
}
