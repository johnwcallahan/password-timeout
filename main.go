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
	passwordTimeout := NewPasswordTimeout("hellogoodbye", 3, []int{5, 10, 30})
	passwordTimeout.Run()
}

type PasswordTimeout struct {
	password                     string
	maxAttempts                  int
	timeoutSequence              []int
	deadline                     time.Time
	attemptCount, currentTimeout int
}

func NewPasswordTimeout(password string, maxAttempts int, timeoutSequence []int) *PasswordTimeout {
	return &PasswordTimeout{
		password:        password,
		timeoutSequence: timeoutSequence,
		maxAttempts:     maxAttempts,
		deadline:        time.Now(),
		attemptCount:    0,
	}
}

func (p PasswordTimeout) Run() {
	for {
		input := p.Prompt()

		// If the user is locked out, prompt them to retry.
		if p.isLockedOut() {
			p.Retry()
			continue
		}

		if p.Test(input) {
			// The password matches; exit the main program loop.
			p.Success()
			return
		}

		p.attemptCount++

		// If the timeout sequence has been exhausted, ban the user.
		if p.currentTimeout > len(p.timeoutSequence)-1 {
			p.Ban()
			return
		}

		// If the maximum number of attempts has been exceeded, set the next timeout.
		if p.attemptCount >= p.maxAttempts {
			p.setTimeout()
			p.Retry()
		}

		// Otherwise, keeping going around the loop (prompt the user again).
	}
}

func (p PasswordTimeout) Prompt() string {
	fmt.Println("")
	fmt.Println("Enter your password")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("An error occurred")
	}
	return strings.TrimSuffix(input, "\n")
}

func (p PasswordTimeout) Test(val string) bool {
	return val == p.password
}

func (p PasswordTimeout) Retry() {
	dur := time.Until(p.deadline).Round(time.Second)
	fmt.Printf("Try again in %s\n\n", dur)
}

func (p PasswordTimeout) Success() {
	fmt.Println("Success!")
}

func (p PasswordTimeout) Ban() {
	fmt.Println("You're banned!")
}

func (p PasswordTimeout) isLockedOut() bool {
	return time.Now().Before(p.deadline)
}

func (p *PasswordTimeout) setTimeout() {
	timeout := p.timeoutSequence[p.currentTimeout]
	p.deadline = time.Now().Add(time.Second * time.Duration(timeout))
	p.currentTimeout++
}
