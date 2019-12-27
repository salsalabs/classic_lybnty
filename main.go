package main

import (
	"fmt"
	"log"
	"os"
	"time"

	godig "github.com/salsalabs/godig/pkg"
	"gopkg.in/alecthomas/kingpin.v2"
)

//Application to create a Last Year But Not This Year  (LYBNTY) report for Salsa Classic.
//The user gets to define the start and end of both "Last Year" and "This Year", making this
//more of a general-purpose report than a typicaly LYBUNTY report.

//Runtime contains the data necessary to run the application.
type Runtime struct {
	LastStart time.Time
	LastEnd   time.Time
	ThisStart time.Time
	ThisEnd   time.Time
	Org       string
	API       *godig.API
}

//NewRuntime build the runtime environment for this app.
func NewRuntime(a *godig.API, org *string, lastStart *string, lastEnd *string, thisStart *string, thisEnd *string) (r Runtime, err error) {
	b := []*string{
		org,
		lastStart,
		lastEnd,
		thisStart,
		thisEnd,
	}
	e := false
	for i, v := range b {
		fmt.Printf("checking parameter %v\n", *v)
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
	if e {
		err := fmt.Errorf("%s", "Unable to continue due to parameter errors")
		return r, err
	}
	r.API = a
	r.Org = *org
	return r, nil
}

func main() {

	t := time.Now()
	y := t.Year()
	v1 := fmt.Sprintf("%4d-01-01", y-1)
	v2 := fmt.Sprintf("%4d-12-31", y-1)
	v3 := fmt.Sprintf("%4d-01-01", y)
	v4 := fmt.Sprintf("%4d-12-31", y)
	p1 := fmt.Sprintf("Last year start date, default is %s", v1)
	p2 := fmt.Sprintf("Last year end date, default is %s", v2)
	p3 := fmt.Sprintf("This year start date, default is %s", v3)
	p4 := fmt.Sprintf("This year end date, default is %s", v4)

	var (
		app       = kingpin.New("LYBNTY Report for Classic", "A command-line app to create an LYBNTY for a Salsa Classic instance")
		login     = app.Flag("login", "YAML file with API token").String()
		org       = app.Flag("org", "Organization name (for output file)").String()
		lastStart = app.Flag("last-year-start", p1).Default(v1).String()
		lastEnd   = app.Flag("last-year-end", p2).Default(v2).String()
		thisStart = app.Flag("this-year-start", p3).Default(v3).String()
		thisEnd   = app.Flag("this-year-end", p4).Default(v4).String()
	)
	app.Parse(os.Args[1:])

	fmt.Printf("login: %s\n", *login)
	a, err := godig.YAMLAuth(*login)
	if err != nil {
		log.Fatalf("Authentication error: %+v\n", err)
	}
	rt, err := NewRuntime(a, org, lastStart, lastEnd, thisStart, thisEnd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Runtime: %+v\n", rt)
}
