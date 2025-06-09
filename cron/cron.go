package cron

import (
	"github.com/robfig/cron/v3"
)

var c *cron.Cron

// InitCron 初始化定时任务管理器
func InitCron() {
	c = cron.New()
	c.Start()
}

// AddJob 添加定时任务
func AddJob(spec string, cmd func()) error {
	_, err := c.AddFunc(spec, cmd)
	return err
}

// StopCron 停止定时任务管理器
func StopCron() {
	if c != nil {
		c.Stop()
	}
}
