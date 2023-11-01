// Copyright 2023 The PipeCD Authors.
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

// Copyright 2013-2018 Docker, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filematcher

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/scanner"
)

type RegexpProvider func(string) (*regexp.Regexp, error)

// PatternMatcher allows checking paths against a list of patterns.
type PatternMatcher struct {
	patterns       []*Pattern
	exclusions     []*Pattern
	regexpProvider RegexpProvider
}

// An Option configures a PatternMatcher.
type Option func(*PatternMatcher)

// WithRegexpProvider sets a custom regexp provider.
func WithRegexpProvider(p RegexpProvider) Option {
	return func(pm *PatternMatcher) {
		pm.regexpProvider = p
	}
}

// NewPatternMatcher creates a new matcher object for specific patterns that can
// be used later to match against patterns against paths.
func NewPatternMatcher(patterns []string, opts ...Option) (*PatternMatcher, error) {
	pm := &PatternMatcher{
		regexpProvider: regexp.Compile,
	}
	for _, opt := range opts {
		opt(pm)
	}
	for _, p := range patterns {
		// Eliminate leading and trailing whitespace.
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		p = filepath.Clean(p)
		negative := false
		if p[0] == '!' {
			if len(p) == 1 {
				return nil, errors.New("illegal exclusion pattern: \"!\"")
			}
			p = p[1:]
			negative = true
		}
		// Do some syntax checking on the pattern.
		// filepath's Match() has some really weird rules that are inconsistent
		// so instead of trying to dup their logic, just call Match() for its
		// error state and if there is an error in the pattern return it.
		// If this becomes an issue we can remove this since its really only
		// needed in the error (syntax) case - which isn't really critical.
		if _, err := filepath.Match(p, "."); err != nil {
			return nil, err
		}
		newp := &Pattern{
			cleanedPattern: p,
			dirs:           strings.Split(p, string(os.PathSeparator)),
		}
		regexp, err := pm.regexpProvider(newp.regexpString())
		if err != nil {
			return nil, filepath.ErrBadPattern
		}
		newp.regexp = regexp
		if negative {
			pm.exclusions = append(pm.exclusions, newp)
		} else {
			pm.patterns = append(pm.patterns, newp)
		}
	}
	return pm, nil
}

// Matches matches path against all the patterns. Matches is not safe to be
// called concurrently.
func (pm *PatternMatcher) Matches(file string) bool {
	matched := matches(file, pm.exclusions)
	if matched {
		return false
	}
	return matches(file, pm.patterns)
}

func (pm *PatternMatcher) MatchesAny(files []string) bool {
	for _, file := range files {
		if pm.Matches(file) {
			return true
		}
	}
	return false
}

func matches(file string, patterns []*Pattern) bool {
	file = filepath.FromSlash(file)
	parentPath := filepath.Dir(file)
	parentPathDirs := strings.Split(parentPath, string(os.PathSeparator))

	for _, pattern := range patterns {
		matched := pattern.regexp.MatchString(file)
		if !matched && parentPath != "." {
			// Check to see if the pattern matches one of our parent dirs.
			if len(pattern.dirs) <= len(parentPathDirs) {
				matched = pattern.regexp.MatchString(strings.Join(parentPathDirs[:len(pattern.dirs)], string(os.PathSeparator)))
			}
		}
		if matched {
			return true
		}
	}
	return false
}

// Exclusions returns array of negative patterns.
func (pm *PatternMatcher) Exclusions() []*Pattern {
	return pm.exclusions
}

// Patterns returns array of active patterns.
func (pm *PatternMatcher) Patterns() []*Pattern {
	return pm.patterns
}

// Pattern defines a single regexp used to filter file paths.
type Pattern struct {
	cleanedPattern string
	dirs           []string
	regexp         *regexp.Regexp
}

func (p *Pattern) String() string {
	return p.cleanedPattern
}

func (p *Pattern) regexpString() string {
	regStr := "^"
	pattern := p.cleanedPattern
	// Go through the pattern and convert it to a regexp.
	// We use a scanner so we can support utf-8 chars.
	var scan scanner.Scanner
	scan.Init(strings.NewReader(pattern))

	sl := string(os.PathSeparator)
	escSL := sl
	if sl == `\` {
		escSL += `\`
	}

	for scan.Peek() != scanner.EOF {
		ch := scan.Next()

		switch ch {
		case '*':
			if scan.Peek() == '*' {
				// Is some flavor of "**".
				scan.Next()

				// Treat **/ as ** so eat the "/".
				if string(scan.Peek()) == sl {
					scan.Next()
				}

				if scan.Peek() == scanner.EOF {
					// Is "**EOF" - to align with .gitignore just accept all.
					regStr += ".*"
				} else {
					// Is "**".
					// Note that this allows for any # of /'s (even 0) because
					// the .* will eat everything, even /'s.
					regStr += "(.*" + escSL + ")?"
				}
			} else {
				// Is "*" so map it to anything but "/".
				regStr += "[^" + escSL + "]*"
			}
		case '?':
			// "?" is any char except "/".
			regStr += "[^" + escSL + "]"
		case '.', '$':
			// Escape some regexp special chars that have no meaning
			// in golang's filepath.Match.
			regStr += `\` + string(ch)
		case '\\':
			// Escape next char. Note that a trailing \ in the pattern
			// will be left alone (but need to escape it).
			if sl == `\` {
				// On windows map "\" to "\\", meaning an escaped backslash,
				// and then just continue because filepath.Match on
				// Windows doesn't allow escaping at all.
				regStr += escSL
				continue
			}
			if scan.Peek() != scanner.EOF {
				regStr += `\` + string(scan.Next())
			} else {
				regStr += `\`
			}
		default:
			regStr += string(ch)
		}
	}
	regStr += "$"
	return regStr
}

// Matches returns true if file matches any of the patterns
// and isn't excluded by any of the subsequent patterns.
func Matches(file string, patterns []string, opts ...Option) (bool, error) {
	pm, err := NewPatternMatcher(patterns, opts...)
	if err != nil {
		return false, err
	}
	file = filepath.Clean(file)

	if file == "." {
		// Don't let them exclude everything, kind of silly.
		return false, nil
	}

	return pm.Matches(file), nil
}
