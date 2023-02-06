package server

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/iotaledger/inx-api-core-v1/pkg/restapi"
	"github.com/iotaledger/inx-api-core-v1/pkg/utxo"
)

func (s *DatabaseServer) receipts(_ echo.Context) (*receiptsResponse, error) {
	receipts := make([]*utxo.ReceiptTuple, 0)
	if err := s.UTXOManager.ForEachReceiptTuple(func(rt *utxo.ReceiptTuple) bool {
		receipts = append(receipts, rt)

		return true
	}); err != nil {
		return nil, errors.WithMessagef(echo.ErrInternalServerError, "unable to retrieve receipts: %s", err)
	}

	return &receiptsResponse{Receipts: receipts}, nil
}

func (s *DatabaseServer) receiptsByMigratedAtIndex(c echo.Context) (*receiptsResponse, error) {
	migratedAt, err := restapi.ParseMilestoneIndexParam(c)
	if err != nil {
		return nil, err
	}

	receipts := make([]*utxo.ReceiptTuple, 0)
	if err := s.UTXOManager.ForEachReceiptTupleMigratedAt(migratedAt, func(rt *utxo.ReceiptTuple) bool {
		receipts = append(receipts, rt)

		return true
	}); err != nil {
		return nil, errors.WithMessagef(echo.ErrInternalServerError, "unable to retrieve receipts for migrated at index %d: %s", migratedAt, err)
	}

	return &receiptsResponse{Receipts: receipts}, nil
}
