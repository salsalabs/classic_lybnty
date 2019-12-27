package lybnty

import (
	"log"
)

//Cache accepts maps of strings, converts them to useful data, then
//caches them.
func (rt *Runtime) Cache(in chan map[string]string, out chan bool) {
	log.Println("Cache: begin")
	for m := range in {
		d := NewDonation(m)
		s := NewSupporter(m)
		data := Data{
			Supporter: s,
			Donation:  d,
		}
		if data.TransactionDate.After(rt.ThisStart) && d.TransactionDate.Before(rt.ThisEnd) {
			rt.DonorThisYear[data.SupporterKey] = true
		}
		if data.TransactionDate.After(rt.LastStart) && d.TransactionDate.Before(rt.LastEnd) {
			rt.DataCache = append(rt.DataCache, data)
		}
	}
	out <- true
	close(out)
	log.Println("Cache: end")
}
