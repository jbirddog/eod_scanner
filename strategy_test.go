package main

import (
	"testing"
	"time"
)

func tc1(sym string,
	d1 time.Time, o1, h1, l1, c1, v1 float64,
	d2 time.Time, o2, h2, l2, c2, v2 float64,
	ml, ms float64,
	r1, r2, r3, r4, r5 float64,
	s float64,
	flags IndicatorFlags,
) *AnalyzedData {
	d := &AnalyzedData{Symbol: sym}

	d.EODData = []*EODData{
		&EODData{sym, d1, o1, h1, l1, c1, v1},
		&EODData{sym, d2, o2, h2, l2, c2, v2},
	}

	i := &d.Indicators
	i.Init()

	i.AvgVolume = 1_000_001.0
	i.AvgClose = 5.01

	i.MACD.Line = ml
	i.MACD.Signal.Value = ms

	setRSIs(i, r1, r2, r3, r4, r5)

	i.SMA20.Value = s
	i.Flags = flags

	return d
}

func tc2(sym string,
	d1 time.Time, o1, h1, l1, c1, v1 float64,
	d2 time.Time, o2, h2, l2, c2, v2 float64,
	d3 time.Time, o3, h3, l3, c3, v3 float64,
	d4 time.Time, o4, h4, l4, c4, v4 float64,
	d5 time.Time, o5, h5, l5, c5, v5 float64,
	ml, ms float64,
	r1, r2, r3, r4, r5 float64,
	s float64,
	flags IndicatorFlags,
) *AnalyzedData {
	d := &AnalyzedData{Symbol: sym}

	d.EODData = []*EODData{
		&EODData{sym, d1, o1, h1, l1, c1, v1},
		&EODData{sym, d2, o2, h2, l2, c2, v2},
		&EODData{sym, d3, o3, h3, l3, c3, v3},
		&EODData{sym, d4, o4, h4, l4, c4, v4},
		&EODData{sym, d5, o5, h5, l5, c5, v5},
	}

	i := &d.Indicators
	i.Init()

	i.AvgVolume = 1_000_001.0
	i.AvgClose = 5.01

	i.MACD.Line = ml
	i.MACD.Signal.Value = ms

	setRSIs(i, r1, r2, r3, r4, r5)

	i.SMA20.Value = s
	i.Flags = flags

	return d
}

type testCaseGen func() *AnalyzedData

func ALGM_08032023() *AnalyzedData {
	return tc2("ALGM",
		Day(2023, 7, 28), 50.23, 50.89, 49.67, 50.84, 1_161_507,
		Day(2023, 7, 31), 51.34, 52.26, 51.16, 51.61, 2_267_876,
		Day(2023, 8, 1), 49.50, 49.64, 45.02, 45.24, 4_354_653,
		Day(2023, 8, 2), 45.00, 45.00, 42.29, 43.22, 3_141_946,
		Day(2023, 8, 3), 42.55, 43.39, 41.92, 43.13, 1_815_615,
		1.53, 2.03,
		56.71, 58.45, 62.38, 64.20, 44.84,
		48.51,
		0b00000000)
}

func APP_01202023() *AnalyzedData {
	return tc1("APP",
		Day(2023, 1, 19), 10.40, 10.62, 10.09, 10.18, 1_856_516.0,
		Day(2023, 1, 20), 10.27, 10.95, 10.13, 10.94, 1_960_167.0,
		-0.2818, -0.4939,
		48.51, 49.47, 47.35, 42.27, 50.61,
		10.31,
		0b00000101)
}

func CRDO_05152023() *AnalyzedData {
	return tc1("CRDO",
		Day(2023, 5, 12), 7.95, 8.06, 7.82, 7.91, 537_400.0,
		Day(2023, 5, 15), 7.99, 8.64, 7.96, 8.61, 1_463_100.0,
		-0.2746, -0.3848,
		35.48, 41.89, 43.11, 42.62, 55.19,
		8.09,
		0b00000001)
}

