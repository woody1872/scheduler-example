package main

func main() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(5).Seconds().Do(func(){ ... })
	
	// strings parse to duration
	s.Every("5m").Do(func(){ ... })
	
	s.Every(5).Days().Do(func(){ ... })
	
	// cron expressions supported
	s.Cron("*/1 * * * *").Do(task) // every minute
	
	// you can start running the scheduler in two different ways:
	// starts the scheduler asynchronously
	s.StartAsync()
	// starts the scheduler and blocks current execution path 
	s.StartBlocking()	
}
