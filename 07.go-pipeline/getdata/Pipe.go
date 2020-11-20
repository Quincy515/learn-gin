package getdata

import "sync"

type InChan chan interface{}                  // InChan 管道入参
type OutChan chan interface{}                 // OutChan 管道数据输出
type CmdFunc func(args ...interface{}) InChan // CmdFunc 读取数据源的普通方法类型
type PipeCmdFunc func(in InChan) OutChan      // PipeCmdFunc 管道方法的类型
type Pipe struct {                            // Pipe 管道的定义
	Cmd     CmdFunc
	PipeCmd PipeCmdFunc
	Count   int
}

func NewPipe() *Pipe {
	return &Pipe{Count: 1}
}

// SetCmd 设置普通取数据方法
func (this *Pipe) SetCmd(c CmdFunc) {
	this.Cmd = c
}

// SetPipeCmd 设置管道的执行方法
func (this *Pipe) SetPipeCmd(c PipeCmdFunc, count int) {
	this.PipeCmd = c
	this.Count = count
}

// Exec 管道的实现
func (this *Pipe) Exec(args ...interface{}) OutChan {
	in := this.Cmd(args)
	out := make(OutChan)
	wg := sync.WaitGroup{}
	for i := 0; i < this.Count; i++ {
		getChan := this.PipeCmd(in)
		wg.Add(1)
		go func(input OutChan) {
			defer wg.Done()
			for v := range input {
				out <- v
			}
		}(getChan)
	}
	go func() {
		defer close(out)
		wg.Wait()
	}()
	return out
}
