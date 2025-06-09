package postgres

import (
	"context"
	"pavlov061356/go-masters-homework/final_task/internal/models"
	"reflect"
	"testing"
)

func TestStorage_Services(t *testing.T) {
	services, err := testDB.GetServices(context.Background())
	if err != nil {
		t.Error(err)
	}
	initialServicesLen := len(services)

	newService := models.Service{
		Name:        "test",
		Description: "test",
	}

	id, err := testDB.NewService(context.Background(), newService)
	if err != nil {
		t.Error(err)
	}

	newService.ID = id

	gotService, err := testDB.GetService(context.Background(), newService.ID)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(gotService, newService) {
		t.Errorf("got %v, want %v", gotService, newService)
	}

	services, err = testDB.GetServices(context.Background())
	if err != nil {
		t.Error(err)
	}

	if len(services) != initialServicesLen+1 {
		t.Fatalf("got %d, want %d", len(services), 1)
	}

	for i := 0; i < len(services); i++ {
		if services[i].ID == id {
			if !reflect.DeepEqual(services[i], newService) {
				t.Errorf("got %v, want %v", services[i], newService)
			}
		}
	}

	newService.Description = "test2"
	err = testDB.UpdateService(context.Background(), newService)
	if err != nil {
		t.Error(err)
	}

	gotService, err = testDB.GetService(context.Background(), newService.ID)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(gotService, newService) {
		t.Errorf("got %v, want %v", gotService, newService)
	}

	err = testDB.DeleteService(context.Background(), newService.ID)
	if err != nil {
		t.Error(err)
	}

	err = testDB.RecomputeServicesScore(context.Background())
	if err != nil {
		t.Error(err)
	}

	recomputeTime, err := testDB.GetLastRecomputeTime(context.Background())
	if err != nil {
		t.Error(err)
	}

	if recomputeTime.IsZero() {
		t.Errorf("got %v, want not zero", recomputeTime)
	}

	t.Log(recomputeTime)
}
