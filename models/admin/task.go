// Copyright 2019 Gitea. All rights reserved.
// SPDX-License-Identifier: MIT

package admin

import (
	"context"
	"fmt"

	"code.gitea.io/gitea/models/db"
	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/modules/timeutil"
	"code.gitea.io/gitea/modules/util"
)

// Task represents a task
type Task struct {
	ID     int64
	DoerID int64            `xorm:"index"` // operator
	Doer   *user_model.User `xorm:"-"`
	// RepoID         int64                  `xorm:"index"`
	// Repo           *repo_model.Repository `xorm:"-"`
	Type           structs.TaskType
	Status         structs.TaskStatus `xorm:"index"`
	StartTime      timeutil.TimeStamp
	EndTime        timeutil.TimeStamp
	PayloadContent string             `xorm:"TEXT"`
	Message        string             `xorm:"TEXT"` // if task failed, saved the error reason, it could be a JSON string of TranslatableMessage or a plain message
	Created        timeutil.TimeStamp `xorm:"created"`
}

func init() {
	db.RegisterModel(new(Task))
}

// TranslatableMessage represents JSON struct that can be translated with a Locale
type TranslatableMessage struct {
	Format string
	Args   []any `json:"omitempty"`
}

// LoadDoer loads do user
func (task *Task) LoadDoer(ctx context.Context) error {
	if task.Doer != nil {
		return nil
	}

	var doer user_model.User
	has, err := db.GetEngine(ctx).ID(task.DoerID).Get(&doer)
	if err != nil {
		return err
	} else if !has {
		return user_model.ErrUserNotExist{
			UID: task.DoerID,
		}
	}
	task.Doer = &doer

	return nil
}

// UpdateCols updates some columns
func (task *Task) UpdateCols(ctx context.Context, cols ...string) error {
	_, err := db.GetEngine(ctx).ID(task.ID).Cols(cols...).Update(task)
	return err
}

// ErrTaskDoesNotExist represents a "TaskDoesNotExist" kind of error.
type ErrTaskDoesNotExist struct {
	ID     int64
	RepoID int64
	Type   structs.TaskType
}

// IsErrTaskDoesNotExist checks if an error is a ErrTaskDoesNotExist.
func IsErrTaskDoesNotExist(err error) bool {
	_, ok := err.(ErrTaskDoesNotExist)
	return ok
}

func (err ErrTaskDoesNotExist) Error() string {
	return fmt.Sprintf("task does not exist [id: %d, repo_id: %d, type: %d]",
		err.ID, err.RepoID, err.Type)
}

func (err ErrTaskDoesNotExist) Unwrap() error {
	return util.ErrNotExist
}

// CreateTask creates a task on database
func CreateTask(ctx context.Context, task *Task) error {
	return db.Insert(ctx, task)
}

// FinishMigrateTask updates database when migrate task finished
func FinishMigrateTask(ctx context.Context, task *Task) error {
	task.Status = structs.TaskStatusFinished
	task.EndTime = timeutil.TimeStampNow()

	_, err := db.GetEngine(ctx).ID(task.ID).Cols("status", "end_time").Update(task)
	return err
}
