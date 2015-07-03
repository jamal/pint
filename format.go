package pint

import (
	"fmt"
	"net/mail"
)

type FormatHandler func(string) (string, error)

var formatHandlers = map[string]FormatHandler{
	"email": formatEmail,
}

// RegisterHandler can be used to register custom field validation. The
// name of the handler will be used in the struct tag for the validate key. For
// example, if you registered a "phone" validate handler, your struct tag would
// be `pint:"phone,validate:phone"`.
func RegisterHandler(name string, handler FormatHandler) {
	formatHandlers[name] = handler
}

func formatEmail(val string) (string, error) {
	address, err := mail.ParseAddress(val)
	if err != nil {
		return "", &ErrValidate{fmt.Sprintf(`"%v" is not a valid email address`, val)}
	}
	return address.Address, nil
}
