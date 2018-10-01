package goeth

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/dneprix/goeth/model"
	"github.com/dneprix/goeth/tx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethereum "github.com/ethereum/go-ethereum/core/types"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

type Service interface {
	SendWei(from common.Address, to common.Address, amount *model.Wei) (*model.Transaction, error)
	GetLast() (txs []*model.Transaction)
	LoadAccounts() error
	LoadBalances()
	LoadHistoryBlocks(num int)
	LoadNewBlocks()
	LoadLast()
}

func NewService(db *gorm.DB, client Client) (Service, error) {
	svc := &service{
		db:       db,
		client:   client,
		accounts: make(map[common.Address]string),
		nblocksC: make(chan *types.Header),
		txLast:   tx.NewLast(db),
	}
	return svc, nil
}

type service struct {
	db       *gorm.DB
	client   Client
	accounts map[common.Address]string
	nblocksC chan *types.Header
	txLast   tx.Last
}

func (s *service) SendWei(from common.Address, to common.Address, amount *model.Wei) (*model.Transaction, error) {
	pass, ok := s.accounts[from]
	if !ok {
		return nil, fmt.Errorf("account is not added to service: %v", from.String())
	}
	txHash, err := s.client.SendTransaction(from, to, amount, pass)
	if err != nil {
		return nil, fmt.Errorf("send tx fail: %v", err)
	}
	tx := &model.Transaction{
		AccountFrom: from,
		AccountTo:   to,
		Amount:      amount,
		Hash:        *txHash,
	}
	s.db.Create(tx)
	return tx, nil
}

func (s *service) GetLast() (txs []*model.Transaction) {
	return s.txLast.View()
}

func (s *service) LoadAccounts() error {
	addrList, err := s.client.AccountsList()
	if err != nil {
		return fmt.Errorf("get accounts: %v", err)
	}
	for _, addr := range addrList {
		err := s.processAccount(addr)
		if err != nil {
			log.Warnf("add account %v error: %v", addr.String(), err)
			continue
		}
	}
	if len(s.accounts) == 0 {
		return fmt.Errorf("no accounts for service")
	}
	return nil
}

func (s *service) processAccount(addr common.Address) error {
	s.accounts[addr] = "410b410b"
	s.client.CheckPassphrase(addr, s.accounts[addr])
	return nil

	var answer string
	fmt.Printf("\nDo you want to add %v account? (y/N): ", addr.String())
	if _, err := fmt.Scan(&answer); err != nil {
		return err
	}
	if strings.ToLower(answer) != "y" {
		return nil
	}
	var bytePass []byte
	var okPass bool
	for i := 0; i < 3; i++ {
		fmt.Print("Enter passphrase: ")
		bytePass, _ = terminal.ReadPassword(0)
		fmt.Print("\n")
		if okPass = s.client.CheckPassphrase(addr, string(bytePass)); okPass {
			break
		}
	}
	if !okPass {
		log.Fatalf("Wrong passphase")
	}

	s.db.Save(&model.Account{
		Address: addr,
	})
	s.accounts[addr] = string(bytePass)
	return nil
}

func (s *service) LoadBalances() {
	for addr, _ := range s.accounts {
		balance, err := s.client.BalanceAt(addr)
		if err != nil {
			log.Errorf("%v balance:", addr.String(), err)
			continue
		}
		s.db.Model(&model.Account{Address: addr}).
			Updates(&model.Account{Balance: balance})
	}
}

func (s *service) LoadNewBlocks() {
	subs, err := s.client.SubscribeNewHead(s.nblocksC)
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln("Subscription started, listening for new blocks")
	for {
		select {
		case h := <-s.nblocksC:
			block, err := s.client.BlockByHash(h.Hash())
			if err != nil {
				log.Warnf("error get block %v: %v", h.Hash().String(), err)
				continue
			}
			s.processBlock(block, nil)
			go s.processConfirmations(block.Number())
			go s.LoadBalances()
		case err := <-subs.Err():
			log.Warnln("error from subscription: ", err)
			log.Warnln("waiting 10 seconds and retrying")
			time.Sleep(time.Second * 10)
			nsubs, err := s.client.SubscribeNewHead(s.nblocksC)
			if err != nil {
				log.Fatal("error trying to resubscribe:", err)
			}
			subs = nsubs
		}
	}
}

func (s *service) LoadHistoryBlocks(num int) {
	lastBlock, err := s.client.BlockByNumber(nil)
	if err != nil {
		log.Warnln("load history error: get last block: %v", err)
		return
	}
	bar := pb.StartNew(num)
	for i := 0; i < num; i++ {
		blockNum := big.NewInt(0).Sub(lastBlock.Number(), big.NewInt(int64(i)))
		block, err := s.client.BlockByNumber(blockNum)
		if err != nil {
			log.Warnf("error get block %v: %v", blockNum.String(), err)
			continue
		}
		s.processBlock(block, lastBlock)
		bar.Increment()
	}
	bar.Finish()
}

func (s *service) processBlock(block *ethereum.Block, lastBlock *ethereum.Block) {
	for index, tx := range block.Transactions() {
		txTo := tx.To()
		if txTo == nil {
			continue
		}
		if _, ok := s.accounts[*txTo]; ok {
			confirmations := decimal.New(0, 0)
			blockNum := decimal.NewFromBigInt(block.Number(), 0)
			if lastBlock != nil {
				confirmations = decimal.NewFromBigInt(lastBlock.Number(), 0).Sub(blockNum)
			}
			txFrom, _ := s.client.TransactionSender(tx, block.Hash(), uint(index))
			s.db.Save(&model.Transaction{
				Hash:          tx.Hash(),
				AccountFrom:   txFrom,
				AccountTo:     *txTo,
				Amount:        model.WeiFromString(tx.Value().String()),
				BlockNum:      blockNum,
				Confirmations: confirmations,
				ChainId:       decimal.NewFromBigInt(tx.ChainId(), 0),
				Gas:           tx.Gas(),
				GasPrice:      decimal.NewFromBigInt(tx.GasPrice(), 0),
				Nonce:         tx.Nonce(),
				Protected:     tx.Protected(),
				CreatedAt:     time.Now(),
			})
		}
	}
}

func (s *service) processConfirmations(lastBlockNum *big.Int) {
	var txs []model.Transaction
	s.db.Where("block_num > 0 AND confirmations < ?", model.TX_CONFIRM_MAX).Find(&txs)
	for _, tx := range txs {
		txBlockNum, err := s.client.BlockNumByTxHash(tx.Hash)
		if err != nil {
			log.Warnf("confirmation get tx %v: %v", tx.Hash.String(), err)
			return
		}
		confirmations := decimal.NewFromBigInt(lastBlockNum, 0).Sub(*txBlockNum)
		if txBlockNum.IsPositive() && tx.Confirmations.LessThan(confirmations) {
			s.db.Model(&model.Transaction{Hash: tx.Hash}).
				Updates(&model.Transaction{
					Confirmations: confirmations,
					BlockNum:      *txBlockNum,
				})
		}
	}
}

func (s *service) LoadLast() {
	accounts := make([]common.Address, 0, len(s.accounts))
	for addr := range s.accounts {
		accounts = append(accounts, addr)
	}
	s.txLast.Load(accounts)
}
