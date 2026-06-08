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
// Method matching is case-insensitive. Path matching is exact.
// Returns the matched route and true if found, nil and false otherwise.
func (m *StaticMatcher) Match(method, path string) (*config.Route, bool) {
	for i := range m.routes {
		if strings.EqualFold(m.routes[i].Method, method) && m.routes[i].Path == path {
			return &m.routes[i], true
		}
	}
	return nil, false
}
