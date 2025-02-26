package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/cloudwego/biz-demo/gomall/app/agent/biz/chat"
	"github.com/cloudwego/biz-demo/gomall/app/agent/biz/rpc"
	"github.com/cloudwego/biz-demo/gomall/app/frontend/types"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/agent"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order"
	rpcproduct "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
)

type QueryOrderService struct {
	ctx context.Context
}

// NewQueryOrderService new QueryOrderService
func NewQueryOrderService(ctx context.Context) *QueryOrderService {
	return &QueryOrderService{ctx: ctx}
}

func (s *QueryOrderService) Run(req *agent.QueryOrderReq) (resp *agent.QueryOrderResp, err error) {
	// 第一步：调用订单服务获取订单列表
	listOrderResp, err := rpc.OrderClient.ListOrder(s.ctx, &order.ListOrderReq{UserId: req.UserId})
	if listOrderResp == nil || len(listOrderResp.Orders) == 0 {
		return nil, errors.New("no matched order")
	}

	// 第二步：将订单列表连接商品信息以获取订单详情
	orders := s.getOrderDetail(listOrderResp)

	// 第三步：根据订单详情和用户问题生成与模型对话的消息
	messages := createMessagesFromTemplate(orders, req.Question)
	klog.Info("messages: %+v\n", messages)

	// 第四步：调用对话模型生成对话结果
	cm := chat.NewArkChatModel(s.ctx, nil)
	result := chat.Generate(s.ctx, cm, messages)
	klog.Info("result: %+v\n\n", result.Content)

	// 第五步：解析对话结果，提取订单ID和查询结果
	orderIDs := extractOrderIDsFromResult(result.Content)
	matchedOrders := filterOrdersByIDs(listOrderResp.Orders, orderIDs)
	responseText := extractResponseTextFromResult(result.Content)

	return &agent.QueryOrderResp{
		Orders:   matchedOrders,
		Response: responseText,
	}, nil
}

// extractOrderIDsFromResult 从 result 中提取订单ID列表
func extractOrderIDsFromResult(result string) []string {
	// 按行分割结果
	lines := strings.Split(result, "\n")

	// 遍历每一行，查找以 "订单ID:" 开头的行
	for _, line := range lines {
		if strings.HasPrefix(line, "订单ID:") {
			// 提取 "订单ID:" 后面的内容
			idPart := strings.TrimSpace(strings.TrimPrefix(line, "订单ID:"))
			// 按逗号分隔，得到订单ID列表
			orderIDs := strings.Split(idPart, ",")
			// 去除每个ID前后的空格
			for i, id := range orderIDs {
				orderIDs[i] = strings.TrimSpace(id)
			}
			return orderIDs
		}
	}

	// 如果没有找到以 "订单ID:" 开头的行，返回空列表
	return nil
}

// filterOrdersByIDs 根据订单ID列表从 orders 中过滤出匹配的订单
func filterOrdersByIDs(orders []*order.Order, orderIDs []string) []*order.Order {
	matchedOrders := make([]*order.Order, 0)
	orderIDSet := make(map[string]struct{})
	for _, id := range orderIDs {
		orderIDSet[id] = struct{}{}
	}

	for _, o := range orders {
		if _, exists := orderIDSet[o.OrderId]; exists {
			matchedOrders = append(matchedOrders, o)
		}
	}
	return matchedOrders
}

// extractResponseTextFromResult 从 result 中提取查询结果
func extractResponseTextFromResult(result string) string {
	// 按行分割结果
	lines := strings.Split(result, "\n")

	// 遍历每一行，查找以 "查询结果:" 开头的行
	for _, line := range lines {
		if strings.HasPrefix(line, "查询结果:") {
			// 提取 "查询结果:" 后面的内容
			responseText := strings.TrimSpace(strings.TrimPrefix(line, "查询结果:"))
			return responseText
		}
	}

	// 如果没有找到以 "查询结果:" 开头的行，返回空字符串
	return ""
}

func (s *QueryOrderService) getOrderDetail(listOrderResp *order.ListOrderResp) []*types.Order {
	var orders []*types.Order
	for _, v := range listOrderResp.Orders {
		var items []types.OrderItem
		var total float32
		for _, vv := range v.OrderItems {
			total += vv.Cost
			i := vv.Item
			productResp, err := rpc.ProductClient.GetProduct(s.ctx, &rpcproduct.GetProductReq{Id: i.ProductId})
			if err != nil {
				klog.Error(err)
			}
			if productResp.Product == nil {
				continue
			}
			p := productResp.Product
			items = append(items, types.OrderItem{
				ProductId:   i.ProductId,        // 商品ID
				Qty:         uint32(i.Quantity), // 商品数量
				ProductName: p.Name,             // 商品名称
				Picture:     p.Picture,          // 商品图片
				Cost:        vv.Cost,            // 商品成本
			})
		}
		// 将订单创建时间戳转换为时间对象，第一个参数是秒，第二个参数是纳秒
		timeObj := time.Unix(int64(v.CreatedAt), 0)
		orders = append(orders, &types.Order{
			Cost:  total,
			Items: items,
			// 订单创建日期，格式化为 "2006-01-02 15:04:05"
			CreatedDate: timeObj.Format("2006-01-02 15:04:05"),
			OrderId:     v.OrderId,
			Consignee:   types.Consignee{Email: v.Email},
		})
	}
	return orders
}

func createTemplate() prompt.ChatTemplate {
	// 创建模板，使用 FString 格式
	return prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage(`
你是一个订单查询助手，能够根据用户的问题从订单详情中提取相关信息。订单详情是一个Order结构体的数组，结构定义如下：

type Consignee struct
    Email         string
    StreetAddress string
    City          string
    State         string
    Country       string
    ZipCode       int32

type Order struct 
    Consignee   Consignee
    OrderId     string
    CreatedDate string
    OrderState  string
    Cost        float32
    Items       []OrderItem

type OrderItem struct 
    ProductId   uint32
    ProductName string
    Picture     string
    Qty         uint32
    Cost        float32

你的任务是根据用户的问题查询订单详情，并返回相关订单ID和查询结果。

请注意，你的输出必须严格遵循以下格式：
1. 第一行是查询结果的文本描述。
2. 第二行以"订单ID"开头，列出所有相关订单ID，用逗号分隔。
例如：
查询结果: 订单状态为已发货。
订单ID: 12345,67890
`),
		// 用户消息模板
		schema.UserMessage("订单详情: {orders}, 我的查询: {question}"),
	)
}

func createMessagesFromTemplate(orders []*types.Order, question string) []*schema.Message {
	template := createTemplate()
	jsonData, err := json.Marshal(orders)
	if err != nil {
		klog.Errorf("序列化失败: %v", err)
	}
	ordersJSON := string(jsonData)
	messages, err := template.Format(context.Background(), map[string]any{
		"orders":   ordersJSON,
		"question": question,
	})
	if err != nil {
		klog.Errorf("format template failed: %v\n", err)
	}
	return messages
}
