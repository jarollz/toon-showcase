package presenter

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"toon-showcase/internal/core/usecase"
)

type progressBar struct {
	total   int64
	done    int64
	stageMu sync.RWMutex
	stage   string
	stopCh  chan struct{}
	stopped chan struct{}
	started bool
	startMu sync.Mutex
}

func NewProgressBar(total int64) usecase.Progress {
	if total <= 0 {
		total = 1
	}
	return &progressBar{
		total:   total,
		stopCh:  make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

func (p *progressBar) SetStage(stage string) {
	p.stageMu.Lock()
	p.stage = stage
	p.stageMu.Unlock()
}

func (p *progressBar) Advance(n int64) {
	if n <= 0 {
		return
	}
	atomic.AddInt64(&p.done, n)
}

func (p *progressBar) Start() {
	p.startMu.Lock()
	defer p.startMu.Unlock()
	if p.started {
		return
	}
	p.started = true
	go p.renderLoop()
}

func (p *progressBar) Finish() {
	p.startMu.Lock()
	started := p.started
	p.startMu.Unlock()
	if !started {
		return
	}
	atomic.StoreInt64(&p.done, p.total)
	close(p.stopCh)
	<-p.stopped
}

func (p *progressBar) renderLoop() {
	ticker := time.NewTicker(120 * time.Millisecond)
	defer ticker.Stop()
	defer close(p.stopped)

	for {
		select {
		case <-ticker.C:
			p.render(false)
		case <-p.stopCh:
			p.render(true)
			return
		}
	}
}

func (p *progressBar) render(final bool) {
	done := atomic.LoadInt64(&p.done)
	total := p.total
	if done > total {
		done = total
	}
	percent := (float64(done) / float64(total)) * 100
	if total == 0 {
		percent = 100
	}
	if final {
		percent = 100
		done = total
	}

	p.stageMu.RLock()
	stage := p.stage
	p.stageMu.RUnlock()

	const barWidth = 30
	filled := int((float64(done) / float64(total)) * barWidth)
	if final {
		filled = barWidth
	}
	if filled > barWidth {
		filled = barWidth
	}
	if filled < 0 {
		filled = 0
	}
	bar := strings.Repeat("#", filled) + strings.Repeat("-", barWidth-filled)

	if stage == "" {
		stage = "processing"
	}

	fmt.Fprintf(os.Stderr, "\rProgress [%s] %6.2f%% (%d/%d) %s", bar, percent, done, total, stage)
	if final {
		fmt.Fprintln(os.Stderr)
	}
}
