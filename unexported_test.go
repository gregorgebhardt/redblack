package redblack

func CheckNoRedRed[V any, T Orderable[V]](t *Tree[V, T]) bool {
	return t.checkNoRedRed()
}

func CheckBlackHeight[V any, T Orderable[V]](t *Tree[V, T]) (uint, bool) {
	return t.checkBlackHeight()
}

func CheckLeftLeaning[V any, T Orderable[V]](t *Tree[V, T]) bool {
	return t.checkLeftLeaning()
}
