package integration

import (
	"base-project/internal/domain"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserEndpoints(t *testing.T) {
	var createdUser domain.User

	// --- Test Create User ---
	t.Run("CreateUser", func(t *testing.T) {
		userData := domain.UserData{
			Name: "John",
		}
		body, _ := json.Marshal(userData)

		resp, err := http.Post(srv.URL+"/user", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		defer resp.Body.Close() // nolint:errcheck

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		err = json.NewDecoder(resp.Body).Decode(&createdUser)
		require.NoError(t, err)

		assert.NotEqual(t, uuid.Nil, createdUser.ID)
		assert.Equal(t, userData.Name, createdUser.Name)
	})

	// --- Test Get User ---
	t.Run("GetUser", func(t *testing.T) {
		require.NotEqual(t, uuid.Nil, createdUser.ID, "CreateUser test must run first and successfully")

		resp, err := http.Get(srv.URL + "/user/" + createdUser.ID.String())
		require.NoError(t, err)
		defer resp.Body.Close() // nolint:errcheck

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var user domain.User
		err = json.NewDecoder(resp.Body).Decode(&user)
		require.NoError(t, err)

		assert.Equal(t, createdUser.ID, user.ID)
		assert.Equal(t, createdUser.Name, user.Name)
	})

	// --- Test Get Users ---
	t.Run("GetUsers", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/user")
		require.NoError(t, err)
		defer resp.Body.Close() // nolint:errcheck

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var users []domain.User
		err = json.NewDecoder(resp.Body).Decode(&users)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, len(users), 1)
	})

	// --- Test Update User ---
	t.Run("UpdateUser", func(t *testing.T) {
		require.NotEqual(t, uuid.Nil, createdUser.ID, "CreateUser test must run first and successfully")

		updatedUserData := domain.UserData{
			Name: "Jane",
		}
		body, _ := json.Marshal(updatedUserData)

		req, err := http.NewRequest(http.MethodPut, srv.URL+"/user/"+createdUser.ID.String(), bytes.NewBuffer(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close() // nolint:errcheck

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var updatedUser domain.User
		err = json.NewDecoder(resp.Body).Decode(&updatedUser)
		require.NoError(t, err)

		assert.Equal(t, createdUser.ID, updatedUser.ID)
		assert.Equal(t, updatedUserData.Name, updatedUser.Name)
	})

	// --- Test Delete User ---
	t.Run("DeleteUser", func(t *testing.T) {
		require.NotEqual(t, uuid.Nil, createdUser.ID, "CreateUser test must run first and successfully")

		req, err := http.NewRequest(http.MethodDelete, srv.URL+"/user/"+createdUser.ID.String(), nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close() // nolint:errcheck

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify user is deleted
		resp, err = http.Get(srv.URL + "/user/" + createdUser.ID.String())
		require.NoError(t, err)
		defer resp.Body.Close() // nolint:errcheck

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err, "failed to read body of not found response")

		type ErrorResponse struct {
			Title string `json:"title"`
		}
		var errorResponse ErrorResponse
		err = json.Unmarshal(bodyBytes, &errorResponse)
		require.NoError(t, err, "failed to unmarshal error response")

		assert.Equal(t, "Not found", errorResponse.Title)
	})
}
