package env

import (
	"os"
	"reflect"
	"testing"
)

type Specification struct {
	Bool   bool
	Int    int
	Float  float32
	String string
}

func TestFill(t *testing.T) {
	var got Specification

	expected := Specification{
		Bool:   true,
		Int:    8080,
		Float:  0.5,
		String: "foo",
	}

	os.Clearenv()
	os.Setenv("ENV_BOOL", "true")
	os.Setenv("ENV_INT", "8080")
	os.Setenv("ENV_FLOAT", "0.5")
	os.Setenv("ENV_STRING", "foo")
	err := Fill("env", &got)
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %+v", expected)
		t.Errorf("got %+v", got)
	}
}
