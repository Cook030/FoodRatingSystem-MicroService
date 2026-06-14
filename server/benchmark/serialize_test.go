package benchmark

import (
	"encoding/json"
	"testing"
	"time"

	"foodRatingSystem/proto/restaurant"
	"foodRatingSystem/shared/model"

	"google.golang.org/protobuf/proto"
)

// 构造一个典型的餐厅数据
func sampleRestaurant() model.Restaurant {
	return model.Restaurant{
		ID:           42,
		Name:         "老北京炸酱面馆",
		Latitude:     39.9042,
		Longitude:    116.4074,
		AverageScore: 4.5,
		Category:     "中餐",
		CreatedAt:    time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		ReviewCount:  128,
	}
}

func sampleProtoRestaurant() *restaurant.RestaurantMessage {
	return &restaurant.RestaurantMessage{
		Id:          42,
		Name:        "老北京炸酱面馆",
		Latitude:    39.9042,
		Longitude:   116.4074,
		AvgScore:    4.5,
		Category:    "中餐",
		ReviewCount: 128,
	}
}

func BenchmarkJSONSerialize(b *testing.B) {
	r := sampleRestaurant()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(r)
	}
}

func BenchmarkProtobufSerialize(b *testing.B) {
	r := sampleProtoRestaurant()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = proto.Marshal(r)
	}
}

func TestSerializeSize(t *testing.T) {
	r := sampleRestaurant()
	pr := sampleProtoRestaurant()

	jsonBytes, _ := json.Marshal(r)
	protoBytes, _ := proto.Marshal(pr)

	t.Logf("单条数据:")
	t.Logf("  JSON 大小:     %d bytes", len(jsonBytes))
	t.Logf("  Protobuf 大小: %d bytes", len(protoBytes))
	t.Logf("  缩减比例:      %.1f%%", (1-float64(len(protoBytes))/float64(len(jsonBytes)))*100)

	// 测试列表场景（10条）
	var list []model.Restaurant
	var protoList []*restaurant.RestaurantMessage
	for i := 0; i < 10; i++ {
		list = append(list, r)
		protoList = append(protoList, pr)
	}

	jsonListBytes, _ := json.Marshal(list)

	resp := &restaurant.NearbyResponse{
		Restaurants: make([]*restaurant.RestaurantWithDistance, 10),
	}
	for i := 0; i < 10; i++ {
		resp.Restaurants[i] = &restaurant.RestaurantWithDistance{
			Id:          pr.Id,
			Name:        pr.Name,
			Latitude:    pr.Latitude,
			Longitude:   pr.Longitude,
			AvgScore:    pr.AvgScore,
			Category:    pr.Category,
			ReviewCount: pr.ReviewCount,
			Distance:    1.234,
		}
	}
	protoListBytes, _ := proto.Marshal(resp)

	t.Logf("\n10条列表:")
	t.Logf("  JSON 大小:     %d bytes", len(jsonListBytes))
	t.Logf("  Protobuf 大小: %d bytes", len(protoListBytes))
	t.Logf("  缩减比例:      %.1f%%", (1-float64(len(protoListBytes))/float64(len(jsonListBytes)))*100)
}
