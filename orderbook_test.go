package main

import (
	"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}
}

func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	buyOrderA := NewOrder(true, 5)
	buyOrderB := NewOrder(true, 10)
	buyOrderC := NewOrder(true, 15)

	l.AddOrderToLimit(buyOrderA)
	l.AddOrderToLimit(buyOrderB)
	l.AddOrderToLimit(buyOrderC)

	l.DeleteOrderFromLimit(buyOrderB)

	fmt.Println(l)
}

func TestPlaceLimitOrder(t *testing.T) {
	ob := NewOrderbook()

	sellOrderA := NewOrder(false, 10)
	sellOrderB := NewOrder(false, 5)

	ob.PlaceLimitOrder(10_000, sellOrderA)
	ob.PlaceLimitOrder(9_000, sellOrderB)

	assert(t, len(ob.Asks), 2)
}

func TestPlaceMarketOrderBuyBulk(t *testing.T) {
	ob := NewOrderbook()

	sellOrderA := NewOrder(false, 20)
	sellOrderB := NewOrder(false, 5)
	sellOrderC := NewOrder(false, 40)
	ob.PlaceLimitOrder(10_000, sellOrderA)
	ob.PlaceLimitOrder(5_000, sellOrderB)
	ob.PlaceLimitOrder(20_000, sellOrderC)

	buyOrderA := NewOrder(true, 10)
	matches := ob.PlaceMarketOrder(buyOrderA)

	assert(t, len(matches), 2)
	assert(t, len(ob.Asks), 2)
	assert(t, ob.AskTotalVolume(), 55.0)
	assert(t, matches[0].Ask, sellOrderB)
	assert(t, matches[0].Bid, buyOrderA)
	assert(t, buyOrderA.IsFilled(), true)

	fmt.Printf("%+v", matches)
}

func TestPlaceMarketOrderSellBulk(t *testing.T) {
	ob := NewOrderbook()

	buyOrderA := NewOrder(true, 5)
	buyOrderB := NewOrder(true, 8)
	buyOrderC := NewOrder(true, 10)
	buyOrderD := NewOrder(true, 1)
	ob.PlaceLimitOrder(10_000, buyOrderA)
	ob.PlaceLimitOrder(9_000, buyOrderB)
	ob.PlaceLimitOrder(5_000, buyOrderC)
	ob.PlaceLimitOrder(5_000, buyOrderD)

	assert(t, ob.BidTotalVolume(), 24.0)

	sellOrderA := NewOrder(false, 24)
	matches := ob.PlaceMarketOrder(sellOrderA)

	assert(t, len(matches), 4)
	assert(t, len(ob.Bids), 0)

	fmt.Printf("%+v", matches)
}
