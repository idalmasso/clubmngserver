package userprops

import (
	"encoding/json"
	"strconv"
)

//BoolUserPropertyValue is a userProperty containing a boolean value
type BoolUserPropertyValue struct{
	BoolUserProperty
	Value bool `json:"value"`
}
//MarshalJSON needed for json-marshalling
func(property *BoolUserPropertyValue)	MarshalJSON()([]byte, error){
	return json.Marshal(*property)
}
//GetValueString get the value formatted as string
func(property *BoolUserPropertyValue)	GetValueString()string{
	return strconv.FormatBool(property.Value)
}
//SetValueString set the value by a string
func(property *BoolUserPropertyValue)	SetValueString(value string) error{
	parsedValue, err:=strconv.ParseBool(value)
	if err!=nil{
		property.Value=parsedValue
	}
	return err
}
