package xmlDeserializer

import (
	"encoding/xml"
	"github.com/beevik/etree"
	"strings"
)

func getChildrenByTagName(root *etree.Element, name string) []*etree.Element {
	var resSlice []*etree.Element

	selectNode := root
	arr := strings.Split(name, ">")
	for i := 0; i < len(arr); i++ {
		subName := arr[i]
		children := selectNode.ChildElements()
		for _, node := range children {
			if node.Tag == subName {
				if i == len(arr)-1 {
					resSlice = append(resSlice, node)
					continue
				}
				selectNode = node
				break
			}
		}
	}

	return resSlice
}

func getChildByTagName(root *etree.Element, name string) *etree.Element {
	selectNode := root
	arr := strings.Split(name, ">")
	for i := 0; i < len(arr); i++ {
		subName := arr[i]
		children := selectNode.ChildElements()
		for _, node := range children {
			if node.Tag == subName {
				if i == len(arr)-1 {
					return node
				}
				selectNode = node
				break
			}
		}
	}

	return nil
}

func unmarshalByElement(root *etree.Element, instance interface{}) error {
	xmlStr := GetElementXml(root, true)
	err := xml.Unmarshal([]byte(xmlStr), instance)
	if err != nil {
		return err
	}

	return nil
}

func checkIsPrefixXmlTag(codeXmlTag string, prefix string) bool {
	facIdx := strings.LastIndex(codeXmlTag, prefix+".")
	if facIdx < 0 {
		return false
	}
	return true
}

func getMapTypeNameFromXmlTag(codeXmlTag string, prefix string) string {
	facIdx := strings.LastIndex(codeXmlTag, prefix+".")
	if facIdx < 0 {
		return ""
	}
	mapTypeName := codeXmlTag[facIdx+len(prefix)+1:]
	return mapTypeName
}

/*
return parent line string and tail string
*/
func parsePrefixXmlTag(codeXmlTag string) ([]string, string) {
	arr := strings.Split(codeXmlTag, ">")
	return arr[:len(arr)-1], arr[len(arr)-1]
}
