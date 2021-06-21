package main

import (
	"encoding/csv"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	//	"fmt"
	"log"
	"os"
	"time"
)

const (
	NumUsers        = 100000
	NumPayments     = 1000000
	MaxPaymentCents = 100000000
	MaxZip          = 99999
	MaxAge          = 120
)

// FYI, this is how test data was generated. It's unlikely that
// you need to modify this, and if you do, you may break the
// expected (hard coded) test values
func main() {
	rand.Seed(0xdeadbeef)

	// Get lists of names and words
	// NOTE don't do multiple passes like this in prod
	namesBytes, err := ioutil.ReadFile("/usr/share/dict/propernames")
	if err != nil {
		log.Fatal(err)
	}
	names := strings.Split(string(namesBytes), "\n")
	numNames := len(names)

	wordsBytes, err := ioutil.ReadFile("/usr/share/dict/words")
	if err != nil {
		log.Fatal(err)
	}
	words := strings.Split(string(wordsBytes), "\n")
	numWords := len(words)

	// Write out random user data
	f, err := os.OpenFile("metrics/users.csv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Error opening users.csv", err)
	}
	w := csv.NewWriter(f)
	for i := 0; i < NumUsers; i++ {
		record := []string{
			strconv.Itoa(i),
			names[rand.Intn(numNames-1)] + " " + names[rand.Intn(numNames-1)],
			strconv.Itoa(rand.Intn(MaxAge)),
			strings.Title(strconv.Itoa(rand.Intn(500)) + " " + words[rand.Intn(numWords-1)] + " St, " + words[rand.Intn(numWords-1)] + "town"),
			strconv.Itoa(rand.Intn(MaxZip)),
		}
		if err := w.Write(record); err != nil {
			log.Fatalln("Error writing csv:", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	// Write out payment data
	f, err = os.OpenFile("metrics/payments.csv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Error opening users.csv", err)
	}
	w = csv.NewWriter(f)
	minDate := time.Date(2010, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	maxDate := time.Date(2021, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	for i := 0; i < NumPayments; i++ {
		// amount in cents, datetime, user id
		record := []string{
			strconv.Itoa(rand.Intn(MaxPaymentCents)),
			time.Unix(rand.Int63n(maxDate-minDate)+minDate, 0).Format(time.RFC3339),
			strconv.Itoa(rand.Intn(NumUsers)),
		}
		if err := w.Write(record); err != nil {
			log.Fatalln("Error writing csv:", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

}
