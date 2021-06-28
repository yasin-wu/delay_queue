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
	//时间类型:0-延迟多少秒执行,1-具体执行时间(时间戳:秒)
	Type DelayType
	//延迟执行时间,单位:秒
	DelayTime int64
	//任务执行参数
	Args []interface{}
}

// DelayType 延迟任务类型
type DelayType int

const (
	DelayTypeDuration DelayType = iota //延迟多少秒执行
	DelayTypeDate                      //具体执行时间(时间戳:秒)
)
