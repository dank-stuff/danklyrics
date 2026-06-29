package api

const docsLink = "https://codeberg.org/dankstuff/danklyrics"

type errorResponse struct {
	Message         string `json:"message"`
	SuggestedAction string `json:"suggested_action,omitempty"`
	DocsLink        string `json:"docs_link,omitempty"`
}
