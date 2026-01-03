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
	BigPicture     *bigPicture
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
	log.Printf("Cache: updating Accounts")
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
	log.Printf("Cache: updating Categories")
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
		// No category budgets exist.
		return nil
	}
	f.cache.mu.Lock()
	log.Printf("Cache: updating CategoryTotals for key %d, %s, %s", key.CategoryID, key.Start, key.End)
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

// invalidateTransactionsCache will invalidate all cached transactions
// lists
func (f *Firefly) invalidateTransactionsCache() {
	f.cache.mu.Lock()
	defer f.cache.mu.Unlock()
	log.Print("Cache: clearing Transactions")
	f.cache.Transactions = nil
}

func (f *Firefly) invalidateAllCaches() {
	f.cache.mu.Lock()
	defer f.cache.mu.Unlock()
	log.Print("Cache: clearing all caches")

	f.cache.Accounts = make([]Account, 0, len(f.cache.Accounts))
	f.cache.BigPicture = nil
	f.cache.Categories = make([]Category, 0, len(f.cache.Categories))
	f.cache.CategoryTotals = map[categoryTotalsKey][]CategoryTotal{}
	f.cache.Transactions = map[transactionsKey][]Transactions{}
}

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
	log.Printf("Cache: updating Transactions for key %d, %s, %s", key.Page, key.Start, key.End)
	if f.cache.Transactions == nil {
		f.cache.Transactions = make(map[transactionsKey][]Transactions)
	}
	f.cache.Transactions[key] = t
	f.cache.mu.Unlock()
	return nil
}

func (f *Firefly) CachedBigPicture() (*bigPicture, error) {
	f.cache.mu.Lock()
	if f.cache.BigPicture == nil {
		f.cache.mu.Unlock()
		err := f.refreshBigPicture()
		if err != nil {
			return nil, err
		}
		f.cache.mu.Lock()
	}
	defer f.cache.mu.Unlock()
	return f.cache.BigPicture, nil
}

func (f *Firefly) refreshBigPicture() error {
	bp, err := f.fetchBigPicture()
	if err != nil {
		return err
	}
	f.cache.mu.Lock()
	log.Printf("Cache: updating Big Picture")
	f.cache.BigPicture = bp
	f.cache.mu.Unlock()
	return nil
}

// RefreshCaches refreshes caches for the current budget and its related data.
// This is intended to be run when lychnos launches.
func (f *Firefly) RefreshCaches(c *categorybudget.CategoryBudgets, b *budget.Budgets) error {
	err := f.refreshCategories()
	if err != nil {
		return fmt.Errorf("failed to refresh categories: %s", err)
	}
	err = f.refreshAccounts()
	if err != nil {
		return fmt.Errorf("failed to refresh accounts: %s", err)
	}
	err = f.refreshBigPicture()
	if err != nil {
		return fmt.Errorf("failed to refresh big picture: %s", err)
	}

	bs, err := b.List()
	if err != nil {
		return fmt.Errorf("failed to list budgets: %s", err)
	}
	cbs, err := c.List()
	if err != nil {
		return fmt.Errorf("failed to list category budgets: %s", err)
	}

	now := time.Now()
	for _, bgt := range bs {
		if now.Before(bgt.Start) || now.After(bgt.End) {
			continue // Only update the cache for the current budget.
		}

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
			intervals := interval.Get(bgt.Start, bgt.End, time.Now().Local().Location())
			for _, i := range intervals {
				go func(i interval.ReportingInterval, cb categorybudget.CategoryBudget) {
					key := categoryTotalsKey{
						CategoryID: cb.Category,
						Start:      i.Start.Local(),
						End:        i.End.Local(),
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

// InvalidateCacheIfAccountBalancesHaveChanged checks if the balances of the
// asset accounts are different than what we currently have in the cache. If
// yes, we refresh all our caches.
//
// Initially, I designed the cache with the assumption that all transaction
// inputs would happen in lychnos. However, it turns out that I often balance
// accounts in firefly-iii directly. This causes the lychnos cache to become
// stale, and then I have to manually restart lychnos. Instead, we can call this
// function on some interval to make sure that our caches are always fresh.
func (f *Firefly) InvalidateCacheIfAccountBalancesHaveChanged(c *categorybudget.CategoryBudgets, b *budget.Budgets) error {
	cachedAccounts, err := f.CachedAccounts()
	if err != nil {
		return err
	}
	freshAssetAccounts, err := f.ListAccounts(AcctTypeAsset)
	if err != nil {
		return err
	}

	// Nested loop isn't ideal, but it's easier for now than changing the data
	// structures. Plus, we won't have more than a few hundred accounts.
	for _, freshAssetAccount := range freshAssetAccounts {
		for _, cachedAccount := range cachedAccounts {
			if freshAssetAccount.ID == cachedAccount.ID && !freshAssetAccount.Attributes.CurrentBalance.Equal(cachedAccount.Attributes.CurrentBalance) {
				f.invalidateAllCaches()
				return f.RefreshCaches(c, b)
			}
		}
	}
	return nil
}

// InvalidateCacheIfCategoriesHaveChanged does the same thing as
// InvalidateCacheIfAccountBalancesHaveChanged, but for categories. This makes
// initial setup easier, or invalidates the cache if you happen to edit the
// categories. Since we validate that you're not creating new categories in the
// lychnos frontend, it's important to keep the cache up-to-date.
func (f *Firefly) InvalidateCacheIfCategoriesHaveChanged(c *categorybudget.CategoryBudgets, b *budget.Budgets) error {
	cachedCategories, err := f.CachedCategories()
	if err != nil {
		return err
	}

	freshCategories, err := f.Categories()
	if err != nil {
		return err
	}

	if len(cachedCategories) != len(freshCategories) {
		f.invalidateAllCaches()
		return f.RefreshCaches(c, b)
	}

	for _, freshCategory := range freshCategories {
		var found bool
		for _, cachedCategory := range cachedCategories {
			if freshCategory.ID == cachedCategory.ID && freshCategory.Name == cachedCategory.Name {
				found = true
				break
			}
		}
		if !found {
			f.invalidateAllCaches()
			return f.RefreshCaches(c, b)
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
			log.Printf("Cache: clearing CategoryTotals for key %d, %s, %s", k.CategoryID, k.Start, k.End)
			delete(f.cache.CategoryTotals, k)
			go func(k categoryTotalsKey) {
				f.refreshCategoryTotals(k)
			}(k)
		}
	}
}
