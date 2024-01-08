package prime_test

import (
	"testing"

	"github.com/pengubco/algorithms/prime"
	"github.com/stretchr/testify/assert"
)

func TestIsPrime(t *testing.T) {
	assert.False(t, prime.IsPrime(-1))
	assert.False(t, prime.IsPrime(0))
	assert.False(t, prime.IsPrime(1))

	primeNumbers := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	for _, v := range primeNumbers {
		assert.True(t, prime.IsPrime(v))
	}

	compositeNumbers := []int{4, 6, 8, 10, 12, 14, 15, 16, 18, 20}
	for _, v := range compositeNumbers {
		assert.False(t, prime.IsPrime(v))
	}
}

func TestNextPrime(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	nextPrimeNumbers := []int{2, 3, 5, 5, 7, 7, 11, 11, 11, 11}
	for i, _ := range numbers {
		assert.Equal(t, nextPrimeNumbers[i], prime.NextPrime(numbers[i]))
	}
}

func TestNextPrime_Big(t *testing.T) {
	assert.Equal(t, 10_007, prime.NextPrime(10_000))
	assert.Equal(t, 100_003, prime.NextPrime(100_000))
}

func TestXxx(t *testing.T) {
	nums := []int{100, 200, 500, 1000, 5000}
	for _, n := range nums {
		assert.Equal(t, 1, prime.NextPrime(n))
	}
}
