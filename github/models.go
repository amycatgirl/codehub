package github

type User struct {
	Username    string  `json:"login"`
	DisplayName *string `json:"name,omitempty"`
	Biography   *string `json:"bio,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	GravatarURL *string `json:"gravatar_url,omitempty"`
	Pronouns    *string `json:"pronouns,omitempty"`
}

type Repository struct {
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Archived    bool      `json:"archived"`
	Forks       int       `json:"forks_count"`
	Stars       int       `json:"stargazers_count"`
	Description *string   `json:"description,omitempty"`
	Topics      *[]string `json:"topics,omitempty"`
	Owner       User      `json:"owner"`
}

type internalFile struct {
	Name     string `json:"name"`
	Content  string `json:"Content"`
	Encoding string `json:"Encoding"`
}

type File struct {
	Name    string
	Content string
}
