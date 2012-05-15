package main

import (
	"fmt"
	"github.com/mattn/go-simplenote"
)

func main() {
	c, err := simplenote.NewClient("your-id", "your-pass")
	if err != nil {
		fmt.Println(err)
		return
	}
	notes, err := c.GetNotes()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, note := range notes {
		err = c.GetNote(&note)
		if err != nil {
			continue
		}
		println(note.Content)
	}
}
