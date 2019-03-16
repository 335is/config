package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalString(t *testing.T) {
	s := "just a normal string"
	d := ""
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, s, d)

	s = ""
	d = "no string off my back"
	err = Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, s, d)
}

func TestUnmarshalBool(t *testing.T) {
	s := "true"
	d := false
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, true, d)

	s = "0"
	d = true
	err = Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, false, d)
}

func TestUnmarshalDuration(t *testing.T) {
	s := "1s"
	d := time.Duration(0)
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, time.Second, d)

	s = "10m"
	d = time.Duration(0)
	err = Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, 10*time.Minute, d)
}

func TestUnmarshalInt(t *testing.T) {
	s := "42"
	d := int(0)
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, int(42), d)

	s = "-127"
	d = int(0)
	err = Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, int(-127), d)

	s = "0x0"
	d = int(0)
	err = Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, int(0), d)
}

func TestUnmarshalInt8(t *testing.T) {
	s := "12"
	d := int8(0)
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, int8(12), d)
}

func TestUnmarshalInt16(t *testing.T) {
	s := "378"
	d := int16(0)
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, int16(378), d)
}

func TestUnmarshalInt32(t *testing.T) {
	s := "12345"
	d := int32(0)
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, int32(12345), d)
}

func TestUnmarshalInt64(t *testing.T) {
	s := "528012385"
	d := int64(0)
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, int64(528012385), d)
}
func TestUnmarshalUint(t *testing.T) {
	s := "77"
	d := uint(0)
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, uint(77), d)

	s = "0x0"
	d = uint(0)
	err = Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, uint(0), d)
}

func TestUnmarshalUint8(t *testing.T) {
	s := "96"
	d := uint8(0)
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, uint8(96), d)
}

func TestUnmarshalUint16(t *testing.T) {
	s := "9347"
	d := uint16(0)
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, uint16(9347), d)
}

func TestUnmarshalUint32(t *testing.T) {
	s := "93848473"
	d := uint32(0)
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, uint32(93848473), d)
}

func TestUnmarshalUint64(t *testing.T) {
	s := "1122334455667788"
	d := uint64(0)
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, uint64(1122334455667788), d)
}

func TestUnmarshalFloat32(t *testing.T) {
	s := "128.231"
	d := float32(0)
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, float32(128.231), d)
}

func TestUnmarshalFloat64(t *testing.T) {
	s := "-5280.768"
	d := float64(0)
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, float64(-5280.768), d)
}

func TestUnmarshalSliceBool(t *testing.T) {
	s := "1, T, false, 0"
	d := []bool{}
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(d))
	assert.ElementsMatch(t, []bool{true, true, false, false}, d)
}

func TestUnmarshalSliceDuration(t *testing.T) {
	s := "2h, 15m, 32m22s"
	d := []time.Duration{}
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(d))
	expected := []time.Duration{
		2 * time.Hour,
		15 * time.Minute,
		(32 * time.Minute) + (22 * time.Second),
	}
	assert.ElementsMatch(t, expected, d)
}

func TestUnmarshalSliceInt(t *testing.T) {
	s := "1, 2, 3, 4"
	d := []int{}
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(d))
	assert.ElementsMatch(t, []int{1, 2, 3, 4}, d)
}

func TestUnmarshalSliceFloat(t *testing.T) {
	s := "3.14159, 2.71828, 0.0, -0.577215, 4"
	d := []float64{}
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, 5, len(d))
	expected := []float64{
		3.14159,
		2.71828,
		0.0,
		-0.577215,
		4.0,
	}
	assert.ElementsMatch(t, expected, d)
}

func TestUnmarshalSliceString(t *testing.T) {
	s := "some string, another string, even more stringy"
	d := []string{}
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(d))
	expected := []string{
		"some string",
		"another string",
		"even more stringy",
	}
	assert.ElementsMatch(t, expected, d)
}

