package main

import (
	"fmt"
	"time"

	"github.com/umahmood/haversine"
)

const SpeedInvalid = 100.0
const SpeedIdle = 10.0

const FareSinglePerKm = 74
const FareDoubleePerKm = 130
const FareDoubleHr = 5
const FareIdlePerHr = 1190
const FareFlag = 130
const FareMin = 347

func CalculateFareInfo(point1 Path, point2 Path) FareInfo {
	p1 := haversine.Coord{Lat: point1.lat, Lon: point1.lng}
	p2 := haversine.Coord{Lat: point2.lat, Lon: point2.lng}
	_, km := haversine.Distance(p2, p1)
	t2 := time.Unix(point2.timestamp, 0)
	dt := t2.Sub(time.Unix(point1.timestamp, 0))
	return FareInfo{distance: km, speed: km / dt.Hours(), time: t2, duration: int(dt.Seconds())}
}

type FareInfo struct {
	distance float64
	speed    float64
	time     time.Time
	duration int
}

func FilterPaths(paths *map[int]Path) []FareInfo {
	id_ride := (*paths)[len(*paths)-1].id_ride
	fmt.Printf("[id_ride %v] Total number of points: %v\n", id_ride, len(*paths))
	var filtered []FareInfo
	last_path := (*paths)[0]
	for i := 1; i < len(*paths); i++ {
		path := (*paths)[i]
		info := CalculateFareInfo(last_path, path)
		if info.speed <= SpeedInvalid {
			filtered = append(filtered, info)
			last_path = path
		}
	}
	fmt.Printf("[id_ride %v] Removed invalid points: %v\n", id_ride, len(*paths)-len(filtered)-1)
	return filtered
}

const TimeFormat = "_2 Jan 2006 15:04:05"

func EstimateFare(infos []FareInfo, id_ride int) int32 {
	start, end := infos[0].time, infos[len(infos)-1].time
	fmt.Printf("[id_ride %v] Ride start: %v end: %v duration: %v\n", id_ride, start.Format(TimeFormat), end.Format(TimeFormat), end.Sub(start))
	var normal_km, double_km float64
	var idle_s int
	for _, info := range infos {
		if info.speed > SpeedIdle {
			time := info.time
			if time.Hour() > FareDoubleHr {
				normal_km += info.distance
			} else if time.Hour() <= FareDoubleHr {
				double_km += info.distance
			}
		} else {
			idle_s += info.duration
		}
	}
	return ApplyPrice(normal_km, double_km, idle_s, id_ride)
}

func ApplyPrice(normal_km, double_km float64, idle_s, id_ride int) int32 {
	idle_h := (float64(idle_s) / 3600)
	fare := FareFlag + normal_km*FareSinglePerKm + double_km*FareDoubleePerKm + idle_h*FareIdlePerHr
	fmt.Printf("[id_ride %v] Calculated €%.2f - single: %.2fkm double: %.2fkm idle: %.2fh\n", id_ride, float32(fare)/100.0, normal_km, double_km, idle_h)
	if fare < FareMin {
		fare = FareMin
	}
	return int32(fare)
}
