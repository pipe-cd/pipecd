// Copyright (C) 2014-2019, Matt Butcher and Matt Farina

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

package semver

// ByNewer is a collection of Version instances and implements the sort
// interface. See the sort package for more details.
// https://golang.org/pkg/sort/
type ByNewer []*Version

// Len returns the length of a collection. The number of Version instances
// on the slice.
func (c ByNewer) Len() int {
	return len(c)
}

// Less checks if one is greater than the other because of sorting by newer version.
func (c ByNewer) Less(i, j int) bool {
	return c[i].GreaterThan(c[j])
}

// Swap is needed for the sort interface to replace the Version objects
// at two different positions in the slice.
func (c ByNewer) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
