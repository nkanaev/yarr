package expirer

import (
	"log"
	"time"

	"github.com/nkanaev/yarr/src/storage"
)

type Expirer struct {
	db              *storage.Storage
	stop            chan bool
	intervalMinutes uint64
}

func NewExpirer(db *storage.Storage) *Expirer {
	return &Expirer{
		db:              db,
		stop:            make(chan bool),
		intervalMinutes: 0,
	}
}

func expire(db *storage.Storage, rateMinutes uint64, stop chan bool) {
	tick := time.NewTicker(time.Minute * time.Duration(rateMinutes) / 2)
	db.ExpireUnreads(rateMinutes)
	log.Printf("expirer %dm: starting", rateMinutes)
	for {
		select {
		case <-tick.C:
			log.Printf("expirer %dm: firing", rateMinutes)
			db.ExpireUnreads(rateMinutes)
		case <-stop:
			log.Printf("expirer %dm: stopping", rateMinutes)
			tick.Stop()
			stop <- true
			return
		}
	}
}

func (e *Expirer) getCheckInterval(globalExpirationPeriod uint64) uint64 {
	minFeedExpirationPeriod, err := e.db.GetMinExpirationPeriod()
	if err != nil {
		log.Fatal(err)
	}

	checkInterval := globalExpirationPeriod
	if *minFeedExpirationPeriod != 0 && *minFeedExpirationPeriod < checkInterval {
		checkInterval = *minFeedExpirationPeriod
	}
	return checkInterval
}

func (e *Expirer) StartUnreadsExpirer(globalExpirationPeriod uint64) {
	checkInterval := e.getCheckInterval(globalExpirationPeriod)
	if checkInterval > 0 {
		e.intervalMinutes = uint64(checkInterval)
		go expire(e.db, e.intervalMinutes, e.stop)
	}
}

func (e *Expirer) SetExpirationRate(globalExpirationPeriod uint64) {
	checkInterval := e.getCheckInterval(globalExpirationPeriod)
	if checkInterval == e.intervalMinutes {
		return
	}
	e.stop <- true
	<-e.stop
	e.intervalMinutes = globalExpirationPeriod
	if checkInterval == 0 {
		return
	}
	go expire(e.db, e.intervalMinutes, e.stop)
}
