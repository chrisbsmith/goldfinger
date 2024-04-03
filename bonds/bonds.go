package bonds

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/itzg/restify"
	"golang.org/x/net/html"

	"github.com/chrisbsmith/goldfinger/config"
)

var (
	// ErrInvalidDataReturned is an error returned when invalid is returned from the server
	ErrInvalidDataReturned = errors.New("invalid data returned from server")
)

// Bond
type Bond struct {
	// Denomination of the Bond
	Denomination string

	// Serial number for the bond
	Serial string

	// Issue Date of the bond
	IssueDate string

	// Series of the bond
	Series string

	//NextAccrual Date of the bond
	NextAccrual string

	// FinalMaturity of the bond
	FinalMaturity string

	// IssuePrice of the bond
	IssuePrice string

	// Interest of the bond
	Interest string

	//InterestRate of the bond
	InterestRate string

	//Value of the bond
	Value string

	// Note of the bond
	Note string
}

type Bonds struct {
	Bonds []Bond
}

// LoadBonds loads the bond information from the config file
func LoadBonds(config *config.Config) (Bonds, error) {
	var bonds Bonds
	log.Printf("Retrieving bond values\n")

	for _, b := range config.Bonds {
		v, err := GetBond(b)
		if err != nil {
			log.Printf("Error getting bond data: %s\n", err.Error())
			return bonds, err
		}
		bonds.Bonds = append(bonds.Bonds, *v)
	}

	return bonds, nil
}

// GetBond queries the Treasury's server to get the bond info
func GetBond(bond config.ConfigBond) (*Bond, error) {
	now := time.Now()

	// API Endpoints
	apiUrl := "https://treasurydirect.gov"
	resource := "/BC/SBCPrice"

	data := url.Values{}
	data.Set("Denomination", strconv.Itoa(bond.Denomination))
	data.Set("SerialNumber", bond.Serial)
	data.Set("IssueDate", bond.IssueDate)
	data.Set("RedemptionDate", now.Format("01/2006"))
	data.Set("Series", bond.Series)
	data.Set("btnAdd.x", "Calculate")
	data.Set("Version", "6")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(r)

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}

	b, err := parseBondData(doc)
	if err != nil {
		return &b, err
	}
	return &b, nil
}

// parseBondData parses the returned HTML to retrieve the bond data
func parseBondData(n *html.Node) (Bond, error) {
	var SavingsBond Bond

	bonddata := restify.FindSubsetByClass(n, "bnddata")

	data := []string{}
	var val string

	for _, nodes := range bonddata {
		n := restify.FindSubsetByClass(nodes, "altrow1")
		for _, ns := range n {
			values := restify.FindSubsetByTagName(ns, "td")
			for _, value := range values {
				if value.FirstChild.Data == "strong" {
					v := restify.FindSubsetByTagName(value, "strong")
					for _, x := range v {
						val = x.FirstChild.Data
					}
				} else if value.FirstChild.Data == "a" {
					v := restify.FindSubsetByTagName(value, "a")
					for _, x := range v {
						val = x.FirstChild.Data
					}
				} else if value.FirstChild.Data == "input" {
					// ignore these cases where it's actuallying finding the tag <input
					continue
				} else {
					val = value.FirstChild.Data
				}

				data = append(data, val)
			}
		}
	}

	// A valid retrieval of data will return 11 rows of data, so if we don't
	// end up with 11, then something is wrong
	if len(data) != 11 {
		return SavingsBond, ErrInvalidDataReturned
	}

	for count := range data {

		// Data is returned as a table but we only get the row data. We know the
		// order of the table columns, so we can use this ugly switch statement
		// to save the data into the proper fields of the struct
		switch count {
		case 0:
			SavingsBond.Serial = strings.TrimSpace(data[count])
		case 1:
			SavingsBond.Series = strings.TrimSpace(data[count])
		case 2:
			SavingsBond.Denomination = strings.TrimSpace(data[count])
		case 3:
			SavingsBond.IssueDate = strings.TrimSpace(data[count])
		case 4:
			SavingsBond.NextAccrual = strings.TrimSpace(data[count])
		case 5:
			SavingsBond.FinalMaturity = strings.TrimSpace(data[count])
		case 6:
			SavingsBond.IssuePrice = strings.TrimSpace(data[count])
		case 7:
			SavingsBond.Interest = strings.TrimSpace(data[count])
		case 8:
			SavingsBond.InterestRate = strings.TrimSpace(data[count])
		case 9:
			SavingsBond.Value = strings.TrimSpace(data[count])
		case 10:
			SavingsBond.Note = strings.TrimSpace(data[count])
		default:
			log.Println("Unknown value")
		}
	}

	return SavingsBond, nil
}

// FindUnmaturedBonds returns bonds that have not yet matured
func (b Bonds) FindUnmaturedBonds() []Bond {
	var r []Bond
	for _, bd := range b.Bonds {
		// Treasury puts "MA" in the note column to designate a matured bond
		if bd.Note != "MA" {
			r = append(r, bd)
		}
	}
	return r
}

// GetTotalPurchasedPrice returns the total price paid for the bonds
func (b Bonds) GetTotalPurchasedPrice() (float64, error) {
	sum := 0.0
	for _, b := range b.Bonds {
		s := strings.Replace(b.IssuePrice, "$", "", 1)

		r, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return 0, err
		}
		sum += r
	}
	return sum, nil
}

// GetTotalValue returns the total value of all bonds
func (b Bonds) GetTotalValue() (float64, error) {
	sum := 0.0
	for _, b := range b.Bonds {
		s := strings.Replace(b.Value, "$", "", 1)

		r, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return 0, err
		}
		sum += r
	}
	return sum, nil
}

// GetTotalInterest returns the total interest accrued by all bonds
func (b Bonds) GetTotalInterest() (float64, error) {
	sum := 0.0
	for _, b := range b.Bonds {
		s := strings.Replace(b.Interest, "$", "", 1)

		r, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return 0, err
		}
		sum += r
	}
	return sum, nil
}
