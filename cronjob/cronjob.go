package cronjob

type CronJob interface {
	Name() string
	Process() error
	IfActive() bool
	IfReboot() bool
}
