package insightcollector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectorMode_EnableChangeFailureRate(t *testing.T) {
	tests := []struct {
		name string
		m    CollectorMode
		want bool
	}{
		{
			name: "return true",
			m: func() CollectorMode {
				m := NewCollectorMode()
				m.Set("C")
				return m
			}(),
			want: true,
		},
		{
			name: "return false",
			m: func() CollectorMode {
				m := NewCollectorMode()
				m.Set("D")
				return m
			}(),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.EnableChangeFailureRate()
			assert.Equal(t, tt.want, got)
		})
	}
}
func TestCollectorMode_EnableDevelopmentFrequency(t *testing.T) {
	tests := []struct {
		name string
		m    CollectorMode
		want bool
	}{
		{
			name: "return true",
			m: func() CollectorMode {
				m := NewCollectorMode()
				m.Set("D")
				return m
			}(),
			want: true,
		},
		{
			name: "return false",
			m: func() CollectorMode {
				m := NewCollectorMode()
				m.Set("C")
				return m
			}(),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.EnableDevelopmentFrequency()
			assert.Equal(t, tt.want, got)
		})
	}
}
func TestCollectorMode_EnableApplicationCount(t *testing.T) {
	tests := []struct {
		name string
		m    CollectorMode
		want bool
	}{
		{
			name: "return true",
			m: func() CollectorMode {
				m := NewCollectorMode()
				m.Set("A")
				return m
			}(),
			want: true,
		},
		{
			name: "return false",
			m: func() CollectorMode {
				m := NewCollectorMode()
				m.Set("D")
				return m
			}(),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.EnableApplicationCount()
			assert.Equal(t, tt.want, got)
		})
	}
}
