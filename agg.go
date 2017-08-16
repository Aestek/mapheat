package main

import (
	"fmt"
	"math"
	"time"
)

type AggElem struct {
	*Evt
	Sub map[string]*Evt
}

const growth = 0.01

func Agg(c <-chan *Evt) <-chan *Evt {
	out := make(chan *Evt)

	var elems []*AggElem
	var total int

	add := func(evt *Evt) {
		total++
		for _, e := range elems {
			dLat := math.Abs(e.Lat - evt.Lat)
			dLng := math.Abs(e.Lng - evt.Lng)

			if dLat > DefaultRad*2 || dLng > DefaultRad*2 {
				continue
			}
			e.Sub[evt.ID] = evt
			e.Rad = DefaultRad + growth*float64(len(e.Sub))
			e.Lat = (e.Lat*float64(len(e.Sub)) + e.Lat) / (float64(len(e.Sub)) + 1)
			e.Lng = (e.Lng*float64(len(e.Sub)) + e.Lng) / (float64(len(e.Sub)) + 1)

			out <- e.Evt
			return
		}

		evt.Rad = DefaultRad + growth
		elems = append(elems, &AggElem{
			Evt: evt,
			Sub: map[string]*Evt{evt.ID: evt},
		})
		out <- evt
	}

	remove := func(evt *Evt) {
		total--
		for i, e := range elems {
			if _, ok := e.Sub[evt.ID]; !ok {
				continue
			}

			delete(e.Sub, evt.ID)
			if len(e.Sub) == 0 {
				elems = append(elems[:i], elems[i+1:]...)

				out <- &Evt{
					ID:     e.ID,
					Action: ActionRemove,
				}
			} else {
				e.Rad = DefaultRad + growth*float64(len(e.Sub))
				out <- e.Evt
			}
			return
		}

		fmt.Println("NOT FOUND")
	}

	go func() {
		for evt := range c {
			if evt.Action == ActionAdd {
				add(evt)
			} else {
				remove(evt)
			}
		}
	}()

	go func() {
		for {
			fmt.Println(len(elems), total)
			time.Sleep(time.Second)
		}
	}()

	return out
}
