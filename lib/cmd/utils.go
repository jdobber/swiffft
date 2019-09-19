package cmd

import "log"

// Check ...
func Check(e error) {
	if e != nil {
		log.Fatalln(e)
		//panic(e)
	}
}
