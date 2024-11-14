package validators

import (
	"net/url"

	"github.com/gofiber/fiber/v2"

	"github.com/thedevsaddam/govalidator"
)

type validator struct{}

func init() {}

// ValidateUUID that the payload is a UUID
func (validator *validator) ValidateUUID(c *fiber.Ctx, name string) url.Values {
	request := map[string]string{
		name: c.Params(name),
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
