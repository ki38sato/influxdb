package influxql_test

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/influxdata/influxdb/influxql"
	"github.com/influxdata/influxdb/pkg/deep"
)

// Ensure that a set of iterators can be merged together, sorted by window and name/tag.
func TestMergeIterator_Float(t *testing.T) {
	inputs := []*FloatIterator{
		{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: 1},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: 3},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: 4},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: 2},
			{Name: "mem", Tags: ParseTags("host=B"), Time: 11, Value: 8},
		}},
		{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: 7},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: 5},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: 6},
			{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: 9},
		}},
		{Points: []influxql.FloatPoint{}},
		{Points: []influxql.FloatPoint{}},
	}

	itr := influxql.NewMergeIterator(FloatIterators(inputs), influxql.IteratorOptions{
		Interval: influxql.Interval{
			Duration: 10 * time.Nanosecond,
		},
		Dimensions: []string{"host"},
		Ascending:  true,
	})
	if a, err := Iterators([]influxql.Iterator{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: 1}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: 3}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: 7}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: 4}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: 2}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: 5}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: 6}},
		{&influxql.FloatPoint{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: 9}},
		{&influxql.FloatPoint{Name: "mem", Tags: ParseTags("host=B"), Time: 11, Value: 8}},
	}) {
		t.Errorf("unexpected points: %s", spew.Sdump(a))
	}

	for i, input := range inputs {
		if !input.Closed {
			t.Errorf("iterator %d not closed", i)
		}
	}
}

// Ensure that a set of iterators can be merged together, sorted by window and name/tag.
func TestMergeIterator_Integer(t *testing.T) {
	inputs := []*IntegerIterator{
		{Points: []influxql.IntegerPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: 1},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: 3},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: 4},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: 2},
			{Name: "mem", Tags: ParseTags("host=B"), Time: 11, Value: 8},
		}},
		{Points: []influxql.IntegerPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: 7},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: 5},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: 6},
			{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: 9},
		}},
		{Points: []influxql.IntegerPoint{}},
	}
	itr := influxql.NewMergeIterator(IntegerIterators(inputs), influxql.IteratorOptions{
		Interval: influxql.Interval{
			Duration: 10 * time.Nanosecond,
		},
		Dimensions: []string{"host"},
		Ascending:  true,
	})

	if a, err := Iterators([]influxql.Iterator{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: 1}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: 3}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: 7}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: 4}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: 2}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: 5}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: 6}},
		{&influxql.IntegerPoint{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: 9}},
		{&influxql.IntegerPoint{Name: "mem", Tags: ParseTags("host=B"), Time: 11, Value: 8}},
	}) {
		t.Errorf("unexpected points: %s", spew.Sdump(a))
	}

	for i, input := range inputs {
		if !input.Closed {
			t.Errorf("iterator %d not closed", i)
		}
	}
}

// Ensure that a set of iterators can be merged together, sorted by window and name/tag.
func TestMergeIterator_String(t *testing.T) {
	inputs := []*StringIterator{
		{Points: []influxql.StringPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: "a"},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: "c"},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: "d"},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: "b"},
			{Name: "mem", Tags: ParseTags("host=B"), Time: 11, Value: "h"},
		}},
		{Points: []influxql.StringPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: "g"},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: "e"},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: "f"},
			{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: "i"},
		}},
		{Points: []influxql.StringPoint{}},
	}
	itr := influxql.NewMergeIterator(StringIterators(inputs), influxql.IteratorOptions{
		Interval: influxql.Interval{
			Duration: 10 * time.Nanosecond,
		},
		Dimensions: []string{"host"},
		Ascending:  true,
	})

	if a, err := Iterators([]influxql.Iterator{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: "a"}},
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: "c"}},
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: "g"}},
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: "d"}},
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: "b"}},
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: "e"}},
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: "f"}},
		{&influxql.StringPoint{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: "i"}},
		{&influxql.StringPoint{Name: "mem", Tags: ParseTags("host=B"), Time: 11, Value: "h"}},
	}) {
		t.Errorf("unexpected points: %s", spew.Sdump(a))
	}

	for i, input := range inputs {
		if !input.Closed {
			t.Errorf("iterator %d not closed", i)
		}
	}
}

// Ensure that a set of iterators can be merged together, sorted by window and name/tag.
func TestMergeIterator_Boolean(t *testing.T) {
	inputs := []*BooleanIterator{
		{Points: []influxql.BooleanPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: true},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: true},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: false},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: false},
			{Name: "mem", Tags: ParseTags("host=B"), Time: 11, Value: true},
		}},
		{Points: []influxql.BooleanPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: true},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: true},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: false},
			{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: false},
		}},
		{Points: []influxql.BooleanPoint{}},
	}
	itr := influxql.NewMergeIterator(BooleanIterators(inputs), influxql.IteratorOptions{
		Interval: influxql.Interval{
			Duration: 10 * time.Nanosecond,
		},
		Dimensions: []string{"host"},
		Ascending:  true,
	})

	if a, err := Iterators([]influxql.Iterator{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.BooleanPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: true}},
		{&influxql.BooleanPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: true}},
		{&influxql.BooleanPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: true}},
		{&influxql.BooleanPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: false}},
		{&influxql.BooleanPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: false}},
		{&influxql.BooleanPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: true}},
		{&influxql.BooleanPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: false}},
		{&influxql.BooleanPoint{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: false}},
		{&influxql.BooleanPoint{Name: "mem", Tags: ParseTags("host=B"), Time: 11, Value: true}},
	}) {
		t.Errorf("unexpected points: %s", spew.Sdump(a))
	}

	for i, input := range inputs {
		if !input.Closed {
			t.Errorf("iterator %d not closed", i)
		}
	}
}

