package pint

import (
	"bytes"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

type testStruct struct {
	String    string
	Int       int
	Int64     int64
	Uint      uint
	Uint64    uint64
	Float     float64
	Bool      bool
	OtherBool bool
}

func TestDecode(t *testing.T) {
	expected := &testStruct{
		"stringvalue",
		-1,
		9223372036854775807,
		4294967295,
		18446744073709551615,
		1.23456789,
		true,
		false,
	}

	data := url.Values{}
	data.Add("String", expected.String)
	data.Add("Int", strconv.FormatInt(int64(expected.Int), 10))
	data.Add("Int64", strconv.FormatInt(expected.Int64, 10))
	data.Add("Uint", strconv.FormatUint(uint64(expected.Uint), 10))
	data.Add("Uint64", strconv.FormatUint(expected.Uint64, 10))
	data.Add("Float", "1.23456789")
	data.Add("Bool", "true")
	data.Add("OtherBool", "0")

	r, _ := http.NewRequest("POST", "/", bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.ParseForm()

	testData := &testStruct{}
	err := Parse(r, testData)
	if err != nil {
		t.Fatalf("Parse returned an unexpected error %v", err)
	}

	if !reflect.DeepEqual(testData, expected) {
		t.Fatalf("Parsed data did not match expected:\n\tGot: %v\n\tExpected: %v", testData, expected)
	}
}
