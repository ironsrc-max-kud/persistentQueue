package services

import (
	"persistentQueue/models"
	"bytes"
	"sync"
)

type MessagingInterface interface {
	Push(m *models.Message, flush bool)
	Pop(queueId string, n int) bytes.Buffer
}

func InitMessaging() *Messaging {
	ser := Messaging{fileQueues: make(map[string]*FileQueue), mutex: sync.RWMutex{}}
	return &ser
}

type Messaging struct {
	mutex      sync.RWMutex
	fileQueues map[string]*FileQueue
}


func (s *Messaging) Push(m *models.Message, flush bool) {
	//Create a single file descriptor for each queue_id
	s.mutex.RLock()
	val, ok := s.fileQueues[m.QueueId]
	s.mutex.RUnlock()

	if !ok{
		s.mutex.Lock()
		if val, ok = s.fileQueues[m.QueueId]; !ok {
			val = NewFileQueue(m.QueueId, 1024*1024*500)
			s.fileQueues[m.QueueId] = val
		}
		s.mutex.Unlock()
	}

	val.Push(m.Data, flush)


}

func (s *Messaging) Pop(queueId string, n int) bytes.Buffer {
	b := bytes.Buffer{}
	if _, ok := s.fileQueues[queueId]; !ok {
		return b
	}
	return s.fileQueues[queueId].Pop(n)
}