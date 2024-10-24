package bitflags

type BitFlag int

func Has[T ~int](flags, wanted T) bool {
	return flags&wanted == wanted
}

func Set[T ~int](flags, newFlags T) T {
	return flags | newFlags
}

func Toggle[T ~int](flags, newFlags T) T {
	return flags ^ newFlags
}

func Clear[T ~int](flags, newFlags T) T {
	return flags &^ newFlags
}
