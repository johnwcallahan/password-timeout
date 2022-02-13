# password-timeout

A toy program that prompts for a password and locks you out at progressive intervals if the password is incorrect.

Just a little experiment with Go's `time` package.

```
Enter your password
1234

Enter your password
2345
Try again in 5s


Enter your password
3456
Try again in 10s
```
