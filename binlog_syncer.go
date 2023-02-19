package main

import "time"

type BinlogSyncer struct {
}

func (b *BinlogSyncer) StartBackup(backupDir string, p Position, timeout time.Duration) error {
	return nil
}
