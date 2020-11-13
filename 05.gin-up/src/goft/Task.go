package goft

import "sync"

// TaskFunc 任务执行的函数
type TaskFunc func(params ...interface{})

// taskList 没有初始化的任务列表
var taskList chan *TaskExecutor

// init 引用包时执行
func init() {
	chlist := getTaskList() // 得到任务列表
	go func() {
		for t := range chlist {
			t.Exec() // 执行任务
		}
	}()
}

var once sync.Once

// getTaskList
func getTaskList() chan *TaskExecutor {
	once.Do(func() { // 单例模式
		// 对 taskList 进行初始化 chan
		taskList = make(chan *TaskExecutor)
	})
	return taskList
}

// TaskExecutor 任务执行者
type TaskExecutor struct {
	f TaskFunc      // 任务中需要执行的函数
	p []interface{} // 任务重需要执行函数的参数
}

func NewTaskExecutor(f TaskFunc, p []interface{}) *TaskExecutor {
	return &TaskExecutor{f: f, p: p}
}

// Exec 执行任务
func (this *TaskExecutor) Exec() {
	this.f(this.p...)
}

// Task 对外的方法，当有任务时，把任务塞入任务列表管道
func Task(f TaskFunc, params ...interface{}) {
	getTaskList() <- NewTaskExecutor(f, params) // 增加任务队列
}
