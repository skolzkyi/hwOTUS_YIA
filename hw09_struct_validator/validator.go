package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

type Field struct {
	Name   string
	RawTag string
}

type Tag struct {
	Name   string
	Params []string
}

type Validator struct {
	Name          string
	SupportedType reflect.Kind
	Method        func(data any, params []any) error
}

var (
	ErrBadInputData        = errors.New("input data is not struct or slice of struct")
	ErrBadTagFormat        = errors.New("tag format is bad")
	ErrUnsupportedType     = errors.New("unsupported type")
	ErrBadTypeAssertion    = errors.New("bad type assertion")
	ErrUnsupportedValidate = errors.New("unsupported validate")
	ErrStrLen              = errors.New("out of limit string length")
	ErrStrRegexp           = errors.New("string does not match the conditions of the regular expression")
	ErrStrMStrings         = errors.New("string is not in multiple strings")
	ErrNumLessMin          = errors.New("number less than min")
	ErrNumGreaterMax       = errors.New("number greater max")
	ErrNumMNums            = errors.New("number is not in multiple nums")
)

func Validate(v interface{}) error {
	var valErrors ValidationErrors
	var err error
	validators := initValidators()
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice:
		if reflect.TypeOf(v).Elem().Kind() == reflect.Struct {
			valueOfInputSlice := reflect.ValueOf(v)
			for i := 0; i < valueOfInputSlice.Len(); i++ {
				valErrors, err = validateData(valueOfInputSlice.Index(i), validators)
				if err != nil {
					return err
				}
			}
		} else {
			return ErrBadInputData
		}
	case reflect.Struct:
		valErrors, err = validateData(v, validators)
		if err != nil {
			return err
		}
	default:
		return ErrBadInputData
	}
	return packErrors(valErrors)
}

func validateData(v any, validators map[string]Validator) (ValidationErrors, error) {
	roadMapOfInputData := make([]Field, 0)
	errors := make(ValidationErrors, 0)
	var i int
	for i < int(reflect.TypeOf(v).NumField()) {
		newField := Field{
			Name:   reflect.TypeOf(v).Field(i).Name,
			RawTag: reflect.TypeOf(v).Field(i).Tag.Get("validate"),
		}
		roadMapOfInputData = append(roadMapOfInputData, newField)
		i++
	}
	for _, curField := range roadMapOfInputData {
		curTags, err := parseTag(curField.RawTag)
		if err != nil {
			return nil, err
		}
		fieldContent := reflect.ValueOf(v).FieldByName(curField.Name)
		for _, curTag := range curTags {
			valErrors, err := validateDataByTag(fieldContent, curField.Name, curTag, validators)
			if err != nil {
				return nil, err
			}
			if len(valErrors) > 0 {
				errors = append(errors, valErrors...)
			}
		}

	}
	return errors, nil
}

func parseTag(RawTag string) ([]Tag, error) {
	var parsedTag Tag
	parsedTags := make([]Tag, 0)
	tags := strings.Split(RawTag, "|")
	for _, tag := range tags {
		tagData := strings.Split(tag, ":")
		if len(tagData) > 2 {
			return parsedTags, ErrBadTagFormat
		}
		parsedTag.Name = tagData[0]
		if len(tagData) > 1 {
			parsedTag.Params = strings.Split(tagData[1], ",")
		}
		parsedTags = append(parsedTags, parsedTag)
	}
	return parsedTags, nil
}

func validateDataByTag(value reflect.Value, fieldname string, tag Tag, validators map[string]Validator) (ValidationErrors, error) {
	valErrors := make(ValidationErrors, 0)
	validator, ok := validators[tag.Name+value.Type().Kind().String()]
	if !ok {
		return nil, ErrBadTypeAssertion
	}
	if value.Type().Kind() != validator.SupportedType {
		return nil, ErrUnsupportedType
	}
	err := validator.Method(value.Interface(), tag.Params)
	if err != nil {
		valError := ValidationError{
			Field: fieldname,
			Err:   err,
		}
		valErrors = append(valErrors, valError)
	}
	return valErrors, nil

}

func initValidators() map[string]Validator {
	validators := make([]Validator, 6)
	validators[0] = Validator{
		Name:          "min",
		SupportedType: reflect.Int,
		Method: func(data any, params []any) error {
			min, ok := params[0].(int)
			if !ok {
				return ErrBadTypeAssertion
			}
			dataInt, ok := data.(int)
			if !ok {
				return ErrBadTypeAssertion
			}
			if dataInt < min {
				return ErrNumLessMin
			}
			return nil
		},
	}
	validators[1] = Validator{
		Name:          "max",
		SupportedType: reflect.Int,
		Method: func(data any, params []any) error {
			max, ok := params[0].(int)
			if !ok {
				return ErrBadTypeAssertion
			}
			dataInt, ok := data.(int)
			if !ok {
				return ErrBadTypeAssertion
			}
			if dataInt > max {
				return ErrNumGreaterMax
			}
			return nil
		},
	}
	validators[2] = Validator{
		Name:          "in",
		SupportedType: reflect.Int,
		Method: func(data any, params []any) error {
			dataInt, ok := data.(int)
			if !ok {
				return ErrBadTypeAssertion
			}
			for _, curParam := range params {
				curParamInt, ok := curParam.(int)
				if !ok {
					return ErrBadTypeAssertion
				}
				if dataInt == curParamInt {
					return nil
				}
			}
			return ErrNumMNums
		},
	}
	validators[3] = Validator{
		Name:          "len",
		SupportedType: reflect.String,
		Method: func(data any, params []any) error {
			length, ok := params[0].(int)
			if !ok {
				return ErrBadTypeAssertion
			}
			dataStr, ok := data.(string)
			if !ok {
				return ErrBadTypeAssertion
			}
			if len(dataStr) == length {
				return nil
			}
			return ErrStrLen
		},
	}
	validators[4] = Validator{
		Name:          "regexp",
		SupportedType: reflect.String,
		Method: func(data any, params []any) error {
			expression, ok := params[0].(string)
			if !ok {
				return ErrBadTypeAssertion
			}
			dataStr, ok := data.(string)
			if !ok {
				return ErrBadTypeAssertion
			}
			regexp, err := regexp.Compile(expression)
			if err != nil {
				return err
			}
			if regexp.MatchString(dataStr) {
				return nil
			} else {
				return ErrStrRegexp
			}
		},
	}
	validators[5] = Validator{
		Name:          "in",
		SupportedType: reflect.String,
		Method: func(data any, params []any) error {
			dataInt, ok := data.(string)
			if !ok {
				return ErrBadTypeAssertion
			}
			for _, curParam := range params {
				curParamInt, ok := curParam.(string)
				if !ok {
					return ErrBadTypeAssertion
				}
				if dataInt == curParamInt {
					return nil
				}
			}
			return ErrStrMStrings
		},
	}
	validatorsMap := make(map[string]Validator)
	for _, curValidator := range validators {
		validatorsMap[curValidator.Name+curValidator.SupportedType.String()] = curValidator
	}
	return validatorsMap
}

func packErrors(input ValidationErrors) error {
	var errRez error
	for _, curValErr := range input {
		errRez = fmt.Errorf("Field - "+curValErr.Field+":%w", curValErr.Err)
	}
	return errRez
}
