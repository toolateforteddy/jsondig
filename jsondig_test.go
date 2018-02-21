package jsondig

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestTwoLevel(t *testing.T) {
	str := `{"foo": {"bar":"baz"}}`
	var data interface{}

	json.NewDecoder(bytes.NewReader([]byte(str))).Decode(&data)

	v, err := JsonDig(data, "foo", "bar")
	if err != nil {
		t.Fatalf("Failed to get foo.bar. Error: %v", err)
		return
	}

	if v != "baz" {
		t.Fatalf("Result wasn't baz. Got %v", v)
	}
}

func TestTwoLevelError(t *testing.T) {
	str := `{"foo": {"foobar":"baz"}}`
	var data interface{}

	json.NewDecoder(bytes.NewReader([]byte(str))).Decode(&data)

	_, err := JsonDig(data, "foo", "bar")
	if err == nil {
		t.Fatal("Failed to fail....")
		return
	}

	expectedErr := `Could not find object at "foo.bar".
Found map[string]interface {}{"foobar":"baz"} at "foo".`
	if err.Error() != expectedErr {
		t.Fatalf("Badly formatted Error recieved. Got:\n%v", err)
	}
}

func TestArrayDive(t *testing.T) {
	str := `{"foo": ["foobar","baz"]}`
	var data interface{}

	json.NewDecoder(bytes.NewReader([]byte(str))).Decode(&data)

	v, err := JsonDig(data, "foo", "1")
	if err != nil {
		t.Fatalf("Failed to get foo[1]. Error: %v", err)
		return
	}

	if v != "baz" {
		t.Fatalf("Result wasn't baz. Got %v", v)
	}
}

func TestArrayError(t *testing.T) {
	str := `{"foo": ["foobar","baz"]}`
	var data interface{}

	json.NewDecoder(bytes.NewReader([]byte(str))).Decode(&data)

	_, err := JsonDig(data, "foo", "5")
	if err == nil {
		t.Fatal("Failed to fail....")
		return
	}

	expectedErr := `Could not find object at "foo.5".
Found []interface {}{"foobar", "baz"} at "foo".`
	if err.Error() != expectedErr {
		t.Fatalf("Badly formatted Error recieved. Got:\n%v", err)
	}

	_, err = JsonDig(data, "foo", "bar")
	if err == nil {
		t.Fatal("Failed to fail....")
		return
	}

	expectedErr = `strconv.Atoi: parsing "bar": invalid syntax`
	if err.Error() != expectedErr {
		t.Fatalf("Badly formatted Error recieved. Got:\n%v", err)
	}
}
