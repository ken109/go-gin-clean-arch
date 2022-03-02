package context

func (c *ctx) Validate(request interface{}) (invalid bool) {
	return c.verr.Validation().Validate(request)
}

func (c *ctx) FieldError(fieldName string, message string) {
	c.verr.Validation().Add(fieldName, message)
}

func (c *ctx) IsInValid() bool {
	return c.verr.Validation().Invalid()
}

func (c *ctx) ValidationError() error {
	return c.verr
}
