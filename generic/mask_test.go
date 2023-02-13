package generic

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenericMask(t *testing.T) {
	assert.Equal(t,
		Mask1[testStruct0](),
		[]reflect.Type{
			typeOf[testStruct0](),
		},
	)

	assert.Equal(t,
		Mask2[testStruct0, testStruct1](),
		[]reflect.Type{
			typeOf[testStruct0](),
			typeOf[testStruct1](),
		},
	)

	assert.Equal(t,
		Mask3[testStruct0, testStruct1, testStruct2](),
		[]reflect.Type{
			typeOf[testStruct0](),
			typeOf[testStruct1](),
			typeOf[testStruct2](),
		},
	)

	assert.Equal(t,
		Mask4[testStruct0, testStruct1, testStruct2, testStruct3](),
		[]reflect.Type{
			typeOf[testStruct0](),
			typeOf[testStruct1](),
			typeOf[testStruct2](),
			typeOf[testStruct3](),
		},
	)

	assert.Equal(t,
		Mask5[testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4](),
		[]reflect.Type{
			typeOf[testStruct0](),
			typeOf[testStruct1](),
			typeOf[testStruct2](),
			typeOf[testStruct3](),
			typeOf[testStruct4](),
		},
	)

	assert.Equal(t,
		Mask6[testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4, testStruct5](),
		[]reflect.Type{
			typeOf[testStruct0](),
			typeOf[testStruct1](),
			typeOf[testStruct2](),
			typeOf[testStruct3](),
			typeOf[testStruct4](),
			typeOf[testStruct5](),
		},
	)

	assert.Equal(t,
		Mask7[testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4, testStruct5, testStruct6](),
		[]reflect.Type{
			typeOf[testStruct0](),
			typeOf[testStruct1](),
			typeOf[testStruct2](),
			typeOf[testStruct3](),
			typeOf[testStruct4](),
			typeOf[testStruct5](),
			typeOf[testStruct6](),
		},
	)

	assert.Equal(t,
		Mask8[testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4, testStruct5, testStruct6, testStruct7](),
		[]reflect.Type{
			typeOf[testStruct0](),
			typeOf[testStruct1](),
			typeOf[testStruct2](),
			typeOf[testStruct3](),
			typeOf[testStruct4](),
			typeOf[testStruct5](),
			typeOf[testStruct6](),
			typeOf[testStruct7](),
		},
	)
}
