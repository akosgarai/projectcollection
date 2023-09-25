package grifts

import (
	"fmt"

	. "github.com/gobuffalo/grift/grift"
)

var _ = Namespace("processor", func() {

	Desc("application", "Task Description")
	Add("application", func(c *Context) error {
		fmt.Println("Hello!")
		return nil
	})

})
