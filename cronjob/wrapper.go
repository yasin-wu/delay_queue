package cronjob

type Wrapper struct {
	period int
	count  int
	phase  []int
	job    CronJob
}

func (c *Wrapper) Process() error {
	return c.job.Process()
}

func (c *Wrapper) name() string {
	return c.job.Name()
}

func (c *Wrapper) ifActive() bool {
	return c.job.IfActive()
}

func (c *Wrapper) ifReboot() bool {
	return c.job.IfReboot()
}
