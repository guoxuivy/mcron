package mcron

import (
	"log"
)

type SchduleActiveCallback func(int)

//定时任务实体
type ScheduleJob struct {
	id       int
	callback SchduleActiveCallback
}

func NewScheduleJob(_id int, _job SchduleActiveCallback) *ScheduleJob {
	instance := &ScheduleJob{
		id:       _id,
		callback: _job,
	}
	return instance
}

//执行入口
func (this *ScheduleJob) Run() {
	log.Println("Invalid callback")
	if nil != this.callback {
		this.callback(this.id)
	} else {
		log.Println("Invalid callback function")
	}
}
