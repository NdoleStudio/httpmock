package validators

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"

	"github.com/thedevsaddam/govalidator"
)

type validator struct{}

const (
	requestHeaders = "requestHeaders"
	requestPath    = "requestPath"
)

func init() {
	govalidator.AddCustomRule(requestHeaders, func(field string, rule string, message string, value interface{}) error {
		input, ok := value.(string)
		if !ok {
			return fmt.Errorf("the %s field must be a string", field)
		}

		if err := json.Unmarshal([]byte(input), &[]map[string]string{}); err != nil {
			return fmt.Errorf("the %s field is not a valid JSON array with schema [{\"key\": \"value\"}]", field)
		}

		return nil
	})

	govalidator.AddCustomRule(requestPath, func(field string, rule string, message string, value interface{}) error {
		input, ok := value.(string)
		if !ok {
			return fmt.Errorf("the %s field must be a string", field)
		}

		if _, err := url.Parse("https://httpmock.dev" + input); err != nil {
			return fmt.Errorf("the %s field must be a valid URI like /post/1", field)
		}

		return nil
	})
}

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
