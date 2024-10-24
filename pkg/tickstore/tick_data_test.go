package tickstore

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTickData(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name  string
		args  args
		want  *TickData[float32]
		panic bool
	}{
		{
			name:  "Size 0",
			args:  args{size: 0},
			panic: true,
		},
		{
			name: "Size 1",
			args: args{size: 1},
			want: &TickData[float32]{
				values: make([]float32, 1),
			},
		},
		{
			name: "Size 2",
			args: args{size: 2},
			want: &TickData[float32]{
				values: make([]float32, 2),
			},
		},
		{
			name: "Size 4",
			args: args{size: 4},
			want: &TickData[float32]{
				values: make([]float32, 4),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panic {
				assert.Panics(t, func() { NewTickData[float32](tt.args.size) })
				return
			}
			got := NewTickData[float32](tt.args.size)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTickData_AddValues(t *testing.T) {
	type args struct {
		values []float32
	}
	tests := []struct {
		tData *TickData[float32]
		args  args
		want  *TickData[float32]
		err   bool
	}{
		{
			tData: &TickData[float32]{values: []float32{-1, -2, -3, -4}},
			args:  args{values: []float32{1}},
			want:  &TickData[float32]{values: []float32{1, -1, -2, -3}},
		},
		{
			tData: &TickData[float32]{values: []float32{1, 0, 0, 0}},
			args:  args{values: []float32{2, 3}},
			want:  &TickData[float32]{values: []float32{2, 3, 1, 0}},
		},
		{
			tData: &TickData[float32]{values: []float32{1, 2, 3, 0}},
			args:  args{values: []float32{4, 5}},
			want:  &TickData[float32]{values: []float32{4, 5, 1, 2}},
		},
		{
			tData: &TickData[float32]{values: []float32{5, 2, 3, 4}},
			args:  args{values: []float32{6, 7, 8}},
			want:  &TickData[float32]{values: []float32{6, 7, 8, 5}},
		},
		{
			tData: &TickData[float32]{values: []float32{5, 6, 7, 8}},
			args:  args{values: []float32{9, 9, 9, 9}},
			want:  &TickData[float32]{values: []float32{9, 9, 9, 9}},
		},
		{
			tData: &TickData[float32]{values: []float32{9, 9, 9, 9}},
			args:  args{values: []float32{1, 2, 3, 4, 5}},
			want:  &TickData[float32]{values: []float32{9, 9, 9, 9}},
			err:   true,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.args.values), func(t *testing.T) {
			oldPointer := &tt.tData.values[0]

			err := tt.tData.AddValues(tt.args.values...)
			if tt.err && err != nil {
				require.Equal(t, tt.want, tt.tData,
					"returned error, but values has been changed")
				return
			}

			if tt.err && err == nil {
				t.Fatalf("AddValues() doesn't return expected error")
			}
			if err != nil {
				t.Fatalf("AddValues() returns unexpected error: %v", err)
			}

			assert.Equal(t, tt.want, tt.tData)

			newPointer := &tt.tData.values[0]
			if newPointer != oldPointer {
				t.Fatalf(
					"Array reallocated. Old ptr: %v; New ptr: %v",
					oldPointer, newPointer,
				)
			}
		})
	}
}

func TestTickData_GetValues(t *testing.T) {
	tests := []struct {
		tData *TickData[float32]
		want  []float32
		err   bool
	}{
		{
			tData: &TickData[float32]{
				values: []float32{1, 0, 0, 0},
			},
			want: []float32{1, 0, 0, 0},
		},
		{
			tData: &TickData[float32]{
				values: []float32{1, 2, 3, 4},
			},
			want: []float32{1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.want), func(t *testing.T) {
			val := tt.tData.values
			oldPointer := &val[0]
			got := tt.tData.GetValues()

			assert.Equal(t, tt.want, got)

			newPointer := &val[0]
			getterPointer := &got[0]
			if newPointer != oldPointer || getterPointer != oldPointer {
				t.Fatalf(
					"Array reallocated. Old ptr: %v; New ptr: %v; Get ptr: %v",
					oldPointer, newPointer, getterPointer,
				)
			}
		})
	}
}

func BenchmarkTickData_AddValues(b *testing.B) {
	td := &TickData[float32]{
		values: make([]float32, 180),
	}

	table := []struct {
		values []float32
	}{
		{values: slices.Repeat([]float32{1}, 1)},
		{values: slices.Repeat([]float32{1}, 5)},
		{values: slices.Repeat([]float32{1, 2, 3, 4, 5}, 1)},
		{values: slices.Repeat([]float32{1}, 10)},
		{values: slices.Repeat([]float32{1, 2, 3, 4, 5}, 2)},
		{values: slices.Repeat([]float32{1}, 60)},
		{values: slices.Repeat([]float32{1, 2, 3, 4, 5}, 12)},
		{values: slices.Repeat([]float32{1}, 90)},
		{values: slices.Repeat([]float32{1}, 180)},
		{values: slices.Repeat([]float32{1, 2, 3, 4, 5}, 36)},
		{values: []float32{2, 3}},
		{values: []float32{2, 3, 4}},
		{values: []float32{2, 3, 4, 5, 3, 4, 5, 3, 4, 5}},
	}

	for idx, v := range table {
		b.Run(fmt.Sprintf("%02d values %d", idx, len(v.values)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				td.AddValues(v.values...)
			}
		})
	}
}
