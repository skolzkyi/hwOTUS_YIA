package hw09structvalidator

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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
		Phones []string `validate:"len:11"`
		// meta json.RawMessage //т.к. на github этот линтер почему-то не отключается через nolint
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
				ID:     "123456789012345678901234567890123456",
				Name:   "Antuan",
				Age:    36,
				Email:  "antuan@mail.ru",
				Role:   "stuff",
				Phones: []string{"89123456781", "89123456782", "89123456783"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "1234567890123456789012345678901234567",
				Name:   "Antuan",
				Age:    36,
				Email:  "antuan@mail.ru",
				Role:   "stuff",
				Phones: []string{"89123456781", "89123456782", "89123456783"},
			},
			expectedErr: ErrStrLen,
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Antuan",
				Age:    36,
				Email:  "antuanmail.ru",
				Role:   "stuff",
				Phones: []string{"89123456781", "89123456782", "89123456783"},
			},
			expectedErr: ErrStrRegexp,
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Antuan",
				Age:    36,
				Email:  "antuan@mail.ru",
				Role:   "noname",
				Phones: []string{"89123456781", "89123456782", "89123456783"},
			},
			expectedErr: ErrStrMStrings,
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Antuan",
				Age:    16,
				Email:  "antuan@mail.ru",
				Role:   "stuff",
				Phones: []string{"89123456781", "89123456782", "89123456783"},
			},
			expectedErr: ErrNumLessMin,
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Antuan",
				Age:    56,
				Email:  "antuan@mail.ru",
				Role:   "stuff",
				Phones: []string{"89123456781", "89123456782", "89123456783"},
			},
			expectedErr: ErrNumGreaterMax,
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Antuan",
				Age:    36,
				Email:  "antuan@mail.ru",
				Role:   "stuff",
				Phones: []string{"89123456781", "8912345678200", "89123456783"},
			},
			expectedErr: ErrStrLen,
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
			tt := tt
			t.Parallel()
			// иначе не сохранить сишнатуру Validate
			//nolint:errorlint
			valErrors := Validate(tt.in).(ValidationErrors)
			if tt.expectedErr == nil {
				require.Len(t, valErrors, 0)
			} else {
				var flag bool
				for _, curValErr := range valErrors {
					if errors.Is(curValErr.Err, tt.expectedErr) {
						flag = true
						break
					}
				}
				require.Truef(t, flag, "actual errors %q", valErrors)
			}
		})
	}
}
