package value_object

import "testing"

func TestNewExecutionStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		wantErr bool
	}{
		{
			name:    "valid in_progress",
			status:  "in_progress",
			wantErr: false,
		},
		{
			name:    "valid completed",
			status:  "completed",
			wantErr: false,
		},
		{
			name:    "valid failed",
			status:  "failed",
			wantErr: false,
		},
		{
			name:    "invalid status",
			status:  "invalid",
			wantErr: true,
		},
		{
			name:    "empty status",
			status:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewExecutionStatus(tt.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewExecutionStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.status {
				t.Errorf("NewExecutionStatus() = %v, want %v", got, tt.status)
			}
		})
	}
}

func TestExecutionStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status ExecutionStatus
		want   bool
	}{
		{
			name:   "in_progress is valid",
			status: ExecutionStatusInProgress,
			want:   true,
		},
		{
			name:   "completed is valid",
			status: ExecutionStatusCompleted,
			want:   true,
		},
		{
			name:   "failed is valid",
			status: ExecutionStatusFailed,
			want:   true,
		},
		{
			name:   "invalid status",
			status: ExecutionStatus("invalid"),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.want {
				t.Errorf("ExecutionStatus.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecutionStatus_StatusChecks(t *testing.T) {
	tests := []struct {
		name         string
		status       ExecutionStatus
		isInProgress bool
		isCompleted  bool
		isFailed     bool
	}{
		{
			name:         "in_progress status",
			status:       ExecutionStatusInProgress,
			isInProgress: true,
			isCompleted:  false,
			isFailed:     false,
		},
		{
			name:         "completed status",
			status:       ExecutionStatusCompleted,
			isInProgress: false,
			isCompleted:  true,
			isFailed:     false,
		},
		{
			name:         "failed status",
			status:       ExecutionStatusFailed,
			isInProgress: false,
			isCompleted:  false,
			isFailed:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsInProgress(); got != tt.isInProgress {
				t.Errorf("ExecutionStatus.IsInProgress() = %v, want %v", got, tt.isInProgress)
			}
			if got := tt.status.IsCompleted(); got != tt.isCompleted {
				t.Errorf("ExecutionStatus.IsCompleted() = %v, want %v", got, tt.isCompleted)
			}
			if got := tt.status.IsFailed(); got != tt.isFailed {
				t.Errorf("ExecutionStatus.IsFailed() = %v, want %v", got, tt.isFailed)
			}
		})
	}
}

func TestExecutionStatus_Equals(t *testing.T) {
	status1, _ := NewExecutionStatus("in_progress")
	status2, _ := NewExecutionStatus("completed")
	status1Copy, _ := NewExecutionStatus("in_progress")

	if !status1.Equals(status1Copy) {
		t.Error("Same statuses should be equal")
	}
	if status1.Equals(status2) {
		t.Error("Different statuses should not be equal")
	}
}
