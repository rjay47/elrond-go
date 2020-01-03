package preprocess

import (
	"github.com/ElrondNetwork/elrond-go/storage"
	"github.com/ElrondNetwork/elrond-go/storage/txcache"
)

// createSortedTransactionsProvider creates a sorted transactions provider for a given cache
func createSortedTransactionsProvider(transactionsPreprocessor *transactions, cache storage.Cacher, cacheKey string) SortedTransactionsProvider {
	txCache, isTxCache := cache.(*txcache.TxCache)
	if isTxCache {
		return newTxCacheToSortedTransactionsProviderAdapter(txCache)
	}

	return newCacheToSortedTransactionsProviderAdapter(transactionsPreprocessor, cache, cacheKey)
}
