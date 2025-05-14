//nolint:exhaustruct
package main

import (
	"reflect"
	"strings"
	"testing"
)

//nolint:paralleltest
func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		inputBody     string
		expected      []string
		errorContains string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{
				"https://blog.boot.dev/path/one",
				"https://other.com/path/one",
			},
		},
		{
			name:     "no URLs",
			inputURL: "https://example.com",
			inputBody: `
<html>
	<body>
		<p>No links here.</p>
	</body>
</html>
`,
			expected: []string{},
		},
		{
			name:     "only relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="/path/one">One</a>
		<a href="/path/two">Two</a>
	</body>
</html>
`,
			expected: []string{
				"https://blog.boot.dev/path/one",
				"https://blog.boot.dev/path/two",
			},
		},
		{
			name:     "only absolute URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="https://example.com/path/one">One</a>
		<a href="https://example.com/path/two">Two</a>
	</body>
</html>
`,
			expected: []string{
				"https://example.com/path/one",
				"https://example.com/path/two",
			},
		},
		{
			name:     "URLs with query params and fragments",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="/path/one?param=value#section">One</a>
		<a href="https://example.com/path/two?param=value#section">Two</a>
	</body>
</html>
`,
			expected: []string{
				"https://blog.boot.dev/path/one?param=value#section",
				"https://example.com/path/two?param=value#section",
			},
		},
		{
			name:     "malformed HTML",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="/path/one">One
		<a href="https://example.com/path/two">Two</a>
	</body>
</html>
`,
			expected: []string{
				"https://blog.boot.dev/path/one",
				"https://example.com/path/two",
			},
		},
		{
			name:     "relative URLs with parent directory",
			inputURL: "https://blog.boot.dev/subdir/",
			inputBody: `
<html>
	<body>
		<a href="../parent">Parent</a>
		<a href="./child">Child</a>
	</body>
</html>
`,
			expected: []string{
				"https://blog.boot.dev/parent",
				"https://blog.boot.dev/subdir/child",
			},
		},
		{
			name:     "URLs with ports",
			inputURL: "https://blog.boot.dev:8080",
			inputBody: `
<html>
	<body>
		<a href="/path/one">One</a>
		<a href="https://example.com:9090/path/two">Two</a>
	</body>
</html>
`,
			expected: []string{
				"https://blog.boot.dev:8080/path/one",
				"https://example.com:9090/path/two",
			},
		},
		{
			name:          "empty input",
			inputURL:      "",
			inputBody:     "",
			expected:      []string{},
			errorContains: "base URL cannot be empty",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil {
				if tc.errorContains != "" &&
					strings.Contains(err.Error(), tc.errorContains) {
					return
				}

				t.Errorf(
					"Test %v - '%s' FAIL: unexpected error: %v",
					i,
					tc.name,
					err,
				)
				return
			}

			if len(actual) > 0 && len(tc.expected) > 0 &&
				!reflect.DeepEqual(actual, tc.expected) {
				t.Errorf(
					"Test %v - %s FAIL: expected URL: %v, actual: %v",
					i,
					tc.name,
					tc.expected,
					actual,
				)
			}
		})
	}
}
