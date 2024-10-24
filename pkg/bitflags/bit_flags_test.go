package bitflags

import (
	"reflect"
	"testing"
)

type TestBitFlag BitFlag

const (
	TestFlagsNone TestBitFlag = 0
	TestFlagsA    TestBitFlag = 1 << (iota - 1)
	TestFlagsB
	TestFlagsC
)

const (
	TestFlagsAB = TestFlagsA | TestFlagsB
)

func TestHas(t *testing.T) {
	type args struct {
		flags TestBitFlag
		check TestBitFlag
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "None has None",
			args: args{flags: TestFlagsNone, check: TestFlagsNone},
			want: true,
		},
		{
			name: "None !has A",
			args: args{flags: TestFlagsNone, check: TestFlagsA},
			want: false,
		},
		{
			name: "AB has A",
			args: args{flags: TestFlagsAB, check: TestFlagsA},
			want: true,
		},
		{
			name: "AB has B",
			args: args{flags: TestFlagsAB, check: TestFlagsB},
			want: true,
		},
		{
			name: "AB has AB",
			args: args{flags: TestFlagsAB, check: TestFlagsA | TestFlagsB},
			want: true,
		},
		{
			name: "AB !has C",
			args: args{flags: TestFlagsAB, check: TestFlagsC},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Has(tt.args.flags, tt.args.check)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet(t *testing.T) {
	type args struct {
		flags    TestBitFlag
		newFlags TestBitFlag
	}
	tests := []struct {
		name string
		args args
		want TestBitFlag
	}{
		{
			name: "None sets nothing",
			args: args{flags: TestFlagsNone, newFlags: TestFlagsNone},
			want: TestFlagsNone,
		},
		{
			name: "Add A to None",
			args: args{flags: TestFlagsNone, newFlags: TestFlagsA},
			want: TestFlagsA,
		},
		{
			name: "Add A to A",
			args: args{flags: TestFlagsA, newFlags: TestFlagsA},
			want: TestFlagsA,
		},
		{
			name: "Add B to A",
			args: args{flags: TestFlagsA, newFlags: TestFlagsB},
			want: TestFlagsAB,
		},
		{
			name: "Add B to AB",
			args: args{flags: TestFlagsAB, newFlags: TestFlagsB},
			want: TestFlagsAB,
		},
		{
			name: "Add C to AB",
			args: args{flags: TestFlagsAB, newFlags: TestFlagsC},
			want: TestFlagsAB | TestFlagsC,
		},
		{
			name: "Add BC to A = ABC",
			args: args{flags: TestFlagsA, newFlags: TestFlagsB | TestFlagsC},
			want: TestFlagsAB | TestFlagsC,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Set(tt.args.flags, tt.args.newFlags)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got: %04b, want: %04b", got, tt.want)
			}
		})
	}
}

func TestToggle(t *testing.T) {
	type args struct {
		flags    TestBitFlag
		newFlags TestBitFlag
	}
	tests := []struct {
		name string
		args args
		want TestBitFlag
	}{
		{
			name: "None to None = None",
			args: args{flags: TestFlagsNone, newFlags: TestFlagsNone},
			want: TestFlagsNone,
		},
		{
			name: "None to A = A",
			args: args{flags: TestFlagsA, newFlags: TestFlagsNone},
			want: TestFlagsA,
		},
		{
			name: "A to None = A",
			args: args{flags: TestFlagsNone, newFlags: TestFlagsA},
			want: TestFlagsA,
		},
		{
			name: "A to A = None",
			args: args{flags: TestFlagsA, newFlags: TestFlagsA},
			want: TestFlagsNone,
		},
		{
			name: "B to A = AB",
			args: args{flags: TestFlagsA, newFlags: TestFlagsB},
			want: TestFlagsAB,
		},
		{
			name: "B to AB = A",
			args: args{flags: TestFlagsAB, newFlags: TestFlagsB},
			want: TestFlagsA,
		},
		{
			name: "C to AB = ABC",
			args: args{flags: TestFlagsAB, newFlags: TestFlagsC},
			want: TestFlagsAB | TestFlagsC,
		},
		{
			name: "AB to C = ABC",
			args: args{flags: TestFlagsC, newFlags: TestFlagsAB},
			want: TestFlagsAB | TestFlagsC,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Toggle(tt.args.flags, tt.args.newFlags)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got: %04b, want: %04b", got, tt.want)
			}
		})
	}
}

func TestClear(t *testing.T) {
	type args struct {
		flags    TestBitFlag
		newFlags TestBitFlag
	}
	tests := []struct {
		name string
		args args
		want TestBitFlag
	}{
		{
			name: "None from None = None",
			args: args{flags: TestFlagsNone, newFlags: TestFlagsNone},
			want: TestFlagsNone,
		},
		{
			name: "A from None = None",
			args: args{flags: TestFlagsNone, newFlags: TestFlagsA},
			want: TestFlagsNone,
		},
		{
			name: "A from A = None",
			args: args{flags: TestFlagsA, newFlags: TestFlagsA},
			want: TestFlagsNone,
		},
		{
			name: "B from A = A",
			args: args{flags: TestFlagsA, newFlags: TestFlagsB},
			want: TestFlagsA,
		},
		{
			name: "B from AB = A",
			args: args{flags: TestFlagsAB, newFlags: TestFlagsB},
			want: TestFlagsA,
		},
		{
			name: "C from ABC = AB",
			args: args{flags: TestFlagsAB | TestFlagsC, newFlags: TestFlagsC},
			want: TestFlagsAB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Clear(tt.args.flags, tt.args.newFlags)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got: %04b, want: %04b", got, tt.want)
			}
		})
	}
}
