package internal

import "github.com/oklog/ulid"

// // Create a new ULID (using current time and default entropy source)
// newUlid := ulid.Make()
// fmt.Println("New ULID:", newUlid)

// // Create a ULID from a specific time
// timestamp := time.Now()
// ulidFromTime := ulid.MustNew(ulid.Timestamp(timestamp), ulid.DefaultEntropy())
// fmt.Println("ULID from time:", ulidFromTime)

// // Parse a ULID string
// ulidString := "01ARZ3NDEKTSV4RRFFQ69G5FAV"
// parsedUlid, err := ulid.Parse(ulidString)
// if err != nil {
// 	panic(err)
// }
// fmt.Println("Parsed ULID:", parsedUlid)

func GenerateUniqueId() string {
	newUlid, err := ulid.New(ulid.Now(), nil)
	if err != nil {
		panic(err)
	}
	return newUlid.String()
}

func ParseUlid(id string) uint64 {
	parsedUlid, err := ulid.Parse(id)
	if err != nil {
		panic(err)
	}

	return parsedUlid.Time()
}
