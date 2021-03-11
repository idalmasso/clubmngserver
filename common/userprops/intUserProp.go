package userprops

import (
	"encoding/json"
	"strconv"
)

//IntUserProperty is a userProperty containing an integer value
type IntUserProperty struct{
	Name string `json:"name"`
	Value int64 `json:"value"`
	Mandatory bool `json:"mandatory"`
}
//MarshalJSON needed for json-marshalling
func(property *IntUserProperty)	MarshalJSON()([]byte, error){
	return json.Marshal(*property)
}
//GetName returns the name of the property
func(property *IntUserProperty)	GetName()string{
	return property.Name
}
//SetName sets the name of the property
func(property *IntUserProperty)	SetName(value string){
	property.Name=value
}
//GetValueString get the value formatted as string
func(property *IntUserProperty)	GetValueString()string{
	return strconv.FormatInt(property.Value,10)
}
//SetValueString set the value by a string
func(property *IntUserProperty)	SetValueString(value string) error{
	parsedValue, err:=strconv.ParseInt(value,10,64)
	if err!=nil{
		property.Value=parsedValue
	}
	return err
}
//GetTypeString returns the type in a string formatted fashion
func(property *IntUserProperty)	GetTypeString()UserPropertyType{
	return UserTypeInt64
}

//IsMandatory should be used ONLY on frontend to decide if a field should be mandatory. no logic on backend here, because if added after would be hard
func(property *IntUserProperty)	IsMandatory()bool {
	return property.Mandatory
}
