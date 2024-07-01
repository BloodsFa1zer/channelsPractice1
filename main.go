//Task 1: Concurrent Prime Number Generator
//Description:

//Write a program that generates prime numbers concurrently.
//One goroutine should generate numbers from 2 upwards and send them to a channel.
//Multiple worker goroutines should receive numbers from this channel, check if they are prime,
//and send the prime numbers to a results channel.
//The main goroutine should collect and print the prime numbers from the results channel.

package main

import (
	"fmt"
	"math"
	"sync"
)

// Function to generate numbers from 2 upwards
func numberGenerator(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

// Function to check if a number is prime
func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// Function to process numbers and send primes to results channel
func primeWorker(ch <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		if isPrime(num) {
			results <- num
		}
	}
}

func main() {
	numCh := make(chan int)
	resultsCh := make(chan int)
	var wg sync.WaitGroup

	// Start number generator
	go numberGenerator(numCh)

	// Start multiple worker goroutines
	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go primeWorker(numCh, resultsCh, &wg)
	}

	// Start a goroutine to close resultsCh when all workers are done
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// Collect and print prime numbers from resultsCh
	go func() {
		for prime := range resultsCh {
			fmt.Println(prime)
		}
	}()

	// Run for a limited time or until a certain number of primes are found
	// For simplicity, we'll stop after finding 10 primes
	foundPrimes := 0
	for prime := range resultsCh {
		fmt.Println(prime)
		foundPrimes++
		if foundPrimes >= 10 {
			close(numCh)
			break
		}
	}

	fmt.Println("Done generating primes")
}
