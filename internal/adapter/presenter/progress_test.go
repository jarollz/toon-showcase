package presenter

import (
	"testing"
	"time"
)

func TestProgressBarLifecycle(t *testing.T) {
	pAny := NewProgressBar(10)
	p, ok := pAny.(*progressBar)
	if !ok {
		t.Fatalf("NewProgressBar returned unexpected type")
	}

	p.SetStage("stage-1")
	p.Advance(0)
	p.Advance(2)
	p.Start()
	p.Start()
	time.Sleep(5 * time.Millisecond)
	p.Finish()

	if !p.started {
		t.Fatalf("progress should be started")
	}
}

func TestProgressBarNoStartFinish(t *testing.T) {
	pAny := NewProgressBar(0)
	p := pAny.(*progressBar)
	p.Finish()
	p.render(true)
}
