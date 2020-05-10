package sdk

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectsGet(t *testing.T) {
	requestAssertions := func(t *testing.T, req *http.Request) {
		t.Helper()

		require.True(t, strings.HasSuffix(req.URL.Path, "/projects/"))
		require.Equal(t, http.MethodGet, req.Method)
		cookieSid, err := req.Cookie("sid")
		require.NoError(t, err)
		require.Equal(t, &http.Cookie{Name: "sid", Value: "my-random-sid"}, cookieSid)
	}

	t.Run("correctly returns projects", func(t *testing.T) {
		responseBody := `[{"_id":"mongo-id-1","name":"Project 1","configurationGitPath":"/clients/path","projectId":"project-1","environments":[{"label":"Development","value":"development","cluster":{"hostname":"127.0.0.1","namespace":"project-1-dev"}}],"pipelines":{"type":"gitlab"}},{"_id":"mongo-id-2","name":"Project 2","configurationGitPath":"/clients/path/configuration","projectId":"project-2","environments":[{"label":"Development","value":"development","cluster":{"hostname":"127.0.0.1","namespace":"project-2-dev"}},{"label":"Production","value":"production","cluster":{"hostname":"127.0.0.1","namespace":"project-2"}}]}]`
		expectedProjects := Projects{
			Project{
				ID:                   "mongo-id-1",
				Name:                 "Project 1",
				ConfigurationGitPath: "/clients/path",
				ProjectID:            "project-1",
				Environments: []Environment{
					{
						Cluster: Cluster{
							Hostname:  "127.0.0.1",
							Namespace: "project-1-dev",
						},
						DisplayName: "Development",
						EnvID:       "development",
					},
				},
				Pipelines: Pipelines{
					Type: "gitlab",
				},
			},
			Project{
				ID:                   "mongo-id-2",
				Name:                 "Project 2",
				ConfigurationGitPath: "/clients/path/configuration",
				ProjectID:            "project-2",
				Environments: []Environment{
					{
						Cluster: Cluster{
							Hostname:  "127.0.0.1",
							Namespace: "project-2-dev",
						},
						DisplayName: "Development",
						EnvID:       "development",
					},
					{
						Cluster: Cluster{
							Hostname:  "127.0.0.1",
							Namespace: "project-2",
						},
						DisplayName: "Production",
						EnvID:       "production",
					},
				},
			},
		}

		s := testCreateResponseServer(t, requestAssertions, responseBody, 200)
		client := testCreateProjectClient(t, s.URL)

		projects, err := client.Get()
		require.NoError(t, err)
		require.Equal(t, expectedProjects, projects)
	})

	t.Run("throws when server respond with 401", func(t *testing.T) {
		responseBody := `{"statusCode":401,"error":"Unauthorized","message":"Unauthorized"}`
		s := testCreateResponseServer(t, requestAssertions, responseBody, 401)
		client := testCreateProjectClient(t, s.URL)

		projects, err := client.Get()
		require.Nil(t, projects)
		require.EqualError(t, err, fmt.Sprintf("GET %s/api/backend/projects/: 401 - %s", s.URL, responseBody))
		require.True(t, errors.Is(err, ErrHTTP))
	})

	t.Run("throws if response body is not as expected", func(t *testing.T) {
		responseBody := `{"statusCode":401,"error":"Unauthorized","message":"Unauthorized"}`
		s := testCreateResponseServer(t, requestAssertions, responseBody, 200)
		client := testCreateProjectClient(t, s.URL)

		projects, err := client.Get()
		require.Nil(t, projects)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrGeneric))
	})
}

