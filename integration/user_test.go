package integration

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Irori235/system-design-2023-v2/internal/handler"

	"github.com/google/uuid"
)

func TestUser(t *testing.T) {
	t.Run("signup user", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "POST", "/api/v1/auth/signup", `{"name":"test_user","password":"pass"}`)
			assert(t, 200, rec.Code)

			res := handler.SignUpResponse{}
			assert(t, nil, json.Unmarshal(rec.Body.Bytes(), &res))
			assert(t, false, uuid.Nil == res.ID)

			userIDMap["user1"] = res.ID
		})

		t.Run("invalid json", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "POST", "/api/v1/auth/signup", `"name":"test_user","password":`)
			assert(t, 400, rec.Code)
		})

		t.Run("invalid request body", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "POST", "/api/v1/auth/signup", `{"password":"pass"}`)
			assert(t, 400, rec.Code)

			rec = doRequest(t, "POST", "/api/v1/auth/signup", `{"name":"test_user"}`)
			assert(t, 400, rec.Code)
		})
	})

	t.Run("signin user", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "POST", "/api/v1/auth/signin", `{"name":"test_user","password":"pass"}`)
			assert(t, 200, rec.Code)

			res := handler.SignInResponse{}
			assert(t, nil, json.Unmarshal(rec.Body.Bytes(), &res))

			jwtMap["user1"] = res.Token
		})

		t.Run("invalid json", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "POST", "/api/v1/auth/signin", `"name":"test_user","password":`)
			assert(t, 400, rec.Code)

			rec = doRequest(t, "POST", "/api/v1/auth/signin", `{"name":"test_user"}`)
			assert(t, 400, rec.Code)

			rec = doRequest(t, "POST", "/api/v1/auth/signin", `{"password":"pass"}`)
			assert(t, 400, rec.Code)
		})

		t.Run("invalid password", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "POST", "/api/v1/auth/signin", `{"name":"test_user","password":"invalid_pass"}`)
			assert(t, 401, rec.Code)
		})
	})

	t.Run("get me", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			jwt := jwtMap["user1"]
			header := map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", jwt),
			}
			rec := doRequest(t, "GET", "/api/v1/users/me", "", header)
			assert(t, 200, rec.Code)

			res := handler.GetMeResponse{}
			assert(t, nil, json.Unmarshal(rec.Body.Bytes(), &res))
			assert(
				t,
				handler.GetMeResponse{
					ID:   userIDMap["user1"],
					Name: "test_user",
				},
				handler.GetMeResponse{
					ID:   res.ID,
					Name: res.Name,
				},
			)
		})

		t.Run("invalid jwt", func(t *testing.T) {
			t.Parallel()

			header := map[string]string{
				"Authorization": "Bearer invalid_jwt",
			}
			rec := doRequest(t, "GET", "/api/v1/users/me", "", header)

			assert(t, 400, rec.Code)
		})
	})

	t.Run("update name", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			jwt := jwtMap["user1"]
			header := map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", jwt),
			}
			rec := doRequest(t, "PATCH", "/api/v1/users/name", `{"name":"updated_name"}`, header)
			assert(t, 200, rec.Code)

			rec2 := doRequest(t, "GET", "/api/v1/users/me", "", header)
			assert(t, 200, rec2.Code)

			res := handler.GetMeResponse{}
			assert(t, nil, json.Unmarshal(rec2.Body.Bytes(), &res))
			assert(
				t,
				handler.GetMeResponse{
					ID:   userIDMap["user1"],
					Name: "updated_name",
				},
				handler.GetMeResponse{
					ID:   res.ID,
					Name: res.Name,
				},
			)
		})
	})

	t.Run("update password", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			jwt := jwtMap["user1"]
			header := map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", jwt),
			}
			rec := doRequest(t, "PATCH", "/api/v1/users/password", `{"password":"updated_pass"}`, header)
			assert(t, 200, rec.Code)

			rec2 := doRequest(t, "POST", "/api/v1/auth/signin", `{"name":"updated_name","password":"updated_pass"}`)
			assert(t, 200, rec2.Code)

			res := handler.SignInResponse{}
			assert(t, nil, json.Unmarshal(rec2.Body.Bytes(), &res))

			jwtMap["user1"] = res.Token
		})
	})

	t.Run("quit user", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "POST", "/api/v1/auth/signup", `{"name":"test_user2","password":"pass"}`)
			assert(t, 200, rec.Code)

			res := handler.SignUpResponse{}
			assert(t, nil, json.Unmarshal(rec.Body.Bytes(), &res))
			assert(t, false, uuid.Nil == res.ID)

			rec2 := doRequest(t, "POST", "/api/v1/auth/signin", `{"name":"test_user2","password":"pass"}`)
			assert(t, 200, rec2.Code)

			res2 := handler.SignInResponse{}
			assert(t, nil, json.Unmarshal(rec2.Body.Bytes(), &res2))

			header := map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", res2.Token),
			}
			rec3 := doRequest(t, "DELETE", "/api/v1/users/quit", "", header)
			assert(t, 200, rec3.Code)

			rec4 := doRequest(t, "POST", "/api/v1/auth/signin", `{"name":"test_user2","password":"pass"}`)
			assert(t, 500, rec4.Code) //  sql: no rows in result set
		})
	})

}
