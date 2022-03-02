package errors

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestValidation_Validate(t *testing.T) {
	type Profile struct {
		Age int `validate:"required,min=0"`
	}

	type User struct {
		Email   string `validate:"required"`
		Profile `inline:"true"`
	}

	type args struct {
		request interface{}
	}
	tests := []struct {
		name        string
		args        args
		wantInvalid bool
		wantError   string
	}{
		{
			name: "default",
			args: args{request: &User{
				Email:   "",
				Profile: Profile{Age: 0},
			}},
			wantInvalid: true,
			wantError:   "validation error: {\"age\":[\"ageは必須フィールドです\"],\"email\":[\"メールアドレスは必須フィールドです\"]}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verr := NewValidation()
			if gotInvalid := verr.Validation().Validate(tt.args.request); gotInvalid != tt.wantInvalid {
				t.Errorf("Validate() = %v, wantInvalid %v", gotInvalid, tt.wantInvalid)
			}
			assert.Equal(t, verr.Error(), tt.wantError)
		})
	}
}
