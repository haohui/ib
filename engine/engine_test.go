package engine

import (
	"reflect"
	"testing"
	"time"
)

func (engine *Handle) expect(t *testing.T, ch chan Reply, expected int64) (Reply, error) {
	for {
		select {
		case <-time.After(engine.timeout):
			return nil, timeout()
		case v := <-ch:
			if v.code() == 0 {
				t.Fatalf("don't know message '%v'", v)
			}

			if v.code() == expected {
				return v, nil
			} else {
				// wrong message received
				t.Logf("received message '%v' of type '%v'\n",
					v, reflect.ValueOf(v).Type())
			}
		}
	}

	return nil, nil
}

func TestConnect(t *testing.T) {
	_, err := Make(1)

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}
}

func TestMarketData(t *testing.T) {
	engine, err := Make(2)

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	req1 := &RequestMarketData{
		Contract: Contract{
			Symbol:       "AAPL",
			SecurityType: "STK",
			Exchange:     "SMART",
			Currency:     "USD",
		},
	}

	id := engine.NextRequestId()
	req1.SetId(id)
	ch := make(chan Reply)
	engine.Subscribe(ch, id)
	defer engine.Unsubscribe(id)

	if err := engine.Send(req1); err != nil {
		t.Fatalf("cannot send market data request: %s", err)
	}

	rep1, err := engine.expect(t, ch, mTickPrice)

	if err != nil {
		t.Fatalf("cannot receive market data: %s", err)
	}

	t.Logf("received packet '%v' of type %v\n", rep1, reflect.ValueOf(rep1).Type())

	if err := engine.Send(&CancelMarketData{id}); err != nil {
		t.Fatalf("cannot send cancel request: %s", err)
	}
}

func TestContractDetails(t *testing.T) {
	engine, err := Make(3)

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	req1 := &RequestContractData{
		Symbol:       "AAPL",
		SecurityType: "STK",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	id := engine.NextRequestId()
	req1.SetId(id)
	ch := make(chan Reply)
	engine.Subscribe(ch, id)
	defer engine.Unsubscribe(id)

	if err := engine.Send(req1); err != nil {
		t.Fatalf("cannot send contract data request: %s", err)
	}

	rep1, err := engine.expect(t, ch, mContractData)

	if err != nil {
		t.Fatalf("cannot receive contract details: %s", err)
	}

	t.Logf("received packet '%v' of type %v\n", rep1, reflect.ValueOf(rep1).Type())

	rep2, err := engine.expect(t, ch, mContractDataEnd)

	if err != nil {
		t.Fatalf("cannot receive end of contract details: %s", err)
	}

	t.Logf("received packet '%v' of type %v\n", rep2, reflect.ValueOf(rep2).Type())
}

func TestOptionChainRequest(t *testing.T) {
	engine, err := Make(4)

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	req1 := &RequestContractData{
		Symbol:       "AAPL",
		SecurityType: "OPT",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	id := engine.NextRequestId()
	req1.SetId(id)
	ch := make(chan Reply)
	engine.Subscribe(ch, id)
	defer engine.Unsubscribe(id)

	if err := engine.Send(req1); err != nil {
		t.Fatalf("cannot send contract data request: %s", err)
	}

	rep1, err := engine.expect(t, ch, mContractDataEnd)

	if err != nil {
		t.Fatalf("cannot receive contract details: %s", err)
	}

	t.Logf("received packet '%v' of type %v\n", rep1, reflect.ValueOf(rep1).Type())
}