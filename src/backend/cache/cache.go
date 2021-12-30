// Package Cache stores results from firefly. Since queries to firefly are slow
// (up to 5 seconds), keep a cache of these requests. Allow the cache to be
// initialized, and selectively updated on-demand. If we are only using this app
// to record new transactions, the cache should always be fresh.
package cache

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/davidschlachter/lychnos/src/backend/budget"
	"github.com/davidschlachter/lychnos/src/backend/categorybudget"
	"github.com/davidschlachter/lychnos/src/backend/firefly"
	"github.com/davidschlachter/lychnos/src/backend/interval"
)

type Cache struct {
	Categories     []firefly.Category
	CategoryTotals map[string][]firefly.CategoryTotal
	mu             sync.Mutex
	f              *firefly.Firefly
}

func New(f *firefly.Firefly) *Cache {
	return &Cache{f: f}
}

func (h *Cache) CachedCategories() ([]firefly.Category, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.Categories == nil {
		err := h.refreshCategories()
		if err != nil {
			return nil, err
		}
	}
	return h.Categories, nil
}

// refreshCategories refreshes the cached Categories. The caller is responsible
// for locking the mutex.
func (h *Cache) refreshCategories() error {
	log.Printf("Updating Categories cache")
	c, err := h.f.Categories()
	if err != nil {
		return err
	}
	h.Categories = c
	return nil
}

func (h *Cache) CachedListCategoryTotals(start, end time.Time) ([]firefly.CategoryTotal, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	key := fmt.Sprintf("%s%s", start, end)
	_, ok := h.CategoryTotals[key]

	if !ok {
		err := h.refreshCategoryTotalsList(start, end, key)
		if err != nil {
			return nil, err
		}
	}
	return h.CategoryTotals[key], nil
}

// refreshCategoryTotalsList refreshes the cached CategoryTotalsList. The caller
// is responsible for locking the mutex.
func (h *Cache) refreshCategoryTotalsList(start, end time.Time, key string) error {
	log.Printf("Updating CategoryTotals cache for: %s, %s", start, end)
	c, err := h.f.ListCategoryTotals(start, end)
	if h.CategoryTotals == nil {
		h.CategoryTotals = make(map[string][]firefly.CategoryTotal)
	}
	h.CategoryTotals[key] = c
	return err
}

func (h *Cache) CachedFetchCategoryTotals(catID int, start, end time.Time) ([]firefly.CategoryTotal, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	key := fmt.Sprintf("%d%s%s", catID, start, end)
	_, ok := h.CategoryTotals[key]

	if !ok {
		err := h.refreshCategoryTotalsFetch(catID, start, end, key)
		if err != nil {
			return nil, err
		}
	}
	return h.CategoryTotals[key], nil
}

// refreshCategoryTotalsFetch refreshes the cached CategoryTotalsList. The caller
// is responsible for locking the mutex.
func (h *Cache) refreshCategoryTotalsFetch(catID int, start, end time.Time, key string) error {
	log.Printf("Updating CategoryTotals cache for: %d, %s, %s", catID, start, end)
	if h.CategoryTotals == nil {
		h.CategoryTotals = make(map[string][]firefly.CategoryTotal)
	}
	c, err := h.f.FetchCategoryTotal(catID, start, end)
	h.CategoryTotals[key] = c
	return err
}

func (h *Cache) RefreshCaches(c *categorybudget.CategoryBudgets, b *budget.Budgets) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	bs, err := b.List()
	if err != nil {
		return fmt.Errorf("failed to list budgets: %s", err)
	}
	cbs, err := c.List()
	if err != nil {
		return fmt.Errorf("failed to list category budgets: %s", err)
	}

	for _, bgt := range bs {
		_ = h.refreshCategoryTotalsList(bgt.Start, bgt.End, fmt.Sprintf("%s%s", bgt.Start, bgt.End))
		for _, cb := range cbs {
			if cb.Budget != bgt.ID {
				continue
			}
			intervals := interval.Get(bgt.Start, bgt.End, time.Now().UTC().Location())
			for _, i := range intervals {
				key := fmt.Sprintf("%d%s%s", cb.Category, i.Start, i.End)
				h.refreshCategoryTotalsFetch(cb.Category, i.Start, i.End, key)
			}
		}
	}

	return nil
}
