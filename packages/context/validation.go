package context

func (c *ctx) Validate(request interface{}) (invalid bool) {
	return c.validationErr.Validation().Validate(request)
}

func (c *ctx) FieldError(fieldName string, message string) {
	c.validationErr.Validation().Add(fieldName, message)
}

func (c *ctx) IsInValid() bool {
	return c.validationErr.Validation().Invalid()
}

func (c *ctx) ValidationError() error {
	return c.validationErr
}
