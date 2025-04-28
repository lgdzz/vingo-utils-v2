// *****************************************************************************
// 作者: lgdz
// 创建时间: 2025/4/16
// 描述：协程池
// *****************************************************************************

package pool

import (
	"context"
	"fmt"
	"github.com/lgdzz/vingo-utils-v2/vingo"
	"sync"
)

// 支持返回结果和错误的任务函数签名
type TaskFunc func(ctx context.Context) Result

type Result struct {
	Data   any     `json:"data"`
	Result any     `json:"result"`
	Error  *string `json:"error"`
}

type GoroutinePool struct {
	maxWorkers int
	tasks      chan TaskFunc
	results    chan Result
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// 创建协程池，支持 context
func NewGoroutinePool(ctx context.Context, maxWorkers int, taskNum int) *GoroutinePool {
	c, cancel := context.WithCancel(ctx)
	return &GoroutinePool{
		maxWorkers: maxWorkers,
		tasks:      make(chan TaskFunc),
		results:    make(chan Result, taskNum),
		ctx:        c,
		cancel:     cancel,
	}
}

// 启动 worker
func (p *GoroutinePool) Run() {
	for i := 0; i < p.maxWorkers; i++ {
		go func() {
			for {
				select {
				case <-p.ctx.Done():
					return
				case task, ok := <-p.tasks:
					if !ok {
						return
					}
					res := task(p.ctx)
					p.results <- res
					p.wg.Done()
				}
			}
		}()
	}
}

// 提交任务
func (p *GoroutinePool) Submit(task TaskFunc) {
	p.wg.Add(1)
	p.tasks <- task
}

// 取消所有任务
func (p *GoroutinePool) Cancel() {
	p.cancel()
}

// 等待所有任务完成并关闭通道，返回所有结果
func (p *GoroutinePool) CloseAndWait() []Result {
	p.wg.Wait()
	close(p.tasks)
	close(p.results)

	var out []Result
	for r := range p.results {
		out = append(out, r)
	}
	return out
}

func BusinessHandle[T any](data T, handle func(object T) any) (res Result) {
	defer func() {
		if err := recover(); err != nil {
			vingo.LogError(fmt.Sprintf("协程池中业务处理错误：：%v", err))
			res = Result{
				Data:   data,
				Result: "fail",
				Error:  vingo.Of(fmt.Sprintf("%v", err)),
			}
		}
	}()

	resData := handle(data)
	if resData != nil {
		res = Result{
			Data:   resData,
			Result: "success",
		}
	} else {
		res = Result{
			Data:   data,
			Result: "success",
		}
	}
	return
}
