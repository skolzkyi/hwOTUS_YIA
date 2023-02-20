package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	caseName       string
	dir            string
	expEnvironment Environment
	expectedError  error
}

func TestReadDir(t *testing.T) {
	cases := createPositiveCases()
	t.Run(cases.caseName, func(t *testing.T) {
		env, err := ReadDir(cases.dir)
		require.NoError(t, err)
		require.EqualValues(t, env, cases.expEnvironment)
	})
	cases = createNegativeCases()
	t.Run(cases.caseName, func(t *testing.T) {
		_, err := ReadDir(cases.dir)
		require.Truef(t, errors.Is(err, cases.expectedError), "actual error %q", err)
	})
}

func createPositiveCases() testCase {
	testEnvironmentMap := make(Environment)
	newEnvironment := EnvValue{
		Value:      "bar",
		NeedRemove: false,
	}
	testEnvironmentMap["BAR"] = newEnvironment
	newEnvironment = EnvValue{
		Value:      "",
		NeedRemove: false,
	}
	testEnvironmentMap["EMPTY"] = newEnvironment
	newEnvironment = EnvValue{
		Value:      "   foo\nwith new line",
		NeedRemove: false,
	}
	testEnvironmentMap["FOO"] = newEnvironment
	newEnvironment = EnvValue{
		Value:      `"hello"`,
		NeedRemove: false,
	}
	testEnvironmentMap["HELLO"] = newEnvironment
	newEnvironment = EnvValue{
		Value:      "",
		NeedRemove: true,
	}
	testEnvironmentMap["UNSET"] = newEnvironment
	positiveTestCase := testCase{
		caseName:       "PositiveTestCase",
		dir:            "./testdata/env",
		expEnvironment: testEnvironmentMap,
		expectedError:  nil,
	}
	return positiveTestCase
}

func createNegativeCases() testCase {
	negativeTestCase := testCase{
		caseName:       "NegativeTestCase",
		dir:            "./testdata/env/BAR",
		expEnvironment: nil,
		expectedError:  ErrNotDirectory,
	}
	return negativeTestCase
}
