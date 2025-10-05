package main

import (
	"fmt"
	"testing"
)

func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	buyOrderA := NewOrder(true, 5)
	buyOrderB := NewOrder(true, 5)

	l.AddOrderToLimit(buyOrderA)
	l.AddOrderToLimit(buyOrderB)

	fmt.Println(l)
}

func TestOrderbook(t *testing.T) {

}
