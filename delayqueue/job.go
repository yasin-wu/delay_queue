package delayqueue

type JobBaseAction interface {
	ID() string
	Execute(args []interface{}) error
}

type jobExecutor struct {
	ID     string
	action JobBaseAction
}

type DelayJob struct {
	ID        string
	DelayTime int64
	Args      []interface{}
}
