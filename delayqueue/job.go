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
	//任务ID
	ID string
	//延迟执行时间,单位:秒
	DelayTime int64
	//任务执行参数
	Args []interface{}
}
