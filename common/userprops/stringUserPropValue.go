package userprops

import "encoding/json"

//StringUserPropertyValue is a userProperty string value
type StringUserPropertyValue struct{
	StringUserProperty
	Value string `json:"value"`

}
//MarshalJSON needed for json-marshalling
func(property *StringUserPropertyValue)	MarshalJSON()([]byte, error){
	return json.Marshal(*property)
}
//GetValueString get the value formatted as string
func(property *StringUserPropertyValue)	GetValueString()string{
	return property.Value
}
//SetValueString set the value by a string
func(property *StringUserPropertyValue)	SetValueString(value string) error{
	property.Value=value
	return nil
}
