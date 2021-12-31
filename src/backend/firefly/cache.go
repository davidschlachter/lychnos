package firefly

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/davidschlachter/lychnos/src/backend/budget"
	"github.com/davidschlachter/lychnos/src/backend/categorybudget"
	"github.com/davidschlachter/lychnos/src/backend/interval"
)

// Since queries to firefly are slow (up to 5 seconds), keep a cache of these
// requests. Allow the cache to be initialized, and selectively updated
// on-demand. If we are only using this app to record new transactions, the
// cache should always be fresh.

type Cache struct {
	Accounts       []Account
	Categories     []Category
	CategoryTotals map[categoryTotalsKey][]CategoryTotal
	Transactions   map[int][]Transactions
	mu             sync.Mutex
}

type categoryTotalsKey struct {
	CategoryID int
	Start      time.Time
	End        time.Time
}

func (f *Firefly) CachedAccounts() ([]Account, error) {
	f.cache.mu.Lock()
	defer f.cache.mu.Unlock()
	if f.cache.Accounts == nil {
		err := f.refreshAccounts()
		if err != nil {
			return nil, err
		}
	}
	return f.cache.Accounts, nil
}

// refreshAccounts refreshes the cached Accounts. The caller is responsible for
// locking the mutex.
func (f *Firefly) refreshAccounts() error {
	log.Printf("Updating Accounts cache")
	c, err := f.ListAccounts("")
	if err != nil {
		return err
	}
	f.cache.Accounts = c
	return nil
}

func (f *Firefly) CachedCategories() ([]Category, error) {
	f.cache.mu.Lock()
	defer f.cache.mu.Unlock()
	if f.cache.Categories == nil {
		err := f.refreshCategories()
		if err != nil {
			return nil, err
		}
	}
	return f.cache.Categories, nil
}

// refreshCategories refreshes the cached Categories. The caller is responsible
// for locking the mutex.
func (f *Firefly) refreshCategories() error {
	log.Printf("Updating Categories cache")
	c, err := f.Categories()
	if err != nil {
		return err
	}
	f.cache.Categories = c
	return nil
}

func (f *Firefly) CachedListCategoryTotals(start, end time.Time) ([]CategoryTotal, error) {
	f.cache.mu.Lock()
	defer f.cache.mu.Unlock()

	key := categoryTotalsKey{
		Start: start,
		End:   end,
	}
	_, ok := f.cache.CategoryTotals[key]

	if !ok {
		err := f.refreshCategoryTotals(key)
		if err != nil {
			return nil, err
		}
	}
	return f.cache.CategoryTotals[key], nil
}

func (f *Firefly) CachedFetchCategoryTotals(catID int, start, end time.Time) ([]CategoryTotal, error) {
	f.cache.mu.Lock()
	defer f.cache.mu.Unlock()

	key := categoryTotalsKey{
		CategoryID: catID,
		Start:      start,
		End:        end,
	}
	_, ok := f.cache.CategoryTotals[key]

	if !ok {
		err := f.refreshCategoryTotals(key)
		if err != nil {
			return nil, err
		}
	}
	return f.cache.CategoryTotals[key], nil
}

// refreshCategoryTotals refreshes a cached CategoryTotalsList. The caller
// is responsible for locking the mutex.
func (f *Firefly) refreshCategoryTotals(key categoryTotalsKey) error {
	var (
		c   []CategoryTotal
		err error
	)
	log.Printf("Updating CategoryTotals cache for: %d, %s, %s", key.CategoryID, key.Start, key.End)
	if f.cache.CategoryTotals == nil {
		f.cache.CategoryTotals = make(map[categoryTotalsKey][]CategoryTotal)
	}
	if key.CategoryID == 0 {
		c, err = f.ListCategoryTotals(key.Start, key.End)
	} else {
		c, err = f.FetchCategoryTotal(key.CategoryID, key.Start, key.End)
	}
	f.cache.CategoryTotals[key] = c
	return err
}

func (f *Firefly) CachedTransactions(page int) ([]Transactions, error) {
	f.cache.mu.Lock()
	defer f.cache.mu.Unlock()

	_, ok := f.cache.Transactions[page]
	if !ok {
		err := f.refreshTransactions(page)
		if err != nil {
			return nil, err
		}
	}
	return f.cache.Transactions[page], nil
}

// refreshTransactions refreshes the cached Transactions. The caller is
// responsible for locking the mutex.
func (f *Firefly) refreshTransactions(page int) error {
	log.Printf("Updating Transactions cache")
	if f.cache.Transactions == nil {
		f.cache.Transactions = make(map[int][]Transactions)
	}
	t, err := f.ListTransactions(1)
	if err != nil {
		return err
	}
	f.cache.Transactions[page] = t
	return nil
}

func (f *Firefly) RefreshCaches(c *categorybudget.CategoryBudgets, b *budget.Budgets) error {
	bs, err := b.List()
	if err != nil {
		return fmt.Errorf("failed to list budgets: %s", err)
	}
	cbs, err := c.List()
	if err != nil {
		return fmt.Errorf("failed to list category budgets: %s", err)
	}

	for _, bgt := range bs {
		go func(bgt budget.Budget) {
			f.cache.mu.Lock()
			defer f.cache.mu.Unlock()
			key := categoryTotalsKey{
				Start: bgt.Start,
				End:   bgt.End,
			}
			_ = f.refreshCategoryTotals(key)
		}(bgt)
		for _, cb := range cbs {
			if cb.Budget != bgt.ID {
				continue
			}
			intervals := interval.Get(bgt.Start, bgt.End, time.Now().UTC().Location())
			for _, i := range intervals {
				go func(i interval.ReportingInterval, cb categorybudget.CategoryBudget) {
					f.cache.mu.Lock()
					defer f.cache.mu.Unlock()
					key := categoryTotalsKey{
						CategoryID: cb.Category,
						Start:      i.Start,
						End:        i.End,
					}
					f.refreshCategoryTotals(key)
				}(i, cb)
			}
		}
	}

	f.refreshCategories()
	f.refreshAccounts()

	return nil
}

// invalidateCategoryCache will invalidate cache entries related to a particular
// category. This should be called after creating a transaction.
func (f *Firefly) invalidateCategoryCache(tgt categoryTotalsKey) {
	f.cache.mu.Lock()
	defer f.cache.mu.Unlock()

	for k := range f.cache.CategoryTotals {
		if k.Start.Year() == tgt.Start.Year() || k.End.Year() == tgt.End.Year() || k.CategoryID == tgt.CategoryID {
			delete(f.cache.CategoryTotals, k)
		}
	}
}

func (f *Firefly) invalidateAccountsCache() {
	f.cache.mu.Lock()
	defer f.cache.mu.Unlock()

	f.cache.Accounts = nil
}

func (f *Firefly) invalidateTransactionsCache() {
	f.cache.mu.Lock()
	defer f.cache.mu.Unlock()

	f.cache.Transactions = nil
}
