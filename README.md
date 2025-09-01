# scheduler

## Example

### Run a Job Daily at a Specific Time

```go
package main

import (
	"log"
	"os"
	"os/signal"
	"scheduler-project/scheduler"
	"syscall"
)

func main() {
	s := scheduler.NewScheduler()

	s.AddJob(scheduler.Daily(8, 30, 0), func() {
		log.Println("Running daily job: Sending morning report...")
	})

	s.Start()
	log.Println("Scheduler started. Daily job scheduled for 8:30 AM. Press Ctrl+C to exit.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	s.Stop()
	log.Println("Scheduler stopped.")
}
```

### Run a Job at the Top of Every Hour

```go
package main

import (
	"log"
	"os"
	"os/signal"
	"scheduler-project/scheduler"
	"syscall"
)

func main() {
	s := scheduler.NewScheduler()

	s.AddJob(scheduler.AtHour(1), func() {
		log.Println("Running hourly job: Clearing temporary cache...")
	})

	s.Start()
	log.Println("Scheduler started. Hourly job is scheduled. Press Ctrl+C to exit.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	s.Stop()
	log.Println("Scheduler stopped.")
}
```

### Run a Job at the Top of Every Two Hours

Runs every second hour (ie 0200, 0400)

```go
package main

import (
	"log"
	"os"
	"os/signal"
	"scheduler-project/scheduler"
	"syscall"
)

func main() {
	s := scheduler.NewScheduler()

	s.AddJob(scheduler.AtHour(2), func() {
		log.Println("Running bi-hourly job: Performing data aggregation...")
	})

	s.Start()
	log.Println("Scheduler started. Bi-hourly job is scheduled. Press Ctrl+C to exit.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	s.Stop()
	log.Println("Scheduler stopped.")
}
```

### Run a Job Every X Minutes

```go
package main

import (
	"log"
	"os"
	"os/signal"
	"scheduler-project/scheduler"
	"syscall"
)

func main() {
	s := scheduler.NewScheduler()

	s.AddJob(scheduler.AtMinute(5), func() {
		log.Println("Running every 5 mins: Checking system health...")
	})

	s.Start()
	log.Println("Scheduler started. Health check will run every 5 minutes. Press Ctrl+C to exit.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	s.Stop()
	log.Println("Scheduler stopped.")
}
```

### Run a Job Using a Cron Expression

```go
package main

import (
	"log"
	"os"
	"os/signal"
	"scheduler-project/scheduler"
	"syscall"
)

func main() {
	s := scheduler.NewScheduler()

	cronSchedule, err := scheduler.Cron("30 * * * *")
	if err != nil {
		log.Fatalf("Invalid cron expression: %v", err)
	}
	s.AddJob(cronSchedule, func() {
		log.Println("Running cron job: Backing up configuration...")
	})

	s.Start()
	log.Println("Scheduler started. Cron job is scheduled. Press Ctrl+C to exit.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	s.Stop()
	log.Println("Scheduler stopped.")
}
```

