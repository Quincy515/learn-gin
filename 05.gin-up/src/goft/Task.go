package goft

import (
	"github.com/robfig/cron/v3"
	"sync"
)

// TaskFunc 任务执行的函数
type TaskFunc func(params ...interface{})

// taskList 没有初始化的任务列表
var taskList chan *TaskExecutor

var (
	once     sync.Once
	onceCron sync.Once
	taskCron *cron.Cron // 定时任务
)

// init 引用包时执行
func init() {
	chlist := getTaskList() // 得到任务列表
	go func() {
		for t := range chlist {
			doTask(t)
		}
	}()
}

// doTask  开协程执行任务
func doTask(t *TaskExecutor) {
	go func() {
		defer func() {
			if t.callback != nil { // 第二个参数运行为 nil
				t.callback() // 简单的回调
			}
		}()
		t.Exec()
	}()
}

// getCronTask 初始化定时任务
func getCronTask() *cron.Cron {
	onceCron.Do(func() {
		taskCron = cron.New(cron.WithSeconds())
	})
	return taskCron
}

// getTaskList 初始化协程任务
func getTaskList() chan *TaskExecutor {
	once.Do(func() { // 单例模式
		// 对 taskList 进行初始化 chan
		taskList = make(chan *TaskExecutor)
	})
	return taskList
}

// TaskExecutor 任务执行者
type TaskExecutor struct {
	f        TaskFunc      // 任务中需要执行的函数
	p        []interface{} // 任务重需要执行函数的参数
	callback func()        // 回调
}

func NewTaskExecutor(f TaskFunc, p []interface{}, callback func()) *TaskExecutor {
	return &TaskExecutor{f: f, p: p, callback: callback}
}

// Exec 执行任务
func (this *TaskExecutor) Exec() {
	this.f(this.p...)
}

// Task 对外的方法，当有任务时，把任务塞入任务列表管道
func Task(f TaskFunc, callback func(), params ...interface{}) {
	if f == nil { // 第一个参数执行函数不允许为 nil
		return
	}
	go func() {
		getTaskList() <- NewTaskExecutor(f, params, callback) // 增加任务队列
	}()
}
