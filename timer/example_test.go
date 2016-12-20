package timer_test

import (
	"fmt"
	"ggs/timer"
	"time"
)

func ExampleTimer() {
	d := timer.NewDispatcher(10)

	// timer 1
	d.AfterFunc(1, func() {
		fmt.Println("My name is ggs")
	})

	// timer 2
	t := d.AfterFunc(1, func() {
		fmt.Println("will not print")
	})
	t.Stop()

	// dispatch
	(<-d.ChanTimer).Cb()

	// Output:
	// My name is ggs
}

func ExampleCronExpr() {
	cronExpr, err := timer.NewCronExpr("0 * * * *")
	if err != nil {
		return
	}

	fmt.Println(cronExpr.Next(time.Date(
		2000, 1, 1,
		20, 10, 5,
		0, time.UTC,
	)))

	// Output:
	// 2000-01-01 21:00:00 +0000 UTC
}

func ExampleCron() {
	d := NewDispatcher(10)

	// cron expr
	cronExpr, err := timer.NewCronExpr("* * * * * *")
	if err != nil {
		return
	}

	// cron
	var c *timer.Cron
	c := d.CronFunc(cronExpr, func() {
		fmt.Println("My name is ggs")
		c.Stop()
	})

	// dispatch
	(<-timer.ChanTimer).Cb()

	// Output:
	// My name is ggs
}
