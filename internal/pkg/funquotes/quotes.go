package funquotes

import (
	"math/rand"
	"time"
)

var quotes []string

func init() {
	quotes = []string{
		"Sir, take a deep breath.",
		"Sir, there are still terabytes of calculations required before an actual flight is possible.",
		"Test complete. Preparing to power down and begin diagnostics.",
		"Commencing automated assembly.",
		"It is a tight fit sir.",
		"Sir, the more you struggle the more this is going to hurt.",
		"You are not authorized to access this area.",
		"For you sir, always.",
		"Sir, we will lose power before we penetrate that shell.",
	}
	//init random generator
	rand.Seed(time.Now().Unix())
}

func GiveMeAQuote() string {
	//give me a quote randomly between {0,n}
	return quotes[rand.Intn(len(quotes))]
}
