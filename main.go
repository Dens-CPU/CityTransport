package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// СТРУКТУРЫ

// Структра Passanger
type Passanger struct {
	FinalStop  string
	WatingTime time.Duration
}

func NewPassanger(finalStop string) Passanger {
	return Passanger{FinalStop: finalStop, WatingTime: time.Duration(3) * time.Second}
}

// Структура BusStop
type BusStop struct {
	Name           string
	Wating         []Passanger
	QueueForArival chan Passanger
}

func NewBusStop(name string) BusStop {
	return BusStop{Name: name, Wating: make([]Passanger, rand.Intn(20)), QueueForArival: make(chan Passanger, rand.Intn(20))}
}

//Структрура Transport

type Transport struct {
	Name     string
	Route    []BusStop
	Capasity chan Passanger
	OutQueue chan Passanger
}

func NewTrasport(name string, route []BusStop) Transport {
	return Transport{Name: name, Route: route, Capasity: make(chan Passanger, 10), OutQueue: make(chan Passanger, 10)}
}

// Создание пассажира
func CreatPassanger(cityBusStops []BusStop) Passanger {
	n := rand.Intn(NumberOfBusStops)
	finalStop := cityBusStops[n].Name
	// fmt.Printf("Пункт назнанчения %s\n", finalStop)
	return NewPassanger(finalStop)
}

// Создание остановки
func CreatBusStop() BusStop {
	var name string
	fmt.Print("Название остановки: ")
	fmt.Scanln(&name)
	return NewBusStop(name)
}

// Создание транспорта
func CreatTransport(cityBusStops []BusStop) Transport {
	var name string
	fmt.Print("Название транспорта:")
	fmt.Scanln(&name)
	MAP := make(map[string]bool)
	//Количество оствновок на маршруте
	n := rand.Intn(NumberOfBusStops) + 2
	route := []BusStop{}
	for i := 0; i < n; i++ {
		j := rand.Intn(4) + 1
		stop := cityBusStops[j].Name
		if !MAP[stop] {
			route = append(route, cityBusStops[j])
			MAP[stop] = true
		}
	}
	return NewTrasport(name, route)
}

var mutex = sync.Mutex{}

// Формирование очереди на остановке
func CreatQueue(cityBusStop []BusStop, wg *sync.WaitGroup) {
	for i := range cityBusStop {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			stop := cityBusStop[i]
			mutex.Lock()
			for j := range stop.Wating {
				stop.Wating[j] = CreatPassanger(cityBusStop)
				// time.Sleep(time.Second)
			}
			mutex.Unlock()
		}(i)
	}
}

// Движение транспорта
func TransportWorking(CityTransport []Transport, wg *sync.WaitGroup) {
	// done := make(chan int)
	for i := range CityTransport {
		//Работа транспорта
		wg.Add(1)
		go func(i int) {
			numberOfStop := 0
			defer wg.Done()
			transport := CityTransport[i]
			for _, stop := range transport.Route {
				time.Sleep(time.Second)
				fmt.Printf("Транспорт %s прибыл на остановку %s\n", transport.Name, stop.Name)
				numberOfStop++
				
				//Высадка из транспорта
				if len(transport.Capasity) != 0 {
					c := len(transport.Capasity)
					for c > 0 {
						passange := <-transport.Capasity
						if passange.FinalStop == stop.Name {
							transport.OutQueue <- passange
							<-transport.OutQueue

						} else {
							transport.Capasity <- passange
						}
						c--
					}
				}
				//Посадку в транспорт
			Boarding:
				for j := 0; j < len(stop.Wating); j++ { //Обход людей на остановке
					for h := numberOfStop; h < len(transport.Route); h++ { // Обход пути транспорта
						// fmt.Printf("h=%d\n", h)
						if len(transport.Capasity) < cap(transport.Capasity) {
							if stop.Wating[j].FinalStop == transport.Route[h].Name && stop.Wating[j].FinalStop != stop.Name {
								passanger := stop.Wating[j]
								transport.Capasity <- passanger
							}
						} else {
							fmt.Println("Транспорт переполнен")

							break Boarding
						}
					}
				}
				fmt.Printf("Людей в транспотре %s = %d\n", transport.Name, len(transport.Capasity))
				fmt.Printf("Вместимость транспорта = %d\n", cap(transport.Capasity))
			}
			fmt.Println(len(transport.Capasity))
			fmt.Println(cap(transport.Capasity))
		}(i)
	}
}

const NumberOfTransport int = 2
const NumberOfBusStops int = 5

func main() {
	var wg sync.WaitGroup
	CityBusStops := make([]BusStop, NumberOfBusStops)
	CityTransport := make([]Transport, NumberOfTransport)
	//Создание остановок
	for i := range CityBusStops {
		CityBusStops[i] = CreatBusStop()
	}
	//Создание Транспорта
	for i := range CityTransport {
		CityTransport[i] = CreatTransport(CityBusStops)
	}
	//Образование очереди на остановках
	CreatQueue(CityBusStops, &wg)
	for i := range CityTransport {
		for _, stop := range CityTransport[i].Route {
			fmt.Printf("Остановка %s на маршркте транспорта %s\n", stop.Name, CityTransport[i].Name)
		}
	}
	TransportWorking(CityTransport, &wg)

	wg.Wait()
	// for _, stop := range CityBusStops {
	// 	fmt.Printf(" Остановка %s\n", stop.Name)
	// 	time.Sleep(time.Second)
	// 	for _, queue := range stop.Wating {
	// 		// time.Sleep(time.Second)
	// 		fmt.Printf("Пункт назначения %s\n", queue.FinalStop)
	// 	}
	// }
	// fmt.Println(CityTransport)
}