func TestUnmarshalMapStringBool(t *testing.T) {
	s := "first:false, second:1, third:True, fourth:FALSE"
	d := map[string]bool{}
	err := Unmarshal(s, &d)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(d))
	expected := map[string]bool{
		"first":  false,
		"second": true,
		"third":  true,
		"fourth": false,
	}
	assert.ObjectsAreEqual(expected, d)
}

// Error test cases

// nil pointer is not allowed
func TestUnmarshalNilPtr(t *testing.T) {
	s := "something"
	var d *string
	err := Unmarshal(s, d)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNilPointer, err)
	assert.Equal(t, (*string)(nil), d)
}

// non-pointer parameter is not settable
func TestUnmarshalNotSettable(t *testing.T) {
	s := "a string lives here"
	d := ""
	err := Unmarshal(s, d)
	assert.NotNil(t, err)
	assert.Equal(t, ErrCannotSetValue, err)
	assert.Equal(t, "", d)
}

// non-parseable bool string
func TestUnmarshalInvalidBool(t *testing.T) {
	s := "Maybe"
	d := false
	err := Unmarshal(s, &d)
	assert.NotNil(t, err)
	assert.Equal(t, false, d)
}

// non-parseable time duration string
func TestUnmarshalInvalidDuration(t *testing.T) {
	s := "yesterday"
	d := time.Duration(0)
	err := Unmarshal(s, &d)
	assert.NotNil(t, err)
	assert.Equal(t, time.Duration(0), d)
}

// non-parseable int string
func TestUnmarshalInvalidInt(t *testing.T) {
	s := "one hundred"
	d := int(0)
	err := Unmarshal(s, &d)
	assert.NotNil(t, err)
	assert.Equal(t, int(0), d)
}

// non-parseable uint string
func TestUnmarshalInvalidUint(t *testing.T) {
	s := "unsigned forty"
	d := uint(0)
	err := Unmarshal(s, &d)
	assert.NotNil(t, err)
	assert.Equal(t, uint(0), d)
}

// non-parseable float string
func TestUnmarshalInvalidFloat(t *testing.T) {
	s := "some precision number"
	d := float64(0.0)
	err := Unmarshal(s, &d)
	assert.NotNil(t, err)
	assert.Equal(t, float64(0.0), d)
}

// non-parseable slice string
func TestUnmarshalInvalidSlice(t *testing.T) {
	s := "1.2.3.4 scooby dont"
	d := []int{}
	err := Unmarshal(s, &d)
	assert.NotNil(t, err)
	assert.ObjectsAreEqual([]int{}, d)
}

// non-parseable map because of invalid key:value pair separator
func TestUnmarshalInvalidMapSeparator(t *testing.T) {
	s := "first:false; second:1"
	d := map[string]bool{}
	err := Unmarshal(s, &d)
	assert.NotNil(t, err)
	assert.Equal(t, ErrBadMap, err)
	assert.ObjectsAreEqual(map[string]bool{}, d)
}

// non-parseable map because of invalid key
func TestUnmarshalInvalidMapKey(t *testing.T) {
	s := "first:false"
	d := map[int]bool{}
	err := Unmarshal(s, &d)
	assert.NotNil(t, err)
	assert.ObjectsAreEqual(map[string]bool{}, d)
}

// non-parseable map because of invalid value
func TestUnmarshalInvalidMapValue(t *testing.T) {
	s := "0:starburst"
	d := map[int]bool{}
	err := Unmarshal(s, &d)
	assert.NotNil(t, err)
	assert.ObjectsAreEqual(map[string]bool{}, d)
}

// type struct is not supported
func TestUnmarshalStruct(t *testing.T) {
	s := "true, time is good, 23"
	type test struct {
		b bool
		s string
		i int
	}
	d := test{}
	err := Unmarshal(s, &d)
	assert.NotNil(t, err)
	assert.Equal(t, ErrUnsupportedType, err)
	expected := test{false, "", 0}
	assert.ElementsMatch(t, expected, d)
}
