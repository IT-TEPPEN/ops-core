package errors

import "fmt"

type CustomError struct {
	feature     ErrorFeature
	layer       ErrorLayer
	no          string // 4-digit number
	kind        string
	messsage    string
	parentError error
}

func NewErrorFromCode(feature ErrorFeature, layer ErrorLayer, errorCode, message string) *CustomError {
	codeNo := errorCode[3:7]

	return &CustomError{
		feature:  feature,
		layer:    layer,
		no:       codeNo,
		kind:     errorCode,
		messsage: message,
	}
}

func (e *CustomError) Error() string {
	if e.parentError != nil {
		return fmt.Sprintf("[%s%s%s] %s\n  - %s", e.feature, e.layer, e.no, e.messsage, e.parentError.Error())
	}

	return fmt.Sprintf("[%s%s%s] %s", e.feature, e.layer, e.no, e.messsage)
}

func (e *CustomError) WithParent(err error) *CustomError {
	return &CustomError{
		feature:     e.feature,
		layer:       e.layer,
		no:          e.no,
		kind:        e.kind,
		messsage:    e.messsage,
		parentError: err,
	}
}

func (e *CustomError) Kind() string {
	return e.kind
}

func (e *CustomError) Message() string {
	return e.messsage
}

func (e *CustomError) Unwrap() error {
	return e.parentError
}

func (e *CustomError) Is(target error) bool {
	if ce, ok := target.(*CustomError); ok {
		return e.kind == ce.kind
	}
	return false
}

func (e *CustomError) Feature() ErrorFeature {
	return e.feature
}

func (e *CustomError) Layer() ErrorLayer {
	return e.layer
}

func (e *CustomError) No() string {
	return e.no
}

func (e *CustomError) Parent() error {
	return e.parentError
}
