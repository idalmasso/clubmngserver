package userprops

import (
	"encoding/json"
	"strconv"
)

//FloatUserPropertyValue is a userProperty containing a float value
type FloatUserPropertyValue struct{
	FloatUserProperty
	Value float64 `json:"value"`
	
}
//MarshalJSON needed for json-marshalling
func(property *FloatUserPropertyValue)	MarshalJSON()([]byte, error){
	return json.Marshal(*property)
}

//GetValueString get the value formatted as string
func(property *FloatUserPropertyValue)	GetValueString()string{
	return strconv.FormatFloat(property.Value, byte('E'),6,64)
}
//SetValueString set the value by a string
func(property *FloatUserPropertyValue)	SetValueString(value string) error{
	parsedValue, err:=strconv.ParseFloat(value,64)
	if err!=nil{
		property.Value=parsedValue
	}
	return err
}
