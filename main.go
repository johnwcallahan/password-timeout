package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	lockout := NewLockout("hellogoodbye", 2, []int{5, 10})
	lockout.Run()
}

type Lockout struct {
	password                     string
	maxAttempts                  int
	timeoutSequence              []int
	deadline                     time.Time
	attemptCount, currentTimeout int
}

func NewLockout(password string, maxAttempts int, timeoutSequence []int) *Lockout {
	return &Lockout{
		password:        password,
		timeoutSequence: timeoutSequence,
		maxAttempts:     maxAttempts,
		deadline:        time.Now(),
		attemptCount:    0,
	}
}

func (l Lockout) Run() {
	for {
		input := l.Prompt()

		// If the user is locked out, prompt them to retry.
		if l.isLockedOut() {
			l.Retry()
			continue
		}

		if l.Test(input) {
			// The password matches; exit the main program loop.
			l.Success()
			return
		}

		l.attemptCount++

		// If the timeout sequence has been exhausted, ban the user.
		if l.currentTimeout > len(l.timeoutSequence)-1 {
			l.Ban()
			return
		}

		// If the maximum number of attempts has been exceeded, set the next timeout.
		if l.attemptCount >= l.maxAttempts {
			l.setLockout()
			l.Retry()
		}

		// Otherwise, keeping going around the loop (prompt the user again).
	}
}

func (l Lockout) Prompt() string {
	fmt.Println("")
	fmt.Println("Enter your password")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("An error occurred")
	}
	return strings.TrimSuffix(input, "\n")
}

func (l Lockout) Test(val string) bool {
	return val == l.password
}

func (l Lockout) Retry() {
	dur := time.Until(l.deadline).Round(time.Second)
	fmt.Printf("Try again in %s\n\n", dur)
}

func (l Lockout) Success() {
	fmt.Println("Success!")
}

func (l Lockout) Ban() {
	fmt.Println("You're locked out!")
}

func (l Lockout) isLockedOut() bool {
	return time.Now().Before(l.deadline)
}

func (l *Lockout) setLockout() {
	timeout := l.timeoutSequence[l.currentTimeout]
	l.deadline = time.Now().Add(time.Second * time.Duration(timeout))
	l.currentTimeout++
}
