package gomockhelpers

import (
	"testing"
)

func TestMatcherFunc(t *testing.T) {
	type FakeStruct struct {
		KeyForTest string
	}

	fakeVal := FakeStruct{
		KeyForTest: "Something",
	}

	var testMatcher MatcherFunc = (func(x interface{}) bool {
		if valForTest, ok := x.(FakeStruct); ok {
			return valForTest.KeyForTest == "Something"
		}

		return false
	})

	t.Run("Matches", func(t *testing.T) {
		if !testMatcher.Matches(fakeVal) {
			t.Error("Expected true, got false")
		}
	})

	t.Run("String", func(t *testing.T) {
		expectedDescription := "matches values captured in a func"
		actualDescription := testMatcher.String()

		if expectedDescription != actualDescription {
			t.Errorf("got %#v, expected %#v",
				actualDescription, expectedDescription)
		}
	})
}
