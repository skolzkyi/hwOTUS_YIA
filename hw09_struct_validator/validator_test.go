package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string //`validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:    "123456789012345678901234567890123456",
				Name:  "Antuan",
				Age:   36,
				Email: "antuan@mail.ru",
				Role:  "stuff",
				//	Phones: "89123456789",
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:    "1234567890123456789012345678901234567",
				Name:  "Antuan",
				Age:   36,
				Email: "antuan@mail.ru",
				Role:  "stuff",
				//Phones: "89123456789",
			},
			expectedErr: ErrStrLen,
		},
		{
			in: User{
				ID:    "123456789012345678901234567890123456",
				Name:  "Antuan",
				Age:   36,
				Email: "antuanmail.ru",
				Role:  "stuff",
				//	Phones: "89123456789",
			},
			expectedErr: ErrStrRegexp,
		},
		{
			in: User{
				ID:    "123456789012345678901234567890123456",
				Name:  "Antuan",
				Age:   36,
				Email: "antuan@mail.ru",
				Role:  "noname",
				//	Phones: "89123456789",
			},
			expectedErr: ErrStrMStrings,
		},
		{
			in: User{
				ID:    "123456789012345678901234567890123456",
				Name:  "Antuan",
				Age:   16,
				Email: "antuan@mail.ru",
				Role:  "stuff",
				//	Phones: "89123456789",
			},
			expectedErr: ErrNumLessMin,
		},
		{
			in: User{
				ID:    "123456789012345678901234567890123456",
				Name:  "Antuan",
				Age:   56,
				Email: "antuan@mail.ru",
				Role:  "stuff",
				//Phones: "89123456789",
			},
			expectedErr: ErrNumGreaterMax,
		},
		{
			in: App{
				Version: "1.2.2",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1.2.2.16",
			},
			expectedErr: ErrStrLen,
		},
		{
			in: Response{
				Code: 114,
				Body: "{}",
			},
			expectedErr: ErrNumMNums,
		},
		{
			in: Response{
				Code: 404,
				Body: "{}",
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			//tt := tt
			//t.Parallel()
			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Truef(t, errors.Is(err, tt.expectedErr), "actual error %q", err)
			}
			//_ = tt
		})
	}
}
