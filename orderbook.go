package main

import (
	"fmt"
	"sort"
	"time"
)

type Match struct {
	Ask        *Order
	Bid        *Order
	SizeFilled float64
	Price      float64
}

type Order struct {
	Size      float64
	Bid       bool
	Limit     *Limit
	Timestamp int64
}

type Limit struct {
	Price       float64
	Orders      []*Order
	TotalVolume float64
}

type Orderbook struct {
	Asks []*Limit
	Bids []*Limit

	AsksLimits map[float64]*Limit
	BidsLimits map[float64]*Limit
}

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{},
	}
}

func NewOrder(bid bool, size float64) *Order {
	return &Order{
		Bid:       bid,
		Size:      size,
		Timestamp: time.Now().Unix(),
	}
}

func (l *Limit) AddOrderToLimit(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size
}

func (l *Limit) DeleteOrderFromLimit(o *Order) {
	for i := 0; i < len(l.Orders); i++ {
		if l.Orders[i] == o {
			l.Orders = append(l.Orders[:i], l.Orders[i+1:]...)
			break
		}
	}
	o.Limit = nil
	l.TotalVolume -= o.Size

	sort.Slice(l.Orders, func(i, j int) bool {
		return l.Orders[i].Timestamp < l.Orders[j].Timestamp
	})
}

func (o *Order) String() string {
	return fmt.Sprintf("[size: %.2f]", o.Size)
}

// func (l *Limit) String() string {
// 	orderStr := ""
// 	for i, o := range l.Orders {
// 		if i > 0 {
// 			orderStr += ", "
// 		}
// 		orderStr += o.String()
// 	}
// 	return fmt.Sprintf("[price: %.2f | orders: [%s] | totalVolume: %.2f]", l.Price, orderStr, l.TotalVolume)
// }

func NewOrderbook() *Orderbook {
	return &Orderbook{
		Asks:       []*Limit{},
		Bids:       []*Limit{},
		AsksLimits: make(map[float64]*Limit),
		BidsLimits: make(map[float64]*Limit),
	}
}

func (ob *Orderbook) PlaceOrder(price float64, o *Order) []Match {
	// 1. try to match the orders using matching logic
	// 2. add the rest of the order to the book
	if o.Size > 0.0 {
		ob.add(price, o)
	}

	return []Match{}
}

func (ob *Orderbook) add(price float64, o *Order) {
	var limit *Limit

	if o.Bid {
		limit = ob.BidsLimits[price]
	} else {
		limit = ob.AsksLimits[price]
	}

	if limit == nil {
		limit = NewLimit(price)
		if o.Bid {
			ob.Bids = append(ob.Bids, limit)
			ob.BidsLimits[price] = limit
		} else {
			ob.Asks = append(ob.Asks, limit)
			ob.AsksLimits[price] = limit
		}
	}

	if o.Bid {
		sort.Slice(ob.Bids, func(i, j int) bool {
			return ob.Bids[i].Price > ob.Bids[j].Price
		})
	} else {
		sort.Slice(ob.Asks, func(i, j int) bool {
			return ob.Asks[i].Price < ob.Asks[j].Price
		})
	}

	limit.AddOrderToLimit(o)

	sort.Slice(limit.Orders, func(i, j int) bool {
		return limit.Orders[i].Timestamp < limit.Orders[j].Timestamp
	})
}
