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

type Point struct {
	x float64
	y float64
}

type Load struct {
	id int
	start Point
	end Point
	loadDistance float64
	depotToStart float64
	depotToEnd float64
	truck *Truck
}

type Truck struct {
	time float64
	loads []*Load
}

type Saving struct {
	load1Id int
	load2Id int
	savedDistance float64
}

const MAX_TIME = 12*60

var loadMap = make(map[int]*Load)
var savings = make([]Saving, 0)
var depot = Point{x: 0, y: 0}
var trucks = make([]*Truck, 0)

func main() {
	args := os.Args

	parseLoadData(args[1])

	generateSavingsList()

	processSavingsList()
}

func pointDistance(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow(p2.x - p1.x, 2) + math.Pow(p2.y - p1.y, 2))
}

func findLoadIndex(loads []*Load, load *Load) int {
	for i, current := range loads {
			if current == load {
					return i
			}
	}
	return -1 // Not found
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

	// Read the file line by line
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, " ")

		if items[0] == "loadNumber" { // Skip the header line
			continue
		}

		id, _ := strconv.Atoi(items[0])
		// Process coordinates in the form of "(-42.51051149928979,-116.19788220835095) (-69.06284568487868,-44.12633704833111)"
		start := strings.Split(items[1][1 : len(items[1])-1], ",") // Remove the parentheses and split by comma
		startX, _ := strconv.ParseFloat(start[0], 64)
		startY, _ := strconv.ParseFloat(start[1], 64)
		startPoint := Point{x: startX, y: startY}

		end := strings.Split(items[2][1 : len(items[2])-1], ",")
		endX, _ := strconv.ParseFloat(end[0], 64)
		endY, _ := strconv.ParseFloat(end[1], 64)
		endPoint := Point{x: endX, y: endY}

		loadDistance := pointDistance(startPoint, endPoint)
		depotToStart := pointDistance(depot, startPoint)
		depotToEnd := pointDistance(depot, endPoint)

		load := &Load{
			id: id,
			start: startPoint,
			end: endPoint,
			loadDistance: loadDistance,
			depotToStart: depotToStart,
			depotToEnd: depotToEnd,
		}
		loadMap[id] = load
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return
	}
}

func generateSavingsList() {
	for i := 1; i <= len(loadMap); i++ {
		for j := 1; j <= len(loadMap); j++ {
			if i != j {

				saving := Saving{
					load1Id: i,
					load2Id: j,
					savedDistance: loadMap[i].depotToEnd + loadMap[j].depotToStart - pointDistance(loadMap[i].end, loadMap[j].start),
				}
				savings = append(savings, saving)
			}
		}
	}

	// Sort the savings list in descending order by savedDistance
	sort.Slice(savings, func(i, j int) bool {
		return savings[i].savedDistance > savings[j].savedDistance
	})
}

func processSavingsList() {
	// Clark and Wright algorithm
	for _, saving := range savings {
		load1 := loadMap[saving.load1Id]
		load2 := loadMap[saving.load2Id]

		if load1.truck == nil && load2.truck == nil  { // a. If both loads are unassigned to a truck

			roundTripTime := load1.depotToStart + load1.loadDistance + pointDistance(load1.end, load2.start) + load2.loadDistance + load2.depotToEnd

			if roundTripTime <= MAX_TIME {
				truck := &Truck{
					time: roundTripTime,
					loads: []*Load{load1, load2},
				}
				trucks = append(trucks, truck)
				load1.truck = truck
				load2.truck = truck
			}
		} else if load1.truck != nil && load2.truck == nil { // b. If load1 is assigned and load2 is unassigned

			load1Index := findLoadIndex(load1.truck.loads, load1)

			// if load1 is last load, attempt to append load2 if maximum time is not reached
			if load1Index == (len(load1.truck.loads) - 1) {
				roundTripTime :=
				load1.truck.time -
				load1.depotToEnd +
				pointDistance(load1.end, load2.start) +
				load2.loadDistance +
				load2.depotToEnd

				if roundTripTime <= MAX_TIME {
					load1.truck.time = roundTripTime
					load1.truck.loads = append(load1.truck.loads, load2)
					load2.truck = load1.truck
				}
			}


		} else if load1.truck == nil && load2.truck != nil { // b. Opposite case, if load1 is unassigned and load2 is assigned

			load2Index := findLoadIndex(load2.truck.loads, load2)

			// if load2 is first load, attempt to prepend load1 if maximum time is not reached
			if load2Index == 0 {
				roundTripTime :=
				load2.truck.time -
				load2.depotToStart +
				pointDistance(load1.end, load2.start) +
				load1.loadDistance +
				load1.depotToStart

				if roundTripTime <= MAX_TIME {
					load2.truck.time = roundTripTime
					load2.truck.loads = append([]*Load{load1}, load2.truck.loads...)
					load1.truck = load2.truck
				}
			}

		} else { // c. both load1 and load2 are assigned, so both routes are merged via a combination of the methods in step b.

			if load1.truck != load2.truck {

				load1Index := findLoadIndex(load1.truck.loads, load1)
				load2Index := findLoadIndex(load2.truck.loads, load2)

				// Ensure that load1 and load2 are not interior load nodes
				if load1Index == (len(load1.truck.loads) - 1) && load2Index == 0 { // load1 is last, and load2 is first
					// Check if linking load1 and load2 together would exceed the maximum time
					roundTripTime :=
					load1.truck.time - // current total running time for all loads on truck1
					load1.depotToEnd + // subtract return leg home of truck1
					pointDistance(load1.end, load2.start) + // add gap from load1 end to load2 start
					load2.truck.time - // current total running time for all loads on truck2
					load2.depotToStart // subtract initial leg of truck2

					if roundTripTime <= MAX_TIME {
						load1.truck.time = roundTripTime
						for _, load := range load2.truck.loads {
							load.truck = load1.truck
						}
					}

				}
			}

		}
	}

	// Assign remaining unassigned loads to new trucks
	for _, load := range loadMap {
		if load.truck == nil {
			truck := &Truck{
				time: load.depotToStart + load.loadDistance + load.depotToEnd,
				loads: []*Load{load},
			}
			load.truck = truck
			trucks = append(trucks, truck)
		}
	}

	// Print the results
	for _, truck := range trucks {
		fmt.Print("[")
		for i, load := range truck.loads {
			if i == len(truck.loads) - 1 {
				fmt.Printf("%d", load.id)
			} else {
				fmt.Printf("%d,", load.id)
			}
		}
		fmt.Print("]")
		fmt.Println()
	}
}
