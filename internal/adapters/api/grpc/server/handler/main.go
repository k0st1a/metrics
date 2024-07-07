// Package handler contains handler for grpc server.
package handler

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/k0st1a/metrics/internal/adapters/api/grpc/protobuf"
	"github.com/k0st1a/metrics/internal/ports"
	"github.com/rs/zerolog/log"
)

type MetricsServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedMetricsServer
	Storage ports.Storage
}

func (s *MetricsServer) GetMetric(ctx context.Context, r *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	log.Printf("GetMetric, Id:%v, Type:%v", r.Id, r.Type)

	switch r.Type {
	case pb.Type_GAUGE:
		v, err := s.Storage.GetGauge(ctx, r.Id)

		switch {
		case errors.Is(err, ports.ErrMetricsNoGauge):
			return nil, status.Errorf(codes.NotFound, "gauge %v not found", r.Id)
		case err != nil:
			return nil, status.Errorf(codes.Internal, "for gauge %v error:%v", r.Id, err)
		case v == nil:
			return nil, status.Errorf(codes.Internal, "no value for gauge %v", r.Id)
		default:
			return &pb.GetMetricResponse{
				Metric: &pb.Metric{
					Id:    r.Id,
					Type:  pb.Type_GAUGE,
					Value: *v,
				},
			}, nil
		}
	case pb.Type_COUNTER:
		v, err := s.Storage.GetCounter(ctx, r.Id)

		switch {
		case errors.Is(err, ports.ErrMetricsNoCounter):
			return nil, status.Errorf(codes.NotFound, "counter %v not found", r.Id)
		case err != nil:
			return nil, status.Errorf(codes.Internal, "for counter %v error:%v", r.Id, err)
		case v == nil:
			return nil, status.Errorf(codes.Internal, "no value for counter %v", r.Id)
		default:
			return &pb.GetMetricResponse{
				Metric: &pb.Metric{
					Id:    r.Id,
					Type:  pb.Type_GAUGE,
					Delta: *v,
				},
			}, nil
		}
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unknown metric type:%v", r.Type)
	}
}
func (s *MetricsServer) StoreMetric(ctx context.Context, r *pb.StoreMetricRequest) (*pb.StoreMetricResponse, error) {
	switch r.Metric.Type {
	case pb.Type_GAUGE:
		log.Printf("Store gauge, name(%v), value(%v)", r.Metric.Id, r.Metric.Value)
		err = h.retry.Retry(ctx, retry.IsConnectionException, func() error {
			return h.storage.StoreGauge(ctx, m.Metric.Id, r.Metric.Value)
		})
		if err != nil {
			log.Error().Err(err).Msg("h.storage.StorageGauge error")
			return nil, status.Errorf(codes.Internal, "for gauge %v error:%v", r.Metric.Id, err)
		}
	case pb.Type_COUNTER:
		log.Printf("Store counter, name(%v), value(%v)", r.Metric.Id, r.Metric.Delta)
		err = h.retry.Retry(ctx, retry.IsConnectionException, func() error {
			//nolint // Не за чем оборачивать ошибку
			return h.storage.StoreCounter(ctx, m.ID, *m.Delta)
		})
		if err != nil {
			log.Error().Err(err).Msg("h.storage.StoreCounter error")
			return nil, status.Errorf(codes.Internal, "for counter %v error:%v", r.Metric.Id, err)
		}
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unknown metric type:%v", r.Type)
	}

	return &pb.StoreMetricResponse{}, nil
}
func (s *MetricsServer) StoreMetricList(ctx context.Context, r *pb.StoreMetricListRequest) (*pb.StoreMetricListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreMetricList not implemented")
}
