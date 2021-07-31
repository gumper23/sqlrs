package sqlrs

import "fmt"

type Resultset struct {
}

func New() *Resultset {
	return &Resultset{}
}

func (*Resultset) Hello(name string) string {
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	return message
}
