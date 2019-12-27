package lybnty

import (
	"log"
	"strconv"
	"strings"
	"time"

	godig "github.com/salsalabs/godig/pkg"
)

//Supporter is the record of supporter information that we're manipulating in this
//application.  Note that typical LYBNTY reports need contact information for
//supporters and not a lot else.  This record follows along in that vein by constraining
//the data being retrieved from the database.
type Supporter struct {
	SupporterKey string
	FirstName    string
	LastName     string
	Email        string
	Street       string
	Street2      string
	City         string
	State        string
	Zip          string
	Country      string
}

//SupporterFields lists the fields that we're using.  Useful for retrieving data
//and putting the data into a CSV.
const SupporterFields = `supporter_KEY
First_Name
Last_Name
Email
Street
Street_2
City
State
Zip
Country`

//NewSupporter creates a supporter record and populates it from a map of strings.  The map
//keys are the field names, the map values are the field values.
func NewSupporter(m map[string]string) (s Supporter) {
	a := strings.Split(SupporterFields, "\n")
	for _, k := range a {
		v, ok := m[k]
		if ok {
			v := strings.TrimSpace(v)
			if len(v) > 0 {
				switch k {
				case "supporter_KEY":
					s.SupporterKey = v
				case "First_Name":
					s.FirstName = v
				case "Last_Name":
					s.LastName = v
				case "Email":
					s.Email = v
				case "Street":
					s.Street = v
				case "Street_2":
					s.Street2 = v
				case "City":
					s.City = v
				case "State":
					s.State = v
				case "Zip":
					s.Zip = v
				case "Country":
					s.Country = v
				}
			}
		}
	}
	return s
}

//Donation is a record of donation information that we're manipulating in this application.
//There's a lot of esoteric stuff in a Classic donation record. We'll use the typically useful
//stuff and disregard the rest.
type Donation struct {
	DonationKey          string
	TransactionDate      time.Time
	TrackingCode         string
	DonationTrackingCode string
	DesignationCode      string
	Result               string
	TransactionType      string
	Amount               float64
}

//DonationFields lists the fields that we're using.  Useful for retrieving data
//and putting the data into a CSV.
const DonationFields = `donation_KEY
Transaction_Date
Tracking_Code
Donation_Tracking_Code
Designation_Code
Result
Transaction_Type
amount`

//NewDonation creates a donation record and populates it from a map of strings.  The map
//keys are the field names, the map values are the field values.
func NewDonation(m map[string]string) (d Donation) {
	a := strings.Split(DonationFields, "\n")
	for _, k := range a {
		v, ok := m[k]
		if ok {
			v := strings.TrimSpace(v)
			if len(v) > 0 {
				switch k {
				case "donation_KEY":
					d.DonationKey = v
				case "Transaction_Date":
					w, err := godig.ClassicTime(v)
					if err != nil {
						log.Fatal(err)
					}
					d.TransactionDate = w
				case "Tracking_Code":
					d.TrackingCode = v
				case "Donation_Tracking_Code":
					d.DonationTrackingCode = v
				case "Designation_Code":
					d.DesignationCode = v
				case "Result":
					d.Result = v
				case "TransactionType":
					d.TransactionType = v
				case "amount":
					x, err := strconv.ParseFloat(v, 64)
					if err != nil {
						log.Fatal(err)
					}
					d.Amount = x
				}
			}
		}
	}
	return d
}

//Data conntains a supporter and a donation and forms a single unit of work.
type Data struct {
	Supporter
	Donation
}
