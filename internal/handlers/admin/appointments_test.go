package admin

import (
	"testing"
	"time"
)

func TestNoofDays(t *testing.T) {
	days := noOfDays(time.Now())
	if days < 28 || days > 31 {
		t.Fatalf("invalid days:%d\t received", days)
	}

}