func TestMergeIterator_Nil(t *testing.T) {
	itr := influxql.NewMergeIterator([]influxql.Iterator{nil}, influxql.IteratorOptions{})
	if itr != nil {
		t.Fatalf("unexpected iterator: %#v", itr)
	}
}

func TestMergeIterator_Cast_Float(t *testing.T) {
	inputs := []influxql.Iterator{
		&IntegerIterator{Points: []influxql.IntegerPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: 1},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: 3},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: 4},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: 2},
			{Name: "mem", Tags: ParseTags("host=B"), Time: 11, Value: 8},
		}},
		&FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: 7},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: 5},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: 6},
			{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: 9},
		}},
	}

	itr := influxql.NewMergeIterator(inputs, influxql.IteratorOptions{
		Interval: influxql.Interval{
			Duration: 10 * time.Nanosecond,
		},
		Dimensions: []string{"host"},
		Ascending:  true,
	})
	if a, err := Iterators([]influxql.Iterator{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: 1}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: 3}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: 7}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: 4}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: 2}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: 5}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: 6}},
		{&influxql.FloatPoint{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: 9}},
		{&influxql.FloatPoint{Name: "mem", Tags: ParseTags("host=B"), Time: 11, Value: 8}},
	}) {
		t.Errorf("unexpected points: %s", spew.Sdump(a))
	}

	for i, input := range inputs {
		switch input := input.(type) {
		case *FloatIterator:
			if !input.Closed {
				t.Errorf("iterator %d not closed", i)
			}
		case *IntegerIterator:
			if !input.Closed {
				t.Errorf("iterator %d not closed", i)
			}
		}
	}
}

// Ensure that a set of iterators can be merged together, sorted by name/tag.
func TestSortedMergeIterator_Float(t *testing.T) {
	inputs := []*FloatIterator{
		{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: 1},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: 3},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: 4},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: 2},
			{Name: "mem", Tags: ParseTags("host=B"), Time: 4, Value: 8},
		}},
		{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: 7},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: 5},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: 6},
			{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: 9},
		}},
		{Points: []influxql.FloatPoint{}},
	}
	itr := influxql.NewSortedMergeIterator(FloatIterators(inputs), influxql.IteratorOptions{
		Interval: influxql.Interval{
			Duration: 10 * time.Nanosecond,
		},
		Dimensions: []string{"host"},
		Ascending:  true,
	})
	if a, err := Iterators([]influxql.Iterator{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: 1}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: 3}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: 7}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: 4}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: 2}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: 5}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: 6}},
		{&influxql.FloatPoint{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: 9}},
		{&influxql.FloatPoint{Name: "mem", Tags: ParseTags("host=B"), Time: 4, Value: 8}},
	}) {
		t.Errorf("unexpected points: %s", spew.Sdump(a))
	}

	for i, input := range inputs {
		if !input.Closed {
			t.Errorf("iterator %d not closed", i)
		}
	}
}

// Ensure that a set of iterators can be merged together, sorted by name/tag.
func TestSortedMergeIterator_Integer(t *testing.T) {
	inputs := []*IntegerIterator{
		{Points: []influxql.IntegerPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: 1},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: 3},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: 4},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: 2},
			{Name: "mem", Tags: ParseTags("host=B"), Time: 4, Value: 8},
		}},
		{Points: []influxql.IntegerPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: 7},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: 5},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: 6},
			{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: 9},
		}},
		{Points: []influxql.IntegerPoint{}},
	}
	itr := influxql.NewSortedMergeIterator(IntegerIterators(inputs), influxql.IteratorOptions{
		Interval: influxql.Interval{
			Duration: 10 * time.Nanosecond,
		},
		Dimensions: []string{"host"},
		Ascending:  true,
	})
	if a, err := Iterators([]influxql.Iterator{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: 1}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: 3}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: 7}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: 4}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: 2}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: 5}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: 6}},
		{&influxql.IntegerPoint{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: 9}},
		{&influxql.IntegerPoint{Name: "mem", Tags: ParseTags("host=B"), Time: 4, Value: 8}},
	}) {
		t.Errorf("unexpected points: %s", spew.Sdump(a))
	}

	for i, input := range inputs {
		if !input.Closed {
			t.Errorf("iterator %d not closed", i)
		}
	}
}

// Ensure that a set of iterators can be merged together, sorted by name/tag.
func TestSortedMergeIterator_String(t *testing.T) {
	inputs := []*StringIterator{
		{Points: []influxql.StringPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: "a"},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: "c"},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: "d"},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: "b"},
			{Name: "mem", Tags: ParseTags("host=B"), Time: 4, Value: "h"},
		}},
		{Points: []influxql.StringPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: "g"},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: "e"},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: "f"},
			{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: "i"},
		}},
		{Points: []influxql.StringPoint{}},
	}
	itr := influxql.NewSortedMergeIterator(StringIterators(inputs), influxql.IteratorOptions{
		Interval: influxql.Interval{
			Duration: 10 * time.Nanosecond,
		},
		Dimensions: []string{"host"},
		Ascending:  true,
	})
	if a, err := Iterators([]influxql.Iterator{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: "a"}},
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: "c"}},
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: "g"}},
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: "d"}},
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: "b"}},
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: "e"}},
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: "f"}},
		{&influxql.StringPoint{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: "i"}},
		{&influxql.StringPoint{Name: "mem", Tags: ParseTags("host=B"), Time: 4, Value: "h"}},
	}) {
		t.Errorf("unexpected points: %s", spew.Sdump(a))
	}

	for i, input := range inputs {
		if !input.Closed {
			t.Errorf("iterator %d not closed", i)
		}
	}
}

