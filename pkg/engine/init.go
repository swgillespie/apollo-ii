package engine

import "sync"

var initializeGuard sync.Once

// Initialize performs all engine-wide initialization that needs
// to occur before engine operations can begin.
func Initialize() {
	initializeGuard.Do(func() {
		initializeAttackTables()
	})
}
