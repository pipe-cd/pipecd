// Copyright 2020 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package insight

// import (
// 	"time"

// 	"github.com/pipe-cd/pipecd/pkg/model"
// )

// func NormalizeTime(from time.Time, step model.InsightStep) time.Time {
// 	var formattedTime time.Time
// 	switch step {
// 	case model.InsightStep_DAILY:
// 		formattedTime = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
// 	case model.InsightStep_WEEKLY:
// 		// Sunday in the week of rangeFrom
// 		sunday := from.AddDate(0, 0, -int(from.Weekday()))
// 		formattedTime = time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 0, 0, 0, 0, time.UTC)
// 	case model.InsightStep_MONTHLY:
// 		formattedTime = time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, time.UTC)
// 	case model.InsightStep_YEARLY:
// 		formattedTime = time.Date(from.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
// 	}
// 	return formattedTime
// }
