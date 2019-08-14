package convert

import (
	"reflect"
	"testing"
)

// convert slice of strings to slice of interfaces
func TestIface(t *testing.T) {
	list := []string{"a", "b", "c"}
	result := Iface(list)
	for i := 0; i < len(list); i++ {
		if reflect.ValueOf(result[i]).String() != list[i] {
			t.Fatal("the interface doesn't equal the input string")
		}
	}
}
