package utils

import (
	"image"
	"reflect"
	"testing"
)

var test_targetRect = image.Rectangle{
	Min: image.Pt(0, 0),
	Max: image.Pt(10, 10),
}

var test_table = []struct {
	testedRect image.Rectangle
	wantedRect image.Rectangle
}{
	{
		testedRect: image.Rectangle{
			Min: image.Pt(0, 0),
			Max: image.Pt(0, 0),
		},
		wantedRect: image.Rectangle{
			Min: image.Pt(0, 0),
			Max: image.Pt(0, 0),
		},
	},
	{
		testedRect: image.Rectangle{
			Min: image.Pt(0, 0),
			Max: image.Pt(1, 1),
		},
		wantedRect: image.Rectangle{
			Min: image.Pt(0, 0),
			Max: image.Pt(1, 1),
		},
	},
	{
		testedRect: image.Rectangle{
			Min: image.Pt(0, 0),
			Max: image.Pt(2, 2),
		},
		wantedRect: image.Rectangle{
			Min: image.Pt(0, 0),
			Max: image.Pt(2, 2),
		},
	},
	{
		testedRect: image.Rectangle{
			Min: image.Pt(0, 0),
			Max: image.Pt(10, 10),
		},
		wantedRect: image.Rectangle{
			Min: image.Pt(0, 0),
			Max: image.Pt(10, 10),
		},
	},
	{
		testedRect: image.Rectangle{
			Min: image.Pt(0, 0),
			Max: image.Pt(10, 1),
		},
		wantedRect: image.Rectangle{
			Min: image.Pt(0, 0),
			Max: image.Pt(10, 1),
		},
	},
	{
		testedRect: image.Rectangle{
			Min: image.Pt(0, 0),
			Max: image.Pt(20, 10),
		},
		wantedRect: image.Rectangle{
			Min: image.Pt(0, 0),
			Max: image.Pt(10, 10),
		},
	},
	{
		testedRect: image.Rectangle{
			Min: image.Pt(0, 0),
			Max: image.Pt(20, 20),
		},
		wantedRect: image.Rectangle{
			Min: image.Pt(0, 0),
			Max: image.Pt(10, 10),
		},
	},
	{
		testedRect: image.Rectangle{
			Min: image.Pt(9, 0),
			Max: image.Pt(20, 20),
		},
		wantedRect: image.Rectangle{
			Min: image.Pt(9, 0),
			Max: image.Pt(10, 10),
		},
	},
	{
		testedRect: image.Rectangle{
			Min: image.Pt(100, 100),
			Max: image.Pt(200, 200),
		},
		wantedRect: image.Rectangle{
			Min: image.Pt(0, 0),
			Max: image.Pt(0, 0),
		},
	},
	{
		testedRect: image.Rectangle{
			Min: image.Pt(-10, 0),
			Max: image.Pt(10, 1),
		},
		wantedRect: image.Rectangle{
			Min: image.Pt(0, 0),
			Max: image.Pt(10, 1),
		},
	},
}

func Test_LimitPointToCanvas(t *testing.T) {
	for _, tt := range test_table {
		t.Run("", func(t *testing.T) {
			result := image.Rectangle{
				Min: LimitPointToRectangle(tt.testedRect.Min, test_targetRect),
				Max: LimitPointToRectangle(tt.testedRect.Max, test_targetRect),
			}
			if result.Empty() {
				result = image.Rectangle{}
			}
			if !reflect.DeepEqual(result, tt.wantedRect) {
				t.Fatalf(
					"Input: %v, Result: %v, want: %v",
					tt.testedRect, result, tt.wantedRect,
				)
			}
		})
	}
}

func Test_LimitRectangleToCanvas(t *testing.T) {
	for _, tt := range test_table {
		t.Run("", func(t *testing.T) {
			result := LimitRectangleToRectangle(
				tt.testedRect,
				test_targetRect,
			)
			if !reflect.DeepEqual(result, tt.wantedRect) {
				t.Fatalf(
					"Input: %v, Result: %v, want: %v",
					tt.testedRect, result, tt.wantedRect,
				)
			}
		})
	}
}

func Benchmark_LimitPointToRectangle(b *testing.B) {
	for _, v := range test_table {
		b.Run("", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// Min in
				LimitPointToRectangle(v.testedRect.Min, test_targetRect)
				// Max in
				LimitPointToRectangle(v.testedRect.Max, test_targetRect)
			}
		})
	}
}

func Benchmark_LimitRectangleToRectangle(b *testing.B) {
	for _, v := range test_table {
		b.Run("", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				LimitRectangleToRectangle(v.testedRect, test_targetRect)
			}
		})
	}
}
