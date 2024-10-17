# Vehicle Routing Problem

This project focuses on solving the Single Depot Vehicle Routing Problem, a combinatorial optimization and integer programming problem. The objective is to determine the optimal set of routes for a fleet of vehicles to traverse in order to deliver to a given set of customers. The solution has been implemented following the Clark Wright Savings Algorithm.

The problem input contains a list of loads. Each load is formatted as an id followed by pickup and dropoff locations in (x,y) floating point coordinates. An example input with four loads is:

```text
loadNumber pickup dropoff
1 (-50.1,80.0) (90.1,12.2)
2 (-24.5,-19.2) (98.5,1.8)
3 (0.3,8.9) (40.9,55.0)
4 (5.3,-61.1) (77.8,-5.4)
```

An example solution to the above problem could be:

```text
[1]
[4,2]
[3]
```

This solution means one driver does load 1; another driver does load 4 followed by load 2; and a final driver does load 3.

## Solution Constraints

- Drivers cannot travel a Euclidean distance greater than `(12 * 60)`.
- Drivers must originate and finish from a central depot located at `(0, 0)`.

## Prerequisites

- [Golang](https://go.dev/doc/install)
- [Python3](https://www.python.org/downloads/)

## Usage

To run the full VRP testing suite:

```bash
python3 evaluateShared.py --cmd "go run main.go" --problemDir trainingProblems
```

To run an individual VRP test from the `/trainingProblems` folder:

```bash
go run main.go trainingProblems/problem1.txt
```

## References

- <https://en.wikipedia.org/wiki/Vehicle_routing_problem>
- <https://web.mit.edu/urban_or_book/www/book/chapter6/6.4.12.html>
- <https://folk.ntnu.no/skoge/prost/proceedings/acc11/data/papers/1317.pdf>
