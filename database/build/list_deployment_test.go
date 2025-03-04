// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestBuild_Engine_ListBuildsForDeployment(t *testing.T) {
	// setup types
	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepoID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetDeployPayload(nil)
	_buildOne.SetSource("https://github.com/github/octocat/deployments/1")

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepoID(1)
	_buildTwo.SetNumber(2)
	_buildTwo.SetDeployPayload(nil)
	_buildTwo.SetSource("https://github.com/github/octocat/deployments/1")

	_deployment := testDeployment()
	_deployment.SetID(1)
	_deployment.SetRepoID(1)
	_deployment.SetURL("https://github.com/github/octocat/deployments/1")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected count query result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the count query
	_mock.ExpectQuery(`SELECT count(*) FROM "builds" WHERE source = $1`).WithArgs("https://github.com/github/octocat/deployments/1").WillReturnRows(_rows)

	// create expected query result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "repo_id", "pipeline_id", "number", "parent", "event", "event_action", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"}).
		AddRow(2, 1, nil, 2, 0, "", "", "", "", 0, 0, 0, 0, "", nil, "", "https://github.com/github/octocat/deployments/1", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0).
		AddRow(1, 1, nil, 1, 0, "", "", "", "", 0, 0, 0, 0, "", nil, "", "https://github.com/github/octocat/deployments/1", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE source = $1 ORDER BY number DESC LIMIT 10`).WithArgs("https://github.com/github/octocat/deployments/1").WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateBuild(context.TODO(), _buildOne)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	_, err = _sqlite.CreateBuild(context.TODO(), _buildTwo)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []*library.Build
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*library.Build{_buildTwo, _buildOne},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*library.Build{_buildTwo, _buildOne},
		},
	}

	filters := map[string]interface{}{}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _, err := test.database.ListBuildsForDeployment(context.TODO(), _deployment, filters, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListBuildsForDeployment for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListBuildsForDeployment for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListBuildsForDeployment for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
