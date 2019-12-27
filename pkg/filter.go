package lybnty

import (
	"log"
)

//Filter accepts maps of strings, converts them to useful data, then
//pushes qualified records downstream.
func (rt Runtime) Filter(in chan map[string]string, out chan Data) {
	log.Println("Filter: begin")
	for m := range in {
		d := NewDonation(m)
		if d.TransactionDate.After(rt.ThisStart) && d.TransactionDate.Before(rt.ThisEnd) {
			log.Printf("Filter: %v is between this start %v and this end %v\n", d.TransactionDate, rt.ThisStart, rt.ThisEnd)
		} else {
			s := NewSupporter(m)
			data := Data{
				Supporter: s,
				Donation:  d,
			}
			out <- data
		}
	}
	close(out)
	log.Println("Filter: end")
}
