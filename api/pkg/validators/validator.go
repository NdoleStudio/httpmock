package validators

import (
	"context"
	"net/url"

	"github.com/thedevsaddam/govalidator"
)

type validator struct{}

func init() {}

// ValidateUUID that the payload is a UUID
func (validator *validator) ValidateUUID(_ context.Context, ID string, name string) url.Values {
	request := map[string]string{
		name: ID,
	}

	v := govalidator.New(govalidator.Options{
		Data: &request,
		Rules: govalidator.MapData{
			name: []string{
				"required",
				"uuid",
			},
		},
	})

	return v.ValidateStruct()
}
