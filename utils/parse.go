package utils

import (
	"errors"
	"fmt"

	"github.com/Jeffail/gabs/v2"
	"github.com/beevik/etree"
)

func ParseUIDFromJson(jsonByte []byte, uidPath string) (uid string, err error) {
	jsonParsed, err := gabs.ParseJSON(jsonByte)
	if err != nil {
		return "", errors.New("parse json err")
	}
	uidInterface := jsonParsed.Path(uidPath).Data()
	switch uidInterface.(type) {
	case float64:
		uid = fmt.Sprintf("%f", uidInterface.(float64))
	case string:
		uid = uidInterface.(string)
	default:
		return "", errors.New("parse json uid failed, must be number or string")
	}
	return uid, nil
}

func ParseUIDFromXml(xmlByte []byte, uidPath string) (uid string, err error) {
	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(xmlByte); err != nil {
		return "", errors.New("parse xml failed")
	}
	id := doc.FindElement(uidPath)
	if id == nil {
		return "", errors.New("parse xml uid failed")
	} else {
		uid = id.Text()
	}
	return uid, nil
}
