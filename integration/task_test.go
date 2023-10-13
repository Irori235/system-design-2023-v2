package integration

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Irori235/system-design-2023-v2/internal/handler"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

func TestTask(t *testing.T) {
	t.Run("setup user", func(t *testing.T) {
		rec := doRequest(t, "POST", "/api/v1/auth/signup", `{"name":"test_user3","password":"pass"}`)
		assert(t, 200, rec.Code)

		res := handler.SignUpResponse{}
		assert(t, nil, json.Unmarshal(rec.Body.Bytes(), &res))
		assert(t, false, uuid.Nil == res.ID)

		userIDMap["user3"] = res.ID

		rec2 := doRequest(t, "POST", "/api/v1/auth/signin", `{"name":"test_user3","password":"pass"}`)
		assert(t, 200, rec2.Code)

		res2 := handler.SignInResponse{}
		assert(t, nil, json.Unmarshal(rec2.Body.Bytes(), &res2))

		jwtMap["user3"] = res2.Token
	})

	t.Run("create task", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			jwt := jwtMap["user3"]
			header := map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", jwt),
			}
			rec := doRequest(t, "POST", "/api/v1/tasks", `{"title":"test_title"}`, header)
			assert(t, 200, rec.Code)

			rec2 := doRequest(t, "POST", "/api/v1/tasks", `{"title":"test_title2"}`, header)
			assert(t, 200, rec2.Code)

			rec3 := doRequest(t, "POST", "/api/v1/tasks", `{"title":"test_title3"}`, header)
			assert(t, 200, rec3.Code)
		})

		t.Run("invalid json", func(t *testing.T) {
			t.Parallel()
			jwt := jwtMap["user3"]
			header := map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", jwt),
			}
			rec := doRequest(t, "POST", "/api/v1/tasks", `"title":"test_title"`, header)
			assert(t, 400, rec.Code)
		})

		t.Run("invalid request body", func(t *testing.T) {
			t.Parallel()
			jwt := jwtMap["user3"]
			header := map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", jwt),
			}
			rec := doRequest(t, "POST", "/api/v1/tasks", "", header)
			assert(t, 400, rec.Code)
		})
	})

	t.Run("get tasks", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			jwt := jwtMap["user3"]
			header := map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", jwt),
			}
			rec := doRequest(t, "GET", "/api/v1/tasks", "", header)
			assert(t, 200, rec.Code)

			res := handler.GetTasksResponse{}
			assert(t, nil, json.Unmarshal(rec.Body.Bytes(), &res))
			assert(t, 3, len(res))

			titles := []string{"test_title", "test_title2", "test_title3"}

			for _, task := range res {
				assert(t, false, uuid.Nil == task.ID)
				assert(t, true, task.UserID == userIDMap["user3"])
				assert(t, true, slices.Contains(titles, task.Title))
				assert(t, false, task.IsDone)
			}

			taskMap["task1"] = res[0]
			taskMap["task2"] = res[1]
			taskMap["task3"] = res[2]
		})
	})

	t.Run("update task", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			jwt := jwtMap["user3"]
			header := map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", jwt),
			}
			rec := doRequest(t, "PUT", "/api/v1/tasks/"+taskMap["task1"].ID.String(), fmt.Sprintf(`{"title":"updated_title","is_done":%t}`, taskMap["task1"].IsDone), header)
			assert(t, 200, rec.Code)

			rec2 := doRequest(t, "GET", "/api/v1/tasks", "", header)
			assert(t, 200, rec2.Code)

			res := handler.GetTasksResponse{}
			assert(t, nil, json.Unmarshal(rec2.Body.Bytes(), &res))
			assert(t, 3, len(res))

			task1res := handler.GetTaskResponse{}
			for _, task := range res {
				if task.ID == taskMap["task1"].ID {
					task1res = handler.GetTaskResponse{
						ID:     task.ID,
						UserID: task.UserID,
						Title:  task.Title,
						IsDone: task.IsDone,
					}

					break
				}
			}

			assert(t,
				handler.GetTaskResponse{
					ID:     taskMap["task1"].ID,
					UserID: taskMap["task1"].UserID,
					Title:  "updated_title",
					IsDone: taskMap["task1"].IsDone,
				},
				task1res,
			)
		})

		t.Run("success2", func(t *testing.T) {
			t.Parallel()

			jwt := jwtMap["user3"]
			header := map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", jwt),
			}
			rec := doRequest(t, "PUT", "/api/v1/tasks/"+taskMap["task2"].ID.String(), fmt.Sprintf(`{"title":"%s","is_done":true}`, taskMap["task2"].Title), header)
			assert(t, 200, rec.Code)

			rec2 := doRequest(t, "GET", "/api/v1/tasks", "", header)
			assert(t, 200, rec2.Code)

			res := handler.GetTasksResponse{}
			assert(t, nil, json.Unmarshal(rec2.Body.Bytes(), &res))
			assert(t, 3, len(res))

			task2res := handler.GetTaskResponse{}
			for _, task := range res {
				if task.ID == taskMap["task2"].ID {
					task2res = handler.GetTaskResponse{
						ID:     task.ID,
						UserID: task.UserID,
						Title:  task.Title,
						IsDone: task.IsDone,
					}

					break
				}
			}

			assert(t,
				handler.GetTaskResponse{
					ID:     taskMap["task2"].ID,
					UserID: taskMap["task2"].UserID,
					Title:  taskMap["task2"].Title,
					IsDone: true,
				},
				task2res,
			)
		})

		t.Run("success3", func(t *testing.T) {
			t.Parallel()

			jwt := jwtMap["user3"]
			header := map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", jwt),
			}
			rec := doRequest(t, "PUT", "/api/v1/tasks/"+taskMap["task3"].ID.String(), `{"title":"updated_title3","is_done":true}`, header)
			assert(t, 200, rec.Code)

			rec2 := doRequest(t, "GET", "/api/v1/tasks", "", header)
			assert(t, 200, rec2.Code)

			res := handler.GetTasksResponse{}
			assert(t, nil, json.Unmarshal(rec2.Body.Bytes(), &res))
			assert(t, 3, len(res))

			task3res := handler.GetTaskResponse{}
			for _, task := range res {
				if task.ID == taskMap["task3"].ID {
					task3res = handler.GetTaskResponse{
						ID:     task.ID,
						UserID: task.UserID,
						Title:  task.Title,
						IsDone: task.IsDone,
					}

					break
				}
			}

			assert(t,
				handler.GetTaskResponse{
					ID:     taskMap["task3"].ID,
					UserID: taskMap["task3"].UserID,
					Title:  "updated_title3",
					IsDone: true,
				},
				task3res,
			)
		})

		t.Run("invalid request body", func(t *testing.T) {
			t.Parallel()

			jwt := jwtMap["user3"]
			header := map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", jwt),
			}
			rec := doRequest(t, "PUT", "/api/v1/tasks/"+taskMap["task3"].ID.String(), ``, header)
			assert(t, 400, rec.Code)
		})

	})

	t.Run("delete task", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			jwt := jwtMap["user3"]
			header := map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", jwt),
			}
			rec := doRequest(t, "DELETE", "/api/v1/tasks/"+taskMap["task1"].ID.String(), "", header)
			assert(t, 200, rec.Code)

			rec2 := doRequest(t, "GET", "/api/v1/tasks", "", header)
			assert(t, 200, rec2.Code)

			res := handler.GetTasksResponse{}
			assert(t, nil, json.Unmarshal(rec2.Body.Bytes(), &res))
			assert(t, 2, len(res))

			for _, task := range res {
				assert(t, false, task.ID == taskMap["task1"].ID)
			}
		})
	})
}
