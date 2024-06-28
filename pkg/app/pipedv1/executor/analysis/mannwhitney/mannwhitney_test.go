// Copyright 2024 The PipeCD Authors.
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

// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mannwhitney

import "testing"

func TestMannWhitneyUTest(t *testing.T) {
	t.Parallel()

	check := func(want, got *MannWhitneyUTestResult) {
		if want.N1 != got.N1 || want.N2 != got.N2 ||
			!aeq(t, want.U, got.U) ||
			want.AltHypothesis != got.AltHypothesis ||
			!aeq(t, want.P, got.P) {
			t.Errorf("want %+v, got %+v", want, got)
		}
	}
	check3 := func(x1, x2 []float64, U float64, pless, pdiff, pgreater float64) {
		want := &MannWhitneyUTestResult{N1: len(x1), N2: len(x2), U: U}

		want.AltHypothesis = LocationLess
		want.P = pless
		got, _ := MannWhitneyUTest(x1, x2, want.AltHypothesis)
		check(want, got)

		want.AltHypothesis = LocationDiffers
		want.P = pdiff
		got, _ = MannWhitneyUTest(x1, x2, want.AltHypothesis)
		check(want, got)

		want.AltHypothesis = LocationGreater
		want.P = pgreater
		got, _ = MannWhitneyUTest(x1, x2, want.AltHypothesis)
		check(want, got)
	}

	s1 := []float64{2, 1, 3, 5}
	s2 := []float64{12, 11, 13, 15}
	s3 := []float64{0, 4, 6, 7} // Interleaved with s1, but no ties
	s4 := []float64{2, 2, 2, 2}
	s5 := []float64{1, 1, 1, 1, 1}

	// Small sample, no ties
	check3(s1, s2, 0, 0.014285714285714289, 0.028571428571428577, 1)
	check3(s2, s1, 16, 1, 0.028571428571428577, 0.014285714285714289)
	check3(s1, s3, 5, 0.24285714285714288, 0.485714285714285770, 0.8285714285714285)

	// Small sample, ties
	// TODO: Check these against some other implementation.
	check3(s1, s1, 8, 0.6285714285714286, 1, 0.6285714285714286)
	check3(s1, s4, 10, 0.8571428571428571, 0.7142857142857143, 0.3571428571428571)
	check3(s1, s5, 17.5, 1, 0, 0.04761904761904767)

	r, err := MannWhitneyUTest(s4, s4, LocationDiffers)
	if err != ErrSamplesEqual {
		t.Errorf("want ErrSamplesEqual, got %+v, %+v", r, err)
	}

	// Large samples.
	l1 := make([]float64, 500)
	for i := range l1 {
		l1[i] = float64(i * 2)
	}
	l2 := make([]float64, 600)
	for i := range l2 {
		l2[i] = float64(i*2 - 41)
	}
	l3 := append([]float64{}, l2...)
	for i := 0; i < 30; i++ {
		l3[i] = l1[i]
	}
	// For comparing with R's wilcox.test:
	// l1 <- seq(0, 499)*2
	// l2 <- seq(0,599)*2-41
	// l3 <- l2; for (i in 1:30) { l3[i] = l1[i] }

	check3(l1, l2, 135250, 0.0024667680407086112, 0.0049335360814172224, 0.9975346930458906)
	check3(l1, l1, 125000, 0.5000436801680628, 1, 0.5000436801680628)
	check3(l1, l3, 134845, 0.0019351907119808942, 0.0038703814239617884, 0.9980659818257166)
}
