package configurations

import (
	"api/internal/commons"
	"fmt"
	"reflect"
	"testing"
)

func TestUpdateActiveConfigurationTable(t *testing.T) {
	type test struct {
		configurationTable commons.ConfigurationTable
		expMessagesFrom    map[commons.Address]<-chan commons.ConfigurationTable
		expMessage         commons.ConfigurationTable
	}

	subscriberChannel1 := make(chan commons.ConfigurationTable)
	subscriberChannel2 := make(chan commons.ConfigurationTable)
	subscriberChannel3 := make(chan commons.ConfigurationTable)

	activeSubscribers = map[commons.Address]chan<- commons.ConfigurationTable{
		"1.2.3.4":      subscriberChannel1,
		"10.0.1.1":     subscriberChannel2,
		"34.172.96.71": subscriberChannel3,
	}

	tests := []test{
		{
			configurationTable: commons.ConfigurationTable{
				{
					commons.ConfigurationRequirement{
						Package:     "com.softinit.iquitos.mainapp",
						CountryCode: "ZZ",
					},
					commons.ConfigurationChance{
						PercentileMin: 0,
						PercentileMax: 100,
						MainSku:       "default",
					},
				},
			},
			expMessagesFrom: map[commons.Address]<-chan commons.ConfigurationTable{
				"1.2.3.4":      subscriberChannel1,
				"10.0.1.1":     subscriberChannel2,
				"34.172.96.71": subscriberChannel3,
			},
			expMessage: commons.ConfigurationTable{
				{
					commons.ConfigurationRequirement{
						Package:     "com.softinit.iquitos.mainapp",
						CountryCode: "ZZ",
					},
					commons.ConfigurationChance{
						PercentileMin: 0,
						PercentileMax: 100,
						MainSku:       "default",
					},
				},
			},
		},
		{
			configurationTable: commons.ConfigurationTable{
				{
					commons.ConfigurationRequirement{
						Package:     "com.softinit.iquitos.mainapp",
						CountryCode: "US",
					},
					commons.ConfigurationChance{
						PercentileMin: 0,
						PercentileMax: 25,
						MainSku:       "rdm_premium_v3_020_trial_7d_monthly",
					},
				},
				{
					commons.ConfigurationRequirement{
						Package:     "com.softinit.iquitos.mainapp",
						CountryCode: "US",
					},
					commons.ConfigurationChance{
						PercentileMin: 25,
						PercentileMax: 50,
						MainSku:       "rdm_premium_v3_030_trial_7d_monthly",
					},
				},
			},
			expMessagesFrom: map[commons.Address]<-chan commons.ConfigurationTable{
				"1.2.3.4":      subscriberChannel1,
				"10.0.1.1":     subscriberChannel2,
				"34.172.96.71": subscriberChannel3,
			},
			expMessage: commons.ConfigurationTable{
				{
					commons.ConfigurationRequirement{
						Package:     "com.softinit.iquitos.mainapp",
						CountryCode: "US",
					},
					commons.ConfigurationChance{
						PercentileMin: 0,
						PercentileMax: 25,
						MainSku:       "rdm_premium_v3_020_trial_7d_monthly",
					},
				},
				{
					commons.ConfigurationRequirement{
						Package:     "com.softinit.iquitos.mainapp",
						CountryCode: "US",
					},
					commons.ConfigurationChance{
						PercentileMin: 25,
						PercentileMax: 50,
						MainSku:       "rdm_premium_v3_030_trial_7d_monthly",
					},
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			UpdateActiveConfigurationTable(tt.configurationTable)
			for _, retMessageFrom := range tt.expMessagesFrom {
				retMessage := <-retMessageFrom
				if !reflect.DeepEqual(tt.expMessage, retMessage) {
					t.Errorf("message mismatch: expected '%v', received '%v'", tt.expMessage, retMessage)
				}
			}
		})
	}
}

func TestAddSubscriber(t *testing.T) {
	type test struct {
		ip                commons.Address
		channel           chan commons.ConfigurationTable
		expMessage        commons.ConfigurationTable
		expActiveChannels []chan commons.ConfigurationTable
		expClosedChannels []chan commons.ConfigurationTable
	}

	subscriberChannel1 := make(chan commons.ConfigurationTable)
	subscriberChannel2 := make(chan commons.ConfigurationTable)
	subscriberChannel3 := make(chan commons.ConfigurationTable)

	activeSubscribers = make(map[commons.Address]chan<- commons.ConfigurationTable)
	activeConfigurationTable = commons.ConfigurationTable{
		{
			commons.ConfigurationRequirement{
				Package:     "com.softinit.iquitos.mainapp",
				CountryCode: "ZZ",
			},
			commons.ConfigurationChance{
				PercentileMin: 0,
				PercentileMax: 100,
				MainSku:       "default",
			},
		},
	}

	tests := []test{
		{
			ip:      "1.2.3.4",
			channel: subscriberChannel1,
			expMessage: commons.ConfigurationTable{
				{
					commons.ConfigurationRequirement{
						Package:     "com.softinit.iquitos.mainapp",
						CountryCode: "ZZ",
					},
					commons.ConfigurationChance{
						PercentileMin: 0,
						PercentileMax: 100,
						MainSku:       "default",
					},
				},
			},
			expActiveChannels: []chan commons.ConfigurationTable{subscriberChannel1},
			expClosedChannels: nil,
		},
		{
			ip:      "10.0.1.1",
			channel: subscriberChannel2,
			expMessage: commons.ConfigurationTable{
				{
					commons.ConfigurationRequirement{
						Package:     "com.softinit.iquitos.mainapp",
						CountryCode: "ZZ",
					},
					commons.ConfigurationChance{
						PercentileMin: 0,
						PercentileMax: 100,
						MainSku:       "default",
					},
				},
			},
			expActiveChannels: []chan commons.ConfigurationTable{subscriberChannel1, subscriberChannel2},
			expClosedChannels: nil,
		},
		{
			ip:      "10.0.1.1",
			channel: subscriberChannel3,
			expMessage: commons.ConfigurationTable{
				{
					commons.ConfigurationRequirement{
						Package:     "com.softinit.iquitos.mainapp",
						CountryCode: "ZZ",
					},
					commons.ConfigurationChance{
						PercentileMin: 0,
						PercentileMax: 100,
						MainSku:       "default",
					},
				},
			},
			expActiveChannels: []chan commons.ConfigurationTable{subscriberChannel1, subscriberChannel3},
			expClosedChannels: []chan commons.ConfigurationTable{subscriberChannel2},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			go AddSubscriber(tt.ip, tt.channel)
			retMessage := <-tt.channel
			if !reflect.DeepEqual(tt.expMessage, retMessage) {
				t.Errorf("message mismatch: expected '%v', received '%v'", tt.expMessage, retMessage)
			}
			for _, expActiveChannel := range tt.expActiveChannels {
				go func() {
					expActiveChannel <- commons.ConfigurationTable{}
				}()
				_, ok := <-expActiveChannel
				if !ok {
					t.Errorf("active channel mismatch: expected '%v', received '%v'", true, ok)
				}
			}
			for _, expClosedChannel := range tt.expClosedChannels {
				select {
				case _, ok := <-expClosedChannel:
					if ok {
						t.Errorf("closed channel mismatch: expected '%v', received '%v'", false, ok)
					}
				default:
					t.Errorf("closed channel mismatch: not returning anything")
				}
			}
		})
	}
}

