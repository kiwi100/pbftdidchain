package transaction

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"
)

// Transaction 表示一个区块链交易
type Transaction struct {
	Hash        string    `json:"hash"`         // 交易哈希
	From        string    `json:"from"`         // 发送地址
	To          string    `json:"to"`           // 接收地址（可选）
	Value       *big.Int  `json:"value"`        // 交易金额（Wei）
	Timestamp   int64     `json:"timestamp"`     // 时间戳
	Status      string    `json:"status"`       // 状态：pending, committed, failed
	BlockSeq    int64     `json:"blockSeq"`     // 区块序列号（共识后分配）
	ConsensusTime time.Duration `json:"consensusTime"` // 共识耗时
}

// NewTransaction 创建新交易
func NewTransaction(hash, from, to string, value *big.Int) *Transaction {
	return &Transaction{
		Hash:      hash,
		From:      from,
		To:        to,
		Value:     value,
		Timestamp: time.Now().Unix(),
		Status:    "pending",
	}
}

// String 返回交易的字符串表示
func (t *Transaction) String() string {
	return fmt.Sprintf("Tx[%s] From[%s] Value[%s]", 
		t.Hash[:16], t.From[:16], t.Value.String())
}

// ToJSON 将交易转换为JSON
func (t *Transaction) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}

// FromJSON 从JSON创建交易
func FromJSON(data []byte) (*Transaction, error) {
	var tx Transaction
	err := json.Unmarshal(data, &tx)
	return &tx, err
}


