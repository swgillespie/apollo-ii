package engine

import "runtime"

type Options struct {
	debugChecks              bool
	moveGenerationBufferSize int
	workQueueBufferSize      int
	workQueueNumGoroutines   int
}

var options = Options{debugChecks: false,
	moveGenerationBufferSize: 128,
	workQueueBufferSize:      1024,
	workQueueNumGoroutines:   runtime.NumCPU()}
