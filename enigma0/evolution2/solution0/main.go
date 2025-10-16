package main

import "fmt"

func main() {
	var rawr a
	rawr = "asdf"
	test(rawr)
}

type a string

func (asdf a) String() string {
	return string(asdf)
}

func test(values ...string) {

}