func GRPN_02062023() *AnalyzedData {
	return tc1("GRPN",
		Day(2023, 2, 3), 9.08, 9.27, 8.73, 8.78, 508_644.0,
		Day(2023, 2, 6), 8.61, 8.69, 8.11, 8.13, 795_608.0,
		0.1932, 0.2142,
		50.19, 55.46, 61.28, 54.26, 47.57,
		8.47,
		0b00000110)
}

func GRPN_06282023() *AnalyzedData {
	return tc1("GRPN",
		Day(2023, 6, 27), 5.21, 5.73, 5.10, 5.56, 1_205_454.0,
		Day(2023, 6, 28), 5.60, 6.00, 5.49, 5.93, 1_152_836.0,
		0.2382, 0.2085,
		57.21, 57.61, 54.05, 57.63, 61.28,
		5.31,
		0b00000011)
}

func RIOT_01052023() *AnalyzedData {
	return tc1("RIOT",
		Day(2023, 1, 4), 3.44, 3.95, 3.38, 3.88, 12_325_300.0,
		Day(2023, 1, 5), 3.86, 4.28, 3.70, 4.22, 14_097_000.0,
		-0.2691, -0.3502,
		31.64, 31.17, 30.84, 46.35, 53.79,
		3.85,
		0b00000011)
}

func RIOT_06272023() *AnalyzedData {
	return tc1("RIOT",
		Day(2023, 6, 26), 11.49, 12.18, 10.72, 10.77, 22_037_414.0,
		Day(2023, 6, 27), 11.07, 11.71, 10.88, 11.65, 25_685_628.0,
		0.0572, -0.0364,
		56.95, 52.84, 55.66, 48.47, 55.09,
		11.01,
		0b00000101)
}

func RIVN_12062022() *AnalyzedData {
	return tc1("RIVN",
		Day(2022, 12, 5), 31.01, 31.24, 29.43, 29.53, 7_560_385.0,
		Day(2022, 12, 6), 29.50, 29.54, 27.43, 27.89, 13_170_726.0,
		-0.8436, -0.7349,
		51.63, 49.97, 49.60, 44.26, 40.09,
		30.97,
		0b000000100)
}

func RIVN_06292023() *AnalyzedData {
	return tc1("RIVN",
		Day(2023, 6, 28), 13.90, 14.87, 13.82, 14.64, 32_296_426.0,
		Day(2023, 6, 29), 14.74, 16.01, 14.61, 16.01, 48_833_726.0,
		0.1716, 0.1408,
		43.35, 42.78, 47.36, 53.12, 61.91,
		14.49,
		0b00000011)
}

func TestStrategies(t *testing.T) {
	testCases := []struct {
		s    string
		posF []testCaseGen
		negF []testCaseGen
	}{
		{
			s: "monthClimb",
			posF: []testCaseGen{
				APP_01202023,
				CRDO_05152023,
				GRPN_06282023,
				RIOT_01052023,
				RIOT_06272023,
				RIVN_06292023,
			},
			negF: []testCaseGen{
				GRPN_02062023,
			},
		},
		{
			s: "monthFall",
			posF: []testCaseGen{
				GRPN_02062023,
				RIVN_12062022,
			},
			negF: []testCaseGen{
				GRPN_06282023,
			},
		},
		{
			s: "fallLevelFall",
			posF: []testCaseGen{
				ALGM_08032023,
			},
			negF: []testCaseGen{
				GRPN_06282023,
			},
		},
	}

	for i, tc := range testCases {
		s, err := StrategyNamed(tc.s)
		if err != nil {
			t.Error(err)
		}

		for j, f := range tc.posF {
			d := f()
			if !s.SignalDetected(d) {
				t.Fatalf("Expected signal '%s' in case %d:%d for %s on %s",
					tc.s,
					i,
					j,
					d.Symbol,
					d.LastDate().Format("01/02/2006"))
			}
		}

		for j, f := range tc.negF {
			d := f()
			if s.SignalDetected(d) {
				t.Fatalf("Unexpected signal '%s' in case %d:%d for %s on %s",
					tc.s,
					i,
					j,
					d.Symbol,
					d.LastDate().Format("01/02/2006"))
			}
		}
	}
}
