package webfinger

const ContentType string = "application/jrd+json"

type Resource struct {
	Subject    string            `json:"subject,omitempty"`
	Aliases    []string          `json:"aliases,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
	Links      []Link            `json:"links"`
}

type Link struct {
	HRef       string             `json:"href"`
	Type       string             `json:"type,omitempty"`
	Rel        string             `json:"rel"`
	Properties map[string]*string `json:"properties,omitempty"`
	Titles     map[string]string  `json:"titles,omitempty"`
}
