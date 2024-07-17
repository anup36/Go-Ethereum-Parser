package scheduler

import (
	"time"

	"eth-parser/parser"
)

func StartScheduler(ep *parser.EthParser) {
	for {
		ep.FetchCurrentBlock()
		currentBlock := ep.GetCurrentBlock()
		ep.ParseBlock(currentBlock)
		time.Sleep(10 * time.Second)
	}
}
