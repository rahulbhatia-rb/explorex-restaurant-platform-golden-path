package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/rahulbhatia-rb/explorex-restaurant-platform-golden-path/internal/policy"
)

func main() {
	path := flag.String("contract", "", "service contract JSON")
	flag.Parse()
	if *path == "" {
		fmt.Fprintln(os.Stderr, "usage: platformctl -contract <file>")
		os.Exit(2)
	}

	b, err := os.ReadFile(*path)
	if err != nil {
		panic(err)
	}

	var c policy.Contract
	if err := json.Unmarshal(b, &c); err != nil {
		panic(err)
	}

	r := policy.Evaluate(c)
	out, _ := json.MarshalIndent(r, "", "  ")
	fmt.Println(string(out))
	if !r.Allowed {
		os.Exit(1)
	}
}
