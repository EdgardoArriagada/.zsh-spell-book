package tui

type Mode int

const (
	ListMode Mode = iota
	AddMode
	DeleteConfirmMode
	ForceDeleteConfirmMode
	SearchMode
)
