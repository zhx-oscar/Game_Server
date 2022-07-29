package main

const (
	CronJobDailyReset = "daily reset job"
)

type _CronJob struct {
	Type string
	Arg  interface{}
}
