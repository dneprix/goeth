package goeth

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/gorm"
)

type Service interface {
	SendEth() (interface{}, error)
	GetLast() (interface{}, error)
	SyncWallets()
}

func NewService(db *gorm.DB, client *ethclient.Client) (Service, error) {
	svc := &service{
		db:     db,
		client: client,
	}
	return svc, nil
}

type service struct {
	db     *gorm.DB
	client *ethclient.Client
}

func (s *service) SendEth() (interface{}, error) {
	return nil, nil
}

func (s *service) GetLast() (interface{}, error) {
	return nil, nil
}

func (s *service) SyncWallets() {
	account := common.HexToAddress("0xa70ee95f624120cf53d97d46c8063d077282de93")
	balance, err := s.client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Println("[ERR]", err)
		return
	}
	log.Println(balance)
}
