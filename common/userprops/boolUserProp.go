package userprops

import (
	"encoding/json"
	"strconv"
)

//BoolUserProperty is a userProperty containing a boolean value
type BoolUserProperty struct{
	Name string `json:"name"`
	Value bool `json:"value"`
	Mandatory bool `json:"mandatory"`
}
//MarshalJSON needed for json-marshalling
func(property *BoolUserProperty)	MarshalJSON()([]byte, error){
	return json.Marshal(*property)
}
//GetName returns the name of the property
func(property *BoolUserProperty)	GetName()string{
	return property.Name
}
//SetName sets the name of the property
func(property *BoolUserProperty)	SetName(value string){
	property.Name=value
}
//GetValueString get the value formatted as string
func(property *BoolUserProperty)	GetValueString()string{
	return strconv.FormatBool(property.Value)
}
//SetValueString set the value by a string
func(property *BoolUserProperty)	SetValueString(value string) error{
	parsedValue, err:=strconv.ParseBool(value)
	if err!=nil{
		property.Value=parsedValue
	}
	return err
}

//GetTypeString returns the type in a string formatted fashion
func(property *BoolUserProperty)	GetTypeString()UserPropertyType{
	return UserTypeBool
}

//IsMandatory should be used ONLY on frontend to decide if a field should be mandatory. no logic on backend here, because if added after would be hard
func(property *BoolUserProperty)	IsMandatory()bool {
	return property.Mandatory
}
