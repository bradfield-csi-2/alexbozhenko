package metrics

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"

	"unsafe"
)

type UserId int
type UserMap map[UserId]*User

type Address struct {
	fullAddress string
	zip         int
}

type Payment struct {
	amount_cents int
	//	time         time.Time
}

type User struct {
	id       UserId
	name     string
	age      uint8
	address  Address
	payments []Payment
}

func AverageAge(users UserMap) float64 {
	average, count := 0.0, 0.0
	for _, u := range users {
		count += 1
		average += (float64(u.age) - average) / count
	}
	return average
}

func AveragePaymentAmount(users UserMap) float64 {
	average, count := 0.0, 0.0
	for _, u := range users {
		for _, p := range u.payments {
			count += 1
			//		amount := float64(p.amount_cents.dollars) + float64(p.amount.cents)/100
			average += (float64(p.amount_cents) - average) / count
		}
	}
	return average / 100
}

// Compute the standard deviation of payment amounts
func StdDevPaymentAmount(users UserMap) float64 {
	mean := AveragePaymentAmount(users)
	squaredDiffs, count := 0.0, 0.0
	for _, u := range users {
		for _, p := range u.payments {
			count += 1
			//amount := float64(p.amount.dollars) + float64(p.amount.cents)/100
			diff := float64(p.amount_cents/100) - mean
			squaredDiffs += diff * diff
		}
	}
	_ = unsafe.Pointer(&count)
	return math.Sqrt(squaredDiffs / count)

}

func LoadData() UserMap {
	f, err := os.Open("users.csv")
	if err != nil {
		log.Fatalln("Unable to read users.csv", err)
	}
	reader := csv.NewReader(f)
	userLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse users.csv as csv", err)
	}

	users := make(UserMap, len(userLines))
	for _, line := range userLines {
		id, _ := strconv.Atoi(line[0])
		name := line[1]
		age, _ := strconv.Atoi(line[2])
		address := line[3]
		zip, _ := strconv.Atoi(line[3])
		users[UserId(id)] = &User{UserId(id), name, uint8(age), Address{address, zip}, []Payment{}}
	}

	f, err = os.Open("payments.csv")
	if err != nil {
		log.Fatalln("Unable to read payments.csv", err)
	}
	reader = csv.NewReader(f)
	paymentLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse payments.csv as csv", err)
	}

	for _, line := range paymentLines {
		userId, _ := strconv.Atoi(line[2])
		paymentCents, _ := strconv.Atoi(line[0])
		//		datetime, _ := time.Parse(time.RFC3339, line[1])
		//datetime := time.Time{}
		users[UserId(userId)].payments = append(users[UserId(userId)].payments, Payment{
			paymentCents,
			//	datetime,
		})
	}

	return users
}
