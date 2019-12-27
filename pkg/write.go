package lybnty

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

//Write accepts Data records, and writes them to a CSV file.
func (rt Runtime) Write(c chan Data) (err error) {
	f, err := os.Create(rt.Filename)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	s := fmt.Sprintf("%v\n%v", SupporterFields, DonationFields)
	a := strings.Split(s, "\n")
	//Headers
	w.Write(a)

	for d := range c {
		var r []string
		x := ""
		for _, k := range a {
			switch k {
			case "supporter_KEY":
				x = d.SupporterKey
			case "First_Name":
				x = d.FirstName
			case "Last_Name":
				x = d.LastName
			case "Email":
				x = d.Email
			case "Street":
				x = d.Street
			case "Street_2":
				x = d.Street2
			case "City":
				x = d.City
			case "State":
				x = d.State
			case "Zip":
				x = d.Zip
			case "Country":
				x = d.Country
			case "donation_KEY":
				x = d.DonationKey
			case "Transaction_Date":
				x = d.TransactionDate.Format("2006-02-01")
			case "Tracking_Code":
				x = d.TrackingCode
			case "Donation_Tracking_Code":
				x = d.DonationTrackingCode
			case "Designation_Code":
				x = d.DesignationCode
			case "Result":
				x = d.Result
			case "TransactionType":
				x = d.TransactionType
			case "amount":
				x = fmt.Sprintf("%.2e", d.Amount)
			}
			r = append(r, x)
		}
		w.Write(r)
	}
	w.Flush()
	err = f.Close()
	return err
}
