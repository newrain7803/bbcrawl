/* This file is part of bbcrawl, ©2020 Jörg Walter
 *  This software is licensed under the "GNU General Public License version 3" */

package download

import "sync"

type DownloadCounter struct {
	mu      *sync.Mutex
	counter uint64 //can count over 9000 downloads
}

func NewDownloadCounter() *DownloadCounter {
	return &DownloadCounter{mu: new(sync.Mutex)}
}

func (dc *DownloadCounter) Count() uint64 {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	dc.counter++
	return dc.counter
}
