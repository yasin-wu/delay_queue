package cronjob

import (
	"fmt"
	"sort"
	"time"

	"github.com/yasin-wu/delay_queue/v2/internal/logger"
)

type Scheduler struct {
	jobs        []*Wrapper
	logger      logger.Logger
	nameSet     map[string]bool
	secondOfDay int
}

func New() *Scheduler {
	return &Scheduler{
		jobs:        nil,
		nameSet:     make(map[string]bool),
		logger:      logger.DefaultLogger,
		secondOfDay: 24 * 60 * 60,
	}
}

func (sche *Scheduler) Register(phase []int, period int, job CronJob) {
	if len(phase) > 1 {
		sort.Slice(phase, func(i, j int) bool {
			return phase[i] < phase[j]
		})
	}
	if _, ok := sche.nameSet[job.Name()]; job.Name() == "" || ok {
		sche.logger.Errorf("CronJob register failed , JobName:%s", job.Name())
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
	sche.logger.Infof("CronJob starting......")
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
		sche.logger.Infof("job %s next time: %v", job.name(), nextTimeInterval)
		if nextTimeInterval >= 0 {
			time.Sleep(time.Duration(nextTimeInterval) * time.Second)
		} else {
			break
		}
		err := func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					sche.logger.Errorf("job %s panic, recover=%v", job.name(), r)
					err = fmt.Errorf("job %s panic, recover=%v", job.name(), r)
				}
			}()
			sche.logger.Infof("job %s start: %v", job.name(), time.Now().Format("2006-01-02 15:04:05"))
			return job.Process()
		}()
		job.count++
		if err != nil && !job.ifReboot() {
			sche.logger.Errorf("job %s execute failed , error:%v", job.name(), err)
			break
		}
	}
}

func (sche *Scheduler) calculateNextTime(phase []int, period int, calCount int) int {
	if len(phase) == 0 {
		if calCount == 0 {
			return 0
		}
		return period
	}
	nowTimePhase := int((time.Now().Unix() + 28800) % 86400)
	i := 0
	for i = 0; i < len(phase); i++ {
		if nowTimePhase < phase[i] {
			break
		}
	}
	if i == len(phase) {
		return period - nowTimePhase + phase[0]
	}
	return phase[i] - nowTimePhase
}

func (sche *Scheduler) validateJob(job *Wrapper) bool {
	if len(job.phase) == 0 {
		return job.period > 0
	}
	return job.period >= sche.secondOfDay && job.period%sche.secondOfDay == 0
}
