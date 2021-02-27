package main

import (
	_ "../ch2/matchers"
	"../ch2/search"
	"log"
	"os"
)

func init()  {
	log.SetOutput(os.Stdout)
}

func main()  {
	search.Run("president")
}
