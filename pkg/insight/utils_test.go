package insight

// import (
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"

// 	"github.com/pipe-cd/pipecd/pkg/model"
// )

// func TestNormalizeTime(t *testing.T) {
// 	type args struct {
// 		from time.Time
// 		step model.InsightStep
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want time.Time
// 	}{
// 		{
// 			name: "formatted correctly with daily",
// 			args: args{
// 				from: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC),
// 				step: model.InsightStep_DAILY,
// 			},
// 			want: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
// 		},
// 		{
// 			name: "formatted correctly with weekly",
// 			args: args{
// 				from: time.Date(2020, 1, 10, 1, 1, 1, 1, time.UTC),
// 				step: model.InsightStep_WEEKLY,
// 			},
// 			want: time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC),
// 		},
// 		{
// 			name: "formatted correctly with monthly",
// 			args: args{
// 				from: time.Date(2020, 1, 7, 1, 1, 1, 1, time.UTC),
// 				step: model.InsightStep_MONTHLY,
// 			},
// 			want: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
// 		},
// 		{
// 			name: "formatted correctly with yearly",
// 			args: args{
// 				from: time.Date(2020, 7, 7, 1, 1, 1, 1, time.UTC),
// 				step: model.InsightStep_YEARLY,
// 			},
// 			want: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := NormalizeTime(tt.args.from, tt.args.step)
// 			assert.Equal(t, got, tt.want)
// 		})
// 	}
// }
