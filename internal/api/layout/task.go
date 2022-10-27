package layout

type task struct {
	ID     string
	IDS    []string // agent id
	Groups []string // groups
	Index  int      // current index
	Err    error
	Done   bool
}

func newTask(taskID string, ids, groups []string) *task {
	return &task{
		ID:     taskID,
		IDS:    ids,
		Groups: groups,
	}
}
