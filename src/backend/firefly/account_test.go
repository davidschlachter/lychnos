package firefly_test

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestListAccounts(t *testing.T) {
	a, err := f.CachedAccounts()
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	if len(a) != 1 {
		t.Fatalf("Got %d Accounts, wanted 1", len(a))
	}
	if a[0].ID != "464" {
		t.Fatalf("Got account ID %s, wanted 464", a[0].ID)
	}
	if a[0].Attributes.Name != "1Password" {
		t.Fatalf("Got account name %s, wanted 1Password", a[0].Attributes.Name)
	}
	if a[0].Attributes.Type != "expense" {
		t.Fatalf("Got account type %s, wanted expense", a[0].Attributes.Type)
	}
	expectedBalance, _ := decimal.NewFromString("53.97")
	if !a[0].Attributes.CurrentBalance.Equals(expectedBalance) {
		t.Fatalf("Got account balance %s, wanted 53.97", a[0].Attributes.CurrentBalance)
	}
}
