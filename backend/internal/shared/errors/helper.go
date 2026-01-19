package errors

import "errors"

func IsCustomError(err error) bool {
	var customError *CustomError
	return errors.As(err, &customError)
}

func AsCustomError(err error) (*CustomError, bool) {
	var customError *CustomError
	ok := errors.As(err, &customError)
	return customError, ok
}

func IsLayer(err error, layer ErrorLayer) bool {
	customErr, ok := AsCustomError(err)
	if !ok {
		return false
	}
	return customErr.layer == layer
}

func IsFeature(err error, feature ErrorFeature) bool {
	customErr, ok := AsCustomError(err)
	if !ok {
		return false
	}
	return customErr.feature == feature
}

func GetErrorChain(err error) []error {
	var chain []error
	currentErr := err
	for currentErr != nil {
		chain = append(chain, currentErr)
		currentErr = errors.Unwrap(currentErr)
	}
	return chain
}
