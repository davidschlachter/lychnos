package firefly

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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
	Categories     []Category
	CategoryTotals map[string][]CategoryTotal
	mu             sync.Mutex
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

	key := fmt.Sprintf("%s%s", start, end)
	_, ok := f.cache.CategoryTotals[key]

	if !ok {
		err := f.refreshCategoryTotalsList(start, end, key)
		if err != nil {
			return nil, err
		}
	}
	return f.cache.CategoryTotals[key], nil
}

// refreshCategoryTotalsList refreshes the cached CategoryTotalsList. The caller
// is responsible for locking the mutex.
func (f *Firefly) refreshCategoryTotalsList(start, end time.Time, key string) error {
	log.Printf("Updating CategoryTotals cache for: %s, %s", start, end)
	c, err := f.ListCategoryTotals(start, end)
	if f.cache.CategoryTotals == nil {
		f.cache.CategoryTotals = make(map[string][]CategoryTotal)
	}
	f.cache.CategoryTotals[key] = c
	return err
}

func (f *Firefly) CachedFetchCategoryTotals(catID int, start, end time.Time) ([]CategoryTotal, error) {
	f.cache.mu.Lock()
	defer f.cache.mu.Unlock()

	key := fmt.Sprintf("%d%s%s", catID, start, end)
	_, ok := f.cache.CategoryTotals[key]

	if !ok {
		err := f.refreshCategoryTotalsFetch(catID, start, end, key)
		if err != nil {
			return nil, err
		}
	}
	return f.cache.CategoryTotals[key], nil
}

// refreshCategoryTotalsFetch refreshes the cached CategoryTotalsList. The caller
// is responsible for locking the mutex.
func (f *Firefly) refreshCategoryTotalsFetch(catID int, start, end time.Time, key string) error {
	log.Printf("Updating CategoryTotals cache for: %d, %s, %s", catID, start, end)
	if f.cache.CategoryTotals == nil {
		f.cache.CategoryTotals = make(map[string][]CategoryTotal)
	}
	c, err := f.FetchCategoryTotal(catID, start, end)
	f.cache.CategoryTotals[key] = c
	return err
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
			_ = f.refreshCategoryTotalsList(bgt.Start, bgt.End, fmt.Sprintf("%s%s", bgt.Start, bgt.End))
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
					key := fmt.Sprintf("%d%s%s", cb.Category, i.Start, i.End)
					f.refreshCategoryTotalsFetch(cb.Category, i.Start, i.End, key)
				}(i, cb)
			}
		}
	}

	return nil
}

// invalidateCategoryCache will invalidate cache entries related to a particular
// category. This should be called after creating a transaction.
func (f *Firefly) invalidateCategoryCache(categoryID string) {
	f.cache.mu.Lock()
	defer f.cache.mu.Unlock()

	thisYear := strconv.Itoa(time.Now().Year())
	for k := range f.cache.CategoryTotals {
		if strings.HasPrefix(k, thisYear) || strings.HasPrefix(k, categoryID) {
			delete(f.cache.CategoryTotals, k)
		}
	}
}
