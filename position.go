package main

type PositionType int

const (
	Long PositionType = iota
	Short
)

// TODO: stringer when looking into go build for other things
func (t PositionType) String() string {
	switch t {
	case Long:
		return "Long"
	case Short:
		return "Short"
	default:
		panic("Unknown PositionType")
	}
}

type Position struct {
	Type     PositionType
	Shares   int
	Entry    float64
	Capitol  float64
	StopLoss float64
}

func PositionFromAnalyzedData(data *AnalyzedData, risk float64, signalType SignalType) *Position {
	entry := data.LastClose()
	stopLoss := data.SMA20.Value

	var positionType PositionType
	var riskPerShare float64

	if signalType == Buy {
		positionType = Long
		riskPerShare = entry - stopLoss
	} else {
		positionType = Short
		riskPerShare = stopLoss - entry
	}

	shares := int(risk / riskPerShare)
	capitol := float64(shares) * entry

	return &Position{
		Type:     positionType,
		Shares:   shares,
		Entry:    entry,
		Capitol:  capitol,
		StopLoss: stopLoss,
	}
}
