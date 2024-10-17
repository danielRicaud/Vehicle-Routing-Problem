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

const MAX_TIME = 12*60
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

func findTruckIndex(trucks []*Truck, truck *Truck) int {
	for i, current := range trucks {
			if current == truck {
					return i
			}
	}
	return -1 // Not found
}

func deleteTruck(trucks []*Truck, truck *Truck) {
	for i, current := range trucks {
			if truck == current {
					// Slice manipulation to remove the element
					trucks[i] = nil
					trucks = append((trucks)[:i], (trucks)[i+1:]...)
					return
			}
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
		// fmt.Println(load)
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
					// Distance(depot, load1end) + Distance(depot, load2start) - Distance(load1end, load2start)
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

	// fmt.Println("Savings list:")
	// for _, saving := range savings {
	// 	fmt.Println(saving)
	// }
}

func processSavingsList() {
	// Clark and Wright savings algorithm
	// trucks := make([]*Truck, 0)
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
					// Check if adding load1 and load2 would exceed the maximum distance
					roundTripTime :=
					load1.truck.time - // current total running time for all loads on truck1
					load1.depotToEnd + // subtract return leg home of truck1
					pointDistance(load1.end, load2.start) + // add gap from load1 end to load2 start
					load2.truck.time - // current total running time for all loads on truck2
					load2.depotToStart // subtract initial leg of truck2

					// load2.loadDistance - // add load2 distance
					// load2.depotToStart + // subtract initial leg of truck2
					// load2.truck.time // add current total running time for all loads on truck2

					if roundTripTime <= MAX_TIME {
						// fmt.Printf("Merging truck1 with load: %v and truck2 with load: %v\n", load1.truck.loads[load1Index].id, load2.truck.loads[load2Index].id)
						load1.truck.time = roundTripTime
						for _, load := range load2.truck.loads {
							load.truck = load1.truck
						}
						// load1.truck.loads = append(load1.truck.loads, load2.truck.loads...)



						// Delete truck2 from truck list
						// trucks = append(trucks[:findTruckIndex(trucks, load2.truck)], trucks[findTruckIndex(trucks, load2.truck)+1:]...)



						// load2.truck = load1.truck
						// TODO: remove empty truck2 if it messes up printing
						// fmt.Println("Deleting a truck	from the list")
						// fmt.Println("Original truck list: %v", trucks)
						// deleteTruck(trucks, load2.truck)
						// fmt.Println("New truck list: %v", trucks)
					}

				}
			}

		}
	}

	// Assign remaining unassigned loads to new trucks
	for _, load := range loadMap {
		if load.truck == nil {
			// fmt.Printf("FOUND ONE")
			truck := &Truck{
				time: load.depotToStart + load.loadDistance + load.depotToEnd,
				loads: []*Load{load},
			}
			load.truck = truck
			trucks = append(trucks, truck)
		}
	}

	// fmt.Println("Trucks and their loads:")
	// fmt.Println("Max time: ", MAX_TIME)
	for _, truck := range trucks {
		// if truck.loads == nil {
		// 	continue
		// }
		// fmt.Printf("Truck with total time %.2f has loads: ", truck.time)
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
