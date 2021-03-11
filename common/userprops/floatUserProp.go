package userprops

import (
	"encoding/json"
	"strconv"
)

//FloatUserProperty is a userProperty containing a float value
type FloatUserProperty struct{
	Name string `json:"name"`
	Value float64 `json:"value"`
	Mandatory bool `json:"mandatory"`
}
//MarshalJSON needed for json-marshalling
func(property *FloatUserProperty)	MarshalJSON()([]byte, error){
	return json.Marshal(*property)
}
//GetName returns the name of the property
func(property *FloatUserProperty)	GetName()string{
	return property.Name
}
//SetName sets the name of the property
func(property *FloatUserProperty)	SetName(value string){
	property.Name=value
}
//GetValueString get the value formatted as string
func(property *FloatUserProperty)	GetValueString()string{
	return strconv.FormatFloat(property.Value, byte('E'),6,64)
}
//SetValueString set the value by a string
func(property *FloatUserProperty)	SetValueString(value string) error{
	parsedValue, err:=strconv.ParseFloat(value,64)
	if err!=nil{
		property.Value=parsedValue
	}
	return err
}

//GetTypeString returns the type in a string formatted fashion
func(property *FloatUserProperty)	GetTypeString()UserPropertyType{
	return UserTypeFloat64
}

//IsMandatory should be used ONLY on frontend to decide if a field should be mandatory. no logic on backend here, because if added after would be hard
func(property *FloatUserProperty)	IsMandatory()bool {
	return property.Mandatory
}
