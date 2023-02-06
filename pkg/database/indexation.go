package database

import (
	"github.com/iotaledger/inx-api-core-v1/pkg/hornet"
	iotago "github.com/iotaledger/iota.go/v2"
)

const (
	IndexationIndexLength = 64
)

// padIndexationIndex returns a padded indexation index.
func padIndexationIndex(index []byte) []byte {
	return append(index, make([]byte, IndexationIndexLength-len(index))...)
}

// IndexMessageIDs returns all known message IDs for the given index.
func (db *Database) IndexMessageIDs(index []byte, maxResults int) (hornet.MessageIDs, error) {
	indexPadded := padIndexationIndex(index)

	iterations := 0
	var messageIDs hornet.MessageIDs
	if err := db.indexationStore.IterateKeys(indexPadded, func(key []byte) bool {
		iterations++

		messageIDs = append(messageIDs, hornet.MessageIDFromSlice(key[IndexationIndexLength:IndexationIndexLength+iotago.MessageIDLength]))

		// stop if maximum amount of iterations reached
		return iterations <= maxResults
	}); err != nil {
		return nil, err
	}

	return messageIDs, nil
}
