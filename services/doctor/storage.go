// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package doctor

import (
	"context"
	"io/fs"

	"code.gitea.io/gitea/models/packages"
	"code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/log"
	packages_module "code.gitea.io/gitea/modules/packages"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/storage"
)

type commonStorageCheckOptions struct {
	storer     storage.ObjectStorage
	isOrphaned func(path string, obj storage.Object, stat fs.FileInfo) (bool, error)
	name       string
}

func commonCheckStorage(logger log.Logger, autofix bool, opts *commonStorageCheckOptions) error {
	totalCount, orphanedCount := 0, 0
	totalSize, orphanedSize := int64(0), int64(0)

	var pathsToDelete []string
	if err := opts.storer.IterateObjects("", func(p string, obj storage.Object) error {
		defer obj.Close()

		totalCount++
		stat, err := obj.Stat()
		if err != nil {
			return err
		}
		totalSize += stat.Size()

		orphaned, err := opts.isOrphaned(p, obj, stat)
		if err != nil {
			return err
		}
		if orphaned {
			orphanedCount++
			orphanedSize += stat.Size()
			if autofix {
				pathsToDelete = append(pathsToDelete, p)
			}
		}
		return nil
	}); err != nil {
		logger.Error("Error whilst iterating %s storage: %v", opts.name, err)
		return err
	}

	if orphanedCount > 0 {
		if autofix {
			var deletedNum int
			for _, p := range pathsToDelete {
				if err := opts.storer.Delete(p); err != nil {
					log.Error("Error whilst deleting %s from %s storage: %v", p, opts.name, err)
				} else {
					deletedNum++
				}
			}
			logger.Info("Deleted %d/%d orphaned %s(s)", deletedNum, orphanedCount, opts.name)
		} else {
			logger.Warn("Found %d/%d (%s/%s) orphaned %s(s)", orphanedCount, totalCount, base.FileSize(orphanedSize), base.FileSize(totalSize), opts.name)
		}
	} else {
		logger.Info("Found %d (%s) %s(s)", totalCount, base.FileSize(totalSize), opts.name)
	}
	return nil
}

type checkStorageOptions struct {
	All      bool
	Avatars  bool
	Packages bool
}

// checkStorage will return a doctor check function to check the requested storage types for "orphaned" stored object/files and optionally delete them
func checkStorage(opts *checkStorageOptions) func(ctx context.Context, logger log.Logger, autofix bool) error {
	return func(ctx context.Context, logger log.Logger, autofix bool) error {
		if err := storage.Init(); err != nil {
			logger.Error("storage.Init failed: %v", err)
			return err
		}

		if opts.Avatars || opts.All {
			if err := commonCheckStorage(logger, autofix,
				&commonStorageCheckOptions{
					storer: storage.Avatars,
					isOrphaned: func(path string, obj storage.Object, stat fs.FileInfo) (bool, error) {
						exists, err := user.ExistsWithAvatarAtStoragePath(ctx, path)
						return !exists, err
					},
					name: "avatar",
				}); err != nil {
				return err
			}
		}

		if opts.Packages || opts.All {
			if !setting.Packages.Enabled {
				logger.Info("Packages isn't enabled (skipped)")
				return nil
			}
			if err := commonCheckStorage(logger, autofix,
				&commonStorageCheckOptions{
					storer: storage.Packages,
					isOrphaned: func(path string, obj storage.Object, stat fs.FileInfo) (bool, error) {
						key, err := packages_module.RelativePathToKey(path)
						if err != nil {
							// If there is an error here then the relative path does not match a valid package
							// Therefore it is orphaned by default
							return true, nil
						}

						exists, err := packages.ExistPackageBlobWithSHA(ctx, string(key))

						return !exists, err
					},
					name: "package blob",
				}); err != nil {
				return err
			}
		}

		return nil
	}
}

func init() {
	Register(&Check{
		Title:                      "Check if there are orphaned storage files",
		Name:                       "storages",
		IsDefault:                  false,
		Run:                        checkStorage(&checkStorageOptions{All: true}),
		AbortIfFailed:              false,
		SkipDatabaseInitialization: false,
		Priority:                   1,
	})

	Register(&Check{
		Title:                      "Check if there are orphaned package blobs in storage",
		Name:                       "storage-packages",
		IsDefault:                  false,
		Run:                        checkStorage(&checkStorageOptions{Packages: true}),
		AbortIfFailed:              false,
		SkipDatabaseInitialization: false,
		Priority:                   1,
	})
}
