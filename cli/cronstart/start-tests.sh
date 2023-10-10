#!/bin/bash

cronSchedule=`date +'%M %H %d %m %u'`
year=`date +'%Y'`
min=`echo $cronSchedule | cut -d " " -f 1`
hour=`echo $cronSchedule | cut -d " " -f 2`
dom=`echo $cronSchedule | cut -d " " -f 3`
month=`echo $cronSchedule | cut -d " " -f 4`
dow=`echo $cronSchedule | cut -d " " -f 5`

tomorrow_dom=$(($dom+1))
tomorrow_dow=$(($dow+1))

echo "Current time: $year-$month-$dom $hour:$min $dow (tomorrow dom: $tomorrow_dom, dow: $tomorrow_dow)"
echo cronSchedule now \"$min $hour $dom $month $dow\"

# in 5 minutes
in5min=$((min+5))
hour5=$hour
if [[ $in5min -gt 54 ]]
then
  #inc hour
  in5min=$(($in5min-60))
  hour5=$(($hour+1))
fi
# in 10 minutes
in10min=$((min+10))
hour10=$hour
if [[ $in10min -gt 59 ]]
then
  #inc hour
  in10min=$(($in10min-60))
  hour10=$(($hour+1))
fi

# Start workflows (just fyi, haven't checked month threasholds for tests)
echo Starting cron workflow for +5 minutes, cronSchedule: \"$in5min $hour5 $dom $month $dow\"
go run main.go --cronschedule "$in5min $hour $dom $month $dow"

echo Starting cron workflow for tomorrow current time, cronSchedule: \"$min $hour $tomorrow_dom $month $tomorrow_dow\"
go run main.go --cronschedule "$min $hour $tomorrow_dom $month $tomorrow_dow"

echo Starting cron workflow for tomorrow +5 mins, cronSchedule: \"$in5min $hour5 $tomorrow_dom $month $tomorrow_dow\"
go run main.go --cronschedule "$in5min $hour $tomorrow_dom $month $tomorrow_dow"

echo Starting cron workflow for tomorrow +10 mins, cronSchedule:  \"$in10min $hour10 $tomorrow_dom $month $tomorrow_dow\"
go run main.go --cronschedule "$in10min $hour $tomorrow_dom $month $tomorrow_dow"

