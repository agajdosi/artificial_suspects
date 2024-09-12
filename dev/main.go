package main

// JUST FOR TESTING NOW
// GONNA BE A HELPER BINARY TO GENERATE AI DESCRIPTIONS ETC...

import (
	"fmt"
	"suspects/database"
)

func main() {
	path := database.GetDataDirPath()
	fmt.Println("Datadir path is:", path)
}
