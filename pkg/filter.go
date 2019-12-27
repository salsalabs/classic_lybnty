package lybnty

import (
	"log"
)

//Filter iterates through the data cache.  Records that do not have a supporter
//key in the donors this year.
func (rt *Runtime) Filter(cached chan bool, out chan Data) {
	log.Println("Filter: begin")
	<-cached
	log.Println("Filter: start filtering to the output")
	log.Printf("Filter: %d cached items\n", len(rt.DataCache))
	log.Printf("Filter: %d donors this year\n", len(rt.DonorThisYear))
	for _, data := range rt.DataCache {
		_, ok := rt.DonorThisYear[data.SupporterKey]
		if !ok {
			out <- data
		}
	}
	close(out)
	log.Println("Filter: end")
}
