package lsp

import "log"

func (s *Server) async(queueName string, f func()) {
	s.jobQueuesLock.Lock()
	queueEntry, ok := s.jobQueues[queueName]
	if !ok {
		queueEntry = jobQueueEntry{
			queue: make(chan func(), 10_000),
		}

		go func() {
			for {
				select {
				case nextF := <-queueEntry.queue:
					log.Println("got next queue", queueName, queueEntry.lastId)
					nextF()
				}
			}

		}()
	}

	queueEntry.lastId++
	s.jobQueues[queueName] = queueEntry
	s.jobQueuesLock.Unlock()

	queueEntry.queue <- func() {
		s.jobQueuesLock.Lock()
		currentId := s.jobQueues[queueName].lastId
		s.jobQueuesLock.Unlock()

		// check and discard what we can, it is not perfect but avoids most
		if currentId != queueEntry.lastId {
			log.Println("IGNORED ", queueName, currentId, queueEntry.lastId)
			return
		}

		f() // this is not guaranteed the last one, but very likely
	}

}
