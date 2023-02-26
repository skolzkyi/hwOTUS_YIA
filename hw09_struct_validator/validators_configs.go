package hw09structvalidator

import (
	"reflect"
	"regexp"
	"strconv"
)

var (
	funcIntMin = func(data string, params []string) error {
		min, err := strconv.Atoi(params[0])
		if err != nil {
			return err
		}
		dataInt, err := strconv.Atoi(data)
		if err != nil {
			return err
		}
		if dataInt < min {
			return ErrNumLessMin
		}
		return nil
	}
	funcIntMax = func(data string, params []string) error {
		max, err := strconv.Atoi(params[0])
		if err != nil {
			return err
		}
		dataInt, err := strconv.Atoi(data)
		if err != nil {
			return err
		}
		if dataInt > max {
			return ErrNumGreaterMax
		}
		return nil
	}
	funcInInt = func(data string, params []string) error {
		dataInt, err := strconv.Atoi(data)
		if err != nil {
			return err
		}
		for _, curParam := range params {
			curParamInt, err := strconv.Atoi(curParam)
			if err != nil {
				return err
			}
			if dataInt == curParamInt {
				return nil
			}
		}
		return ErrNumMNums
	}
	funcLenStr = func(data string, params []string) error {
		length, err := strconv.Atoi(params[0])
		if err != nil {
			return err
		}
		if len(data) == length {
			return nil
		}
		return ErrStrLen
	}
	funcRegexpStr = func(data string, params []string) error {
		expression := params[0]
		regexp, err := regexp.Compile(expression)
		if err != nil {
			return err
		}
		if regexp.MatchString(data) {
			return nil
		}
		return ErrStrRegexp
	}
	funcInStr = func(data string, params []string) error {
		for _, curParam := range params {
			if data == curParam {
				return nil
			}
		}
		return ErrStrMStrings
	}
)

func initValidators() map[string]Validator {
	validators := make([]Validator, 6)
	validators[0] = Validator{
		Name:          "min",
		SupportedType: reflect.Int,
		Method:        funcIntMin,
	}
	validators[1] = Validator{
		Name:          "max",
		SupportedType: reflect.Int,
		Method:        funcIntMax,
	}
	validators[2] = Validator{
		Name:          "in",
		SupportedType: reflect.Int,
		Method:        funcInInt,
	}
	validators[3] = Validator{
		Name:          "len",
		SupportedType: reflect.String,
		Method:        funcLenStr,
	}
	validators[4] = Validator{
		Name:          "regexp",
		SupportedType: reflect.String,
		Method:        funcRegexpStr,
	}
	validators[5] = Validator{
		Name:          "in",
		SupportedType: reflect.String,
		Method:        funcInStr,
	}
	validatorsMap := make(map[string]Validator)
	for _, curValidator := range validators {
		validatorsMap[curValidator.Name+curValidator.SupportedType.String()] = curValidator
	}

	return validatorsMap
}
