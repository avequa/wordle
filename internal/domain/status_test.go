package domain

import "testing"

func TestStatusIsFinished(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   bool
	}{
		{"игра не начата", StatusUnknown, false},
		{"игра идет", StatusInProgress, false},
		{"победа", StatusWin, true},
		{"поражение", StatusLose, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsFinished(); got != tt.want {
				t.Fatalf("получено %v, ожидалось %v", got, tt.want)
			}
		})
	}
}
