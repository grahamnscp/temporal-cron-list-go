package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"

	"github.com/aptible/supercronic/cronexpr"

	"cronlist"
	u "cronlist/utils"
)

const wkflDuration = 30

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

	// Minute 0-59 - Hour 0-23 - DoM 1-31 - Month 1-12 - DoW 0-6 (Sun-Sat, or 7 for Sun) - Localtime!
	//cronSchedule := "03 14 26 9 2" // Local time
	cronSchedule, err := parseCLIArgs(os.Args[1:])
	if err != nil {
		log.Fatalf("Parameter --cronschedule <cron schedule string> is required")
	}

	nextTime := cronexpr.MustParse(*cronSchedule).Next(time.Now())
	log.Printf("StartCronList: cronSchedule: %s, next time: %v", *cronSchedule, nextTime)

	startSeconds, _ := strconv.Atoi(fmt.Sprintf("%.0f", nextTime.Sub(time.Now()).Seconds()))
	endSeconds := startSeconds + wkflDuration

	utcSchedule, _ := decHour(*cronSchedule) // take hour off for UTC TZ
	log.Printf("StartCronList: utc cronSchedule: %s, startSeconds: %d, endSeconds: %d", utcSchedule, startSeconds, endSeconds)

	taskQueue := os.Getenv("TASK_QUEUE")
	workflowOptions := client.StartWorkflowOptions{
		ID:                       "cronwkfl-" + thisid,
		TaskQueue:                taskQueue,
		CronSchedule:             utcSchedule,
		WorkflowExecutionTimeout: (time.Duration(endSeconds) * time.Second),
	}

	message := uuid.New()
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, cronlist.CronListWorkflow, message)
	if err != nil {
		log.Fatalln("StartCronList: Unable to execute workflow", err)
	}
	log.Println("StartCronList: Started cron workflow", "WorkflowID:", we.GetID(), "RunID:", we.GetRunID(), "Schedule:", *cronSchedule)
}

/*
/* Functions
*/
func parseCLIArgs(args []string) (*string, error) {

	set := flag.NewFlagSet("cronliststart", flag.ExitOnError)
	cronScheduleStr := set.String("cronschedule", "", "cron schedule string")
	if err := set.Parse(args); err != nil {
		return nil, fmt.Errorf("failed parsing args: %w", err)
	} else if *cronScheduleStr == "" {
		return nil, fmt.Errorf("--cronschedule argument is required")
	}
	return cronScheduleStr, nil
}

// remove hour for BST to UTC
func decHour(s string) (string, error) {
	var m, h, d, n, w int
	cronFormat := "%d %d %d %d %d"
	if _, err := fmt.Sscanf(s, cronFormat, &m, &h, &d, &n, &w); err != nil {
		return "", err
	}
	h--
	return fmt.Sprintf(cronFormat, m, h, d, n, w), nil
}
