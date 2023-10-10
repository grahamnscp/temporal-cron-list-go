package main

import (
	"log"
  "os"
  "fmt"
  "context"
  "encoding/json"
  "bytes"
  "time"
  "math"
  "strings"

	"go.temporal.io/sdk/client"
  "go.temporal.io/api/workflowservice/v1"
  commonpb "go.temporal.io/api/common/v1"
  historypb "go.temporal.io/api/history/v1"
  enumspb "go.temporal.io/api/enums/v1"
  "go.temporal.io/sdk/converter"

  "github.com/aptible/supercronic/cronexpr"

	u "cronlist/utils"
)

var log_level = strings.ToLower(os.Getenv("LOG_LEVEL"))

func main() {

  log.Printf("CronList: Started")

	clientOptions, err := u.LoadClientOptions()
	if err != nil {
		log.Fatalln("CronList: Failed to load Temporal Cloud environment:", err)
	}

	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("CronList: Unable to create client", err)
	}
	defer c.Close()


  // list running cron workflows
  namespace := os.Getenv("TEMPORAL_NAMESPACE")

  //query := "CustomStringField='ACTIVE-CRON' and CloseTime is null" // search attribute not set until workflow runs
  // CloseTime is null or ExecutionStatus='Running' 
  query := "WorkflowType='CronListWorkflow' and CloseTime is null"

  log.Printf("CronList: Query: '%s'", query)

  var exec *commonpb.WorkflowExecution
  var nextPageToken []byte
  for hasMore := true; hasMore; hasMore = len(nextPageToken) > 0 {
    resp, err := c.ListWorkflow(context.Background(), &workflowservice.ListWorkflowExecutionsRequest{
      Namespace:     namespace,
      PageSize:      10,
      NextPageToken: nextPageToken,
      Query:         query,
    })
    if err != nil {
      log.Fatal("CronList: ListWorkflows returned an error,", err)
    }
    log.Printf("CronList: Executions: %d", len(resp.Executions))

    for i := range resp.Executions {
      exec = resp.Executions[i].Execution
      //log.Printf("CronList: Execution: WorkflowId: %v, RunId: %v\n", exec.WorkflowId, exec.RunId)

      // using execution get workflow history
      history, err := getHistory(c, context.Background(), exec)
      if err != nil {
        log.Printf("CronList: getHistory error: %v", err)
      }
      //historyEventCount := len(history)
      //fmt.Printf("Event History Count: %d\n", historyEventCount)

      // examine event history
      for h := range history {
        //fmt.Printf("history[%d]:\n", h)

        pc := converter.NewJSONPayloadConverter()
        payload, _ := pc.ToPayload(history[h])
        jsondatastr := string(payload.Data)

        // Unmarshall json data into a map container for decoded the JSON structure into
        var c map[string]interface{}
        err := json.Unmarshal([]byte(jsondatastr), &c)
        if err != nil {
          panic(err)
        }
        event_type_num := int32(c["event_type"].(float64))

        if log_level == "debug" {
          event_type := string(enumspb.EventType_name[event_type_num])
          fmt.Printf("Workflow Event %d: %s\n", event_type_num, event_type)

          // dump event as formatted json output
          data := &bytes.Buffer{}
          if err := json.Indent(data, []byte(jsondatastr), "", "  "); err != nil {
            panic(err)
          }
          fmt.Println(data.String())
        }

        // Pull specific attributes..
        if event_type_num == enumspb.EventType_value["WorkflowExecutionStarted"] {
          attributes := c["Attributes"].(map[string]interface{})["workflow_execution_started_event_attributes"]
          cron_schedule := fmt.Sprintf("%v", attributes.(map[string]interface{})["cron_schedule"])
          workflow_type := attributes.(map[string]interface{})["workflow_type"].(map[string]interface{})["name"]
          backoff_dur := time.Duration(math.Round(attributes.(map[string]interface{})["first_workflow_task_backoff"].(float64)))
          runTime := cronexpr.MustParse(cron_schedule).Next(time.Now())

          fmt.Printf("Execution: WorkflowId: %v, RunId: %v\n", exec.WorkflowId, exec.RunId)
          fmt.Printf("  workflow_type: %s, cron_schedule: '%s', first_workflow_task_backoff: %v, cronTime: '%s'\n", workflow_type, cron_schedule, backoff_dur, runTime)
        }

      }
      fmt.Printf("-----\n")
    }
    nextPageToken = resp.NextPageToken
  }
}

/* getHistory */
func getHistory(c client.Client, ctx context.Context, execution *commonpb.WorkflowExecution) ([]*historypb.HistoryEvent, error) {

  iter := c.GetWorkflowHistory(ctx,
    execution.GetWorkflowId(),
    execution.GetRunId(),
    false,
    enumspb.HISTORY_EVENT_FILTER_TYPE_ALL_EVENT)
  var events []*historypb.HistoryEvent
  for iter.HasNext() {
    event, err := iter.Next()
    if err != nil {
      return nil, err
    }
    events = append(events, event)
  }
  return events, nil
}
