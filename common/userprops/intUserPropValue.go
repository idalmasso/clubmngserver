package userprops

import (
	"encoding/json"
	"strconv"
)

//IntUserPropertyValue is a userProperty value containing an integer
type IntUserPropertyValue struct{
	*IntUserProperty
	Value int64 `json:"value"`
}
//MarshalJSON needed for json-marshalling
func(property *IntUserPropertyValue)	MarshalJSON()([]byte, error){
	return json.Marshal(*property)
}
//GetValueString get the value formatted as string
func(property *IntUserPropertyValue)	GetValueString()string{
	return strconv.FormatInt(property.Value,10)
}
//SetValueString set the value by a string
func(property *IntUserPropertyValue)	SetValueString(value string) error{
	parsedValue, err:=strconv.ParseInt(value,10,64)
	if err!=nil{
		property.Value=parsedValue
	}
	return err
}
