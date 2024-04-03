package bonds

import (
	"errors"
	"reflect"
	"testing"

	"github.com/chrisbsmith/goldfinger/config"
)

func TestLoadBonds(t *testing.T) {

	scenarios := []struct {
		name                   string
		config                 *config.Config // value to pass as the config parameter in LoadBonds
		expectedBonds          Bonds
		expectedUnmaturedCount int
		expectedError          error
	}{
		{
			// Don't use an expected bond for an unmatured test because the values will be different
			// every month
			name: "non-matured-bond",
			config: &config.Config{
				Bonds: []config.ConfigBond{
					{
						Denomination: 50,
						Serial:       "abcdef",
						IssueDate:    "01/2000",
						Series:       "EE",
					},
				},
			},
			expectedUnmaturedCount: 1,
			expectedError:          nil,
		},
		{
			name: "matured-bond",
			config: &config.Config{
				Bonds: []config.ConfigBond{
					{
						Denomination: 50,
						Serial:       "abcdef",
						IssueDate:    "01/1990",
						Series:       "EE",
					},
				},
			},
			expectedBonds: Bonds{
				[]Bond{
					{
						Denomination:  "$50",
						Serial:        "abcdef",
						IssueDate:     "01/1990",
						Series:        "EE",
						NextAccrual:   "",
						FinalMaturity: "01/2020",
						IssuePrice:    "$25.00",
						Interest:      "$78.68",
						InterestRate:  "",
						Value:         "$103.68",
						Note:          "MA",
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "two-matured-bond",
			config: &config.Config{
				Bonds: []config.ConfigBond{
					{
						Denomination: 50,
						Serial:       "abcdef",
						IssueDate:    "01/1990",
						Series:       "EE",
					},
					{
						Denomination: 50,
						Serial:       "ghijk",
						IssueDate:    "01/1990",
						Series:       "EE",
					},
				},
			},
			expectedBonds: Bonds{
				[]Bond{
					{
						Denomination:  "$50",
						Serial:        "abcdef",
						IssueDate:     "01/1990",
						Series:        "EE",
						NextAccrual:   "",
						FinalMaturity: "01/2020",
						IssuePrice:    "$25.00",
						Interest:      "$78.68",
						InterestRate:  "",
						Value:         "$103.68",
						Note:          "MA",
					},
					{
						Denomination:  "$50",
						Serial:        "ghijk",
						IssueDate:     "01/1990",
						Series:        "EE",
						NextAccrual:   "",
						FinalMaturity: "01/2020",
						IssuePrice:    "$25.00",
						Interest:      "$78.68",
						InterestRate:  "",
						Value:         "$103.68",
						Note:          "MA",
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "invalid-bond",
			config: &config.Config{
				Bonds: []config.ConfigBond{
					{
						Denomination: 50,
						Serial:       "abcdef",
						IssueDate:    "01/2000",
						Series:       "AA",
					},
				},
			},
			expectedError: ErrInvalidDataReturned,
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {

			bonds, err := LoadBonds(scenario.config)
			if err != nil {
				if !errors.Is(err, scenario.expectedError) {
					t.Errorf("[%s] expected error %v, got %v", scenario.name, scenario.expectedError, err)
					return
				} else if errors.Is(err, scenario.expectedError) {
					return
				} else {
					t.Fatalf("[%s] failed to retrieve bond information: %v", scenario.name, err)
				}
			}

			// If we're expecting unmatured bonds, let's check to see if it matches
			if scenario.expectedUnmaturedCount > 0 {
				c := bonds.FindUnmaturedBonds()
				if len(c) != scenario.expectedUnmaturedCount {
					t.Errorf("[%s] expected unmature count %d, got %d", scenario.name, scenario.expectedUnmaturedCount, len(c))
				}
				return
			}

			len := len(bonds.Bonds)
			for i := 0; i < len; i++ {
				bond := bonds.Bonds[i]
				expectedBond := scenario.expectedBonds.Bonds[i]

				bondValues := reflect.ValueOf(bond)
				// fields := bondValues.Type()
				expectedBondValues := reflect.ValueOf(expectedBond)

				// First verify that the number of fields returned matches what is expected
				if bondValues.NumField() != expectedBondValues.NumField() {
					t.Errorf("[%s] number of fields in the bond (%d) doesn't match expected (%d)", scenario.name, bondValues.NumField(), expectedBondValues.NumField())
				}

				// Compare the bonds on a field by field basis
				for i := 0; i < bondValues.NumField(); i++ {
					if bondValues.Field(i).Interface().(string) != expectedBondValues.Field(i).Interface().(string) {
						t.Errorf("[%s] value (%s) does not match expected (%s)", scenario.name, bondValues.Field(i), expectedBondValues.Field(i))
					}
				}
			}
		})
	}
}
