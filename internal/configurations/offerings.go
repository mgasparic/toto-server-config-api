package configurations

import (
	"api/internal/commons"
	"log"
	"math/rand"
	"sync"
	"time"
)

const defaultCountryCode = commons.CountryCode("ZZ")

var (
	randomizedQueues map[commons.ConfigurationRequirement]<-chan commons.MainSku
	done             chan struct{}
	mux              sync.RWMutex
)

func UpdateConfigurationTable(configurationTable commons.ConfigurationTable) {
	mux.Lock()
	defer mux.Unlock()

	newRandomizedQueues := make(map[commons.ConfigurationRequirement]<-chan commons.MainSku)
	newRules := make(map[commons.ConfigurationRequirement][]commons.ConfigurationChance)
	newDone := make(chan struct{})

	for _, rule := range configurationTable {
		newRules[rule.ConfigurationRequirement] = append(newRules[rule.ConfigurationRequirement], rule.ConfigurationChance)
	}
	for ruleRequirement, ruleChances := range newRules {
		randomizedQueue := make(chan commons.MainSku, 100)
		newRandomizedQueues[ruleRequirement] = randomizedQueue
		go func(rq chan<- commons.MainSku, rcs []commons.ConfigurationChance) {
			for {
				percentile := commons.Percentile(rand.Intn(100) + 1)
				for _, rc := range rcs {
					if percentile > rc.PercentileMin && percentile <= rc.PercentileMax {
						select {
						case randomizedQueue <- rc.MainSku:
						case <-newDone:
							return
						}
						break
					}
				}
			}
		}(randomizedQueue, ruleChances)
	}

	if done != nil {
		close(done)
	}
	done = newDone
	randomizedQueues = newRandomizedQueues
}

func RandomMainSku(configurationRequirement commons.ConfigurationRequirement) (commons.MainSku, bool) {
	if randomizedQueues == nil {
		log.Fatal("configuration table has to be instantiated first")
	}

	randomizedQueue, ok := randomizedQueues[configurationRequirement]
	if !ok {
		configurationRequirement.CountryCode = defaultCountryCode
		randomizedQueue, ok = randomizedQueues[configurationRequirement]
		if !ok {
			return "", false
		}
	}

	select {
	case sku := <-randomizedQueue:
		return sku, true
	case <-time.After(time.Second):
		return "", false
	}
}
