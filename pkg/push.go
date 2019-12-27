package lybnty

import (
	"fmt"
	"strings"

	godig "github.com/salsalabs/godig/pkg"
)

//Push reads supporter and donation records and pushes them downstream.  Note that
//the records are maps of strings.  Downstream will convert them to useful data.
func (rt Runtime) Push(d chan map[string]string) (err error) {
	t := rt.API.Donation()
	offset := int32(0)
	count := 500
	lastStart := rt.LastStart.Format(godig.DateFormat)
	lastStart = fmt.Sprintf("Transaction_Date >= %v", lastStart)
	lastEnd := rt.LastEnd.Format(godig.DateFormat)
	lastEnd = fmt.Sprintf("Transaction_Date <= %v", lastEnd)
	conditions := []string{"RESULT IN 0,-1&",
		lastStart,
		lastEnd,
	}

	for count == 500 {
		c := strings.Join(conditions, "&condition=")
		fmt.Printf("Conditions are '%v'\n", c)
		m, err := t.ManyMap(offset, count, c)
		if err != nil {
			return err
		}
		for _, r := range m {
			d <- r
		}
		count = len(m)
	}
	return err
}
