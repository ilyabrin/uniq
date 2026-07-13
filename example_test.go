package uniq_test

import (
	"fmt"

	"github.com/ilyabrin/uniq"
)

func ExampleNew() {
	id := uniq.New(10, 3)
	fmt.Println(len(id)) // e.g. "aKbcDefGhi"
	// Output: 10
}

func ExampleToken() {
	token, err := uniq.Token(32)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(token)) // e.g. "kX9mPq2RtY7wBn4ZcV1sAj8dHf3GuLe0"
	// Output: 32
}

func ExampleMustToken() {
	apiKey := uniq.MustToken(64)
	fmt.Println(len(apiKey))
	// Output: 64
}

func ExampleGenerate() {
	id, err := uniq.Generate(20, uniq.WithDigits())
	if err != nil {
		panic(err)
	}
	fmt.Println(len(id)) // e.g. "x7hq2mzk91tbv4wnp0ce"
	// Output: 20
}

func ExampleWithUppercase() {
	id := uniq.MustGenerate(12, uniq.WithUppercase(4))
	upper := 0
	for _, c := range id {
		if c >= 'A' && c <= 'Z' {
			upper++
		}
	}
	fmt.Println(upper) // e.g. "abQcdEfGhIjk" has exactly 4
	// Output: 4
}

func ExampleWithAlphabet() {
	// A 6-digit confirmation code, e.g. for SMS.
	code := uniq.MustGenerate(6, uniq.WithAlphabet("0123456789"))
	fmt.Println(len(code)) // e.g. "482910"
	// Output: 6
}