// Ensure that a set of iterators can be merged together, sorted by name/tag.
func TestSortedMergeIterator_Boolean(t *testing.T) {
	inputs := []*BooleanIterator{
		{Points: []influxql.BooleanPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: true},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: true},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: false},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: false},
			{Name: "mem", Tags: ParseTags("host=B"), Time: 4, Value: true},
		}},
		{Points: []influxql.BooleanPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: true},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: true},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: false},
			{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: true},
		}},
		{Points: []influxql.BooleanPoint{}},
	}
	itr := influxql.NewSortedMergeIterator(BooleanIterators(inputs), influxql.IteratorOptions{
		Interval: influxql.Interval{
			Duration: 10 * time.Nanosecond,
		},
		Dimensions: []string{"host"},
		Ascending:  true,
	})
	if a, err := Iterators([]influxql.Iterator{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.BooleanPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: true}},
		{&influxql.BooleanPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: true}},
		{&influxql.BooleanPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: true}},
		{&influxql.BooleanPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: false}},
		{&influxql.BooleanPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: false}},
		{&influxql.BooleanPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: true}},
		{&influxql.BooleanPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: false}},
		{&influxql.BooleanPoint{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: true}},
		{&influxql.BooleanPoint{Name: "mem", Tags: ParseTags("host=B"), Time: 4, Value: true}},
	}) {
		t.Errorf("unexpected points: %s", spew.Sdump(a))
	}

	for i, input := range inputs {
		if !input.Closed {
			t.Errorf("iterator %d not closed", i)
		}
	}
}

func TestSortedMergeIterator_Nil(t *testing.T) {
	itr := influxql.NewSortedMergeIterator([]influxql.Iterator{nil}, influxql.IteratorOptions{})
	if itr != nil {
		t.Fatalf("unexpected iterator: %#v", itr)
	}
}

func TestSortedMergeIterator_Cast_Float(t *testing.T) {
	inputs := []influxql.Iterator{
		&IntegerIterator{Points: []influxql.IntegerPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: 1},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: 3},
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: 4},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: 2},
			{Name: "mem", Tags: ParseTags("host=B"), Time: 4, Value: 8},
		}},
		&FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: 7},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: 5},
			{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: 6},
			{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: 9},
		}},
	}

	itr := influxql.NewSortedMergeIterator(inputs, influxql.IteratorOptions{
		Interval: influxql.Interval{
			Duration: 10 * time.Nanosecond,
		},
		Dimensions: []string{"host"},
		Ascending:  true,
	})
	if a, err := Iterators([]influxql.Iterator{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: 1}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 12, Value: 3}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 20, Value: 7}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30, Value: 4}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 1, Value: 2}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 11, Value: 5}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 13, Value: 6}},
		{&influxql.FloatPoint{Name: "mem", Tags: ParseTags("host=A"), Time: 25, Value: 9}},
		{&influxql.FloatPoint{Name: "mem", Tags: ParseTags("host=B"), Time: 4, Value: 8}},
	}) {
		t.Errorf("unexpected points: %s", spew.Sdump(a))
	}

	for i, input := range inputs {
		switch input := input.(type) {
		case *FloatIterator:
			if !input.Closed {
				t.Errorf("iterator %d not closed", i)
			}
		case *IntegerIterator:
			if !input.Closed {
				t.Errorf("iterator %d not closed", i)
			}
		}
	}
}

// Ensure limit iterators work with limit and offset.
func TestLimitIterator_Float(t *testing.T) {
	input := &FloatIterator{Points: []influxql.FloatPoint{
		{Name: "cpu", Time: 0, Value: 1},
		{Name: "cpu", Time: 5, Value: 3},
		{Name: "cpu", Time: 10, Value: 5},
		{Name: "mem", Time: 5, Value: 3},
		{Name: "mem", Time: 7, Value: 8},
	}}

	itr := influxql.NewLimitIterator(input, influxql.IteratorOptions{
		Limit:  1,
		Offset: 1,
	})

	if a, err := Iterators([]influxql.Iterator{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Time: 5, Value: 3}},
		{&influxql.FloatPoint{Name: "mem", Time: 7, Value: 8}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}

	if !input.Closed {
		t.Error("iterator not closed")
	}
}

