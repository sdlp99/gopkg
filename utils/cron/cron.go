package cron

import (
	"github.com/robfig/cron/v3"
	"github.com/sdlp99/sdpkg/utils/logger"
)

var (
	parserD cron.Parser
	cronD   *cron.Cron
)

func init() {
	parserD = cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
	)
	cronD = cron.New(cron.WithParser(parserD))
	cronD.Start()
}

func AddCron(schedule string, cmd func()) {
	_, err := cronD.AddFunc(schedule, cmd)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}
}

func ExitCron() {
	cronD.Stop()
}
