// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/types/library"
)

func TestGithub_OrgAccess_Admin(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/orgs/:org/memberships/:username", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/org_admin.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	want := "admin"

	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.OrgAccess(context.TODO(), u, "github")

	if resp.Code != http.StatusOK {
		t.Errorf("OrgAccess returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("OrgAccess returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("OrgAccess is %v, want %v", got, want)
	}
}

func TestGithub_OrgAccess_Member(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/orgs/:org/memberships/:username", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/org_member.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	want := "member"

	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.OrgAccess(context.TODO(), u, "github")

	if resp.Code != http.StatusOK {
		t.Errorf("OrgAccess returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("OrgAccess returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("OrgAccess is %v, want %v", got, want)
	}
}

func TestGithub_OrgAccess_NotFound(t *testing.T) {
	// setup mock server
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup types
	want := ""

	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.OrgAccess(context.TODO(), u, "github")

	if err == nil {
		t.Errorf("OrgAccess should have returned err")
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("OrgAccess is %v, want %v", got, want)
	}
}

func TestGithub_OrgAccess_Pending(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/orgs/:org/memberships/:username", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/org_pending.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	want := ""

	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.OrgAccess(context.TODO(), u, "github")

	if resp.Code != http.StatusOK {
		t.Errorf("OrgAccess returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("OrgAccess returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("OrgAccess is %v, want %v", got, want)
	}
}

func TestGithub_OrgAccess_Personal(t *testing.T) {
	// setup mock server
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup types
	want := "admin"

	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.OrgAccess(context.TODO(), u, "foo")

	if err != nil {
		t.Errorf("OrgAccess returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("OrgAccess is %v, want %v", got, want)
	}
}

func TestGithub_RepoAccess_Admin(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/repo_admin.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	want := "admin"

	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.RepoAccess(context.TODO(), u, u.GetToken(), "github", "octocat")

	if resp.Code != http.StatusOK {
		t.Errorf("RepoAccess returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("RepoAccess returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("RepoAccess is %v, want %v", got, want)
	}
}

func TestGithub_RepoAccess_NotFound(t *testing.T) {
	// setup mock server
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup types
	want := ""

	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.RepoAccess(context.TODO(), u, u.GetToken(), "github", "octocat")

	if err == nil {
		t.Errorf("RepoAccess should have returned err")
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("RepoAccess is %v, want %v", got, want)
	}
}

func TestGithub_TeamAccess_Admin(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/user/teams", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/team_admin.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	want := "admin"

	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.TeamAccess(context.TODO(), u, "github", "octocat")

	if resp.Code != http.StatusOK {
		t.Errorf("TeamAccess returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("TeamAccess returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("TeamAccess is %v, want %v", got, want)
	}
}

func TestGithub_TeamAccess_NoAccess(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/user/teams", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/team_admin.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	want := ""

	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.TeamAccess(context.TODO(), u, "github", "baz")

	if resp.Code != http.StatusOK {
		t.Errorf("TeamAccess returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("TeamAccess returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("TeamAccess is %v, want %v", got, want)
	}
}

func TestGithub_TeamAccess_NotFound(t *testing.T) {
	// setup mock server
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	// setup types
	want := ""

	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.TeamAccess(context.TODO(), u, "github", "octocat")

	if err == nil {
		t.Errorf("TeamAccess should have returned err")
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("TeamAccess is %v, want %v", got, want)
	}
}

func TestGithub_TeamList(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/user/teams", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/team_admin.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	want := []string{"Justice League", "octocat"}

	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.ListUsersTeamsForOrg(context.TODO(), u, "github")

	if resp.Code != http.StatusOK {
		t.Errorf("TeamAccess returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("TeamAccess returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("TeamAccess is %v, want %v", got, want)
	}
}
