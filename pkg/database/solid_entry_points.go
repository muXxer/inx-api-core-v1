package database

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/iotaledger/inx-api-core-v1/pkg/hornet"
	"github.com/iotaledger/inx-api-core-v1/pkg/milestone"
)

type SolidEntryPoints struct {
	entryPointsMap map[string]milestone.Index
}

func newSolidEntryPoints() *SolidEntryPoints {
	return &SolidEntryPoints{
		entryPointsMap: make(map[string]milestone.Index),
	}
}

func (s *SolidEntryPoints) Add(messageID hornet.MessageID, milestoneIndex milestone.Index) {
	messageIDMapKey := messageID.ToMapKey()
	if _, exists := s.entryPointsMap[messageIDMapKey]; !exists {
		s.entryPointsMap[messageIDMapKey] = milestoneIndex
	}
}

func solidEntryPointsFromBytes(solidEntryPointsBytes []byte) (*SolidEntryPoints, error) {
	s := newSolidEntryPoints()

	bytesReader := bytes.NewReader(solidEntryPointsBytes)

	var err error

	solidEntryPointsCount := len(solidEntryPointsBytes) / (32 + 4)
	for i := 0; i < solidEntryPointsCount; i++ {
		messageIDBuf := make([]byte, 32)
		var msIndex uint32

		err = binary.Read(bytesReader, binary.LittleEndian, messageIDBuf)
		if err != nil {
			return nil, fmt.Errorf("solidEntryPoints: %w", err)
		}

		err = binary.Read(bytesReader, binary.LittleEndian, &msIndex)
		if err != nil {
			return nil, fmt.Errorf("solidEntryPoints: %w", err)
		}
		s.Add(hornet.MessageIDFromSlice(messageIDBuf), milestone.Index(msIndex))
	}

	return s, nil
}

func (db *Database) readSolidEntryPoints() (*SolidEntryPoints, error) {
	value, err := db.snapshotStore.Get([]byte("solidEntryPoints"))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to retrieve solid entry points", err)
	}

	points, err := solidEntryPointsFromBytes(value)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to convert solid entry points", err)
	}

	return points, nil
}

func (db *Database) loadSolidEntryPoints() error {
	solidEntryPoints, err := db.readSolidEntryPoints()
	if err != nil {
		return err
	}
	db.solidEntryPoints = solidEntryPoints

	return nil
}
