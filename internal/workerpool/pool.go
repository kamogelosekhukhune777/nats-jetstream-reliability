package workerpool

import "sync"

// 5. Worker pool consumers
type Pool struct {
	wg sync.WaitGroup
}

func (p *Pool) Run(workers int, fn func()) {
	for i := 0; i < workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			fn()
		}()
	}
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
