// package main

// import (
// 	"fmt"
// 	"parkinglot/models"
// 	"parkinglot/services"
// 	"parkinglot/strategies"
// 	"time"
// )

// func main() {
// 	// 1. Build the lot: Floor(number, bikes, cars, trucks)
// 	floors := []*models.Floor{
// 		models.NewFloor(1, 2, 2, 2),
// 		models.NewFloor(2, 2, 1, 2),
// 	}

// 	// 2. Pick strategies (swap these to change behaviour, nothing else changes)
// 	lot := services.NewParkingLotService(
// 		"AMB Mall", floors, strategies.NearestSpotStrategy{},
// 		strategies.NewHourlyStrategy())

// 	lot.ShowAvailability()
// 	fmt.Println()

// 	// 3. Vehicles arrive at 10:00
// 	start := time.Date(2026, 11, 11, 10, 0, 0, 0, time.Local)

// 	_ = park(lot, models.NewVehicle("KA-01-1111", models.Bike), start)
// 	carTicket := park(lot, models.NewVehicle("KA-02-2222", models.Car), start)
// 	_ = park(lot, models.NewVehicle("KA-03-3333", models.Car), start)
// 	_ = park(lot, models.NewVehicle("KA-04-4444", models.Car), start)
// 	_ = park(lot, models.NewVehicle("KA-05-5555", models.Truck), start)

// 	// Lot has only 3 car spots, so this one fails
// 	park(lot, models.NewVehicle("KA-06-6666", models.Car), start)
// 	// Same car trying to park twice
// 	park(lot, models.NewVehicle("KA-02-2222", models.Car), start)

// 	fmt.Println()
// 	lot.ShowAvailability()
// 	fmt.Println()

// 	// 4. Car leaves after 2h30m -> rounds up to 3h -> 3 * 20 = 60
// 	unpark(lot, carTicket.ID, start.Add(2*time.Hour+30*time.Minute))
// 	unpark(lot, "T99", start.Add(3*time.Hour)) // invalid ticket

// 	// 5. The waiting car can now get the freed spot
// 	park(lot, models.NewVehicle("KA-06-6666", models.Car), start.Add(165*time.Minute))

// }

// func park(lot *services.ParkingLotService, v *models.Vehicle, at time.Time) *models.Ticket {
// 	ticket, err := lot.Park(v, at)
// 	if err != nil {
// 		fmt.Printf("[FAIL] Park %s %s: %v\n", v.Type, v.Number, err)
// 		return nil
// 	}
// 	fmt.Printf("[OK]   Parked %s %s at %s | Ticket %s\n", v.Type, v.Number, ticket.Spot.ID, ticket.ID)
// 	return ticket
// }

// func unpark(lot *services.ParkingLotService, ticketID string, at time.Time) {
// 	ticket, err := lot.Unpark(ticketID, at)
// 	if err != nil {
// 		fmt.Printf("[FAIL] Unpark %s: %v\n", ticketID, err)
// 		return
// 	}
// 	fmt.Printf("[OK]   Unparked %s %s from %s | Duration %v | Fee Rs.%.0f\n",
// 		ticket.Vehicle.Type, ticket.Vehicle.Number, ticket.Spot.ID, ticket.Duration(), ticket.Fee)
// }

package main

import (
	"fmt"
	"sync"
	"time"

	"parkinglot/models"
	"parkinglot/services"
	"parkinglot/strategies"
)

