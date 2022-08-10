package dqueue

import "github.com/yasin-wu/delay_queue/v2/pkg"

type JobExecutor struct {
	ID     string
	Action pkg.JobBaseAction
}
