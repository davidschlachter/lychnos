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
	Transactions   map[transactionsKey][]Transactions
	mu             sync.Mutex
}

type categoryTotalsKey struct {
	CategoryID int
	Start      time.Time
	End        time.Time
}

type transactionsKey struct {
	Page  int
	Start string
	End   string
}

func (f *Firefly) CachedAccounts() ([]Account, error) {
	f.cache.mu.Lock()
	if f.cache.Accounts == nil {
		f.cache.mu.Unlock()
		err := f.refreshAccounts()
		if err != nil {
			return nil, err
		}
		f.cache.mu.Lock()
	}
	defer f.cache.mu.Unlock()
	return f.cache.Accounts, nil
}

func (f *Firefly) refreshAccounts() error {
	c, err := f.ListAccounts("")
	if err != nil {
		return err
	}
	f.cache.mu.Lock()
	defer f.cache.mu.Unlock()
	log.Printf("Updating Accounts cache")
	f.cache.Accounts = c
	return nil
}

func (f *Firefly) CachedCategories() ([]Category, error) {
	f.cache.mu.Lock()
	if f.cache.Categories == nil {
		f.cache.mu.Unlock()
		err := f.refreshCategories()
		if err != nil {
			return nil, err
		}
		f.cache.mu.Lock()
	}
	defer f.cache.mu.Unlock()
	return f.cache.Categories, nil
}

// refreshCategories refreshes the cached Categories. The caller is responsible
// for locking the mutex.
func (f *Firefly) refreshCategories() error {
	c, err := f.Categories()
	if err != nil {
		return err
	}
	f.cache.mu.Lock()
	defer f.cache.mu.Unlock()
	log.Printf("Updating Categories cache")
	f.cache.Categories = c
	return nil
}

func (f *Firefly) CachedListCategoryTotals(start, end time.Time) ([]CategoryTotal, error) {
	f.cache.mu.Lock()
	key := categoryTotalsKey{
		Start: start,
		End:   end,
	}
	_, ok := f.cache.CategoryTotals[key]
	if !ok {
		f.cache.mu.Unlock()
		err := f.refreshCategoryTotals(key)
		if err != nil {
			return nil, err
		}
		f.cache.mu.Lock()
	}
	defer f.cache.mu.Unlock()
	return f.cache.CategoryTotals[key], nil
}

func (f *Firefly) CachedFetchCategoryTotals(catID int, start, end time.Time) ([]CategoryTotal, error) {
	f.cache.mu.Lock()
	key := categoryTotalsKey{
		CategoryID: catID,
		Start:      start,
		End:        end,
	}
	_, ok := f.cache.CategoryTotals[key]
	if !ok {
		f.cache.mu.Unlock()
		err := f.refreshCategoryTotals(key)
		if err != nil {
			return nil, err
		}
		f.cache.mu.Lock()
	}
	defer f.cache.mu.Unlock()
	return f.cache.CategoryTotals[key], nil
}

func (f *Firefly) refreshCategoryTotals(key categoryTotalsKey) error {
	var (
		c   []CategoryTotal
		err error
	)
	f.cache.mu.Lock()
	if f.cache.CategoryTotals == nil {
		f.cache.CategoryTotals = make(map[categoryTotalsKey][]CategoryTotal)
	}
	f.cache.mu.Unlock()
	if key.CategoryID == 0 {
		c, err = f.ListCategoryTotals(key.Start, key.End)
	} else {
		c, err = f.FetchCategoryTotal(key.CategoryID, key.Start, key.End)
	}
	if err != nil {
		return fmt.Errorf("could not update category totals cache: %s", err)
	}
	if key.CategoryID != 0 && len(c) != 1 {
		return fmt.Errorf("got %d category totals, wanted 1 for key %d, %s, %s", len(c), key.CategoryID, key.Start, key.End)
	}
	if key.CategoryID == 0 && len(c) == 0 {
		return fmt.Errorf("got 0 category totals, wanted at least one for key %d, %s, %s", key.CategoryID, key.Start, key.End)
	}
	f.cache.mu.Lock()
	log.Printf("Updating CategoryTotals cache for: %d, %s, %s", key.CategoryID, key.Start, key.End)
	f.cache.CategoryTotals[key] = c
	f.cache.mu.Unlock()
	return nil
}

func (f *Firefly) CachedTransactions(key transactionsKey) ([]Transactions, error) {
	f.cache.mu.Lock()
	_, ok := f.cache.Transactions[key]
	if !ok {
		f.cache.mu.Unlock()
		err := f.refreshTransactions(key)
		if err != nil {
			return nil, err
		}
		f.cache.mu.Lock()
	}
	defer f.cache.mu.Unlock()
	return f.cache.Transactions[key], nil
}

// refreshTransactions refreshes the cached Transactions. The caller is
// responsible for locking the mutex.
func (f *Firefly) refreshTransactions(key transactionsKey) error {
	f.cache.mu.Lock()
	if f.cache.Transactions == nil {
		f.cache.Transactions = make(map[transactionsKey][]Transactions)
	}
	f.cache.mu.Unlock()
	t, err := f.ListTransactions(key)
	if err != nil {
		return err
	}
	f.cache.mu.Lock()
	log.Printf("Updating Transactions cache for page %d, start %s, end %s", key.Page, key.Start, key.End)
	f.cache.Transactions[key] = t
	f.cache.mu.Unlock()
	return nil
}

func (f *Firefly) RefreshCaches(c *categorybudget.CategoryBudgets, b *budget.Budgets) error {
	f.refreshCategories()
	f.refreshAccounts()

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
			key := categoryTotalsKey{
				Start: bgt.Start,
				End:   bgt.End,
			}
			err := f.refreshCategoryTotals(key)
			if err != nil {
				log.Printf("Failed to seed category totals cache: %s", err)
			}
		}(bgt)
		for _, cb := range cbs {
			if cb.Budget != bgt.ID {
				continue
			}
			intervals := interval.Get(bgt.Start, bgt.End, time.Now().UTC().Location())
			for _, i := range intervals {
				go func(i interval.ReportingInterval, cb categorybudget.CategoryBudget) {
					key := categoryTotalsKey{
						CategoryID: cb.Category,
						Start:      i.Start,
						End:        i.End,
					}
					err := f.refreshCategoryTotals(key)
					if err != nil {
						log.Printf("Failed to seed category totals cache: %s", err)
					}
				}(i, cb)
			}
		}
	}

	return nil
}

// refreshCategoryTxnCache will invalidate cache entries related to a particular
// category and time. This should be called after creating a transaction.
func (f *Firefly) refreshCategoryTxnCache(tgt categoryTotalsKey) {
	f.cache.mu.Lock()
	defer f.cache.mu.Unlock()

	for k := range f.cache.CategoryTotals {
		if (k.Start.Year() == tgt.Start.Year() && (k.CategoryID == 0 || k.CategoryID == tgt.CategoryID)) ||
			(k.End.Year() == tgt.End.Year() && (k.CategoryID == 0 || k.CategoryID == tgt.CategoryID)) {
			log.Printf("Clearing CategoryTotals cache for: %d, %s, %s", k.CategoryID, k.Start, k.End)
			delete(f.cache.CategoryTotals, k)
			go func(k categoryTotalsKey) {
				f.refreshCategoryTotals(k)
			}(k)
		}
	}
}
