package main

import (
	"fmt"
	"sync"
)

func Producer(limit int, q chan<- int) {
	defer wg.Done()

	for i := 0; i <= limit; i++ {
		q <- i
	}

	close(q)
}

func Consumer(q <-chan int) {
	defer wg.Done()

	for v := range q {
		fmt.Println(v)
	}
}

/*
*	Pub/sub:
*		Publisher: publish message to Message brokers
*		Message brokers: contains group of messages called topics/channels
*		Subcriber: subcribe to topic, get notify when new messages arrive at subcribed topics/channels
 */

type Publisher interface {
	Publish()
}

type Subcriber interface {
	Subcribe()
}

type TopicKeeper struct {
	Topics map[string]chan any
}

type Pub struct {
	Topic chan any
	Publisher
}

type Sub struct {
	Topic chan any
	Subcriber
}

var (
	mutex = new(sync.Mutex)
	wg    = new(sync.WaitGroup)
)

func (t *TopicKeeper) Register(p *Pub, topicName string) {
	mutex.Lock()
	if t.Topics == nil {
		t.Topics = make(map[string]chan any)
	}

	topic, ok := t.Topics[topicName]

	if !ok {
		topic = make(chan any, 0)
		t.Topics[topicName] = topic
	}
	mutex.Unlock()

	p.Topic = topic
}

func (t *TopicKeeper) Listen(s *Sub, topicName string) {
	mutex.Lock()
	if t.Topics == nil {
		t.Topics = make(map[string]chan any)
	}

	topic, ok := t.Topics[topicName]

	if !ok {
		topic = make(chan any, 10)
		t.Topics[topicName] = make(chan any, 0)
	}
	mutex.Unlock()

	s.Topic = topic
}

func (p *Pub) Publish(c func(arg ...any)) {
	c()
}

func (s *Sub) Subcribe(c func(arg ...any)) {
	c()
}

func main() {
	t := &TopicKeeper{}
	p := &Pub{}
	s1 := &Sub{}
	s2 := &Sub{}

	go t.Register(p, "student")
	go t.Listen(s1, "student")
	go t.Listen(s2, "student")

	wg.Add(1)
	go p.Publish(func(arg ...any) {
		for i := 0; i < 10; i++ {
			p.Topic <- i
		}
		close(p.Topic)
		wg.Done()
	})

	go s1.Subcribe(func(arg ...any) {
		for {
			select {
			case i := <-s1.Topic:
				fmt.Println(i)
			default:
				fmt.Println("Waiting")
			}
		}
	})

	go s2.Subcribe(func(arg ...any) {
		select {
		case i := <-s2.Topic:
			fmt.Println(i)
		default:
			fmt.Println("Waiting")
		}
	})

	wg.Wait()
}
