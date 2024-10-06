// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user

// func TestDashboardPagination(t *testing.T) {
// 	ctx, _ := contexttest.MockContext(t, "/", contexttest.MockContextOption{Render: templates.HTMLRenderer()})
// 	page := context.NewPagination(10, 3, 1, 3)

// 	setting.AppSubURL = "/SubPath"
// 	out, err := ctx.RenderToHTML("base/paginate", map[string]any{"Link": setting.AppSubURL, "Page": page})
// 	assert.NoError(t, err)
// 	assert.Contains(t, out, `<a class=" item navigation" href="/SubPath/?page=2">`)

// 	setting.AppSubURL = ""
// 	out, err = ctx.RenderToHTML("base/paginate", map[string]any{"Link": setting.AppSubURL, "Page": page})
// 	assert.NoError(t, err)
// 	assert.Contains(t, out, `<a class=" item navigation" href="/?page=2">`)
// }
