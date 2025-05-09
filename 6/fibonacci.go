package main

func fibonacci(n int) int {
	if n < 0 {
		return 0
	} else if n <= 1 {
		return n
	} else {
		return fibonacci(n-1) + fibonacci(n-2)
	}
}
