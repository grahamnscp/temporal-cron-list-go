package cronlist

import (
	"context"
  "log"

  u "cronlist/utils"
)

func CronActivity (ctx context.Context, message string) (string, error) {

  log.Printf("%sCronActivity:%s Started with message: %s %s\n", u.ColorGreen, u.ColorBlue, message, u.ColorReset)

  // Activity action..
  log.Printf("%sCronActivity:%s Action for message: %s %s\n", u.ColorGreen, u.ColorBlue, message, u.ColorReset)

  log.Printf("%sCronActivity:%s Complete.%s\n", u.ColorGreen, u.ColorBlue, u.ColorReset)

	return "Message sent", nil
}

