package patient

import (
	"testing"
	"time"
)

func TestStartEndTime(t *testing.T) {
	tests := []struct {
		name   string
		expErr bool
		inTime time.Time
	}{
		{
			name:   "valid time",
			inTime: time.Now(),
		},
		{
			name:   "in valid time",
			inTime: time.Time{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := startTime(test.inTime)
			if test.expErr {
				if err == nil {
					t.Fatal("expected error got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("no err expected\tgot:%v\n", err)
				}
			}
			t.Logf("starttime:%s\n", got.Format(time.DateTime))
			if test.expErr && !got.IsZero() {
				t.Fatalf("expected zero time\tgot value:%s\n", got.Format(time.DateTime))

			}
			got1, err1 := endTime(test.inTime)
			if test.expErr {
				if err1 == nil {
					t.Fatal("expected error got nil")
				}
			} else {
				if err1 != nil {
					t.Fatalf("no err expected\tgot:%v\n", err)
				}
			}
			t.Logf("endTime:%s\n", got1.Format(time.DateTime))
			if test.expErr && !got1.IsZero() {
				t.Fatalf("expected zero time\tgot value:%s\n", got1.Format(time.DateTime))

			}
		})
	}

}
