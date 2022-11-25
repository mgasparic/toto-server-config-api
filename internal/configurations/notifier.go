package configurations

import (
	"api/internal/commons"
	"log"
	"sync"
)

var (
	activeConfigurationTable commons.ConfigurationTable
	activeSubscribers        map[commons.Address]chan<- commons.ConfigurationTable
	activeSubscribersMux     sync.RWMutex
)

func UpdateActiveConfigurationTable(configurationTable commons.ConfigurationTable) {
	activeSubscribersMux.Lock()
	defer activeSubscribersMux.Unlock()
	if activeSubscribers == nil {
		activeSubscribers = make(map[commons.Address]chan<- commons.ConfigurationTable)
	}
	activeConfigurationTable = configurationTable
	for _, channel := range activeSubscribers {
		go func(ch chan<- commons.ConfigurationTable) {
			ch <- activeConfigurationTable
		}(channel)
	}
}

func AddSubscriber(address commons.Address, channel chan<- commons.ConfigurationTable) {
	if activeSubscribers == nil {
		log.Fatal("configuration table has to be instantiated first")
	}

	activeSubscribersMux.Lock()
	defer activeSubscribersMux.Unlock()
	closeChannel(address)
	activeSubscribers[address] = channel
	channel <- activeConfigurationTable
}

func DelSubscriber(address commons.Address) {
	activeSubscribersMux.Lock()
	defer activeSubscribersMux.Unlock()
	closeChannel(address)
	delete(activeSubscribers, address)
}

func closeChannel(address commons.Address) {
	ch, ok := activeSubscribers[address]
	if ok {
		close(ch)
	}
}
