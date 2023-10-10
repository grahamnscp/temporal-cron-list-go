package main

import (
  "context"
  "log"
  "fmt"
  "math/rand"
  "os"

  "github.com/pborman/uuid"
  "go.temporal.io/sdk/client"

  "cronlist"
  u "cronlist/utils"
)

func main() {

  thisid := fmt.Sprint(rand.Intn(99999))
  log.Printf("StartCronList: Message: %s", thisid)

  clientOptions, err := u.LoadClientOptions()
  if err != nil {
    log.Fatalln("StartCronList: Failed to load Temporal Cloud environment:", err)
  }

  c, err := client.Dial(clientOptions)
  if err != nil {
    log.Fatalln("StartCronList: Unable to create client", err)
  }
  defer c.Close()

  taskQueue := os.Getenv("TASK_QUEUE")
  workflowOptions := client.StartWorkflowOptions{
    ID:        "cronwkfl-" + thisid,
    TaskQueue: taskQueue,
  }

  message := uuid.New()
  we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, cronlist.CronListWorkflow, message)
  if err != nil {
    log.Fatalln("StartCronList: Unable to execute workflow", err)
  }
  log.Println("StartCronList: Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
