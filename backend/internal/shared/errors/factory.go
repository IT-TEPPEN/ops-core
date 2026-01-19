package errors

import "fmt"

type ErrorFactory struct {
	feature         ErrorFeature
	layer           ErrorLayer
	registeredCodes map[string]bool
}

func NewFactory(feature ErrorFeature, layer ErrorLayer) *ErrorFactory {
	return &ErrorFactory{
		feature:         feature,
		layer:           layer,
		registeredCodes: make(map[string]bool),
	}
}

func (f *ErrorFactory) NewError(errorCode, message string) *CustomError {
	return NewErrorFromCode(f.feature, f.layer, errorCode, message)
}

func (f *ErrorFactory) ValidateErrorCode(errorCode string) error {
	if len(errorCode) != 7 {
		return fmt.Errorf("Invalid error code length: expected 7 characters, got %d (%s)", len(errorCode), errorCode)
	}

	codeFeature := ErrorFeature(errorCode[0:2])
	codeLayer := ErrorLayer(errorCode[2:3])
	codeNo := errorCode[3:7]

	if codeFeature != f.feature {
		return fmt.Errorf("Mismatched feature in error code: expected %s, got %s (%s)", f.feature, codeFeature, errorCode)
	}

	if codeLayer != f.layer {
		return fmt.Errorf("Mismatched layer in error code: expected %s, got %s (%s)", f.layer, codeLayer, errorCode)
	}

	for _, c := range codeNo {
		if c < '0' || c > '9' {
			return fmt.Errorf("Invalid characters in error code number: expected digits only, got %s (%s)", codeNo, errorCode)
		}
	}

	return nil
}

func (f *ErrorFactory) MustBuilder(errorCode string) func(message string) *CustomError {
	if err := f.ValidateErrorCode(errorCode); err != nil {
		panic(err)
	}

	if _, exists := f.registeredCodes[errorCode]; exists {
		panic(fmt.Sprintf("Error code already registered: %s", errorCode))
	}

	f.registeredCodes[errorCode] = true

	return func(message string) *CustomError {
		return f.NewError(errorCode, message)
	}
}
