// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToCorrectPageSize(t *testing.T) {
	assert.EqualValues(t, 30, ToCorrectPageSize(0))
	assert.EqualValues(t, 30, ToCorrectPageSize(-10))
	assert.EqualValues(t, 20, ToCorrectPageSize(20))
	assert.EqualValues(t, 50, ToCorrectPageSize(100))
}
