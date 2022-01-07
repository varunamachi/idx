package httpw

import (
	"fmt"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/varunamachi/sause/pkg/auth"
)

type noopWriter struct{}

func (nw *noopWriter) Write(b []byte) (int, error) {
	return 0, nil
}

type Endpoint struct {
	Method     string
	Path       string
	Category   string
	Desc       string
	Version    string
	Role       auth.Role
	Permission string
	Route      *echo.Route
	Handler    echo.HandlerFunc
}

func (ep *Endpoint) NeedsAuth() bool {
	return ep.Permission != "" && ep.Role != auth.None && ep.Role != ""
}

type Server struct {
	categories    map[string][]*Endpoint
	endpoints     []*Endpoint
	root          *echo.Echo
	printer       io.Writer
	userRetriever auth.UserRetrieverFunc
}

func NewServer(printer io.Writer, userGetter auth.UserRetrieverFunc) *Server {
	if printer == nil {
		printer = &noopWriter{}
	}
	return &Server{
		categories:    make(map[string][]*Endpoint),
		endpoints:     make([]*Endpoint, 0, 100),
		root:          echo.New(),
		printer:       printer,
		userRetriever: userGetter,
	}
}

func (s *Server) AddEndpoints(ep ...*Endpoint) *Server {
	s.endpoints = append(s.endpoints, ep...)
	return s
}

func (s *Server) Start(port uint32) error {
	s.configure()
	s.Print()
	return s.root.Start(fmt.Sprintf(":%d", port))
}

func (s *Server) configure() {

	type groupPair struct {
		versionGrp *echo.Group
		inGrp      *echo.Group
	}
	groups := map[string]*groupPair{}

	authMw := getAuthMiddleware(s)
	for _, ep := range s.endpoints {
		ep := ep

		grp := groups[ep.Version]
		if grp == nil {
			grp = &groupPair{}
			grp.versionGrp = s.root.Group("api/" + ep.Version)
			grp.inGrp = grp.versionGrp.Group("in")
			grp.inGrp.Use(authMw)
		}

		if ep.NeedsAuth() {
			ep.Route = grp.inGrp.Add(
				ep.Method, ep.Path, ep.Handler, getAccessMiddleware(ep, s))

		} else {
			ep.Route = grp.versionGrp.Add(
				ep.Method, ep.Path, ep.Handler, getAccessMiddleware(ep, s))
		}

		if _, found := s.categories[ep.Category]; !found {
			s.categories[ep.Category] = make([]*Endpoint, 0, 20)
		}
		s.categories[ep.Category] = append(s.categories[ep.Category], ep)
	}
}

func (s *Server) Print() {
	for _, ep := range s.endpoints {
		cat := ep.Category
		if len(cat) > 14 {
			cat = ep.Category[:14]
		}
		fmt.Fprintf(s.printer,
			"%-3s %-5s %-40s %-15s %s\n",
			ep.Version, ep.Route.Method, ep.Route.Path, cat, ep.Desc)
	}
	fmt.Print("\n\n")
}
