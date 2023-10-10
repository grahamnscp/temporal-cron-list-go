package cronlist

import (
	"time"
  "fmt"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

  u "cronlist/utils"
)

func CronListWorkflow(ctx workflow.Context, message string) (string, error) {

  logger := workflow.GetLogger(ctx)
  logger.Info(u.ColorGreen, "CronList-Workflow:", u.ColorReset, "Started", "-", workflow.GetInfo(ctx).WorkflowExecution.ID)

  // RetryPolicy specifies how to automatically handle retries if an Activity fails.
  activityretrypolicy := &temporal.RetryPolicy{
    InitialInterval:     time.Second,
    BackoffCoefficient:  2.0,
		MaximumInterval:     time.Minute,         // Short MaxiumInterval!
    MaximumAttempts:     10, 
  }

  // ActivityOptions
  activityoptions := workflow.ActivityOptions{
    StartToCloseTimeout: time.Minute,         // Timeout options specify when to automatically timeout Activity functions.
		HeartbeatTimeout:    2 * time.Second,     // Short timeout to make activity fail over very fast if unresponsive
    RetryPolicy:         activityretrypolicy, // Temporal retries failed Activities by default.
  }
	ctx = workflow.WithActivityOptions(ctx, activityoptions)

  // Set search attribute status to ACTIVE
  _ = u.UpcertSearchAttribute(ctx, "CustomStringField", "ACTIVE-CRON")

  var activityOutput string

  activityErr := workflow.ExecuteActivity(ctx, CronActivity, message).Get(ctx, &activityOutput)

	if activityErr != nil {
    // Set search attribute status to FAILED
    _ = u.UpcertSearchAttribute(ctx, "CustomStringField", "FAILED-CRON")
    logger.Info(u.ColorGreen, "CronList-Workflow:", u.ColorReset, "Failed", u.ColorRed, "CronActivity returned failure:", activityErr, u.ColorReset)
    return activityOutput, fmt.Errorf("CronActivity: failed for message: %s", message)
	}

  // All done
  // Set search attribute status to COMPLETE
  _ = u.UpcertSearchAttribute(ctx, "CustomStringField", "COMPLETE-CRON")
  logger.Info(u.ColorGreen, "CronList-Workflow:", u.ColorReset, "Complete", "-", workflow.GetInfo(ctx).WorkflowExecution.ID)
	return activityOutput, nil
}

