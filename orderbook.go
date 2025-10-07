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

func NewOrderbook() *Orderbook {
	return &Orderbook{
		Asks:       []*Limit{},
		Bids:       []*Limit{},
		AsksLimits: make(map[float64]*Limit),
		BidsLimits: make(map[float64]*Limit),
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

func (l *Limit) fillOrder(o1, o2 *Order) Match {
	var (
		bid        *Order
		ask        *Order
		sizeFilled float64
	)

	if o1.Bid {
		bid = o1
		ask = o2
	} else {
		bid = o2
		ask = o1
	}

	if o1.Size >= o2.Size {
		o1.Size -= o2.Size
		sizeFilled = o2.Size
		o2.Size = 0.0
	} else {
		o2.Size -= o1.Size
		sizeFilled = o1.Size
		o1.Size = 0.0
	}

	return Match{
		Ask:        ask,
		Bid:        bid,
		SizeFilled: sizeFilled,
		Price:      l.Price,
	}

}

func (l *Limit) Fill(o *Order) []Match {
	var (
		matches        []Match
		ordersToDelete []*Order
	)

	for _, order := range l.Orders {
		match := l.fillOrder(o, order)
		matches = append(matches, match)

		l.TotalVolume -= match.SizeFilled

		if order.IsFilled() {
			ordersToDelete = append(ordersToDelete, order)
		}

		if o.IsFilled() {
			break
		}
	}

	for _, order := range ordersToDelete {
		l.DeleteOrderFromLimit(order)
	}

	return matches
}

// func (o *Order) String() string {
// 	return fmt.Sprintf("[size: %.2f]", o.Size)
// }

func (o *Order) IsFilled() bool {
	return o.Size == 0.0
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

func (ob *Orderbook) BidTotalVolume() float64 {
	totalVolume := 0.0

	for i := 0; i < len(ob.Bids); i++ {
		totalVolume += ob.Bids[i].TotalVolume
	}

	return totalVolume
}

func (ob *Orderbook) AskTotalVolume() float64 {
	totalVolume := 0.0

	for i := 0; i < len(ob.Asks); i++ {
		totalVolume += ob.Asks[i].TotalVolume
	}

	return totalVolume
}

func (ob *Orderbook) clearLimit(bid bool, l *Limit) {
	if bid {
		delete(ob.BidsLimits, l.Price)
		for i := 0; i < len(ob.Bids); i++ {
			if ob.Bids[i] == l {
				ob.Bids = append(ob.Bids[:i], ob.Bids[i+1:]...)
			}
		}
	} else {
		delete(ob.AsksLimits, l.Price)
		for i := 0; i < len(ob.Asks); i++ {
			if ob.Asks[i] == l {
				ob.Asks = append(ob.Asks[:i], ob.Asks[i+1:]...)
			}
		}
	}
}

func (ob *Orderbook) PlaceLimitOrder(price float64, o *Order) {
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

func (ob *Orderbook) PlaceMarketOrder(o *Order) []Match {
	matches := []Match{}

	if o.Bid {
		if o.Size > ob.AskTotalVolume() {
			panic(fmt.Errorf("not enough volume [size: %.2f] for market order [size: %.2f]", ob.AskTotalVolume(), o.Size))
		}
		var toClear []*Limit
		for _, limit := range ob.Asks {
			limitMatches := limit.Fill(o)
			matches = append(matches, limitMatches...)

			if len(limit.Orders) == 0 {
				toClear = append(toClear, limit)
			}
			if o.IsFilled() {
				break
			}
		}
		for _, l := range toClear {
			ob.clearLimit(false, l)
		}
	} else {
		if o.Size > ob.BidTotalVolume() {
			panic(fmt.Errorf("not enough volume [size: %.2f] for market order [size: %.2f]", ob.BidTotalVolume(), o.Size))
		}
		var toClear []*Limit
		for _, limit := range ob.Bids {
			limitMatches := limit.Fill(o)
			matches = append(matches, limitMatches...)

			if len(limit.Orders) == 0 {
				toClear = append(toClear, limit)
			}
			if o.IsFilled() {
				break
			}
		}
		for _, l := range toClear {
			ob.clearLimit(true, l)
		}
	}

	return matches
}
