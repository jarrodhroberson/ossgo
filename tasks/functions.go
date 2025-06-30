package tasks

import (
	"context"
	"fmt"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"

	errs "github.com/jarrodhroberson/ossgo/errors"
	"github.com/jarrodhroberson/ossgo/functions/must"
	"github.com/jarrodhroberson/ossgo/gcp"
	"github.com/jarrodhroberson/ossgo/timestamp"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog/log"
)

func CreateQueuePath(queueId string) string {
	projectId := must.Must(gcp.ProjectId())
	locationId := must.Must(gcp.Region())
	return fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectId, locationId, queueId)
}

func DisallowDuplicates(name string) CreateTaskRequestOption {
	return func(ctr *cloudtaskspb.CreateTaskRequest) {
		// Task name must be formatted: "projects/<PROJECT_ID>/locations/<LOCATION_ID>/queues/<QUEUE_ID>/tasks/<TASK_ID>"
		ctr.Task.Name = fmt.Sprintf("%s/tasks/%s", ctr.Parent, name)
	}
}

// CreateTask
// Deprecated: use Create instead
func CreateTask(ctx context.Context, req *cloudtaskspb.CreateTaskRequest) (*cloudtaskspb.Task, error) {
	return Create(ctx, req)
}

// Create creates a new Cloud Task using the provided request. It establishes a new client connection,
// submits the task creation request, and handles any potential errors including duplicate tasks.
//
// Parameters:
//   - ctx: The context.Context for the request
//   - req: The CreateTaskRequest containing task configuration
//
// Returns:
//   - *cloudtaskspb.Task: The created task if successful
//   - error: An error if the operation fails, including specific handling for already existing tasks
func Create(ctx context.Context, req *cloudtaskspb.CreateTaskRequest) (*cloudtaskspb.Task, error) {

	c, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer func(c *cloudtasks.Client) {
		err = c.Close()
		if err != nil {
			log.Error().Stack().Err(err).Msg(err.Error())
		}
	}(c)

	createdTask, err := c.CreateTask(ctx, req)
	if err != nil {
		if err.Error() == "AlreadyExists" {
			return nil, errs.NotCreatedError.Wrap(err, "task already exists %s", req.Task.Name)
		} else {
			return nil, fmt.Errorf("cloudtasks.CreateTask: %v", err)
		}
	}
	return createdTask, nil
}

// Get retrieves a Cloud Task by its name. It establishes a new client connection,
// creates a GetTaskRequest with basic response view, and fetches the task details.
//
// Parameters:
//   - ctx: The context.Context for the request
//   - name: The name of the task to retrieve, formatted as:
//     "projects/<PROJECT_ID>/locations/<LOCATION_ID>/queues/<QUEUE_ID>/tasks/<TASK_ID>"
//
// Returns:
//   - *cloudtaskspb.Task: The retrieved task if successful
//   - error: An error if the operation fails or the task cannot be found
func Get(ctx context.Context, name string) (*cloudtaskspb.Task, error) {
	c, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer func(c *cloudtasks.Client) {
		err = c.Close()
		if err != nil {
			log.Error().Stack().Err(err).Msg(err.Error())
		}
	}(c)

	req := &cloudtaskspb.GetTaskRequest{
		Name:         name,
		ResponseView: cloudtaskspb.Task_BASIC,
	}

	t, err := client.GetTask(ctx, req)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Delete removes a Cloud Task by its name. It establishes a new client connection
// and attempts to delete the specified task.
//
// Parameters:
//   - ctx: The context.Context for the request
//   - name: The name of the task to delete, formatted as:
//     "projects/<PROJECT_ID>/locations/<LOCATION_ID>/queues/<QUEUE_ID>/tasks/<TASK_ID>"
//
// Returns:
//   - error: An error if the operation fails or nil if the deletion is successful
func Delete(ctx context.Context, name string) error {
	c, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return err
	}
	defer func(c *cloudtasks.Client) {
		err = c.Close()
		if err != nil {
			log.Error().Stack().Err(err).Msg(err.Error())
		}
	}(c)

	req := &cloudtaskspb.DeleteTaskRequest{
		Name: name,
	}

	return client.DeleteTask(ctx, req)
}

// WaitForTaskCompletion polls a Cloud Task until it completes successfully or reaches the maximum retry attempts.
// It checks the task's status by making repeated GetTask requests with a fixed delay between attempts.
//
// Parameters:
//   - ctx: The context.Context for the request, which can be used to cancel the waiting operation
//   - client: The Cloud Tasks client instance to use for making requests
//   - task: The task to monitor for completion
//
// The function considers a task complete when its LastAttempt has a ResponseTime and
// the ResponseStatus code is 200. It will retry up to 3 times with a 10-second delay
// between attempts.
//
// Returns:
//   - error: nil if the task completes successfully, an error if the operation fails
//     or times out after maximum retries
//
// Deprecated: this is a naive blocking implementation and should be replaced with a more sophisticated way using goroutines and channels
func WaitForTaskCompletion(ctx context.Context, client *cloudtasks.Client, task *cloudtaskspb.Task, maxWaitDuration time.Duration) error {
	if maxWaitDuration <= 0 {
		return errorx.IllegalArgument.New("maxWaitDuration must be positive: %d", maxWaitDuration)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, maxWaitDuration)
	defer cancel()

	resultChan := ListenForTaskCompletion(timeoutCtx, client, task)

	select {
	case <-timeoutCtx.Done():
		if timeoutCtx.Err() == context.DeadlineExceeded {
			return errs.DurationExceeded.New("maxWaitDuration exceeded waiting for task completion after %s", timestamp.HumanReadableDuration(maxWaitDuration))
		}
		return timeoutCtx.Err()
	case result := <-resultChan:
		return result.Error
	}
}

// TaskResult represents the result of a task completion check
type TaskResult struct {
	Task  *cloudtaskspb.Task
	Error error
}

// ListenForTaskCompletion monitors a Cloud Task's completion status asynchronously using
// channels and goroutines. It asynchronously monitors the task status 5 seconds after the next scheduled attempt and sends the
// result through a channel when the task either completes successfully or encounters an error.
//
// Parameters:
//   - ctx: The context.Context for managing the monitoring lifecycle
//   - client: The Cloud Tasks client instance
//   - task: The task to monitor
//
// Returns:
//   - <-chan TaskResult: A channel that will receive the task completion result
func ListenForTaskCompletion(ctx context.Context, client *cloudtasks.Client, task *cloudtaskspb.Task) <-chan TaskResult {
	resultChan := make(chan TaskResult, 1)

	go func() {
		defer close(resultChan)
		ticker := time.NewTicker(task.ScheduleTime.AsTime().Sub(time.Now()) + time.Second*5)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				resultChan <- TaskResult{Task: task, Error: ctx.Err()}
				return
			case <-ticker.C:
				getTaskRequest := &cloudtaskspb.GetTaskRequest{
					Name:         task.Name,
					ResponseView: cloudtaskspb.Task_BASIC,
				}

				t, err := client.GetTask(ctx, getTaskRequest)
				if err != nil {
					resultChan <- TaskResult{Task: task, Error: err}
					return
				}

				if t.LastAttempt != nil && t.LastAttempt.ResponseTime != nil && t.LastAttempt.ResponseStatus.Code == int32(200) {
					resultChan <- TaskResult{Task: t, Error: nil}
					return
				} else {
					ticker.Reset(task.ScheduleTime.AsTime().Sub(time.Now()) + t.CreateTime.AsTime().Sub(t.LastAttempt.ResponseTime.AsTime()))
				}
			}
		}
	}()

	return resultChan
}
