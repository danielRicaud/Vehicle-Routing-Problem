package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

const MAX_DISTANCE = 12*60
type Point struct {
	x float64
	y float64
}

type Load struct {
	id int
	start Point
	end Point
	distance float64
	totalDistance float64
	assigned bool
}

type Truck struct {
	distance float64
	loads []Load
}

type Saving struct {
	load1Id int
	load2Id int
	savedDistance float64
}

var loadMap = make(map[int]*Load)
var savings = make([]Saving, 0)
var depot = Point{x: 0, y: 0}

func main() {
	args := os.Args

	parseLoadData(args[1])

	generateSavingsList()
}

func pointDistance(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow(p2.x - p1.x, 2) + math.Pow(p2.y - p1.y, 2))
}

func generateSavingsList() {
	for i := 1; i <= len(loadMap); i++ {
		for j := 1; j <= len(loadMap); j++ {
			if i != j {

				saving := Saving{
					load1Id: i,
					load2Id: j,
					// Distance(depot, load1end) + Distance(depot, load2end) - Distance(load1end, load2start)
					savedDistance: loadMap[i].totalDistance + loadMap[j].totalDistance - pointDistance(loadMap[i].end, loadMap[j].start),
				}
				savings = append(savings, saving)
			}
		}
	}

	// Sort the savings list in descending order by savedDistance
	sort.Slice(savings, func(i, j int) bool {
		return savings[i].savedDistance > savings[j].savedDistance
	})

	fmt.Println("Savings list:")
	for _, saving := range savings {
		fmt.Println(saving)
	}

}


func parseLoadData(filePath string)  {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Add each load to the loadMap so they can be referenced by loadNumber, format coordinates into Point structs
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, " ")

		if items[0] == "loadNumber" { // Skip the header line
			continue
		}

		id, _ := strconv.Atoi(items[0])
		// Process coordinates in the form of "(-42.51051149928979,-116.19788220835095) (-69.06284568487868,-44.12633704833111)"
		start := strings.Split(items[1][1 : len(items[1])-1], ",")
		startX, _ := strconv.ParseFloat(start[0], 64)
		startY, _ := strconv.ParseFloat(start[1], 64)
		startPoint := Point{x: startX, y: startY}

		end := strings.Split(items[2][1 : len(items[2])-1], ",")
		endX, _ := strconv.ParseFloat(end[0], 64)
		endY, _ := strconv.ParseFloat(end[1], 64)
		endPoint := Point{x: endX, y: endY}

		loadDistance := pointDistance(startPoint, endPoint)
		totalLoadDistance := pointDistance(depot, startPoint) + loadDistance

		load := &Load{
			id: id,
			start: startPoint,
			end: endPoint,
			distance: loadDistance,
			totalDistance: totalLoadDistance,
			assigned: false,
		}
		loadMap[id] = load
		fmt.Println(load)
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return
	}
}
