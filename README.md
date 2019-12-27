# classic_lybnty
Last Year But Not This Year (LYBNTY) for Salsa Classic

# Summary
Classic LYBNTY is an application that reads supporters and donations, then returns
information about supporters who donated last year but did not donate this year.  The
results are stored in a CSV file for consumption later.

# Prerequisites
1. Set up a go directory in your home directory.  Use this as a guide.
```
$HOME
    + go
        + src
        + bin
        + package
```
1. Install Go.  There is some [nice installation documentation here](https://golang.org/doc/install) that shows how to do that.

# Installation
Type these commands from a console window.  If you are using Windows, then sorry, you'll have to fake it.
```
go get github.com/salsalabs/classic_lybnty
cd ~/go/github.com/salsalabs/classic_lybnty
go get ./...
go install
```
If you don't get errors, then the installation is done.  You'll have an executable file named `classic_lybnty` 
(or `classic_lybnty.exe`) in ~/go/bin.

# Logging into Classic
This application expects a YAML file with the credentials that the Classic API needs.  Here's a sample.
```yaml
host: salsa4.salsalabs.com
email: truly.bogus@your.org
password: your-salsa-api-password
```
The file is *required*, so go ahead and create it somewhere.

# Usage
This application was designed to accept a starting and ending date for both the "last year" and the "this year"
in "LYBNTY".  The defaults are
* Last year is the calendar year for last year
* This year is the calendar year for this year

Here's the usage that you get with `--help`:
```text
usage: classic_lybnty [<flags>]

A command-line app to create an LYBNTY for a Salsa Classic instance

Flags:
  --help                        Show context-sensitive help (also try --help-long and --help-man).
  --login=LOGIN                 YAML file with API token
  --org=ORG                     Organization name (for output file)
  --last-year-start="2018-01-01"  
                                Last year start date, default is 2018-01-01
  --last-year-end="2018-12-31"  Last year end date, default is 2018-12-31
  --this-year-start="2019-01-01"  
                                This year start date, default is 2019-01-01
  --this-year-end="2019-12-31"  This year end date, default is 2019-12-31
  --api-verbose                 Makes the app noisy by showing API calls and results
```
Where:

| Argument | Description |
| --- | --- |
|LOGIN|A yaml file with Classic API credentials. For example, `bobcat.yaml`.|
|ORG|An organization name.  Be sure to quote the name if it has spaces.  For example, `"Bob the Bobcat Refuge"`|

Don't use --api-verbose, K?  It's there for the developer, is very, very messy and, well, don't use it.

# Output

The application reads all donations, filters out the ones that match the criteria, then writes the results
to a CSV file.  The CSV file contains `ORG` and the year from `--this-year-start`.  For example, running this 
application in 2019 using `Bob the Bobcat Refuge`

```classic_lybnty --login bobcat.yaml --org "Bob the Bobcat Refuge"```

stores the matching records in 

`LYBNTY Bob the Bobcat Refuge 2019.csv`.

Likewise, running the same report for 2018

```classic_lybnty --login bobcat.yaml --org "Bob the Bobcat Refuge" --last-year-start 2018-05-01```

stores matching records in 

`LYBNTY Bob the Bobcat Refuge 2018.csv`.

The file contains a header line that shows the fields for each line.  For example,

```
supporter_KEY,First_Name,Last_Name,Email,Street,Street_2,City,State,Zip,Country,donation_KEY,Transaction_Date,Tracking_Code,Donation_Tracking_Code,Designation_Code,Result,Transaction_Type,amount
```

Each line after that records supporter and donation information for a single donation.

```
53770542,Mongo,Baker,mongo@baker.bizi,1213 Grunt and Smash,,Someplace Ugly,WV,20015,,BR-549,2016-11-11,Smash2016,,,,,350.00
```

#Questions?  Comments?
Use the [Issues](https://github.com/salsalabs/classic_lybnty/issues) link at the top of this page.  Don't bother the
nice folks at Salsalabs Support.  It's nesting time, and they tend to bite if you get too close.
