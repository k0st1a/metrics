// Package client for create grpc client.
package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/k0st1a/metrics/internal/adapters/api/grpc/protobuf"
	"github.com/k0st1a/metrics/internal/application/agent/config"
	"github.com/k0st1a/metrics/internal/pkg/rawmetric"
	"github.com/rs/zerolog/log"
)

type Client struct {
	Client pb.MetricsClient
}

func New(cfg *config.Config) (*Client, error) {
	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(cfg.ServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("create grpc client error:%w", err)
	}
	defer conn.Close()

	// получаем переменную интерфейсного типа UsersClient,
	// через которую будем отправлять сообщения
	c := pb.NewMetricsClient(conn)

	return &Client{
		Client: c,
	}, nil
}

func (c *Client) DoBatch(r []rawmetric.Info) error {
	log.Printf("DoBatch:%+v", r)

	m := Raw2MetricList(r)

	resp, err := c.Client.StoreMetricList(context.Background(), &pb.StoreMetricListRequest{Metrics: m})
	if err != nil {
		return fmt.Errorf("store metric list error:%w", err)
	}

	log.Printf("resp:%v", resp)

	return nil
}

// Raw2MetricList - преобразование списка метрики из "сырого" формата в "окончательный" формат
// для отправки по grpc.
func Raw2MetricList(r []rawmetric.Info) []*pb.Metric {
	l := make([]*pb.Metric, len(r))
	i := 0
	for _, v := range r {
		m, err := Raw2Metric(v)
		if err != nil {
			log.Error().Err(err).Msg("Raw2Metric error")
			continue
		}
		l[i] = m
		i++
	}

	return l
}

// Raw2Metric - преобразование метрики из "сырого" формате в "окончательный" формат
// для отправки по grpc.
func Raw2Metric(r rawmetric.Info) (*pb.Metric, error) {
	if r.Type == rawmetric.Gauge {
		v, ok := r.Value.(float64)
		if !ok {
			return nil, fmt.Errorf("for gauge(%+v) value type not float64", r)
		}

		return &pb.Metric{
			Id:    r.Name,
			Type:  pb.Type_GAUGE,
			Value: v,
		}, nil
	}

	if r.Type == rawmetric.Counter {
		v, ok := r.Value.(int64)
		if !ok {
			return nil, fmt.Errorf("for counter(%+v) value type not int64", r)
		}

		return &pb.Metric{
			Id:    r.Name,
			Type:  pb.Type_COUNTER,
			Delta: v,
		}, nil
	}

	return nil, fmt.Errorf("unknown metric type(%+v)", r)
}
