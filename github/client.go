package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	token string
	httpc *http.Client
}

func (c Client) get(url string, accept *string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	c.addRegularHeaders(req)
	if accept != nil {
		req.Header.Add("Accept", *accept)
	}

	return c.httpc.Do(req)
}

func (c Client) addRegularHeaders(request *http.Request) {
	request.Header.Add("User-Agent", "github:amycatgirl/github-ss")
	request.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	// FIXME: REPLACE THE FUCKING TOKEN DON'T LEAK IT DUMBASS
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
}

func (c Client) QueryRepositories(user string) (*[]Repository, error) {
	resp, err := c.get(fmt.Sprintf("https://api.github.com/users/%s/repos?sort=updated&type=all", user), nil)
	if err != nil {
		return nil, err
	}

	var repoList []Repository
	if err := json.NewDecoder(resp.Body).Decode(&repoList); err != nil {
		return nil, fmt.Errorf("Failed to decode: %w", err)
	}

	return &repoList, nil
}

func (c Client) QuerySingleRepository(user string, repository string) (*Repository, error) {
	resp, err := c.get(fmt.Sprintf("https://api.github.com/repos/%s/%s", user, repository), nil)
	if err != nil {
		return nil, err
	}

	var repo Repository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, fmt.Errorf("Failed to decode: %w", err)
	}

	return &repo, nil
}

func (c Client) QueryProfile(user string) (*User, error) {
	resp, err := c.get(fmt.Sprintf("https://api.github.com/users/%s", user), nil)
	if err != nil {
		return nil, err
	}

	var profile User
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("Failed to decode: %w", err)
	}

	return &profile, nil
}

func (c Client) QueryReadmeFromRepo(owner string, repo string) (*string, error) {
	content_type := "application/vnd.github.html+json"
	resp, err := c.get(fmt.Sprintf("https://api.github.com/repos/%s/%s/readme", owner, repo), &content_type)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, nil
	}

	buf := new(strings.Builder)
	n, err := io.Copy(buf, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to copy strbuf: %v", err)
	}
	fmt.Printf("%d bytes read", n)

	content := buf.String()
	return &content, nil
}

type Args struct {
	Token string
}

func New(args *Args) (*Client, error) {
	httpc := http.Client{}
	if args == nil {
		return nil, fmt.Errorf("I NEED A TOKEN TO FUNCTION, DUMBASS")
	}

	client := Client{
		httpc: &httpc,
		token: args.Token,
	}

	return &client, nil
}