func TestDelSubscriber(t *testing.T) {
	type test struct {
		ip                commons.Address
		expActiveChannels []chan commons.ConfigurationTable
		expClosedChannels []chan commons.ConfigurationTable
	}

	subscriberChannel1 := make(chan commons.ConfigurationTable)
	subscriberChannel2 := make(chan commons.ConfigurationTable)
	subscriberChannel3 := make(chan commons.ConfigurationTable)

	activeSubscribers = map[commons.Address]chan<- commons.ConfigurationTable{
		"1.2.3.4":      subscriberChannel1,
		"10.0.1.1":     subscriberChannel2,
		"34.172.96.71": subscriberChannel3,
	}

	tests := []test{
		{
			ip:                "1.2.3.4",
			expActiveChannels: []chan commons.ConfigurationTable{subscriberChannel2, subscriberChannel3},
			expClosedChannels: []chan commons.ConfigurationTable{subscriberChannel1},
		},
		{
			ip:                "10.0.1.1",
			expActiveChannels: []chan commons.ConfigurationTable{subscriberChannel3},
			expClosedChannels: []chan commons.ConfigurationTable{subscriberChannel1, subscriberChannel2},
		},
		{
			ip:                "34.172.96.71",
			expActiveChannels: []chan commons.ConfigurationTable{},
			expClosedChannels: []chan commons.ConfigurationTable{subscriberChannel1, subscriberChannel2, subscriberChannel3},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			DelSubscriber(tt.ip)
			for _, expActiveChannel := range tt.expActiveChannels {
				go func() {
					expActiveChannel <- commons.ConfigurationTable{}
				}()
				_, ok := <-expActiveChannel
				if !ok {
					t.Errorf("active channel mismatch: expected '%v', received '%v'", true, ok)
				}
			}
			for _, expClosedChannel := range tt.expClosedChannels {
				select {
				case _, ok := <-expClosedChannel:
					if ok {
						t.Errorf("closed channel mismatch: expected '%v', received '%v'", false, ok)
					}
				default:
					t.Errorf("closed channel mismatch: not returning anything")
				}
			}
		})
	}
}
