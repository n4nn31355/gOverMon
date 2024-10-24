package utils

import (
	"cmp"
	"image"
	"slices"
)

func SliceMinMax[S ~[]T, T cmp.Ordered](slice S) (sMin T, sMax T) {
	return slices.Min(slice), slices.Max(slice)
}

func LimitPointToRectangle(point image.Point, target image.Rectangle) image.Point {
	return image.Pt(
		max(min(point.X, target.Max.X), target.Min.X),
		max(min(point.Y, target.Max.Y), target.Min.Y),
	)
}

func LimitRectangleToRectangle(rect image.Rectangle, target image.Rectangle) image.Rectangle {
	return target.Intersect(rect)
}
