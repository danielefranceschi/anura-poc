// Copyright 2015 The Gogs Authors. All rights reserved.
// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package models

import (
	"fmt"

	"code.gitea.io/gitea/modules/util"
)

// ErrUserOwnPackages notifies that the user (still) owns the packages.
type ErrUserOwnPackages struct {
	UID int64
}

// IsErrUserOwnPackages checks if an error is an ErrUserOwnPackages.
func IsErrUserOwnPackages(err error) bool {
	_, ok := err.(ErrUserOwnPackages)
	return ok
}

func (err ErrUserOwnPackages) Error() string {
	return fmt.Sprintf("user still has ownership of packages [uid: %d]", err.UID)
}

// ErrDeleteLastAdminUser represents a "DeleteLastAdminUser" kind of error.
type ErrDeleteLastAdminUser struct {
	UID int64
}

// IsErrDeleteLastAdminUser checks if an error is a ErrDeleteLastAdminUser.
func IsErrDeleteLastAdminUser(err error) bool {
	_, ok := err.(ErrDeleteLastAdminUser)
	return ok
}

func (err ErrDeleteLastAdminUser) Error() string {
	return fmt.Sprintf("can not delete the last admin user [uid: %d]", err.UID)
}

// ErrUpdateTaskNotExist represents a "UpdateTaskNotExist" kind of error.
type ErrUpdateTaskNotExist struct {
	UUID string
}

// IsErrUpdateTaskNotExist checks if an error is a ErrUpdateTaskNotExist.
func IsErrUpdateTaskNotExist(err error) bool {
	_, ok := err.(ErrUpdateTaskNotExist)
	return ok
}

func (err ErrUpdateTaskNotExist) Error() string {
	return fmt.Sprintf("update task does not exist [uuid: %s]", err.UUID)
}

func (err ErrUpdateTaskNotExist) Unwrap() error {
	return util.ErrNotExist
}

// ErrFilenameInvalid represents a "FilenameInvalid" kind of error.
type ErrFilenameInvalid struct {
	Path string
}

// IsErrFilenameInvalid checks if an error is an ErrFilenameInvalid.
func IsErrFilenameInvalid(err error) bool {
	_, ok := err.(ErrFilenameInvalid)
	return ok
}

func (err ErrFilenameInvalid) Error() string {
	return fmt.Sprintf("path contains a malformed path component [path: %s]", err.Path)
}

func (err ErrFilenameInvalid) Unwrap() error {
	return util.ErrInvalidArgument
}

// ErrFilePathInvalid represents a "FilePathInvalid" kind of error.
type ErrFilePathInvalid struct {
	Message string
	Path    string
	Name    string
}

// IsErrFilePathInvalid checks if an error is an ErrFilePathInvalid.
func IsErrFilePathInvalid(err error) bool {
	_, ok := err.(ErrFilePathInvalid)
	return ok
}

func (err ErrFilePathInvalid) Error() string {
	if err.Message != "" {
		return err.Message
	}
	return fmt.Sprintf("path is invalid [path: %s]", err.Path)
}

func (err ErrFilePathInvalid) Unwrap() error {
	return util.ErrInvalidArgument
}

// ErrFilePathProtected represents a "FilePathProtected" kind of error.
type ErrFilePathProtected struct {
	Message string
	Path    string
}

// IsErrFilePathProtected checks if an error is an ErrFilePathProtected.
func IsErrFilePathProtected(err error) bool {
	_, ok := err.(ErrFilePathProtected)
	return ok
}

func (err ErrFilePathProtected) Error() string {
	if err.Message != "" {
		return err.Message
	}
	return fmt.Sprintf("path is protected and can not be changed [path: %s]", err.Path)
}

func (err ErrFilePathProtected) Unwrap() error {
	return util.ErrPermissionDenied
}
