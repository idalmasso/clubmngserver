package common

import "fmt"

//NotFoundError is  an error for not found an object
type NotFoundError struct{
	ID string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s - Not found", e.ID)
}
//BadRequestParametersError is an error for params missing or bad
type BadRequestParametersError struct{
	Message string
}
func (e BadRequestParametersError) Error() string {
	return fmt.Sprintf("%s - Not found", e.Message)
}
//AlreadyExistsError is an error for handling duplicate errors
type AlreadyExistsError struct{
	ID string
}

func (e AlreadyExistsError) Error() string {
	return fmt.Sprintf("%s - Not found", e.ID)
}
