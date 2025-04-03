package firefly_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidschlachter/lychnos/src/backend/firefly"
	"github.com/shopspring/decimal"
)

func TestListAccounts(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/accounts/?type=expense", nil)
	f.HandleAccount(w, req)

	var a []firefly.Account
	json.NewDecoder(w.Body).Decode(&a)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Status code = %d, want %d\n", w.Result().StatusCode, http.StatusOK)
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
	if !a[0].Attributes.CurrentBalance.Equal(expectedBalance) {
		t.Fatalf("Got account balance %s, wanted 53.97", a[0].Attributes.CurrentBalance)
	}
}