func main() {
	// 1. Build the lot: 11 spots in total (4 bikes, 3 cars, 4 trucks)
	floors := []*models.Floor{
		models.NewFloor(1, 2, 2, 2),
		models.NewFloor(2, 2, 1, 2),
	}
	lot := services.NewParkingLotService("AMB Mall", floors,
		strategies.NearestSpotStrategy{}, strategies.NewHourlyStrategy())

	// 2. Create gates. All of them share the SAME lot.
	entryGates := []*services.EntryGate{
		services.NewEntryGate("E1", lot),
		services.NewEntryGate("E2", lot),
		services.NewEntryGate("E3", lot),
	}
	exitGates := []*services.ExitGate{
		services.NewExitGate("X1", lot),
		services.NewExitGate("X2", lot),
	}

	lot.ShowAvailability()
	start := time.Date(2026, 11, 11, 10, 0, 0, 0, time.Local)

	// 3. 15 vehicles (5 bikes, 5 cars, 5 trucks) but only 11 spots -> 4 must fail
	types := []models.VehicleType{models.Bike, models.Car, models.Truck}
	queues := make([][]*models.Vehicle, len(entryGates)) // one queue per entry gate
	for i := 1; i <= 15; i++ {
		v := models.NewVehicle(fmt.Sprintf("KA-01-%04d", i), types[(i-1)/5])
		gateIdx := (i - 1) % len(entryGates) // round robin: E1, E2, E3, E1, ...
		queues[gateIdx] = append(queues[gateIdx], v)
	}

	// ---------------- Phase 1: all entry gates work at the same time ----------------
	fmt.Println("\n=== Phase 1: 15 vehicles enter through 3 gates at once ===")
	var wg sync.WaitGroup
	ticketsCh := make(chan *models.Ticket, 15) // goroutines send their tickets here

	for i, gate := range entryGates {
		wg.Add(1) // ALWAYS Add before starting the goroutine
		go func(g *services.EntryGate, queue []*models.Vehicle) {
			defer wg.Done() // tell the WaitGroup "this gate is finished"
			for _, v := range queue {
				time.Sleep(10 * time.Millisecond) // a real barrier takes time to open
				t, err := g.Enter(v, start)
				if err != nil {
					fmt.Printf("[%s] FAIL  %-5s %s: %v\n", g.ID, v.Type, v.Number, err)
					continue
				}
				fmt.Printf("[%s] OK    %-5s %s -> %s (%s)\n", g.ID, v.Type, v.Number, t.Spot.ID, t.ID)
				ticketsCh <- t
			}
		}(gate, queues[i])
	}
	wg.Wait()        // block until all 3 gates call Done()
	close(ticketsCh) // no more senders, so it's safe to close

	// Verify: no spot was given to two vehicles
	var tickets []*models.Ticket
	usedSpots := make(map[string]bool)
	// Listening from the channel
	for t := range ticketsCh {
		if usedSpots[t.Spot.ID] {
			fmt.Println("BUG! spot given twice:", t.Spot.ID)
		}
		usedSpots[t.Spot.ID] = true
		tickets = append(tickets, t)
	}
	fmt.Printf("\nParked %d of 15 vehicles | unique spots used: %d\n", len(tickets), len(usedSpots))
	lot.ShowAvailability()

	// ---------------- Phase 2: same car shows up at two gates at the same moment ----------------
	fmt.Println("\n=== Phase 2: same car at E1 and E2 at the same time ===")
	// Free one car spot first so the car CAN park (otherwise both just fail with "no spot")
	for _, t := range tickets {
		if t.Vehicle.Type == models.Car {
			exitGates[0].Exit(t.ID, start.Add(time.Hour))
			fmt.Printf("(freed %s by letting %s leave)\n", t.Spot.ID, t.Vehicle.Number)
			tickets = remove(tickets, t.ID)
			break
		}
	}
	sameCar := models.NewVehicle("KA-99-9999", models.Car)
	for _, gate := range entryGates[:2] {
		wg.Add(1)
		go func(g *services.EntryGate) {
			defer wg.Done()
			t, err := g.Enter(sameCar, start.Add(time.Hour))
			if err != nil {
				fmt.Printf("[%s] FAIL  %s: %v\n", g.ID, sameCar.Number, err)
				return
			}
			fmt.Printf("[%s] OK    %s -> %s (%s)\n", g.ID, sameCar.Number, t.Spot.ID, t.ID)
			tickets = append(tickets, t) // safe: only the ONE winner reaches this line
		}(gate)
	}
	wg.Wait()

	// ---------------- Phase 3: everyone leaves through 2 exit gates (worker pool) ----------------
	fmt.Println("\n=== Phase 3: everyone leaves through 2 exit gates ===")
	exitQueue := make(chan *models.Ticket, len(tickets))
	for _, t := range tickets {
		exitQueue <- t
	}
	close(exitQueue) // exit gates stop when the queue is empty

	collected := make([]float64, len(exitGates)) // each gate writes ONLY its own index -> no race
	for i, gate := range exitGates {
		wg.Add(1)
		go func(idx int, g *services.ExitGate) {
			defer wg.Done()
			for t := range exitQueue { // whichever gate is free takes the next ticket
				time.Sleep(10 * time.Millisecond) // driver pays, barrier opens
				done, err := g.Exit(t.ID, start.Add(2*time.Hour))
				if err != nil {
					fmt.Printf("[%s] FAIL  %s: %v\n", g.ID, t.ID, err)
					continue
				}
				fmt.Printf("[%s] EXIT  %-5s %s from %s | Fee Rs.%.0f\n",
					g.ID, done.Vehicle.Type, done.Vehicle.Number, done.Spot.ID, done.Fee)
				collected[idx] += done.Fee
			}
		}(i, gate)
	}
	wg.Wait()

	total := 0.0
	for i, g := range exitGates {
		fmt.Printf("Gate %s collected Rs.%.0f\n", g.ID, collected[i])
		total += collected[i]
	}
	fmt.Printf("Total collected: Rs.%.0f\n\n", total)
	lot.ShowAvailability()
}

func remove(tickets []*models.Ticket, id string) []*models.Ticket {
	for i, t := range tickets {
		if t.ID == id {
			return append(tickets[:i], tickets[i+1:]...)
		}
	}
	return tickets
}

// Important things -
// Entry is not a worker pool, it is one go routine per gate with each tnrey gate with it's own fixed queue. Not a pool.
// Exit is a worker pool with 2 number of workers with one shared job queue(channel).
// Also go routines doesn't block, if you write a go routine and inside that go routine you have written soem for loop or some big/lot of code
// and then outside of that go routine if you write another go routine, then we don't wait for second go routine to start until first go routine finishes everything.
// We first start first go routine and then right away we start second go routine and we wait in wg.Wait() method
// for both go routines to process their logic and return.
// In this case as well in entry case we have written am entry gate for loop with size 3 and in each iteration we
// sponned a go routine, so we don't wait for the first go routine to complete everything before sponning the second goroutine,
// We spon the first go routine and it continue to process without blocking.. so again the next for loop iteratin starts right away
// and it starts second go routine and it continues to process without blocking.. so again the next for loop iteration starts right away
// and it starts third go routine. We wait at wg.Wait() for all go routines to finish and return.
