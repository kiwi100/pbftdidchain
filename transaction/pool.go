package transaction

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/big"
	"os"
	"sync"
	"time"
)

// Pool 交易池，管理待处理的交易
type Pool struct {
	mu            sync.RWMutex
	pendingTxs    []*Transaction      // 待处理交易队列
	committedTxs  map[string]*Transaction // 已提交交易（按hash索引）
	allTxs        []*Transaction        // 所有交易（用于统计）
	maxPoolSize   int                 // 最大池大小
	totalLoaded   int                 // 总加载交易数
	startTime     time.Time           // 开始时间（用于统计）
}

// NewPool 创建新的交易池
func NewPool(maxSize int) *Pool {
	return &Pool{
		pendingTxs:   make([]*Transaction, 0),
		committedTxs: make(map[string]*Transaction),
		allTxs:       make([]*Transaction, 0),
		maxPoolSize:  maxSize,
		startTime:    time.Now(),
	}
}

// LoadFromCSV 从CSV文件加载交易
func (p *Pool) LoadFromCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	
	// 读取表头
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}
	fmt.Printf("CSV Header: %v\n", header)

	// 读取数据
	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}

		// 解析交易数据
		// CSV格式: transactionHash,from,,,,,value
		if len(record) < 7 {
			continue
		}

		hash := record[0]
		from := record[1]
		valueStr := record[6]

		// 解析金额（科学计数法格式，如 6.00E+22）
		value, err := parseScientificNotation(valueStr)
		if err != nil {
			fmt.Printf("Warning: failed to parse value for tx %s: %v\n", hash, err)
			continue
		}

		tx := NewTransaction(hash, from, "", value)
		p.AddTransaction(tx)
		count++

		if count%10000 == 0 {
			fmt.Printf("Loaded %d transactions...\n", count)
		}
	}

	p.totalLoaded = count
	fmt.Printf("Successfully loaded %d transactions from CSV\n", count)
	return nil
}

// parseScientificNotation 解析科学计数法字符串为big.Int
func parseScientificNotation(s string) (*big.Int, error) {
	// 处理科学计数法，如 "6.00E+22"
	var floatVal big.Float
	_, _, err := floatVal.Parse(s, 10)
	if err != nil {
		return nil, err
	}

	// 转换为big.Int（Wei单位）
	var intVal big.Int
	floatVal.Int(&intVal)
	return &intVal, nil
}

// AddTransaction 添加交易到池中
func (p *Pool) AddTransaction(tx *Transaction) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 检查是否已存在
	if _, exists := p.committedTxs[tx.Hash]; exists {
		return
	}

	// 检查池大小
	if len(p.pendingTxs) >= p.maxPoolSize {
		// 移除最旧的交易
		p.pendingTxs = p.pendingTxs[1:]
	}

	p.pendingTxs = append(p.pendingTxs, tx)
	p.allTxs = append(p.allTxs, tx)
}

// GetNextTransaction 获取下一个待处理交易
func (p *Pool) GetNextTransaction() *Transaction {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.pendingTxs) == 0 {
		return nil
	}

	tx := p.pendingTxs[0]
	p.pendingTxs = p.pendingTxs[1:]
	return tx
}

// GetNextBatch 获取一批交易（用于批量处理）
func (p *Pool) GetNextBatch(size int) []*Transaction {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.pendingTxs) == 0 {
		return nil
	}

	batchSize := size
	if len(p.pendingTxs) < batchSize {
		batchSize = len(p.pendingTxs)
	}

	batch := make([]*Transaction, batchSize)
	copy(batch, p.pendingTxs[:batchSize])
	p.pendingTxs = p.pendingTxs[batchSize:]

	return batch
}

// CommitTransaction 标记交易为已提交
func (p *Pool) CommitTransaction(tx *Transaction, blockSeq int64, consensusTime time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	tx.Status = "committed"
	tx.BlockSeq = blockSeq
	tx.ConsensusTime = consensusTime
	p.committedTxs[tx.Hash] = tx
}

// GetPendingCount 获取待处理交易数量
func (p *Pool) GetPendingCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.pendingTxs)
}

// GetCommittedCount 获取已提交交易数量
func (p *Pool) GetCommittedCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.committedTxs)
}

// GetTotalLoaded 获取总加载交易数
func (p *Pool) GetTotalLoaded() int {
	return p.totalLoaded
}

// GetStatistics 获取统计信息
func (p *Pool) GetStatistics() *PoolStatistics {
	p.mu.RLock()
	defer p.mu.RUnlock()

	elapsed := time.Since(p.startTime)
	
	stats := &PoolStatistics{
		TotalLoaded:      p.totalLoaded,
		PendingCount:     len(p.pendingTxs),
		CommittedCount:   len(p.committedTxs),
		ElapsedTime:      elapsed,
		Throughput:       float64(len(p.committedTxs)) / elapsed.Seconds(),
	}

	// 计算平均共识时延
	if len(p.committedTxs) > 0 {
		var totalTime time.Duration
		count := 0
		for _, tx := range p.committedTxs {
			if tx.ConsensusTime > 0 {
				totalTime += tx.ConsensusTime
				count++
			}
		}
		if count > 0 {
			stats.AvgConsensusLatency = totalTime / time.Duration(count)
		}
	}

	return stats
}

// PoolStatistics 交易池统计信息
type PoolStatistics struct {
	TotalLoaded          int           `json:"totalLoaded"`
	PendingCount         int           `json:"pendingCount"`
	CommittedCount       int           `json:"committedCount"`
	ElapsedTime          time.Duration `json:"elapsedTime"`
	Throughput           float64       `json:"throughput"`           // TPS (Transactions Per Second)
	AvgConsensusLatency  time.Duration `json:"avgConsensusLatency"`  // 平均共识时延
}