// Ensure limit iterators work with limit and offset.
func TestLimitIterator_Integer(t *testing.T) {
	input := &IntegerIterator{Points: []influxql.IntegerPoint{
		{Name: "cpu", Time: 0, Value: 1},
		{Name: "cpu", Time: 5, Value: 3},
		{Name: "cpu", Time: 10, Value: 5},
		{Name: "mem", Time: 5, Value: 3},
		{Name: "mem", Time: 7, Value: 8},
	}}

	itr := influxql.NewLimitIterator(input, influxql.IteratorOptions{
		Limit:  1,
		Offset: 1,
	})

	if a, err := Iterators([]influxql.Iterator{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.IntegerPoint{Name: "cpu", Time: 5, Value: 3}},
		{&influxql.IntegerPoint{Name: "mem", Time: 7, Value: 8}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}

	if !input.Closed {
		t.Error("iterator not closed")
	}
}

// Ensure limit iterators work with limit and offset.
func TestLimitIterator_String(t *testing.T) {
	input := &StringIterator{Points: []influxql.StringPoint{
		{Name: "cpu", Time: 0, Value: "a"},
		{Name: "cpu", Time: 5, Value: "b"},
		{Name: "cpu", Time: 10, Value: "c"},
		{Name: "mem", Time: 5, Value: "d"},
		{Name: "mem", Time: 7, Value: "e"},
	}}

	itr := influxql.NewLimitIterator(input, influxql.IteratorOptions{
		Limit:  1,
		Offset: 1,
	})

	if a, err := Iterators([]influxql.Iterator{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.StringPoint{Name: "cpu", Time: 5, Value: "b"}},
		{&influxql.StringPoint{Name: "mem", Time: 7, Value: "e"}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}

	if !input.Closed {
		t.Error("iterator not closed")
	}
}

// Ensure limit iterators work with limit and offset.
func TestLimitIterator_Boolean(t *testing.T) {
	input := &BooleanIterator{Points: []influxql.BooleanPoint{
		{Name: "cpu", Time: 0, Value: true},
		{Name: "cpu", Time: 5, Value: false},
		{Name: "cpu", Time: 10, Value: true},
		{Name: "mem", Time: 5, Value: false},
		{Name: "mem", Time: 7, Value: true},
	}}

	itr := influxql.NewLimitIterator(input, influxql.IteratorOptions{
		Limit:  1,
		Offset: 1,
	})

	if a, err := Iterators([]influxql.Iterator{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.BooleanPoint{Name: "cpu", Time: 5, Value: false}},
		{&influxql.BooleanPoint{Name: "mem", Time: 7, Value: true}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}

	if !input.Closed {
		t.Error("iterator not closed")
	}
}

// Ensure auxilary iterators can be created for auxilary fields.
func TestFloatAuxIterator(t *testing.T) {
	itr := influxql.NewAuxIterator(
		&FloatIterator{Points: []influxql.FloatPoint{
			{Time: 0, Value: 1, Aux: []interface{}{float64(100), float64(200)}},
			{Time: 1, Value: 2, Aux: []interface{}{float64(500), math.NaN()}},
		}},
		influxql.IteratorOptions{Aux: []influxql.VarRef{{Val: "f0", Type: influxql.Float}, {Val: "f1", Type: influxql.Float}}},
	)

	itrs := []influxql.Iterator{
		itr,
		itr.Iterator("f0", influxql.Unknown),
		itr.Iterator("f1", influxql.Unknown),
		itr.Iterator("f0", influxql.Unknown),
	}
	itr.Start()

	if a, err := Iterators(itrs).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{
			&influxql.FloatPoint{Time: 0, Value: 1, Aux: []interface{}{float64(100), float64(200)}},
			&influxql.FloatPoint{Time: 0, Value: float64(100)},
			&influxql.FloatPoint{Time: 0, Value: float64(200)},
			&influxql.FloatPoint{Time: 0, Value: float64(100)},
		},
		{
			&influxql.FloatPoint{Time: 1, Value: 2, Aux: []interface{}{float64(500), math.NaN()}},
			&influxql.FloatPoint{Time: 1, Value: float64(500)},
			&influxql.FloatPoint{Time: 1, Value: math.NaN()},
			&influxql.FloatPoint{Time: 1, Value: float64(500)},
		},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure limit iterator returns a subset of points.
func TestLimitIterator(t *testing.T) {
	itr := influxql.NewLimitIterator(
		&FloatIterator{Points: []influxql.FloatPoint{
			{Time: 0, Value: 0},
			{Time: 1, Value: 1},
			{Time: 2, Value: 2},
			{Time: 3, Value: 3},
		}},
		influxql.IteratorOptions{
			Limit:     2,
			Offset:    1,
			StartTime: influxql.MinTime,
			EndTime:   influxql.MaxTime,
		},
	)

	if a, err := (Iterators{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Time: 1, Value: 1}},
		{&influxql.FloatPoint{Time: 2, Value: 2}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

func TestFillIterator_DST_Start_GroupByDay_Ascending(t *testing.T) {
	start := time.Date(2000, 4, 1, 0, 0, 0, 0, LosAngeles)
	end := time.Date(2000, 4, 5, 0, 0, 0, 0, LosAngeles).Add(-time.Nanosecond)
	itr := influxql.NewFillIterator(
		&FloatIterator{Points: []influxql.FloatPoint{{Time: start.UnixNano(), Value: 0}}},
		nil,
		influxql.IteratorOptions{
			StartTime: start.UnixNano(),
			EndTime:   end.UnixNano(),
			Interval: influxql.Interval{
				Duration: 24 * time.Hour,
			},
			Location:  LosAngeles,
			Ascending: true,
		},
	)

	if a, err := (Iterators{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Time: start.UnixNano(), Value: 0}},
		{&influxql.FloatPoint{Time: start.Add(24 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.Add(47 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.Add(71 * time.Hour).UnixNano(), Nil: true}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

func TestFillIterator_DST_Start_GroupByDay_Descending(t *testing.T) {
	start := time.Date(2000, 4, 1, 0, 0, 0, 0, LosAngeles)
	end := time.Date(2000, 4, 5, 0, 0, 0, 0, LosAngeles).Add(-time.Nanosecond)
	itr := influxql.NewFillIterator(
		&FloatIterator{Points: []influxql.FloatPoint{{Time: start.UnixNano(), Value: 0}}},
		nil,
		influxql.IteratorOptions{
			StartTime: start.UnixNano(),
			EndTime:   end.UnixNano(),
			Interval: influxql.Interval{
				Duration: 24 * time.Hour,
			},
			Location:  LosAngeles,
			Ascending: false,
		},
	)

	if a, err := (Iterators{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Time: start.Add(71 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.Add(47 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.Add(24 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.UnixNano(), Value: 0}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

func TestFillIterator_DST_Start_GroupByHour_Ascending(t *testing.T) {
	start := time.Date(2000, 4, 2, 0, 0, 0, 0, LosAngeles)
	end := start.Add(4*time.Hour - time.Nanosecond)
	itr := influxql.NewFillIterator(
		&FloatIterator{Points: []influxql.FloatPoint{{Time: start.UnixNano(), Value: 0}}},
		nil,
		influxql.IteratorOptions{
			StartTime: start.UnixNano(),
			EndTime:   end.UnixNano(),
			Interval: influxql.Interval{
				Duration: time.Hour,
			},
			Location:  LosAngeles,
			Ascending: true,
		},
	)

	if a, err := (Iterators{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Time: start.UnixNano(), Value: 0}},
		{&influxql.FloatPoint{Time: start.Add(1 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.Add(2 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.Add(3 * time.Hour).UnixNano(), Nil: true}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

func TestFillIterator_DST_Start_GroupByHour_Descending(t *testing.T) {
	start := time.Date(2000, 4, 2, 0, 0, 0, 0, LosAngeles)
	end := start.Add(4*time.Hour - time.Nanosecond)
	itr := influxql.NewFillIterator(
		&FloatIterator{Points: []influxql.FloatPoint{{Time: start.UnixNano(), Value: 0}}},
		nil,
		influxql.IteratorOptions{
			StartTime: start.UnixNano(),
			EndTime:   end.UnixNano(),
			Interval: influxql.Interval{
				Duration: time.Hour,
			},
			Location:  LosAngeles,
			Ascending: false,
		},
	)

	if a, err := (Iterators{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Time: start.Add(3 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.Add(2 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.Add(1 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.UnixNano(), Value: 0}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

func TestFillIterator_DST_End_GroupByDay_Ascending(t *testing.T) {
	start := time.Date(2000, 10, 28, 0, 0, 0, 0, LosAngeles)
	end := time.Date(2000, 11, 1, 0, 0, 0, 0, LosAngeles).Add(-time.Nanosecond)
	itr := influxql.NewFillIterator(
		&FloatIterator{Points: []influxql.FloatPoint{{Time: start.UnixNano(), Value: 0}}},
		nil,
		influxql.IteratorOptions{
			StartTime: start.UnixNano(),
			EndTime:   end.UnixNano(),
			Interval: influxql.Interval{
				Duration: 24 * time.Hour,
			},
			Location:  LosAngeles,
			Ascending: true,
		},
	)

	if a, err := (Iterators{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Time: start.UnixNano(), Value: 0}},
		{&influxql.FloatPoint{Time: start.Add(24 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.Add(49 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.Add(73 * time.Hour).UnixNano(), Nil: true}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

func TestFillIterator_DST_End_GroupByDay_Descending(t *testing.T) {
	start := time.Date(2000, 10, 28, 0, 0, 0, 0, LosAngeles)
	end := time.Date(2000, 11, 1, 0, 0, 0, 0, LosAngeles).Add(-time.Nanosecond)
	itr := influxql.NewFillIterator(
		&FloatIterator{Points: []influxql.FloatPoint{{Time: start.UnixNano(), Value: 0}}},
		nil,
		influxql.IteratorOptions{
			StartTime: start.UnixNano(),
			EndTime:   end.UnixNano(),
			Interval: influxql.Interval{
				Duration: 24 * time.Hour,
			},
			Location:  LosAngeles,
			Ascending: false,
		},
	)

	if a, err := (Iterators{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Time: start.Add(73 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.Add(49 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.Add(24 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.UnixNano(), Value: 0}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

func TestFillIterator_DST_End_GroupByHour_Ascending(t *testing.T) {
	start := time.Date(2000, 10, 29, 0, 0, 0, 0, LosAngeles)
	end := start.Add(4*time.Hour - time.Nanosecond)
	itr := influxql.NewFillIterator(
		&FloatIterator{Points: []influxql.FloatPoint{{Time: start.UnixNano(), Value: 0}}},
		nil,
		influxql.IteratorOptions{
			StartTime: start.UnixNano(),
			EndTime:   end.UnixNano(),
			Interval: influxql.Interval{
				Duration: time.Hour,
			},
			Location:  LosAngeles,
			Ascending: true,
		},
	)

	if a, err := (Iterators{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Time: start.UnixNano(), Value: 0}},
		{&influxql.FloatPoint{Time: start.Add(1 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.Add(2 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.Add(3 * time.Hour).UnixNano(), Nil: true}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

func TestFillIterator_DST_End_GroupByHour_Descending(t *testing.T) {
	start := time.Date(2000, 10, 29, 0, 0, 0, 0, LosAngeles)
	end := start.Add(4*time.Hour - time.Nanosecond)
	itr := influxql.NewFillIterator(
		&FloatIterator{Points: []influxql.FloatPoint{{Time: start.UnixNano(), Value: 0}}},
		nil,
		influxql.IteratorOptions{
			StartTime: start.UnixNano(),
			EndTime:   end.UnixNano(),
			Interval: influxql.Interval{
				Duration: time.Hour,
			},
			Location:  LosAngeles,
			Ascending: false,
		},
	)

	if a, err := (Iterators{itr}).ReadAll(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Time: start.Add(3 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.Add(2 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.Add(1 * time.Hour).UnixNano(), Nil: true}},
		{&influxql.FloatPoint{Time: start.UnixNano(), Value: 0}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Iterators is a test wrapper for iterators.
type Iterators []influxql.Iterator

// Next returns the next value from each iterator.
// Returns nil if any iterator returns a nil.
func (itrs Iterators) Next() ([]influxql.Point, error) {
	a := make([]influxql.Point, len(itrs))
	for i, itr := range itrs {
		switch itr := itr.(type) {
		case influxql.FloatIterator:
			fp, err := itr.Next()
			if fp == nil || err != nil {
				return nil, err
			}
			a[i] = fp
		case influxql.IntegerIterator:
			ip, err := itr.Next()
			if ip == nil || err != nil {
				return nil, err
			}
			a[i] = ip
		case influxql.StringIterator:
			sp, err := itr.Next()
			if sp == nil || err != nil {
				return nil, err
			}
			a[i] = sp
		case influxql.BooleanIterator:
			bp, err := itr.Next()
			if bp == nil || err != nil {
				return nil, err
			}
			a[i] = bp
		default:
			panic(fmt.Sprintf("iterator type not supported: %T", itr))
		}
	}
	return a, nil
}

// ReadAll reads all points from all iterators.
func (itrs Iterators) ReadAll() ([][]influxql.Point, error) {
	var a [][]influxql.Point

	// Read from every iterator until a nil is encountered.
	for {
		points, err := itrs.Next()
		if err != nil {
			return nil, err
		} else if points == nil {
			break
		}
		a = append(a, influxql.Points(points).Clone())
	}

	// Close all iterators.
	influxql.Iterators(itrs).Close()

	return a, nil
}

func TestIteratorOptions_Window_Interval(t *testing.T) {
	opt := influxql.IteratorOptions{
		Interval: influxql.Interval{
			Duration: 10,
		},
	}

	start, end := opt.Window(4)
	if start != 0 {
		t.Errorf("expected start to be 0, got %d", start)
	}
	if end != 10 {
		t.Errorf("expected end to be 10, got %d", end)
	}
}

func TestIteratorOptions_Window_Offset(t *testing.T) {
	opt := influxql.IteratorOptions{
		Interval: influxql.Interval{
			Duration: 10,
			Offset:   8,
		},
	}

	start, end := opt.Window(14)
	if start != 8 {
		t.Errorf("expected start to be 8, got %d", start)
	}
	if end != 18 {
		t.Errorf("expected end to be 18, got %d", end)
	}
}

func TestIteratorOptions_Window_Default(t *testing.T) {
	opt := influxql.IteratorOptions{
		StartTime: 0,
		EndTime:   60,
	}

	start, end := opt.Window(34)
	if start != 0 {
		t.Errorf("expected start to be 0, got %d", start)
	}
	if end != 61 {
		t.Errorf("expected end to be 61, got %d", end)
	}
}

func TestIteratorOptions_Window_Location(t *testing.T) {
	now := time.Date(2000, 4, 2, 12, 14, 15, 0, LosAngeles)
	opt := influxql.IteratorOptions{
		Location: LosAngeles,
		Interval: influxql.Interval{
			Duration: 24 * time.Hour,
		},
	}

	start, end := opt.Window(now.UnixNano())
	if exp := time.Date(2000, 4, 2, 0, 0, 0, 0, LosAngeles).UnixNano(); start != exp {
		t.Errorf("expected start to be %d, got %d", exp, start)
	}
	if exp := time.Date(2000, 4, 3, 0, 0, 0, 0, LosAngeles).UnixNano(); end != exp {
		t.Errorf("expected end to be %d, got %d", exp, end)
	}
	if got, exp := time.Duration(end-start), 23*time.Hour; got != exp {
		t.Errorf("expected duration to be %s, got %s", exp, got)
	}
}

func TestIteratorOptions_SeekTime_Ascending(t *testing.T) {
	opt := influxql.IteratorOptions{
		StartTime: 30,
		EndTime:   60,
		Ascending: true,
	}

	time := opt.SeekTime()
	if time != 30 {
		t.Errorf("expected time to be 30, got %d", time)
	}
}

func TestIteratorOptions_SeekTime_Descending(t *testing.T) {
	opt := influxql.IteratorOptions{
		StartTime: 30,
		EndTime:   60,
		Ascending: false,
	}

	time := opt.SeekTime()
	if time != 60 {
		t.Errorf("expected time to be 60, got %d", time)
	}
}

func TestIteratorOptions_DerivativeInterval_Default(t *testing.T) {
	opt := influxql.IteratorOptions{}
	expected := influxql.Interval{Duration: time.Second}
	actual := opt.DerivativeInterval()
	if actual != expected {
		t.Errorf("expected derivative interval to be %v, got %v", expected, actual)
	}
}

func TestIteratorOptions_DerivativeInterval_GroupBy(t *testing.T) {
	opt := influxql.IteratorOptions{
		Interval: influxql.Interval{
			Duration: 10,
			Offset:   2,
		},
	}
	expected := influxql.Interval{Duration: 10}
	actual := opt.DerivativeInterval()
	if actual != expected {
		t.Errorf("expected derivative interval to be %v, got %v", expected, actual)
	}
}

func TestIteratorOptions_DerivativeInterval_Call(t *testing.T) {
	opt := influxql.IteratorOptions{
		Expr: &influxql.Call{
			Name: "mean",
			Args: []influxql.Expr{
				&influxql.VarRef{Val: "value"},
				&influxql.DurationLiteral{Val: 2 * time.Second},
			},
		},
		Interval: influxql.Interval{
			Duration: 10,
			Offset:   2,
		},
	}
	expected := influxql.Interval{Duration: 2 * time.Second}
	actual := opt.DerivativeInterval()
	if actual != expected {
		t.Errorf("expected derivative interval to be %v, got %v", expected, actual)
	}
}

func TestIteratorOptions_ElapsedInterval_Default(t *testing.T) {
	opt := influxql.IteratorOptions{}
	expected := influxql.Interval{Duration: time.Nanosecond}
	actual := opt.ElapsedInterval()
	if actual != expected {
		t.Errorf("expected elapsed interval to be %v, got %v", expected, actual)
	}
}

func TestIteratorOptions_ElapsedInterval_GroupBy(t *testing.T) {
	opt := influxql.IteratorOptions{
		Interval: influxql.Interval{
			Duration: 10,
			Offset:   2,
		},
	}
	expected := influxql.Interval{Duration: time.Nanosecond}
	actual := opt.ElapsedInterval()
	if actual != expected {
		t.Errorf("expected elapsed interval to be %v, got %v", expected, actual)
	}
}

func TestIteratorOptions_ElapsedInterval_Call(t *testing.T) {
	opt := influxql.IteratorOptions{
		Expr: &influxql.Call{
			Name: "mean",
			Args: []influxql.Expr{
				&influxql.VarRef{Val: "value"},
				&influxql.DurationLiteral{Val: 2 * time.Second},
			},
		},
		Interval: influxql.Interval{
			Duration: 10,
			Offset:   2,
		},
	}
	expected := influxql.Interval{Duration: 2 * time.Second}
	actual := opt.ElapsedInterval()
	if actual != expected {
		t.Errorf("expected elapsed interval to be %v, got %v", expected, actual)
	}
}

// Ensure iterator options can be marshaled to and from a binary format.
func TestIteratorOptions_MarshalBinary(t *testing.T) {
	opt := &influxql.IteratorOptions{
		Expr: MustParseExpr("count(value)"),
		Aux:  []influxql.VarRef{{Val: "a"}, {Val: "b"}, {Val: "c"}},
		Interval: influxql.Interval{
			Duration: 1 * time.Hour,
			Offset:   20 * time.Minute,
		},
		Dimensions: []string{"region", "host"},
		Fill:       influxql.NumberFill,
		FillValue:  float64(100),
		Condition:  MustParseExpr(`foo = 'bar'`),
		StartTime:  1000,
		EndTime:    2000,
		Ascending:  true,
		Limit:      100,
		Offset:     200,
		SLimit:     300,
		SOffset:    400,
		Dedupe:     true,
	}

	// Marshal to binary.
	buf, err := opt.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	// Unmarshal back to an object.
	var other influxql.IteratorOptions
	if err := other.UnmarshalBinary(buf); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&other, opt) {
		t.Fatalf("unexpected options: %s", spew.Sdump(other))
	}
}

// Ensure iterator can be encoded and decoded over a byte stream.
func TestIterator_EncodeDecode(t *testing.T) {
	var buf bytes.Buffer

	// Create an iterator with several points & stats.
	itr := &FloatIterator{
		Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: 0},
			{Name: "mem", Tags: ParseTags("host=B"), Time: 1, Value: 10},
		},
		stats: influxql.IteratorStats{
			SeriesN: 2,
			PointN:  0,
		},
	}

	// Encode to the buffer.
	enc := influxql.NewIteratorEncoder(&buf)
	enc.StatsInterval = 100 * time.Millisecond
	if err := enc.EncodeIterator(itr); err != nil {
		t.Fatal(err)
	}

	// Decode from the buffer.
	dec := influxql.NewReaderIterator(&buf, influxql.Float, itr.Stats())

	// Initial stats should exist immediately.
	fdec := dec.(influxql.FloatIterator)
	if stats := fdec.Stats(); !reflect.DeepEqual(stats, influxql.IteratorStats{SeriesN: 2, PointN: 0}) {
		t.Fatalf("unexpected stats(initial): %#v", stats)
	}

	// Read both points.
	if p, err := fdec.Next(); err != nil {
		t.Fatalf("unexpected error(0): %#v", err)
	} else if !reflect.DeepEqual(p, &influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0, Value: 0}) {
		t.Fatalf("unexpected point(0); %#v", p)
	}
	if p, err := fdec.Next(); err != nil {
		t.Fatalf("unexpected error(1): %#v", err)
	} else if !reflect.DeepEqual(p, &influxql.FloatPoint{Name: "mem", Tags: ParseTags("host=B"), Time: 1, Value: 10}) {
		t.Fatalf("unexpected point(1); %#v", p)
	}
	if p, err := fdec.Next(); err != nil {
		t.Fatalf("unexpected error(eof): %#v", err)
	} else if p != nil {
		t.Fatalf("unexpected point(eof); %#v", p)
	}
}

// IteratorCreator is a mockable implementation of SelectStatementExecutor.IteratorCreator.
type IteratorCreator struct {
	CreateIteratorFn  func(m *influxql.Measurement, opt influxql.IteratorOptions) (influxql.Iterator, error)
	FieldDimensionsFn func(m *influxql.Measurement) (fields map[string]influxql.DataType, dimensions map[string]struct{}, err error)
}

func (ic *IteratorCreator) CreateIterator(m *influxql.Measurement, opt influxql.IteratorOptions) (influxql.Iterator, error) {
	return ic.CreateIteratorFn(m, opt)
}

func (ic *IteratorCreator) FieldDimensions(m *influxql.Measurement) (fields map[string]influxql.DataType, dimensions map[string]struct{}, err error) {
	return ic.FieldDimensionsFn(m)
}

func (ic *IteratorCreator) MapType(m *influxql.Measurement, field string) influxql.DataType {
	f, d, err := ic.FieldDimensions(m)
	if err != nil {
		return influxql.Unknown
	}

	if typ, ok := f[field]; ok {
		return typ
	}
	if _, ok := d[field]; ok {
		return influxql.Tag
	}
	return influxql.Unknown
}

// Test implementation of influxql.FloatIterator
type FloatIterator struct {
	Points []influxql.FloatPoint
	Closed bool
	stats  influxql.IteratorStats
}

func (itr *FloatIterator) Stats() influxql.IteratorStats { return itr.stats }
func (itr *FloatIterator) Close() error                  { itr.Closed = true; return nil }

// Next returns the next value and shifts it off the beginning of the points slice.
func (itr *FloatIterator) Next() (*influxql.FloatPoint, error) {
	if len(itr.Points) == 0 || itr.Closed {
		return nil, nil
	}

	v := &itr.Points[0]
	itr.Points = itr.Points[1:]
	return v, nil
}

func FloatIterators(inputs []*FloatIterator) []influxql.Iterator {
	itrs := make([]influxql.Iterator, len(inputs))
	for i := range itrs {
		itrs[i] = influxql.Iterator(inputs[i])
	}
	return itrs
}

// Test implementation of influxql.IntegerIterator
type IntegerIterator struct {
	Points []influxql.IntegerPoint
	Closed bool
	stats  influxql.IteratorStats
}

func (itr *IntegerIterator) Stats() influxql.IteratorStats { return itr.stats }
func (itr *IntegerIterator) Close() error                  { itr.Closed = true; return nil }

// Next returns the next value and shifts it off the beginning of the points slice.
func (itr *IntegerIterator) Next() (*influxql.IntegerPoint, error) {
	if len(itr.Points) == 0 || itr.Closed {
		return nil, nil
	}

	v := &itr.Points[0]
	itr.Points = itr.Points[1:]
	return v, nil
}

func IntegerIterators(inputs []*IntegerIterator) []influxql.Iterator {
	itrs := make([]influxql.Iterator, len(inputs))
	for i := range itrs {
		itrs[i] = influxql.Iterator(inputs[i])
	}
	return itrs
}

// Test implementation of influxql.StringIterator
type StringIterator struct {
	Points []influxql.StringPoint
	Closed bool
	stats  influxql.IteratorStats
}

func (itr *StringIterator) Stats() influxql.IteratorStats { return itr.stats }
func (itr *StringIterator) Close() error                  { itr.Closed = true; return nil }

// Next returns the next value and shifts it off the beginning of the points slice.
func (itr *StringIterator) Next() (*influxql.StringPoint, error) {
	if len(itr.Points) == 0 || itr.Closed {
		return nil, nil
	}

	v := &itr.Points[0]
	itr.Points = itr.Points[1:]
	return v, nil
}

func StringIterators(inputs []*StringIterator) []influxql.Iterator {
	itrs := make([]influxql.Iterator, len(inputs))
	for i := range itrs {
		itrs[i] = influxql.Iterator(inputs[i])
	}
	return itrs
}

// Test implementation of influxql.BooleanIterator
type BooleanIterator struct {
	Points []influxql.BooleanPoint
	Closed bool
	stats  influxql.IteratorStats
}

func (itr *BooleanIterator) Stats() influxql.IteratorStats { return itr.stats }
func (itr *BooleanIterator) Close() error                  { itr.Closed = true; return nil }

// Next returns the next value and shifts it off the beginning of the points slice.
func (itr *BooleanIterator) Next() (*influxql.BooleanPoint, error) {
	if len(itr.Points) == 0 || itr.Closed {
		return nil, nil
	}

	v := &itr.Points[0]
	itr.Points = itr.Points[1:]
	return v, nil
}

func BooleanIterators(inputs []*BooleanIterator) []influxql.Iterator {
	itrs := make([]influxql.Iterator, len(inputs))
	for i := range itrs {
		itrs[i] = influxql.Iterator(inputs[i])
	}
	return itrs
}
