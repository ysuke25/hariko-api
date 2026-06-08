package matcher

import (
	"strings"

	"github.com/ysuke25/hariko/internal/config"
)

// RouteMatcher defines an interface for matching HTTP requests to configured routes.
type RouteMatcher interface {
	Match(method, path string) (*config.Route, bool)
}

// StaticMatcher implements RouteMatcher with exact path matching.
type StaticMatcher struct {
	routes []config.Route
}

// NewStaticMatcher creates a new StaticMatcher with the given routes.
func NewStaticMatcher(routes []config.Route) *StaticMatcher {
	return &StaticMatcher{routes: routes}
}

// Match finds a route that matches the given method and path.
// Method matching is case-insensitive. Path matching supports exact matches
// and path parameters (e.g., /api/users/:id).
// Exact matches take priority over parameterized matches.
// Returns the matched route and true if found, nil and false otherwise.
func (m *StaticMatcher) Match(method, path string) (*config.Route, bool) {
	var paramMatch *config.Route

	for i := range m.routes {
		if !strings.EqualFold(m.routes[i].Method, method) {
			continue
		}

		if m.routes[i].Path == path {
			return &m.routes[i], true
		}

		if paramMatch == nil && matchPath(m.routes[i].Path, path) {
			paramMatch = &m.routes[i]
		}
	}

	if paramMatch != nil {
		return paramMatch, true
	}
	return nil, false
}

// matchPath checks if a pattern with path parameters matches a given path.
// Pattern segments starting with ':' are treated as wildcards that match any non-empty value.
func matchPath(pattern, path string) bool {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return false
	}

	for i := range patternParts {
		if strings.HasPrefix(patternParts[i], ":") {
			if pathParts[i] == "" {
				return false
			}
			continue
		}
		if patternParts[i] != pathParts[i] {
			return false
		}
	}

	return true
}
