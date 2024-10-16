package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

const MAX_DISTANCE = 12*60
type Point struct {
	x float64
	y float64
}

type Load struct {
	id int
	start string
	end string
	distance float64
	assigned float64
}

type Truck struct {
	distance float64
	loads []Load
}

func main() {
	args := os.Args

	parseFile(args[1])
}

func pointDistance(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow(p2.x - p1.x, 2) + math.Pow(p2.y - p1.y, 2))
}


func parseFile(filePath string)  {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate over each line
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, " ")

		if items[0] == "loadNumber" {
			continue
		}
		fmt.Println(items)
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return
	}

}
