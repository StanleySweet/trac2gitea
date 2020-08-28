// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

// AddMilestone adds a milestone to Gitea, returns id of created milestone
func (accessor *DefaultAccessor) AddMilestone(milestone *Milestone) (int64, error) {
	_, err := accessor.db.Exec(`
		INSERT INTO
			milestone(repo_id,name,content,is_closed,deadline_unix,closed_date_unix)
			SELECT $1,$2,$3,$4,$5,$6 WHERE
				NOT EXISTS (SELECT * FROM milestone WHERE repo_id = $1 AND name = $2)`,
		accessor.repoID, milestone.Name, milestone.Description, milestone.Closed, milestone.DueTime, milestone.ClosedTime)
	if err != nil {
		err = errors.Wrapf(err, "adding milestone %s", milestone.Name)
		return -1, err
	}

	var milestoneID int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&milestoneID)
	if err != nil {
		err = errors.Wrapf(err, "retrieving id of new milestone %s", milestone.Name)
		return -1, err
	}

	return milestoneID, nil
}

// GetMilestoneID gets the ID of a named milestone - returns -1 if no such milestone
func (accessor *DefaultAccessor) GetMilestoneID(name string) (int64, error) {
	var milestoneID int64 = -1
	err := accessor.db.QueryRow(`
		SELECT id FROM milestone WHERE name = $1
		`, name).Scan(&milestoneID)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving id of milestone %s", name)
		return -1, err
	}

	return milestoneID, nil
}

// GetMilestoneURL gets the URL for accessing a given milestone
func (accessor *DefaultAccessor) GetMilestoneURL(milestoneID int64) string {
	repoURL := accessor.getUserRepoURL()
	return fmt.Sprintf("%s/milestone/%d", repoURL, milestoneID)
}
