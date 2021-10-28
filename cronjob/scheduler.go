package cronjob

import (
	"fmt"
	"sort"
	"time"

	"github.com/yasin-wu/delay_queue/logger"
)

type Scheduler struct {
	jobs    []*Wrapper
	logger  logger.Logger
	nameSet map[string]bool
}

func New() *Scheduler {
	return &Scheduler{
		jobs:    nil,
		nameSet: make(map[string]bool),
		logger:  logger.DefaultLogger,
	}
}

func (sche *Scheduler) Register(phase []int, period int, job CronJob) {
	if len(phase) > 1 {
		sort.Slice(phase, func(i, j int) bool {
			return phase[i] < phase[j]
		})
	}
	jobName := job.Name()
	if _, ok := sche.nameSet[jobName]; jobName == "" || ok {
		sche.logger.ErrorF("CronJob register failed , JobName:%s", jobName)
	} else {
		sche.jobs = append(sche.jobs, &Wrapper{
			job:    job,
			phase:  phase,
			period: period,
		})
		sche.nameSet[job.Name()] = true
	}
}

func (sche *Scheduler) Start() {
	sche.logger.InfoF("CronJob starting......")
	for i := 0; i < len(sche.jobs); i++ {
		if sche.jobs[i].ifActive() && sche.validateJob(sche.jobs[i]) {
			go sche.run(sche.jobs[i])
		}
	}
}

func (sche *Scheduler) SetLogger(logger logger.Logger) {
	sche.logger = logger
}

func (sche *Scheduler) run(job *Wrapper) {
	for {
		nextTimeInterval := sche.calculateNextTime(job.phase, job.period, job.count)
		if nextTimeInterval >= 0 {
			time.Sleep(time.Duration(nextTimeInterval) * time.Second)
		} else {
			break
		}
		err := func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("job %s panic, recover=%v", job.name(), r)
				}
			}()
			return job.Process()
		}()
		job.count++
		if err != nil && !job.ifReboot() {
			break
		}
	}
}

func (sche *Scheduler) calculateNextTime(phase []int, period int, calCount int) int {
	if len(phase) == 0 {
		if calCount == 0 {
			return 0
		} else {
			return period
		}
	} else {
		nowTimePhase := int((time.Now().Unix() + 28800) % 86400)
		var i = 0
		for i = 0; i < len(phase); i++ {
			if nowTimePhase < phase[i] {
				break
			}
		}
		if i == len(phase) {
			return period - nowTimePhase + phase[0]
		} else {
			return phase[i] - nowTimePhase
		}
	}
}

func (sche *Scheduler) validateJob(job *Wrapper) bool {
	if len(job.phase) == 0 {
		return job.period > 0
	} else {
		return job.period >= SecondOfDay && job.period%SecondOfDay == 0
	}
}
