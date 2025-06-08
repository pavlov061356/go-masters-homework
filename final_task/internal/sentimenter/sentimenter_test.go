package sentimenter

import (
	"context"
	"pavlov061356/go-masters-homework/final_task/internal/config"
	"pavlov061356/go-masters-homework/final_task/internal/models"
	"testing"
)

func TestSentimenter_GetReviewSentiment(t *testing.T) {
	type fields struct {
		cfg config.Sentimenter
	}
	type args struct {
		ctx    context.Context
		review models.Review
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "test",
			fields: fields{
				cfg: config.Sentimenter{
					Model:   "llama2:7b",
					Addr:    "http://localhost:11434",
					Timeout: 10,
				},
			},
			args: args{
				ctx: context.Background(),
				review: models.Review{
					ID:    1,
					Score: 1,
				},
			},
			want: 3,
		},
		{
			name: "positive",
			fields: fields{
				cfg: config.Sentimenter{
					Model:   "mistral",
					Addr:    "http://localhost:11434",
					Timeout: 10,
				},
			},
			args: args{
				ctx: context.Background(),
				review: models.Review{
					ID:      1,
					Content: "Товар отличный. Заказывал уже несколько раз, каждый раз очень доволен. Рекомендую!",
					Score:   1,
				},
			},
			want: 1,
		},
		{
			name: "negative",
			fields: fields{
				cfg: config.Sentimenter{
					Model:   "mistral",
					Addr:    "http://localhost:11434",
					Timeout: 10,
				},
			},
			args: args{
				ctx: context.Background(),
				review: models.Review{
					ID:      1,
					Content: "Товар так себе, пришёл повреждённый. Не советую!",
					Score:   1,
				},
			},
			want: 3,
		},
		{
			name: "neutral",
			fields: fields{
				cfg: config.Sentimenter{
					Model:   "mistral",
					Addr:    "http://localhost:11434",
					Timeout: 10,
				},
			},
			args: args{
				ctx: context.Background(),
				review: models.Review{
					ID:      1,
					Content: "Качество среднее, ничего хорошего сказать нельзя, но и ничего плохого",
					Score:   1,
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(tt.fields.cfg)
			if err != nil {
				t.Errorf("New() error = %v", err)
				return
			}
			got, err := s.GetReviewSentiment(tt.args.ctx, tt.args.review)
			if err != nil {
				t.Errorf("Sentimenter.GetReviewSentiment() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("Sentimenter.GetReviewSentiment() = %v, want %v", got, tt.want)
			}
		})
	}
}