func TestGetProjectByID(t *testing.T) {
	projectRequestAssertions := func(t *testing.T, req *http.Request) {
		t.Helper()

		require.True(t, strings.HasSuffix(req.URL.Path, "/projects/"))
		require.Equal(t, http.MethodGet, req.Method)
		cookieSid, err := req.Cookie("sid")
		require.NoError(t, err)
		require.Equal(t, &http.Cookie{Name: "sid", Value: "my-random-sid"}, cookieSid)
	}

	t.Run("Error creating request for projectId fetch", func(t *testing.T) {
		client := testCreateClient(t, "this-url-does-not-exist")
		project, err := getProjectByID(client, "project1")
		require.Nil(t, project)
		require.EqualError(t, err, fmt.Sprintf("BaseURL must have a trailing slash, but \"this-url-does-not-exist\" does not"))
	})

	t.Run("Network error occurs during projectId fetch", func(t *testing.T) {
		responseBody := `{"statusCode":401,"error":"Unauthorized","message":"Unauthorized"}`
		s := testCreateResponseServer(t, projectRequestAssertions, responseBody, 401)
		defer s.Close()

		client := testCreateClient(t, s.URL)
		project, err := getProjectByID(client, "project1")
		require.Nil(t, project)
		require.EqualError(t, err, fmt.Sprintf("GET %s/api/backend/projects/: 401 - %s", s.URL, responseBody))
		require.True(t, errors.Is(err, ErrHTTP))
	})

	t.Run("Generic error occurs during projectId fetch (malformed data, _id should be a string)", func(t *testing.T) {
		responseBody := `[{"_id":9876543,"name":"Project 1","configurationGitPath":"/clients/path","projectId":"project-1","environments":[{"label":"Development","value":"development","cluster":{"hostname":"127.0.0.1","namespace":"project-1-dev"}}],"pipelines":{"type":"gitlab"}},{"_id":"mongo-id-2","name":"Project 2","configurationGitPath":"/clients/path/configuration","projectId":"project-2","environments":[{"label":"Development","value":"development","cluster":{"hostname":"127.0.0.1","namespace":"project-2-dev"}},{"label":"Production","value":"production","cluster":{"hostname":"127.0.0.1","namespace":"project-2"}}]}]`
		s := testCreateResponseServer(t, projectRequestAssertions, responseBody, 200)
		defer s.Close()

		client := testCreateClient(t, s.URL)
		project, err := getProjectByID(client, "project1")
		require.Nil(t, project)
		require.EqualError(t, err, fmt.Sprintf("%s: json: cannot unmarshal number into Go struct field Project._id of type string", ErrGeneric))
		require.True(t, errors.Is(err, ErrGeneric))
	})

	t.Run("Error projectID not found", func(t *testing.T) {
		responseBody := `[{"_id":"mongo-id-1","name":"Project 1","configurationGitPath":"/clients/path","projectId":"project-1","environments":[{"label":"Development","value":"development","cluster":{"hostname":"127.0.0.1","namespace":"project-1-dev"}}],"pipelines":{"type":"gitlab"}},{"_id":"mongo-id-2","name":"Project 2","configurationGitPath":"/clients/path/configuration","projectId":"project-2","environments":[{"label":"Development","value":"development","cluster":{"hostname":"127.0.0.1","namespace":"project-2-dev"}},{"label":"Production","value":"production","cluster":{"hostname":"127.0.0.1","namespace":"project-2"}}]}]`
		s := testCreateResponseServer(t, projectRequestAssertions, responseBody, 200)
		defer s.Close()

		client := testCreateClient(t, s.URL)
		project, err := getProjectByID(client, "project1")
		require.Nil(t, project)
		require.EqualError(t, err, fmt.Sprintf("%s: project1", ErrProjectNotFound))
		require.True(t, errors.Is(err, ErrProjectNotFound))
	})

	t.Run("Returns desired project", func(t *testing.T) {
		responseBody := `[{"_id":"mongo-id-1","name":"Project 1","configurationGitPath":"/clients/path","projectId":"project-1","environments":[{"label":"Development","value":"development","cluster":{"hostname":"127.0.0.1","namespace":"project-1-dev"}}],"pipelines":{"type":"gitlab"}},{"_id":"mongo-id-2","name":"Project 2","configurationGitPath":"/clients/path/configuration","projectId":"project-2","environments":[{"label":"Development","value":"development","cluster":{"hostname":"127.0.0.1","namespace":"project-2-dev"}},{"label":"Production","value":"production","cluster":{"hostname":"127.0.0.1","namespace":"project-2"}}],"pipelines":{"type":"gitlab"}}]`
		s := testCreateResponseServer(t, projectRequestAssertions, responseBody, 200)
		defer s.Close()

		client := testCreateClient(t, s.URL)
		project, err := getProjectByID(client, "project-2")
		require.NoError(t, err)
		require.Equal(t, &Project{
			ID:                   "mongo-id-2",
			Name:                 "Project 2",
			ConfigurationGitPath: "/clients/path/configuration",
			ProjectID:            "project-2",
			Environments: []Environment{{
				EnvID:       "development",
				DisplayName: "Development",
				Cluster: Cluster{
					Hostname:  "127.0.0.1",
					Namespace: "project-2-dev",
				},
			}, {
				EnvID:       "production",
				DisplayName: "Production",
				Cluster: Cluster{
					Hostname:  "127.0.0.1",
					Namespace: "project-2",
				},
			}},
			Pipelines: Pipelines{Type: "gitlab"},
		},
			project)
	})
}

func testCreateProjectClient(t *testing.T, url string) IProjects {
	t.Helper()
	return ProjectsClient{
		JSONClient: testCreateClient(t, url),
	}
}
