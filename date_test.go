package abstract_test

import (
	"testing"

	"github.com/maxbolgarin/abstract"
)

func TestRange(t *testing.T) {
	testCases := []struct {
		id     string
		input1 string
		input2 string
		res    int
	}{
		{
			id:     "1",
			input1: "2020-01-01",
			input2: "2020-01-01",
			res:    0,
		},
		{
			id:     "2",
			input1: "2020-01-01",
			input2: "2020-01-02",
			res:    1,
		},
		{
			id:     "3",
			input1: "2020-01-02",
			input2: "2020-01-01",
			res:    1,
		},
		{
			id:     "4",
			input1: "2020-01-01",
			input2: "2019-01-01",
			res:    365,
		},
		{
			id:     "5",
			input1: "2020-02-01",
			input2: "2020-01-02",
			res:    30,
		},
		{
			id:     "6",
			input1: "2020-04-02",
			input2: "2020-01-01",
			res:    92,
		},
	}

	for _, test := range testCases {
		i1, err := abstract.ParseDate(test.input1)
		if err != nil {
			t.Fatal(err)
		}
		i2, err := abstract.ParseDate(test.input2)
		if err != nil {
			t.Fatal(err)
		}

		result := i1.Range(i2)
		if result != test.res {
			t.Errorf("%s -> expected %v, got %v", test.id, test.res, result)
		}
	}
}
