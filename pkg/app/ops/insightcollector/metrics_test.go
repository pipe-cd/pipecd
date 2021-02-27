package insightcollector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectorMetrics_IsEnabled(t *testing.T) {
	tests := []struct {
		name string
		m    CollectorMetrics
		args CollectorMetrics
		want bool
	}{
		{
			name: "return true",
			m: func() CollectorMetrics {
				m := NewCollectorMetrics()
				m.Enable(ChangeFailureRate)
				return m
			}(),
			args: ChangeFailureRate,
			want: true,
		},
		{
			name: "return false",
			m: func() CollectorMetrics {
				m := NewCollectorMetrics()
				m.Enable(DevelopmentFrequency)
				return m
			}(),
			args: ChangeFailureRate,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.IsEnabled(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
