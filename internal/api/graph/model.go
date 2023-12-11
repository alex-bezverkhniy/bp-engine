package graph

type (
	StatusNode struct {
		Name         string
		NextStatuses []StatusNode
	}
	ProcessModel struct {
		Name       string
		GraphModel StatusNode
	}
)
