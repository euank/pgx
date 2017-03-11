package pgtype_test

import (
	"reflect"
	"testing"

	"github.com/jackc/pgx/pgtype"
)

func TestXIDTranscode(t *testing.T) {
	testSuccessfulTranscode(t, "xid", []interface{}{
		pgtype.XID{Uint: 42, Status: pgtype.Present},
		pgtype.XID{Status: pgtype.Null},
	})
}

func TestXIDConvertFrom(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result pgtype.XID
	}{
		{source: uint32(1), result: pgtype.XID{Uint: 1, Status: pgtype.Present}},
	}

	for i, tt := range successfulTests {
		var r pgtype.XID
		err := r.ConvertFrom(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if r != tt.result {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestXIDAssignTo(t *testing.T) {
	var ui32 uint32
	var pui32 *uint32

	simpleTests := []struct {
		src      pgtype.XID
		dst      interface{}
		expected interface{}
	}{
		{src: pgtype.XID{Uint: 42, Status: pgtype.Present}, dst: &ui32, expected: uint32(42)},
		{src: pgtype.XID{Status: pgtype.Null}, dst: &pui32, expected: ((*uint32)(nil))},
	}

	for i, tt := range simpleTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Interface(); dst != tt.expected {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	pointerAllocTests := []struct {
		src      pgtype.XID
		dst      interface{}
		expected interface{}
	}{
		{src: pgtype.XID{Uint: 42, Status: pgtype.Present}, dst: &pui32, expected: uint32(42)},
	}

	for i, tt := range pointerAllocTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Elem().Interface(); dst != tt.expected {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	errorTests := []struct {
		src pgtype.XID
		dst interface{}
	}{
		{src: pgtype.XID{Status: pgtype.Null}, dst: &ui32},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}
}