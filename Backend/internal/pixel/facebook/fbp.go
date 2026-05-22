package facebook

import (
	"fmt"
	"math/rand"
	"time"
)

// EnsureFBP returns existing fbp or generates first-party browser id (Meta format).
func EnsureFBP(fbp string) string {
	fbp = trimSpace(fbp)
	if fbp != "" {
		return fbp
	}
	return fmt.Sprintf("fb.1.%d.%d", time.Now().Unix(), rand.Int63n(1_000_000_000_000))
}
