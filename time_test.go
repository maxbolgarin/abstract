package abstract_test

import (
	"testing"
	"time"

	"github.com/maxbolgarin/abstract"
)

// текущее время (допустим 00 15) меньше времени начала дня (допустим 04 00)
// 1. 00:10 -- меньше текущего, меньше начала -- было недавно -- 2 приоритет
// 2. ??:?? -- меньше текущего, больше начала -- невозможно
// 3. 00:20 -- больше текущего, меньше начала -- скоро будет -- 3 приоритет
// 4. 04:10 -- больше текущего, больше начала -- было давно, в самом верху -- 1 приоритет

// текущее время (допустим 18 00) больше времени начала дня (допустим 04 00)
// 1. 02:00 -- меньше текущего, меньше начала -- сегодня еще будет, в самом низу -- 4 приоритет
// 2. 14:00 -- меньше текущего, больше начала -- было недавно --  2 приоритет
// 3. ??:?? -- больше текущего, меньше начала -- невозможно
// 4. 18:30 -- больше текущего, больше начала -- через 30 минут - 3 приоритет

func TestGetTimeSortingPriority(t *testing.T) {
	var (
		dayStart  = abstract.NewTime(4, 0)
		nowBefore = abstract.NewTime(0, 15)
		nowAfter  = abstract.NewTime(18, 0)
	)

	testCases := []struct {
		id           string
		entered, now abstract.Time
		result       abstract.SortingPriority
	}{
		{
			id:      "b1",
			entered: abstract.NewTime(0, 10),
			now:     nowBefore,
			result:  abstract.BeforePriority,
		},
		{
			id:     "b2",
			now:    nowBefore,
			result: abstract.BeforePriority,
		},
		{
			id:      "b3",
			entered: abstract.NewTime(0, 20),
			now:     nowBefore,
			result:  abstract.AfterPriority,
		},
		{
			id:      "b4",
			entered: abstract.NewTime(4, 10),
			now:     nowBefore,
			result:  abstract.LongAgoPriority,
		},
		{
			id:      "a1",
			entered: abstract.NewTime(2, 0),
			now:     nowAfter,
			result:  abstract.NotSoonPriority,
		},
		{
			id:      "a2",
			entered: abstract.NewTime(14, 0),
			now:     nowAfter,
			result:  abstract.BeforePriority,
		},
		{
			id:     "a3",
			now:    nowAfter,
			result: abstract.NotSoonPriority,
		},
		{
			id:      "a4",
			entered: abstract.NewTime(18, 30),
			now:     nowAfter,
			result:  abstract.AfterPriority,
		},
		{
			id:      "eq",
			entered: dayStart,
			now:     nowAfter,
			result:  abstract.BeforePriority,
		},
		{
			id:      "eq2",
			entered: dayStart,
			now:     dayStart,
			result:  abstract.BeforePriority,
		},
	}

	for _, test := range testCases {
		result := abstract.GetTimeSortingPriority(test.entered, test.now, dayStart)
		if test.result != result {
			t.Errorf("%s -> expected %d, got %d", test.id, test.result, result)
		}
	}
}

func TestParseUTCOffset(t *testing.T) {
	utcTime := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	testCases := []struct {
		id     string
		input  string
		result time.Time
		isErr  bool
	}{
		{
			id:    "1",
			input: "",
			isErr: true,
		},
		{
			id:     "2",
			input:  "0",
			result: utcTime,
		},
		{
			id:     "3",
			input:  "+0",
			result: utcTime,
		},
		{
			id:     "4",
			input:  "-0",
			result: utcTime,
		},
		{
			id:     "5",
			input:  "+0 0",
			result: utcTime,
		},
		{
			id:     "6",
			input:  "-0 0",
			result: utcTime,
		},
		{
			id:    "7",
			input: "+",
			isErr: true,
		},
		{
			id:     "8",
			input:  "1",
			result: utcTime.In(time.FixedZone("", getOffset(1, 0, 1))),
		},
		{
			id:     "9",
			input:  "+3",
			result: utcTime.In(time.FixedZone("", getOffset(3, 0, 1))),
		},
		{
			id:     "10",
			input:  "-3",
			result: utcTime.In(time.FixedZone("", getOffset(3, 0, -1))),
		},
		{
			id:     "11",
			input:  "3 30",
			result: utcTime.In(time.FixedZone("", getOffset(3, 30, 1))),
		},
		{
			id:     "12",
			input:  "+3 30",
			result: utcTime.In(time.FixedZone("", getOffset(3, 30, 1))),
		},
		{
			id:     "13",
			input:  "-3 30",
			result: utcTime.In(time.FixedZone("", getOffset(3, 30, -1))),
		},
		{
			id:     "14",
			input:  "-3:30",
			result: utcTime.In(time.FixedZone("", getOffset(3, 30, -1))),
		},
		{
			id:     "15",
			input:  "+14",
			result: utcTime.In(time.FixedZone("", getOffset(14, 0, 1))),
		},
		{
			id:    "16",
			input: "-14",
			isErr: true,
		},
		{
			id:    "17",
			input: "15",
			isErr: true,
		},
		{
			id:    "18",
			input: "3 15",
			isErr: true,
		},
		{
			id:    "19",
			input: "2 30",
			isErr: true,
		},
		{
			id:    "20",
			input: "-4 30",
			isErr: true,
		},
		{
			id:    "21",
			input: "+3 31",
			isErr: true,
		},
		{
			id:    "20",
			input: "-12 45",
			isErr: true,
		},
		{
			id:     "21",
			input:  "+12 45",
			result: utcTime.In(time.FixedZone("", getOffset(12, 45, 1))),
		},
		{
			id:    "22",
			input: "+13 45",
			isErr: true,
		},
		{
			id:    "23",
			input: "+13 45 33",
			isErr: true,
		},
		{
			id:    "24",
			input: "+d",
			isErr: true,
		},
		{
			id:    "25",
			input: "13 a",
			isErr: true,
		},
	}

	for _, test := range testCases {
		tz, err := abstract.ParseUTCOffset(test.input)
		if err != nil {
			if !test.isErr {
				t.Errorf("%s -> unexpected error %s", test.id, err)
			}
			continue
		}

		if !test.result.Equal(utcTime.In(tz)) {
			t.Errorf("%s -> expected %v, got %v", test.id, test.result, utcTime.In(tz))
		}
	}
}

func getOffset(hours, minutes, sign int) int {
	return sign*hours*60*60 + sign*minutes*60
}
