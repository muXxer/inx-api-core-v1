package daemon

const (
	PriorityDisconnectINX = iota // no dependencies
	PriorityStopDatabase
	PriorityStopDatabaseAPI
	PriorityStopDatabaseAPIINX
	PriorityStopPrometheus
)
