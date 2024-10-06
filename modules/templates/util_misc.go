// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package templates

import (
	"html/template"
	"mime"
	"path/filepath"
	"strconv"
	"strings"

	"code.gitea.io/gitea/modules/svg"

	"github.com/editorconfig/editorconfig-core-go/v2"
)

func sortArrow(normSort, revSort, urlSort string, isDefault bool) template.HTML {
	// if needed
	if len(normSort) == 0 || len(urlSort) == 0 {
		return ""
	}

	if len(urlSort) == 0 && isDefault {
		// if sort is sorted as default add arrow tho this table header
		if isDefault {
			return svg.RenderHTML("octicon-triangle-down", 16)
		}
	} else {
		// if sort arg is in url test if it correlates with column header sort arguments
		// the direction of the arrow should indicate the "current sort order", up means ASC(normal), down means DESC(rev)
		if urlSort == normSort {
			// the table is sorted with this header normal
			return svg.RenderHTML("octicon-triangle-up", 16)
		} else if urlSort == revSort {
			// the table is sorted with this header reverse
			return svg.RenderHTML("octicon-triangle-down", 16)
		}
	}
	// the table is NOT sorted with this header
	return ""
}

type remoteAddress struct {
	Address  string
	Username string
	Password string
}

func filenameIsImage(filename string) bool {
	mimeType := mime.TypeByExtension(filepath.Ext(filename))
	return strings.HasPrefix(mimeType, "image/")
}

func tabSizeClass(ec *editorconfig.Editorconfig, filename string) string {
	if ec != nil {
		def, err := ec.GetDefinitionForFilename(filename)
		if err == nil && def.TabWidth >= 1 && def.TabWidth <= 16 {
			return "tab-size-" + strconv.Itoa(def.TabWidth)
		}
	}
	return "tab-size-4"
}
