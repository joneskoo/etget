package energiatili_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/joneskoo/etget/energiatili"
)

func mustTime(t time.Time, err error) time.Time {
	if err != nil {
		panic(err)
	}
	return t
}

func TestModel(t *testing.T) {
	var report energiatili.ConsumptionReport
	err := json.Unmarshal([]byte(sampleJSONData), &report)
	if err != nil {
		t.Fatalf("ERROR parsing JSON structure: %s", err)
	}

	cases := []struct {
		in   int
		want energiatili.Record
	}{
		{0, energiatili.Record{Value: 0, Timestamp: mustTime(time.Parse(time.RFC3339, "2012-08-02T20:00:00Z"))}},
		{1, energiatili.Record{Value: 2.646, Timestamp: mustTime(time.Parse(time.RFC3339, "2014-09-02T21:00:00Z"))}},
	}

	equal := func(a, b energiatili.Record) bool {
		if a.Value != b.Value {
			return false
		}
		if !a.Timestamp.Equal(b.Timestamp) {
			return false
		}
		return true
	}

	points, err := report.Records()
	for _, c := range cases {
		got := points[c.in]
		want := c.want
		if !equal(got, want) {
			t.Errorf("points[%d] = %v, want %v", c.in, got, want)
		}
	}
}

var sampleJSONData = `
{
  "IsValid": true,
  "Message": null,
  "HasTemperatureSeries": true,
  "HasReactivePowerSeries": false,
  "PowerUnit": "kWh",
  "DataInterval": {
    "Duration": "1650.23:00:00",
    "Start": "/Date(1325368800000)/",
    "Stop": "/Date(1468011600000)/",
    "StartValue": "/Date(1325368800000)/",
    "StopValue": "/Date(1468011600000)/",
    "TotalYears": 5,
    "TotalDays": 1650.9583333333333,
    "TotalHours": 1650.9583333333333
  },
  "Hours": {
    "PMax": {
      "Item1": "/Date(1450522800000)/",
      "Item2": 10.434
    },
    "QMax": null,
    "Step": {
      "TimeZoneInfo": {
        "Id": "FLE Standard Time",
        "DisplayName": "(UTC+02:00) Helsinki, Kyiv, Riga, Sofia, Tallinn, Vilnius",
        "StandardName": "FLE Standard Time",
        "DaylightName": "FLE Daylight Time",
        "BaseUtcOffset": "02:00:00",
        "AdjustmentRules": [
          {
            "DateStart": "/Date(-62135596800000)/",
            "DateEnd": "/Date(253402214400000)/",
            "DaylightDelta": "01:00:00",
            "DaylightTransitionStart": {
              "TimeOfDay": "/Date(-62135586000000)/",
              "Month": 3,
              "Week": 5,
              "Day": 1,
              "DayOfWeek": 0,
              "IsFixedDateRule": false
            },
            "DaylightTransitionEnd": {
              "TimeOfDay": "/Date(-62135582400000)/",
              "Month": 10,
              "Week": 5,
              "Day": 1,
              "DayOfWeek": 0,
              "IsFixedDateRule": false
            }
          }
        ],
        "SupportsDaylightSavingTime": true
      },
      "Type": 4,
      "StepLength": "01:00:00",
      "Start": "/Date(1325368800000)/",
      "Stop": "/Date(1468011600000)/",
      "StepCount": 39624
    },
    "Consumptions": [
      {
        "TariffTimeZoneId": 2,
        "TariffTimeZoneName": "Yö",
        "TariffTimeZoneDescription": "Yö",
        "Series": {
          "ReadingCounter": 0,
          "Name": "Yö",
          "Resolution": "Hour",
          "Data": [
            [
              1343948400000,
              0
            ],
            [
              1409702400000,
              2.646
            ],
            [
              1409706000000,
              1.856
            ],
            [
              1409709600000,
              0.547
            ],
            [
              1409713200000,
              0.635
            ],
            [
              1409716800000,
              0.827
            ],
            [
              1409720400000,
              0.358
            ]
          ],
          "Start": "/Date(1325368800000)/",
          "Stop": "/Date(1468011600000)/",
          "Step": {
            "TimeZoneInfo": {
              "Id": "FLE Standard Time",
              "DisplayName": "(UTC+02:00) Helsinki, Kyiv, Riga, Sofia, Tallinn, Vilnius",
              "StandardName": "FLE Standard Time",
              "DaylightName": "FLE Daylight Time",
              "BaseUtcOffset": "02:00:00",
              "AdjustmentRules": [
                {
                  "DateStart": "/Date(-62135596800000)/",
                  "DateEnd": "/Date(253402214400000)/",
                  "DaylightDelta": "01:00:00",
                  "DaylightTransitionStart": {
                    "TimeOfDay": "/Date(-62135586000000)/",
                    "Month": 3,
                    "Week": 5,
                    "Day": 1,
                    "DayOfWeek": 0,
                    "IsFixedDateRule": false
                  },
                  "DaylightTransitionEnd": {
                    "TimeOfDay": "/Date(-62135582400000)/",
                    "Month": 10,
                    "Week": 5,
                    "Day": 1,
                    "DayOfWeek": 0,
                    "IsFixedDateRule": false
                  }
                }
              ],
              "SupportsDaylightSavingTime": true
            },
            "Type": 4,
            "StepLength": "01:00:00",
            "Start": "/Date(1325368800000)/",
            "Stop": "/Date(1468011600000)/",
            "StepCount": 39624
          },
          "DataCount": 39624,
          "Type": "energy2",
          "Unit": "kWh"
        }
      }
    ],
    "Prices": [],
    "Consumption": null,
    "SalesConsumption": null,
    "SalesConsumptionSecondary": null,
    "PeakPower": null,
    "TemperatureCorrectedConsumption": null,
    "DegreeDayValues": null,
    "Temperature": {
      "ReadingCounter": 0,
      "Name": "N/A",
      "Resolution": "Hour",
      "Data": [
        [
          1325376000000,
          -4.5
        ],
        [
          1325379600000,
          -5.1
        ],
        [
          1325383200000,
          -5.9
        ],
        [
          1325386800000,
          -6.3
        ]
      ],
      "Start": "/Date(1325368800000)/",
      "Stop": "/Date(1468011600000)/",
      "Step": {
        "TimeZoneInfo": {
          "Id": "FLE Standard Time",
          "DisplayName": "(UTC+02:00) Helsinki, Kyiv, Riga, Sofia, Tallinn, Vilnius",
          "StandardName": "FLE Standard Time",
          "DaylightName": "FLE Daylight Time",
          "BaseUtcOffset": "02:00:00",
          "AdjustmentRules": [
            {
              "DateStart": "/Date(-62135596800000)/",
              "DateEnd": "/Date(253402214400000)/",
              "DaylightDelta": "01:00:00",
              "DaylightTransitionStart": {
                "TimeOfDay": "/Date(-62135586000000)/",
                "Month": 3,
                "Week": 5,
                "Day": 1,
                "DayOfWeek": 0,
                "IsFixedDateRule": false
              },
              "DaylightTransitionEnd": {
                "TimeOfDay": "/Date(-62135582400000)/",
                "Month": 10,
                "Week": 5,
                "Day": 1,
                "DayOfWeek": 0,
                "IsFixedDateRule": false
              }
            }
          ],
          "SupportsDaylightSavingTime": true
        },
        "Type": 4,
        "StepLength": "01:00:00",
        "Start": "/Date(1325368800000)/",
        "Stop": "/Date(1468011600000)/",
        "StepCount": 39624
      },
      "DataCount": 39624,
      "Type": "temperature",
      "Unit": null
    },
    "ReactivePower": null,
    "ReactivePowerProduction": null,
    "ActivePower": {
      "ReadingCounter": 0,
      "Name": null,
      "Resolution": "Hour",
      "Data": [],
      "Start": "/Date(1325368800000)/",
      "Stop": "/Date(1468011600000)/",
      "Step": {
        "TimeZoneInfo": {
          "Id": "FLE Standard Time",
          "DisplayName": "(UTC+02:00) Helsinki, Kyiv, Riga, Sofia, Tallinn, Vilnius",
          "StandardName": "FLE Standard Time",
          "DaylightName": "FLE Daylight Time",
          "BaseUtcOffset": "02:00:00",
          "AdjustmentRules": [
            {
              "DateStart": "/Date(-62135596800000)/",
              "DateEnd": "/Date(253402214400000)/",
              "DaylightDelta": "01:00:00",
              "DaylightTransitionStart": {
                "TimeOfDay": "/Date(-62135586000000)/",
                "Month": 3,
                "Week": 5,
                "Day": 1,
                "DayOfWeek": 0,
                "IsFixedDateRule": false
              },
              "DaylightTransitionEnd": {
                "TimeOfDay": "/Date(-62135582400000)/",
                "Month": 10,
                "Week": 5,
                "Day": 1,
                "DayOfWeek": 0,
                "IsFixedDateRule": false
              }
            }
          ],
          "SupportsDaylightSavingTime": true
        },
        "Type": 4,
        "StepLength": "01:00:00",
        "Start": "/Date(1325368800000)/",
        "Stop": "/Date(1468011600000)/",
        "StepCount": 39624
      },
      "DataCount": 39624,
      "Type": "activePower",
      "Unit": "kW"
    },
    "ActivePowerProduction": null,
    "HeatingPowerMeterReadingData": null,
    "HeatingWaterMeterReadingData": null,
    "HeatingReadingBasicPriceData": null,
    "HeatingReadingEnergyPriceData": null,
    "WaterFlow": null,
    "WaterFlowInTemperature": null,
    "WaterFlowOutTemperature": null,
    "WaterMeterReadingData": null,
    "Prices": [],
    "Consumption": null,
    "SalesConsumption": null,
    "SalesConsumptionSecondary": null,
    "PeakPower": null,
    "TemperatureCorrectedConsumption": null,
    "DegreeDayValues": {
      "ReadingCounter": 0,
      "Name": null,
      "Resolution": "Day",
      "Data": [
        [
          1325376000000,
          22.9625
        ],
        [
          1325462400000,
          21.554166666666667
        ],
        [
          1325548800000,
          18.5
        ],
        [
          1325635200000,
          18.104166666666668
        ],
        [
          1325721600000,
          19.883333333333333
        ]
      ],
      "Start": "/Date(1325368800000)/",
      "Stop": "/Date(1468011600000)/",
      "Step": {
        "Type": 3,
        "TimeZoneInfo": {
          "Id": "FLE Standard Time",
          "DisplayName": "(UTC+02:00) Helsinki, Kyiv, Riga, Sofia, Tallinn, Vilnius",
          "StandardName": "FLE Standard Time",
          "DaylightName": "FLE Daylight Time",
          "BaseUtcOffset": "02:00:00",
          "AdjustmentRules": [
            {
              "DateStart": "/Date(-62135596800000)/",
              "DateEnd": "/Date(253402214400000)/",
              "DaylightDelta": "01:00:00",
              "DaylightTransitionStart": {
                "TimeOfDay": "/Date(-62135586000000)/",
                "Month": 3,
                "Week": 5,
                "Day": 1,
                "DayOfWeek": 0,
                "IsFixedDateRule": false
              },
              "DaylightTransitionEnd": {
                "TimeOfDay": "/Date(-62135582400000)/",
                "Month": 10,
                "Week": 5,
                "Day": 1,
                "DayOfWeek": 0,
                "IsFixedDateRule": false
              }
            }
          ],
          "SupportsDaylightSavingTime": true
        },
        "Start": "/Date(1325368800000)/",
        "Stop": "/Date(1468011600000)/",
        "StepCount": 1652
      },
      "DataCount": 1652,
      "Type": null,
      "Unit": null
    },
    "ReactivePower": null,
    "ReactivePowerProduction": null,
    "ActivePowerProduction": null,
    "HeatingPowerMeterReadingData": null,
    "HeatingWaterMeterReadingData": null,
    "HeatingReadingBasicPriceData": null,
    "HeatingReadingEnergyPriceData": null,
    "WaterFlow": null,
    "WaterFlowInTemperature": null,
    "WaterFlowOutTemperature": null,
    "WaterMeterReadingData": null,
    "ConsumptionStatuses": {
      "ReadingCounter": 0,
      "Name": "Status",
      "Resolution": "Day",
      "Data": [
        [
          1325376000000,
          255
        ],
        [
          1325462400000,
          255
        ],
        [
          1325548800000,
          255
        ],
        [
          1325635200000,
          255
        ]
      ],
      "Start": "/Date(1325368800000)/",
      "Stop": "/Date(1468011600000)/",
      "Step": {
        "Type": 3,
        "TimeZoneInfo": {
          "Id": "FLE Standard Time",
          "DisplayName": "(UTC+02:00) Helsinki, Kyiv, Riga, Sofia, Tallinn, Vilnius",
          "StandardName": "FLE Standard Time",
          "DaylightName": "FLE Daylight Time",
          "BaseUtcOffset": "02:00:00",
          "AdjustmentRules": [
            {
              "DateStart": "/Date(-62135596800000)/",
              "DateEnd": "/Date(253402214400000)/",
              "DaylightDelta": "01:00:00",
              "DaylightTransitionStart": {
                "TimeOfDay": "/Date(-62135586000000)/",
                "Month": 3,
                "Week": 5,
                "Day": 1,
                "DayOfWeek": 0,
                "IsFixedDateRule": false
              },
              "DaylightTransitionEnd": {
                "TimeOfDay": "/Date(-62135582400000)/",
                "Month": 10,
                "Week": 5,
                "Day": 1,
                "DayOfWeek": 0,
                "IsFixedDateRule": false
              }
            }
          ],
          "SupportsDaylightSavingTime": true
        },
        "Start": "/Date(1325368800000)/",
        "Stop": "/Date(1468011600000)/",
        "StepCount": 1652
      },
      "DataCount": 1652,
      "Type": "status",
      "Unit": "quality"
    },
    "Production": null,
    "OwnAreaComparisionGroupData": null,
    "ForeignAreaComparisionGroupData": null
  },
  "NetworkPriceList": {
    "PriceListName": "Aikasähkö 3X025",
    "ActiveTarificationId": 2,
    "Tarifications": [
      {
        "TarificationId": 2,
        "StartTime": "/Date(1245963600000)/",
        "EndTime": "/Date(1343941200000)/"
      },
      {
        "TarificationId": 2,
        "StartTime": "/Date(1409691600000)/",
        "EndTime": null
      },
      {
        "TarificationId": 2,
        "StartTime": "/Date(1409691600000)/",
        "EndTime": null
      },
      {
        "TarificationId": 2,
        "StartTime": "/Date(1409691600000)/",
        "EndTime": null
      }
    ],
    "SingleTariffBasicPrices": [],
    "SingleEnergyPrices": [],
    "TimeBasedTariffBasicPrices": [
      {
        "ProductComponentTypeCode": "AIKAPERUSMAKSU",
        "StartTime": "/Date(1325368800000)/",
        "EndTime": "/Date(1343941200000)/",
        "PriceWithVat": 10.998333333333333,
        "PriceNoVat": 8.94175
      },
      {
        "ProductComponentTypeCode": "AIKAPERUSMAKSU",
        "StartTime": "/Date(1409691600000)/",
        "EndTime": "/Date(1420063200000)/",
        "PriceWithVat": 11.953333333333333,
        "PriceNoVat": 9.63975
      },
      {
        "ProductComponentTypeCode": "AIKAPERUSMAKSU",
        "StartTime": "/Date(1464728400000)/",
        "EndTime": "/Date(1420063200000)/",
        "PriceWithVat": 15.425833333333333,
        "PriceNoVat": 12.440166666666666
      },
      {
        "ProductComponentTypeCode": "AIKAPERUSMAKSU",
        "StartTime": "/Date(1420063200000)/",
        "EndTime": "/Date(1420668000000)/",
        "PriceWithVat": 11.953333333333333,
        "PriceNoVat": 9.63975
      },
      {
        "ProductComponentTypeCode": "AIKAPERUSMAKSU",
        "StartTime": "/Date(1464728400000)/",
        "EndTime": "/Date(1420668000000)/",
        "PriceWithVat": 15.425833333333333,
        "PriceNoVat": 12.440166666666666
      },
      {
        "ProductComponentTypeCode": "AIKAPERUSMAKSU",
        "StartTime": "/Date(1420668000000)/",
        "EndTime": "/Date(1464728400000)/",
        "PriceWithVat": 11.953333333333333,
        "PriceNoVat": 9.63975
      },
      {
        "ProductComponentTypeCode": "AIKAPERUSMAKSU",
        "StartTime": "/Date(1464728400000)/",
        "EndTime": "/Date(4102437600000)/",
        "PriceWithVat": 15.425833333333333,
        "PriceNoVat": 12.440166666666666
      }
    ],
    "TimeBasedEnergyDayPrices": [
      {
        "ProductComponentTypeCode": "AIKAPÄIVÄ",
        "StartTime": "/Date(1325368800000)/",
        "EndTime": "/Date(1343941200000)/",
        "PriceWithVat": 2.36,
        "PriceNoVat": 1.919
      },
      {
        "ProductComponentTypeCode": "AIKAPÄIVÄ",
        "StartTime": "/Date(1409691600000)/",
        "EndTime": "/Date(1420063200000)/",
        "PriceWithVat": 2.57,
        "PriceNoVat": 2.073
      },
      {
        "ProductComponentTypeCode": "AIKAPÄIVÄ",
        "StartTime": "/Date(1464728400000)/",
        "EndTime": "/Date(1420063200000)/",
        "PriceWithVat": 2.79,
        "PriceNoVat": 2.25
      },
      {
        "ProductComponentTypeCode": "AIKAPÄIVÄ",
        "StartTime": "/Date(1420063200000)/",
        "EndTime": "/Date(1420668000000)/",
        "PriceWithVat": 2.57,
        "PriceNoVat": 2.073
      },
      {
        "ProductComponentTypeCode": "AIKAPÄIVÄ",
        "StartTime": "/Date(1464728400000)/",
        "EndTime": "/Date(1420668000000)/",
        "PriceWithVat": 2.79,
        "PriceNoVat": 2.25
      },
      {
        "ProductComponentTypeCode": "AIKAPÄIVÄ",
        "StartTime": "/Date(1420668000000)/",
        "EndTime": "/Date(1464728400000)/",
        "PriceWithVat": 2.57,
        "PriceNoVat": 2.073
      },
      {
        "ProductComponentTypeCode": "AIKAPÄIVÄ",
        "StartTime": "/Date(1464728400000)/",
        "EndTime": "/Date(4102437600000)/",
        "PriceWithVat": 2.79,
        "PriceNoVat": 2.25
      }
    ],
    "TimeBasedEnergyNightPrices": [
      {
        "ProductComponentTypeCode": "AIKAYÖ",
        "StartTime": "/Date(1325368800000)/",
        "EndTime": "/Date(1343941200000)/",
        "PriceWithVat": 1.45,
        "PriceNoVat": 1.179
      },
      {
        "ProductComponentTypeCode": "AIKAYÖ",
        "StartTime": "/Date(1409691600000)/",
        "EndTime": "/Date(1420063200000)/",
        "PriceWithVat": 1.57,
        "PriceNoVat": 1.266
      },
      {
        "ProductComponentTypeCode": "AIKAYÖ",
        "StartTime": "/Date(1464728400000)/",
        "EndTime": "/Date(1420063200000)/",
        "PriceWithVat": 1.71,
        "PriceNoVat": 1.379
      },
      {
        "ProductComponentTypeCode": "AIKAYÖ",
        "StartTime": "/Date(1420063200000)/",
        "EndTime": "/Date(1420668000000)/",
        "PriceWithVat": 1.57,
        "PriceNoVat": 1.266
      },
      {
        "ProductComponentTypeCode": "AIKAYÖ",
        "StartTime": "/Date(1464728400000)/",
        "EndTime": "/Date(1420668000000)/",
        "PriceWithVat": 1.71,
        "PriceNoVat": 1.379
      },
      {
        "ProductComponentTypeCode": "AIKAYÖ",
        "StartTime": "/Date(1420668000000)/",
        "EndTime": "/Date(1464728400000)/",
        "PriceWithVat": 1.57,
        "PriceNoVat": 1.266
      },
      {
        "ProductComponentTypeCode": "AIKAYÖ",
        "StartTime": "/Date(1464728400000)/",
        "EndTime": "/Date(4102437600000)/",
        "PriceWithVat": 1.71,
        "PriceNoVat": 1.379
      }
    ],
    "SeasonTariffBasicPrices": [],
    "SeasonEnergyWinterdayPrices": [],
    "SeasonEnergyOtherPrices": [],
    "PowerTariffBasicPrices": [],
    "PowerEnergyDayPrices": [],
    "PowerEnergyOtherTimePrices": [],
    "HeatingBasicPrices": [],
    "HeatingEnergyPrices": [],
    "GasBasicPrices": [],
    "GasEnergyPrices": [],
    "CleanWaterBasicPrices": [],
    "WasteWaterBasicPrices": [],
    "CleanWaterUsagePrices": [],
    "WasteWaterUsagePrices": [],
    "RainWaterBasicPrices": [],
    "RainWaterUsagePrices": [],
    "PowerEnergyWinterPrices": [],
    "PeakPowerPrices": [],
    "ReactivePowerPrices": [],
    "IsSpotPriceList": false,
    "SpotPrices": null
  },
  "EnergyTaxes": {
    "Code": "ELECTRICITYTAXCLASS1",
    "Name": "Sähkön veroluokka 1",
    "Taxes": [
      {
        "TotalTaxWithVAT": 2.09469,
        "TotalTaxNoVAT": 1.703,
        "StartDate": "/Date(1293840000000)/",
        "EndDate": "/Date(1356912000000)/"
      },
      {
        "TotalTaxWithVAT": 1.07726,
        "TotalTaxNoVAT": 0.883,
        "StartDate": "/Date(1199145600000)/",
        "EndDate": "/Date(1293753600000)/"
      },
      {
        "TotalTaxWithVAT": 2.11172,
        "TotalTaxNoVAT": 1.703,
        "StartDate": "/Date(1356998400000)/",
        "EndDate": "/Date(1388448000000)/"
      },
      {
        "TotalTaxWithVAT": 2.35972,
        "TotalTaxNoVAT": 1.903,
        "StartDate": "/Date(1388534400000)/",
        "EndDate": "/Date(1419984000000)/"
      },
      {
        "TotalTaxWithVAT": 2.79372,
        "TotalTaxNoVAT": 2.253,
        "StartDate": "/Date(1420070400000)/",
        "EndDate": null
      }
    ]
  },
  "Vat": [
    {
      "Tax": 0.22,
      "Start": "/Date(-62135596800000)/",
      "Stop": "/Date(1277931600000)/"
    },
    {
      "Tax": 0.23,
      "Start": "/Date(1277931600000)/",
      "Stop": "/Date(1356991200000)/"
    },
    {
      "Tax": 0.24,
      "Start": "/Date(1356991200000)/",
      "Stop": "/Date(4102437600000)/"
    }
  ],
  "EnergyConsumptionCorrections": [],
  "UseReadings": false,
  "HourlyDataAvailable": true,
  "UtilityType": 1,
  "CustomerType": 1,
  "LowTemperature": -35,
  "HighTemperature": 35,
  "HasProviderHourlyValuesSeparateLoadingSupport": false,
  "HasHourlyValues": true,
  "FullConsumptionLoaded": true,
  "MeteringPoint": null,
  "PriceLimitPriceList": null,
  "PriceGuaranteedPriceList": null
}
`
