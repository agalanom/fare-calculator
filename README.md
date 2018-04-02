# Fare Calculator

In this exercise you are asked to create a **Fare Estimation script**. Its input will consist of a list of tuples of  the form `(id_ride, lat, lng, timestamp)` representing the position of the taxi-cab during a ride.

Two consecutive tuples belonging to the same ride form a segment **S**. For each segment, we define the elapsed time **Δt** as the absolute difference of the segment endpoint timestamps and the distance covered **Δs** as the haversine distance of the segment endpoint coordinates.

Your first task is to filter the data on the input so that invalid entries are discarded before the fare estimation process begins. Filtering should be performed as follows: consecutive tuples **p1**, **p2** should be used to calculate the segment’s speed **U**. If **U > 100km/h**, **p2** should be removed from the set.

Once the data has been filtered you should proceed in estimating the fare by **aggregating the individual estimates** for the ride segments using rules tabulated below: 

State | Applicable when | Fare amount
---|---|---
|MOVING (U > 10km/h) | Time (05:00, 00:00] | 0.74 per km
||Time (00:00, 05:00] | 1.30 per km
|IDLE (U <= 10km/h) | Always | 11.90 per hour of idle time

At the start of each ride, a standard ‘flag’ amount of 1.30 is charged to the ride’s fare. The minimum ride fare should be **at least** 3.47.

## Input Data
The sample data file contains one record per line (comma separated values). The input data is guaranteed to contain **continuous** row blocks for each individual ride (i.e. the data will not be multiplexed). In addition, the data is also pre-sorted for you in ascending timestamp order.

## Deliverable
A Golang 1.8+ program that processes input as provided by the sample data, filters out invalid points and produces an output comma separated value text file with the following format:

```id_ride, fare_estimate```

## Instructions

Build binary

```bash
go build
```

Run

```bash
./fare-calculator paths.csv
```

Test

```bash
go test
```
