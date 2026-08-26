package tui

type Mode int

const (
	ListMode Mode = iota
	AddMode
	DeleteConfirmMode
	ForceDeleteConfirmMode
	SearchMode
	DeletingMode
	CreatingMode
)

type EditorDoneMsg struct {
	TmpFile string
	Err     error
}
