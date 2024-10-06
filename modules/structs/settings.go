// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package structs

// GeneralUISettings contains global ui settings exposed by API
type GeneralUISettings struct {
	DefaultTheme     string   `json:"default_theme"`
	AllowedReactions []string `json:"allowed_reactions"`
	CustomEmojis     []string `json:"custom_emojis"`
}

// GeneralAPISettings contains global api settings exposed by it
type GeneralAPISettings struct {
	MaxResponseItems       int   `json:"max_response_items"`
	DefaultPagingNum       int   `json:"default_paging_num"`
	DefaultGitTreesPerPage int   `json:"default_git_trees_per_page"`
	DefaultMaxBlobSize     int64 `json:"default_max_blob_size"`
}
