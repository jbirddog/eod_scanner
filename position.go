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

type Position interface {
	Type() PositionType
	Shares() int
	Entry() float64
	Capital() float64
	StopLoss() float64
}

type DefaultPosition struct {
	_type     PositionType
	_shares   int
	_entry    float64
	_capital  float64
	_stopLoss float64
}

func (p *DefaultPosition) Type() PositionType {
	return p._type
}

func (p *DefaultPosition) Shares() int {
	return p._shares
}

func (p *DefaultPosition) Entry() float64 {
	return p._entry
}

func (p *DefaultPosition) Capital() float64 {
	return p._capital
}

func (p *DefaultPosition) StopLoss() float64 {
	return p._stopLoss
}

func PositionFromAnalyzedData(data *AnalyzedData, risk float64, signalType SignalType) Position {
	entry := data.LastClose()
	stopLoss := data.Indicators.SMA20.Value

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
	capital := float64(shares) * entry

	return &DefaultPosition{
		_type:     positionType,
		_shares:   shares,
		_entry:    entry,
		_capital:  capital,
		_stopLoss: stopLoss,
	}
}
