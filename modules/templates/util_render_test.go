// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package templates

import (
	"context"
	"os"
	"strings"
	"testing"

	"code.gitea.io/gitea/models/unittest"
	"code.gitea.io/gitea/modules/markup"

	"github.com/stretchr/testify/assert"
)

func testInput() string {
	s := `  space @mention-user<SPACE><SPACE>
/just/a/path.bin
https://example.com/file.bin
[local link](file.bin)
[remote link](https://example.com)
[[local link|file.bin]]
[[remote link|https://example.com]]
![local image](image.jpg)
![remote image](https://example.com/image.jpg)
[[local image|image.jpg]]
[[remote link|https://example.com/image.jpg]]
https://example.com/user/repo/compare/88fc37a3c0a4dda553bdcfc80c178a58247f42fb...12fc37a3c0a4dda553bdcfc80c178a58247f42fb#hash
com 88fc37a3c0a4dda553bdcfc80c178a58247f42fb...12fc37a3c0a4dda553bdcfc80c178a58247f42fb pare
https://example.com/user/repo/commit/88fc37a3c0a4dda553bdcfc80c178a58247f42fb
com 88fc37a3c0a4dda553bdcfc80c178a58247f42fb mit
:+1:
mail@domain.com
@mention-user test
#123
  space<SPACE><SPACE>
`
	return strings.ReplaceAll(s, "<SPACE>", " ")
}

var testMetas = map[string]string{
	"user":     "user13",
	"repo":     "repo11",
	"repoPath": "../../tests/gitea-repositories-meta/user13/repo11.git/",
	"mode":     "comment",
}

func TestMain(m *testing.M) {
	unittest.InitSettings()
	markup.Init(&markup.ProcessorHelper{
		IsUsernameMentionable: func(ctx context.Context, username string) bool {
			return username == "mention-user"
		},
	})
	os.Exit(m.Run())
}

func TestRenderMarkdownToHtml(t *testing.T) {
	expected := `<p>space <a href="/mention-user" rel="nofollow">@mention-user</a><br/>
/just/a/path.bin
<a href="https://example.com/file.bin" rel="nofollow">https://example.com/file.bin</a>
<a href="/file.bin" rel="nofollow">local link</a>
<a href="https://example.com" rel="nofollow">remote link</a>
<a href="/file.bin" rel="nofollow">local link</a>
<a href="https://example.com" rel="nofollow">remote link</a>
<a href="/image.jpg" target="_blank" rel="nofollow noopener"><img src="/image.jpg" alt="local image"/></a>
<a href="https://example.com/image.jpg" target="_blank" rel="nofollow noopener"><img src="https://example.com/image.jpg" alt="remote image"/></a>
<a href="/image.jpg" rel="nofollow"><img src="/image.jpg" title="local image" alt="local image"/></a>
<a href="https://example.com/image.jpg" rel="nofollow"><img src="https://example.com/image.jpg" title="remote link" alt="remote link"/></a>
<a href="https://example.com/user/repo/compare/88fc37a3c0a4dda553bdcfc80c178a58247f42fb...12fc37a3c0a4dda553bdcfc80c178a58247f42fb#hash" rel="nofollow"><code>88fc37a3c0...12fc37a3c0 (hash)</code></a>
com 88fc37a3c0a4dda553bdcfc80c178a58247f42fb...12fc37a3c0a4dda553bdcfc80c178a58247f42fb pare
<a href="https://example.com/user/repo/commit/88fc37a3c0a4dda553bdcfc80c178a58247f42fb" rel="nofollow"><code>88fc37a3c0</code></a>
com 88fc37a3c0a4dda553bdcfc80c178a58247f42fb mit
<span class="emoji" aria-label="thumbs up">👍</span>
<a href="mailto:mail@domain.com" rel="nofollow">mail@domain.com</a>
<a href="/mention-user" rel="nofollow">@mention-user</a> test
#123
space</p>
`
	assert.Equal(t, expected, string(RenderMarkdownToHtml(context.Background(), testInput())))
}

func TestUserMention(t *testing.T) {
	rendered := RenderMarkdownToHtml(context.Background(), "@no-such-user @mention-user @mention-user")
	assert.EqualValues(t, `<p>@no-such-user <a href="/mention-user" rel="nofollow">@mention-user</a> <a href="/mention-user" rel="nofollow">@mention-user</a></p>`, strings.TrimSpace(string(rendered)))
}
