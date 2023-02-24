package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	builder := strings.Builder{}
	for _, curValErr := range v {
		builder.WriteString("[Field: ")
		builder.WriteString(curValErr.Field)
		builder.WriteString(" Error: ")
		builder.WriteString(curValErr.Err.Error())
		builder.WriteString("]; ")
	}
	return builder.String()
}

type Field struct {
	Name    string
	RawTag  string
	IsSlice bool
}

type Tag struct {
	Name   string
	Params []string
}

type Validator struct {
	Name          string
	SupportedType reflect.Kind
	Method        func(data string, params []string) error
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
	// т.к. обрабатываются только string и int  и их слайсы
	//nolint:exhaustive
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
	fmt.Println()
	return valErrors
}

func validateData(v any, validators map[string]Validator) (ValidationErrors, error) {
	roadMapOfInputData := make([]Field, 0)
	errors := make(ValidationErrors, 0)
	var i int
	for i < reflect.TypeOf(v).NumField() {
		var isSliceFlag bool
		if reflect.TypeOf(v).Field(i).Type.Kind() == reflect.Slice {
			isSliceFlag = true
		}
		newField := Field{
			Name:    reflect.TypeOf(v).Field(i).Name,
			RawTag:  reflect.TypeOf(v).Field(i).Tag.Get("validate"),
			IsSlice: isSliceFlag,
		}
		if newField.RawTag != "" {
			roadMapOfInputData = append(roadMapOfInputData, newField)
		}
		i++
	}
	for _, curField := range roadMapOfInputData {
		curTags, err := parseTag(curField.RawTag)
		if err != nil {
			return nil, err
		}
		fieldContent := reflect.ValueOf(v).FieldByName(curField.Name)
		if curField.IsSlice {
			for i := 0; i < fieldContent.Len(); i++ {
				valErrors, err := validateDataByTag(fieldContent.Index((i)), curField.Name, curTags, validators)
				if err != nil {
					return nil, err
				}
				valErrors = wrapIndex(valErrors, i)
				errors = append(errors, valErrors...)
			}
		} else {
			valErrors, err := validateDataByTag(fieldContent, curField.Name, curTags, validators)
			if err != nil {
				return nil, err
			}
			errors = append(errors, valErrors...)
		}
	}
	return errors, nil
}

func parseTag(rawTag string) ([]Tag, error) {
	var parsedTag Tag
	parsedTags := make([]Tag, 0)
	tags := strings.Split(rawTag, "|")
	for _, tag := range tags {
		tagData := strings.Split(tag, ":")
		if len(tagData) > 2 {
			return parsedTags, ErrBadTagFormat
		}
		parsedTag.Name = tagData[0]
		if len(tagData) > 1 {
			tagParams := make([]string, 0)
			paramsFromTagData := strings.Split(tagData[1], ",")
			for _, curtextParam := range paramsFromTagData {
				tagParams = append(tagParams, strings.ReplaceAll(curtextParam, " ", ""))
			}
			parsedTag.Params = tagParams
		}
		parsedTags = append(parsedTags, parsedTag)
	}
	return parsedTags, nil
}

func wrapIndex(input ValidationErrors, index int) ValidationErrors {
	for i := range input {
		input[i].Err = fmt.Errorf("Index of error data - "+strconv.Itoa((index))+":%w", input[i].Err)
	}
	return input
}

func validateDataByTag(v reflect.Value, fn string, tags []Tag, vd map[string]Validator) (ValidationErrors, error) {
	errors := make(ValidationErrors, 0)
	for _, curTag := range tags {
		valErrors, err := validateElDataByTag(v, fn, curTag, vd)
		if err != nil {
			return nil, err
		}
		if len(valErrors) > 0 {
			errors = append(errors, valErrors...)
		}
	}
	return errors, nil
}

func validateElDataByTag(v reflect.Value, fn string, tag Tag, vd map[string]Validator) (ValidationErrors, error) {
	valErrors := make(ValidationErrors, 0)
	validator, ok := vd[tag.Name+v.Type().Kind().String()]
	if !ok {
		return nil, ErrBadTypeAssertion
	}
	if v.Type().Kind() != validator.SupportedType {
		return nil, ErrUnsupportedType
	}
	var tempValue string
	// т.к. обрабатываются только string и int  и их слайсы
	//nolint:exhaustive
	switch v.Type().Kind() {
	case reflect.Int:
		tempInt := v.Int()
		tempValue = strconv.Itoa(int(tempInt))
	case reflect.String:
		tempValue = v.String()
	default:
		return nil, ErrUnsupportedType
	}
	err := validator.Method(tempValue, tag.Params)
	if err != nil {
		valError := ValidationError{
			Field: fn,
			Err:   err,
		}
		valErrors = append(valErrors, valError)
	}
	return valErrors, nil
}
