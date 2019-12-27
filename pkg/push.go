package lybnty

import (
	"log"
)

//Push reads supporter and donation records and pushes them downstream.  Note that
//the records are maps of strings.  Downstream will convert them to useful data.
func (rt *Runtime) Push(d chan map[string]string) (err error) {
	t := rt.API.NewTable("supporter(supporter_KEY)donation")
	offset := int32(0)
	count := 500
	c := "donation.RESULT IN 0,-1"
	for count == 500 {
		m, err := t.LeftJoinMap(offset, count, c)
		if err != nil {
			return err
		}
		for _, r := range m {
			d <- r
		}
		count = len(m)
		log.Printf("Push: read %d from offset %d\n", count, offset)
		offset += int32(count)
	}
	close(d)
	return err
}
