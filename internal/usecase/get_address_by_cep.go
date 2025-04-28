package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/Ma-Leal/weather/internal/dto"
	"github.com/Ma-Leal/weather/internal/entity"
)

type GetAddressByCEP struct {
	Url    string
	client *http.Client
	tracer trace.Tracer
}

func NewGetAddressByCEP(url string) *GetAddressByCEP {
	return &GetAddressByCEP{
		Url: url,
		client: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
		tracer: otel.Tracer("GetAddressByCEP"),
	}
}

func (uc *GetAddressByCEP) Execute(ctx context.Context, cep string) (*entity.Address, error) {
	ctx, span := uc.tracer.Start(ctx, "GetAddressByCEP.Execute")
	defer span.End()

	c, err := entity.NewCEP(cep)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid CEP format")
		return nil, err
	}

	fullURL := fmt.Sprintf("%s/%s/json/", uc.Url, c.Number)
	span.SetAttributes(
		attribute.String("http.url", fullURL),
		attribute.String("cep", c.Number),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create request")
		return nil, err
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := uc.client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "HTTP request failed")
		return nil, err
	}
	defer resp.Body.Close()

	span.SetAttributes(
		attribute.Int("http.status_code", resp.StatusCode),
	)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to read response body")
		return nil, err
	}

	var r dto.ViaCepResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to unmarshal response")
		return nil, errors.New("unmarshal failed")
	}

	if r.City == "" {
		err := entity.ErrCEPNotFound
		span.RecordError(err)
		span.SetStatus(codes.Error, "CEP not found")
		return nil, err
	}

	add := entity.NewAddress(c, r.City)

	span.SetAttributes(
		attribute.String("address.city", add.City),
	)

	return add, nil
}
