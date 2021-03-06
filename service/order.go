package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	
	"goshop/admin-api/pkg/grpc/gclient"
	"goshop/admin-api/pkg/utils"
	
	"github.com/shinmigo/pb/memberpb"
	
	"github.com/gin-gonic/gin"
	"github.com/shinmigo/pb/orderpb"
)

type Order struct {
	*gin.Context
}

type ListOrderRes struct {
	Total  uint64         `json:"total"`
	Orders []*orderDetail `json:"orders"`
}

type orderDetail struct {
	orderpb.OrderDetail
	Member *memberpb.MemberDetail `json:"member"`
}

func NewOrder(ctx *gin.Context) *Order {
	return &Order{Context: ctx}
}

func (o *Order) Index(req *orderpb.ListOrderReq) (*ListOrderRes, error) {
	var (
		orderList  *orderpb.ListOrderRes
		err        error
		orders     = make([]*orderDetail, 0, req.PageSize)
		jsonBytes  []byte
		memberIds  []uint64
		ctx        context.Context
		cancel     context.CancelFunc
		memberMaps = make(map[uint64]*memberpb.MemberDetail)
	)
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	if orderList, err = gclient.OrderClient.GetOrderList(ctx, req); err != nil {
		return nil, err
	}
	cancel()
	
	if orderList.Total == 0 {
		return &ListOrderRes{
			Total:  0,
			Orders: nil,
		}, nil
	}
	
	for _, ord := range orderList.Orders {
		memberIds = append(memberIds, ord.MemberId)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	var memberList *memberpb.ListMemberRes
	if memberList, err = gclient.Member.GetMemberList(ctx, &memberpb.GetMemberReq{
		Page:     1,
		PageSize: uint64(len(memberIds)),
	}); err == nil {
		for _, memberDetail := range memberList.Members {
			memberMaps[memberDetail.MemberId] = memberDetail
		}
	}
	cancel()
	
	fmt.Println(orderList.Orders)
	if jsonBytes, err = json.Marshal(orderList.Orders); err != nil {
		return nil, err
	}
	json.Unmarshal(jsonBytes, &orders)
	fmt.Println(orders)
	for _, ord := range orders {
		ord.Member = memberMaps[ord.MemberId]
	}
	
	return &ListOrderRes{
		Total:  orderList.Total,
		Orders: orders,
	}, err
}

func (o *Order) Status(storeId uint64) (*orderpb.ListOrderStatusRes, error) {
	var (
		listOrderStatusRes *orderpb.ListOrderStatusRes
		currentStatus      []uint64
		statuses           = []uint64{
			1, 2, 3, 4, 5, 6, 7,
		}
		err error
	)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	if listOrderStatusRes, err = gclient.OrderClient.GetOrderStatus(ctx, &orderpb.GetOrderStatusReq{
		StoreId: storeId,
	}); err != nil {
		return nil, err
	}
	cancel()
	
	var total uint64
	for _, statistics := range listOrderStatusRes.OrderStatistics {
		currentStatus = append(currentStatus, statistics.OrderStatus)
		total = total + statistics.Count
	}
	
	for _, status := range statuses {
		if utils.InSliceUint64(status, currentStatus) == false {
			listOrderStatusRes.OrderStatistics = append(listOrderStatusRes.OrderStatistics, &orderpb.ListOrderStatusRes_OrderStatistics{
				OrderStatus: status,
				Count:       0,
			})
		}
	}
	
	// 加个全部
	listOrderStatusRes.OrderStatistics = append(listOrderStatusRes.OrderStatistics, &orderpb.ListOrderStatusRes_OrderStatistics{
		OrderStatus: 0,
		Count:       total,
	})
	
	return listOrderStatusRes, nil
}
