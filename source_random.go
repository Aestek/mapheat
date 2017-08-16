package main

import (
	"math/rand"
	"strconv"
	"time"
)

func randomSource() <-chan *Evt {
	evts := make(chan *Evt)

	go func() {
		time.Sleep(3 * time.Second)
		for i := 0; ; i++ {
			lngMin := -4.498975
			lngMax := 7.750420
			latMin := 42.567788
			latMax := 50.953201

			id := strconv.Itoa(i)

			func(id string) {
				time.AfterFunc(100*time.Second, func() {
					evts <- &Evt{
						Action: ActionRemove,
						ID:     id,
					}
				})
			}(id)

			evts <- &Evt{
				Action: ActionAdd,
				ID:     id,
				Lng:    rand.Float64()*(lngMax-lngMin) + lngMin,
				Lat:    rand.Float64()*(latMax-latMin) + latMin,
				Rad:    DefaultRad,
			}

			time.Sleep(1 * time.Millisecond)
		}
	}()

	return evts
}
