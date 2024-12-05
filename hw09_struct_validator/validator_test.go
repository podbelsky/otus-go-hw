package hw09structvalidator

import (
	"encoding/json"
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

func TestValidateStruct(t *testing.T) {
	for _, tt := range []struct {
		name           string
		input          interface{}
		expectedErrors []error
	}{
		{
			"str len rule valid",
			App{"12345"},
			nil,
		},
		{
			"str len rule field invalid",
			App{"123456"},
			[]error{ErrStrLengthRuleIsInvalid},
		},
		{
			"int in rule with multiple values valid",
			Response{200, "123"},
			nil,
		},
		{
			"int in rule with multiple values not valid",
			Response{123123123213, "yo"},
			[]error{ErrIntInRuleIsInvalid},
		},
		{
			"no validation rules",
			Token{
				[]byte{1, 2, 3, 4, 5},
				[]byte{1, 2, 3, 4, 5},
				[]byte{1, 2, 3, 4, 5},
			},
			nil,
		},
		{
			"str len and str in rule is invalid at the same time",
			User{
				"123",
				"name",
				20,
				"kek@mail.ru",
				"not_admin",
				[]string{"89261232323", "89263212121", "89261242424"},
				[]byte{1, 2, 3, 4, 5},
			},
			[]error{ErrStrLengthRuleIsInvalid, ErrStrInRuleIsInvalid},
		},
		{
			"len and max rule not valid",
			User{
				"123",
				"name",
				600,
				"kek@mail.ru",
				"admin",
				[]string{"89261232323", "89263212121", "89261242424"},
				[]byte{1, 2, 3, 4, 5},
			},
			[]error{ErrStrLengthRuleIsInvalid, ErrIntMaxRuleIsInvalid},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := ValidateStruct(tt.input)
			if tt.expectedErrors == nil {
				require.NoError(t, err)

				return
			}

			var actualErrors ValidationErrors
			require.ErrorAs(t, err, &actualErrors)
			for i, e := range actualErrors {
				require.ErrorIs(t, e.Err, tt.expectedErrors[i])
			}
		})
	}
}
