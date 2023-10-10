package main

import (
  "log"
  "os"

  "go.temporal.io/sdk/client"
  "go.temporal.io/sdk/worker"

  "cronlist"
  u "cronlist/utils"
)

func main() {
  log.Printf("%sGo worker starting..%s", u.ColorGreen, u.ColorReset)

  // The client and worker are heavyweight objects that should be created once per process.
  clientOptions, err := u.LoadClientOptions()
  if err != nil {
    log.Fatalln("Failed to load Temporal Cloud environment:", err)
  }

  log.Println("Go worker connecting to server..")

  c, err := client.Dial(clientOptions)
  if err != nil {
    log.Fatalln("Unable to create client", err)
  }
  defer c.Close()

  taskQueue := os.Getenv("TASK_QUEUE")
  log.Println("Go worker initialising..")
  w := worker.New(c, taskQueue, worker.Options{})

  log.Println("Registering for workflow and activites..")
  w.RegisterWorkflow(cronlist.CronListWorkflow)
  w.RegisterActivity(cronlist.CronActivity)

  log.Printf("%sGo worker listening on %s task queue..%s", u.ColorGreen, "CronListTQ", u.ColorReset)
  err = w.Run(worker.InterruptCh())
  if err != nil {
    log.Fatalln("Unable to start worker", err)
  }

  log.Printf("%sGo worker stopped.%s", u.ColorGreen, u.ColorReset)
}
