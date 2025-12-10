package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/amycatgirl/codehub/github"
	custom_template "github.com/amycatgirl/codehub/template"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data any, c echo.Context) error {
	err := t.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		fmt.Printf("Error while rendering template %s: %v\n", name, err)
	}
	return err
}

type Server struct {
	gh    *github.Client
	echo  *echo.Echo
	httpd *http.Server
}

type Args struct {
	Addr        string
	GithubToken string
}

func (s *Server) Serve() {
	s.addRoutes()
	s.httpd.ListenAndServe()
}

func handleLanding(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

func (s *Server) handleRepo(c echo.Context) error {
	user := c.Param("user")
	repo := c.Param("repo")
	info, err := s.gh.QuerySingleRepository(user, repo)

	if err != nil {
		fmt.Printf("%v\n", err)
		return fmt.Errorf("Failed to query repo %s/%s: %w", user, repo, err)
	}

	if info == nil {
		return fmt.Errorf("Missing repository metadata, somehow")
	}

	readme, err := s.gh.QueryReadmeFromRepo(user, repo)
	if err != nil {
		fmt.Printf("Error while querying readme: %v", err)
		return err
	}

	fmt.Printf("Rendering repo with info: %v\n", info)

	return c.Render(http.StatusOK, "repo.html", &GithubRepo{
		Repository: *info,
		Readme:     readme,
	})
}

type GithubRepo struct {
	github.Repository
	Readme *string
}

type GithubUser struct {
	github.User
	Repositories []github.Repository
}

func (s *Server) handleUser(c echo.Context) error {
	user := c.Param("user")
	profile, err := s.gh.QueryProfile(user)
	if err != nil {
		return fmt.Errorf("Failed to query profile for user %s: %w", user, err)
	}

	repos, err := s.gh.QueryRepositories(user)
	if err != nil {
		return fmt.Errorf("Failed to query repos for user %s: %w", user, err)
	}

	return c.Render(http.StatusOK, "user.html", GithubUser{
		User:         *profile,
		Repositories: *repos,
	})
}

func (s *Server) addRoutes() {
	s.echo.GET("/", handleLanding)
	s.echo.Static("/static", "public")
	s.echo.GET("/:user/:repo", s.handleRepo)
	s.echo.GET("/:user", s.handleUser)
}

func New(args Args) (*Server, error) {
	t := &Template{
		templates: template.Must(custom_template.RegisterTemplatesWithSanitizer("templates/*.html")),
	}

	e := echo.New()
	e.Renderer = t
	e.Use(middleware.Logger())
	httpd := http.Server{
		Addr:    args.Addr,
		Handler: e,
	}

	ghc, err := github.New(args.GithubToken)
	if err != nil {
		return nil, err
	}

	s := Server{
		echo:  e,
		httpd: &httpd,
		gh:    ghc,
	}
	return &s, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file, see error below:")
		panic(err)
	}

	token := os.Getenv("GH_API_TOKEN")
	s, err := New(Args{
		Addr:        ":8080",
		GithubToken: token,
	})
	if err != nil {
		panic(err)
	}

	s.Serve()
}
