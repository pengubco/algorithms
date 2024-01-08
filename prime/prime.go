package prime

import "math"

// IsPrime returns true iff the given number is a prime number.
func IsPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n == 2 || n == 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	upperBound := math.Sqrt(float64(n)) + 1
	for i := 5; i <= int(upperBound); i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// NextPrime generates the smallest prime number larger than n.
func NextPrime(n int) int {
	for i := n + 1; ; i++ {
		if IsPrime(i) {
			return i
		}
	}
}
