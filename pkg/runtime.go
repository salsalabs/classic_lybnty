package lybnty

import (
	"fmt"
	"log"
	"time"

	godig "github.com/salsalabs/godig/pkg"
)

//Runtime contains the data necessary to run the application.
type Runtime struct {
	LastStart     time.Time
	LastEnd       time.Time
	ThisStart     time.Time
	ThisEnd       time.Time
	Org           string
	API           *godig.API
	Filename      string
	DonorThisYear map[string]bool
	DataCache     []Data
}

//NewRuntime build the runtime environment for this app.
func NewRuntime(a *godig.API, org *string, lastStart *string, lastEnd *string, thisStart *string, thisEnd *string) (rt *Runtime, err error) {
	b := []*string{
		org,
		lastStart,
		lastEnd,
		thisStart,
		thisEnd,
	}
	e := false
	r := Runtime{}
	for i, v := range b {
		if v == nil || len(*v) == 0 {
			switch i {
			case 0:
				log.Println("--org is a required parameter")
				e = true
			case 1:
				log.Println("--last-year-start is a required parameter")
				e = true
			case 2:
				log.Println("--last-year-end is a required parameter")
				e = true
			case 3:
				log.Println("--this-year-start is a required parameter")
				e = true
			case 4:
				log.Println("--this-year-end is a required parameter")
				e = true
			}
		}
	}
	if !e {
		b := b[1:5]
		for i, v := range b {
			d1, err := time.Parse(godig.DateFormat, *v)
			switch i {
			case 0:
				if err != nil {
					log.Println("--last-year-start must be formatted as YYYY-MM-DD!")
					e = true
				} else {
					r.LastStart = d1
				}
			case 1:
				if err != nil {
					log.Println("--last-year-end must be formatted as YYYY-MM-DD!")
					e = true
				} else {
					r.LastEnd = d1
				}
			case 2:
				d1, err := time.Parse(godig.DateFormat, *v)
				if err != nil {
					log.Println("--this-year-start must be formatted as YYYY-MM-DD!")
					e = true
				} else {
					r.ThisStart = d1
				}
			case 3:
				if err != nil {
					log.Println("--this-year-end must be formatted as YYYY-MM-DD!")
					e = true
				} else {
					r.ThisEnd = d1
				}
			}
		}
	}
	rt = &r
	if e {
		err := fmt.Errorf("%s", "Unable to continue due to parameter errors")
		return rt, err
	}
	r.API = a
	r.Org = *org
	y := r.ThisStart.Format("2006")
	r.Filename = fmt.Sprintf("LYBNTY %v %v.csv", r.Org, y)
	r.DonorThisYear = make(map[string]bool)
	return rt, nil
}
