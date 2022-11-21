package gin

import (
	"fmt"
	"github.com/gofrs/uuid"
	"time"
)

type rsTemplate struct {
	correlation *uuid.UUID
	status      int
	size        int
	latency     time.Duration
	body        string
}

func (t *rsTemplate) String() string {
	return fmt.Sprintf("%v %s %d %d bytes %dms", t.correlation, response, t.status, t.size, t.latency.Milliseconds())
}

func (t *rsTemplate) fullString() string {
	if t.size > 0 {
		return fmt.Sprintf("%s, body: %s", t.String(), t.body)
	}
	return fmt.Sprintf("%s, %s", t.String(), noBody)
}
