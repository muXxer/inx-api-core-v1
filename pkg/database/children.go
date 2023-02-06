package database

import (
	"github.com/iotaledger/inx-api-core-v1/pkg/hornet"
	iotago "github.com/iotaledger/iota.go/v2"
)

// ChildrenMessageIDs returns the message IDs of the children of the given message.
func (db *Database) ChildrenMessageIDs(messageID hornet.MessageID, maxResults int) (hornet.MessageIDs, error) {
	var childrenMessageIDs hornet.MessageIDs

	iterations := 0
	if err := db.childrenStore.IterateKeys(messageID, func(key []byte) bool {
		iterations++

		childrenMessageIDs = append(childrenMessageIDs, hornet.MessageIDFromSlice(key[iotago.MessageIDLength:iotago.MessageIDLength+iotago.MessageIDLength]))

		// stop if maximum amount of iterations reached
		return iterations <= maxResults
	}); err != nil {
		return nil, err
	}

	return childrenMessageIDs, nil
}
