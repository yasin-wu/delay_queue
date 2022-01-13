package delayqueue

/**
 * @author: yasin
 * @date: 2022/1/13 11:06
 * @description: 延迟服务
 */
type JobBaseAction interface {
	ID() string                       //服务唯一ID
	Execute(args []interface{}) error //服务执行动作
}

type jobExecutor struct {
	ID     string
	action JobBaseAction
}

/**
 * @author: yasin
 * @date: 2022/1/13 13:23
 * @description: 延迟任务
 */
type DelayJob struct {
	ID        string        //任务ID
	Type      DelayType     //时间类型:0-延迟多少秒执行,1-具体执行时间(时间戳:秒)
	DelayTime int64         //延迟执行时间,单位:秒
	Args      []interface{} //任务执行参数

}

// DelayType 延迟任务类型
type DelayType int

const (
	DelayTypeDuration DelayType = iota //延迟多少秒执行
	DelayTypeDate                      //具体执行时间(时间戳:秒)
)
