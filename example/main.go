package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jamal/pint"
)

// Login example
type Login struct {
	Username string `pint:"username"`
	Password string `pint:"password"`
}

// User example
type User struct {
	Name        string  `pint:"name"`
	Email       string  `pint:"email,format:email"`
	Age         int     `pint:"age,min:13,max:99"`
	Phone       string  `pint:"phone"`
	BigNum      int64   `pint:"big_num"`
	FloatNum    float32 `pint:"float_num"`
	UintNum     uint    `pint:"uint_num"`
	Admin       bool    `pint:"active"`
	Deleted     bool    `pint:"deleted"`
	EmptyString string  `pint:"empty_string,omitempty"`
}

func main() {
	data := url.Values{}
	data.Add("username", "john")
	data.Add("password", "password123")

	r, _ := http.NewRequest("POST", "/", bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.ParseForm()

	login := Login{}
	err := pint.Parse(r, &login)
	fmt.Println(login, err)

	data = url.Values{}
	data.Add("name", "John Doe")
	data.Add("email", "jdoe@example.com")
	data.Add("age", "13")
	data.Add("phone", "5551231234")
	data.Add("big_num", "9223372036854775807")
	data.Add("float_num", "1.23456789")
	data.Add("uint_num", "132")
	data.Add("active", "true")
	data.Add("deleted", "0")

	r, _ = http.NewRequest("POST", "/", bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.ParseForm()

	user := User{}
	err = pint.Parse(r, &user)
	fmt.Println(user, err)
}
