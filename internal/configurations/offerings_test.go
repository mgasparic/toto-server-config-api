package configurations

import (
	"api/internal/commons"
	"fmt"
	"testing"
)

func TestUpdateConfigurationTable(t *testing.T) {
	type test struct {
		configurationTable                      commons.ConfigurationTable
		expMainSkusForConfigurationRequirements map[commons.ConfigurationRequirement][]commons.MainSku
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
			expMainSkusForConfigurationRequirements: map[commons.ConfigurationRequirement][]commons.MainSku{
				{
					Package:     "com.softinit.iquitos.mainapp",
					CountryCode: "ZZ",
				}: {"default"},
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
				{
					commons.ConfigurationRequirement{
						Package:     "com.softinit.iquitos.mainapp",
						CountryCode: "US",
					},
					commons.ConfigurationChance{
						PercentileMin: 50,
						PercentileMax: 75,
						MainSku:       "rdm_premium_v3_100_trial_7d_yearly",
					},
				},
				{
					commons.ConfigurationRequirement{
						Package:     "com.softinit.iquitos.mainapp",
						CountryCode: "US",
					},
					commons.ConfigurationChance{
						PercentileMin: 75,
						PercentileMax: 100,
						MainSku:       "rdm_premium_v3_150_trial_7d_yearly",
					},
				},
				{
					commons.ConfigurationRequirement{
						Package:     "com.softinit.iquitos.mainapp",
						CountryCode: "ZZ",
					},
					commons.ConfigurationChance{
						PercentileMin: 0,
						PercentileMax: 100,
						MainSku:       "rdm_premium_v3_050_trial_7d_yearly",
					},
				},
			},
			expMainSkusForConfigurationRequirements: map[commons.ConfigurationRequirement][]commons.MainSku{
				{
					Package:     "com.softinit.iquitos.mainapp",
					CountryCode: "US",
				}: {"rdm_premium_v3_020_trial_7d_monthly",
					"rdm_premium_v3_030_trial_7d_monthly",
					"rdm_premium_v3_100_trial_7d_yearly",
					"rdm_premium_v3_150_trial_7d_yearly"},
				{
					Package:     "com.softinit.iquitos.mainapp",
					CountryCode: "ZZ",
				}: {"rdm_premium_v3_050_trial_7d_yearly"},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			UpdateConfigurationTable(tt.configurationTable)
		Requirements:
			for requirement, expMainSkus := range tt.expMainSkusForConfigurationRequirements {
				retMainSku := <-randomizedQueues[requirement]
				for _, expMainSku := range expMainSkus {
					if expMainSku == retMainSku {
						continue Requirements
					}
				}
				t.Errorf("main skus mismatch: expected '%v', received '%s'", expMainSkus, retMainSku)
			}
		})
	}
}

func TestRandomMainSku(t *testing.T) {
	type test struct {
		configurationRequirement commons.ConfigurationRequirement
		expMainSkus              []commons.MainSku
		expOk                    bool
	}

	UpdateConfigurationTable(commons.ConfigurationTable{
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
		{
			commons.ConfigurationRequirement{
				Package:     "com.softinit.iquitos.mainapp",
				CountryCode: "US",
			},
			commons.ConfigurationChance{
				PercentileMin: 50,
				PercentileMax: 75,
				MainSku:       "rdm_premium_v3_100_trial_7d_yearly",
			},
		},
		{
			commons.ConfigurationRequirement{
				Package:     "com.softinit.iquitos.mainapp",
				CountryCode: "US",
			},
			commons.ConfigurationChance{
				PercentileMin: 75,
				PercentileMax: 100,
				MainSku:       "rdm_premium_v3_150_trial_7d_yearly",
			},
		},
		{
			commons.ConfigurationRequirement{
				Package:     "com.softinit.iquitos.mainapp",
				CountryCode: "ZZ",
			},
			commons.ConfigurationChance{
				PercentileMin: 0,
				PercentileMax: 100,
				MainSku:       "rdm_premium_v3_050_trial_7d_yearly",
			},
		},
	})

	tests := []test{
		{
			configurationRequirement: commons.ConfigurationRequirement{
				Package:     "com.softinit.iquitos.newapp",
				CountryCode: "SI",
			},
			expMainSkus: []commons.MainSku{""},
			expOk:       false,
		},
		{
			configurationRequirement: commons.ConfigurationRequirement{
				Package:     "com.softinit.iquitos.mainapp",
				CountryCode: "SI",
			},
			expMainSkus: []commons.MainSku{"rdm_premium_v3_050_trial_7d_yearly"},
			expOk:       true,
		},
		{
			configurationRequirement: commons.ConfigurationRequirement{
				Package:     "com.softinit.iquitos.mainapp",
				CountryCode: "US",
			},
			expMainSkus: []commons.MainSku{"rdm_premium_v3_020_trial_7d_monthly",
				"rdm_premium_v3_030_trial_7d_monthly",
				"rdm_premium_v3_100_trial_7d_yearly",
				"rdm_premium_v3_150_trial_7d_yearly"},
			expOk: true,
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			t.Parallel()
			retMainSku, retOk := RandomMainSku(tt.configurationRequirement)
			expMainSkusContainRetMainSku := false
			for _, expMainSku := range tt.expMainSkus {
				if expMainSku == retMainSku {
					expMainSkusContainRetMainSku = true
					break
				}
			}
			if !expMainSkusContainRetMainSku {
				t.Errorf("main skus mismatch: expected '%v', received '%s'", tt.expMainSkus, retMainSku)
			}
			if tt.expOk != retOk {
				t.Errorf("ok mismatch: expected '%v', received '%v'", tt.expOk, retOk)
			}
		})
	}
}
