# pint

Pint is a tiny Go library that helps you handle input from HTTP requests.

## Installation

go get -u github.com/jamal/pint

## Usage

```go
type UserRegistration struct {
    FirstName string `pint:"first_name"`
    LastName string `pint:"last_name"`
    Email string `pint:"email,format:email"`
    Password string `pint:"password"`
    Age int `pint:"age,min:13,max:199"`
}

func UserRegistrationHandler(w http.ResponseWriter, r *http.Request) {
    user := &UserRegistration{}
    err := pint.Parse(r, user)
    if err == pint.ErrValidate {
        fmt.Fprintf(w, "Validation error: %s", er.String())
        return
    }

    // ...
}
```

### Custom Format handler

You can register a custom format handler using `RegisterHandler`.

```go
type User struct {
    Name string `pint:"name"`
    Phone string `pint:"phone,format:phone"`
}

func formatPhone(val string) (string, error) {
	num, err := libphonenumber.Parse(val, "US")
	if err != nil {
		return "", &pint.ErrValidate{fmt.Sprintf("%s is not a valid phone number", val)}
	}
	return libphonenumber.Format(num, libphonenumber.INTERNATIONAL), nil
}

func init() {
    pint.RegisterHandler("phone", formatPhone)
}
```

## License

Released under the [MIT License](https://github.com/jamal/pint/blob/master/License).
