package generic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenericMask(t *testing.T) {
	assert.Equal(t,
		T[testStruct0](),
		Comp(typeOf[testStruct0]()),
	)

	assert.Equal(t,
		T1[testStruct0](),
		[]Comp{
			typeOf[testStruct0](),
		},
	)

	assert.Equal(t,
		T2[testStruct0, testStruct1](),
		[]Comp{
			typeOf[testStruct0](),
			typeOf[testStruct1](),
		},
	)

	assert.Equal(t,
		T3[testStruct0, testStruct1, testStruct2](),
		[]Comp{
			typeOf[testStruct0](),
			typeOf[testStruct1](),
			typeOf[testStruct2](),
		},
	)

	assert.Equal(t,
		T4[testStruct0, testStruct1, testStruct2, testStruct3](),
		[]Comp{
			typeOf[testStruct0](),
			typeOf[testStruct1](),
			typeOf[testStruct2](),
			typeOf[testStruct3](),
		},
	)

	assert.Equal(t,
		T5[testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4](),
		[]Comp{
			typeOf[testStruct0](),
			typeOf[testStruct1](),
			typeOf[testStruct2](),
			typeOf[testStruct3](),
			typeOf[testStruct4](),
		},
	)

	assert.Equal(t,
		T6[testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4, testStruct5](),
		[]Comp{
			typeOf[testStruct0](),
			typeOf[testStruct1](),
			typeOf[testStruct2](),
			typeOf[testStruct3](),
			typeOf[testStruct4](),
			typeOf[testStruct5](),
		},
	)

	assert.Equal(t,
		T7[testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4, testStruct5, testStruct6](),
		[]Comp{
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
		T8[testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4, testStruct5, testStruct6, testStruct7](),
		[]Comp{
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

	assert.Equal(t,
		T9[testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4, testStruct5, testStruct6, testStruct7,
			testStruct8](),
		[]Comp{
			typeOf[testStruct0](),
			typeOf[testStruct1](),
			typeOf[testStruct2](),
			typeOf[testStruct3](),
			typeOf[testStruct4](),
			typeOf[testStruct5](),
			typeOf[testStruct6](),
			typeOf[testStruct7](),
			typeOf[testStruct8](),
		},
	)

	assert.Equal(t,
		T10[testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4, testStruct5, testStruct6, testStruct7,
			testStruct8, testStruct9](),
		[]Comp{
			typeOf[testStruct0](),
			typeOf[testStruct1](),
			typeOf[testStruct2](),
			typeOf[testStruct3](),
			typeOf[testStruct4](),
			typeOf[testStruct5](),
			typeOf[testStruct6](),
			typeOf[testStruct7](),
			typeOf[testStruct8](),
			typeOf[testStruct9](),
		},
	)

	assert.Equal(t,
		T11[testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4, testStruct5, testStruct6, testStruct7,
			testStruct8, testStruct9, testStruct10](),
		[]Comp{
			typeOf[testStruct0](),
			typeOf[testStruct1](),
			typeOf[testStruct2](),
			typeOf[testStruct3](),
			typeOf[testStruct4](),
			typeOf[testStruct5](),
			typeOf[testStruct6](),
			typeOf[testStruct7](),
			typeOf[testStruct8](),
			typeOf[testStruct9](),
			typeOf[testStruct10](),
		},
	)

	assert.Equal(t,
		T12[testStruct0, testStruct1, testStruct2, testStruct3,
			testStruct4, testStruct5, testStruct6, testStruct7,
			testStruct8, testStruct9, testStruct10, testStruct11](),
		[]Comp{
			typeOf[testStruct0](),
			typeOf[testStruct1](),
			typeOf[testStruct2](),
			typeOf[testStruct3](),
			typeOf[testStruct4](),
			typeOf[testStruct5](),
			typeOf[testStruct6](),
			typeOf[testStruct7](),
			typeOf[testStruct8](),
			typeOf[testStruct9](),
			typeOf[testStruct10](),
			typeOf[testStruct11](),
		},
	)
}
