package postgres

import (
	"context"
	"pavlov061356/go-masters-homework/final_task/internal/models"
	"reflect"
	"testing"
)

func TestStorage_Users(t *testing.T) {
	newUser := models.User{
		Email:    "test@email.email",
		Password: "test123",
		Username: "test",
	}

	id, err := testDB.NewUser(context.Background(), newUser)
	if err != nil {
		t.Fatalf("failed to create new user: %v", err)
	}

	newUser.ID = id

	gotUser, err := testDB.GetUser(context.Background(), id)

	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if !reflect.DeepEqual(newUser, gotUser) {
		t.Errorf("got %v, want %v", gotUser, newUser)
	}

	gotUser, err = testDB.GetUserByEmail(context.Background(), newUser.Email)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if !reflect.DeepEqual(newUser, gotUser) {
		t.Errorf("got %v, want %v", gotUser, newUser)
	}

	newUser.Password = "test123456789"

	if err = testDB.UpdateUser(context.Background(), newUser); err != nil {
		t.Fatalf("failed to update user: %v", err)
	}

	gotUser, err = testDB.GetUser(context.Background(), id)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if !reflect.DeepEqual(newUser, gotUser) {
		t.Errorf("got %v, want %v", gotUser, newUser)
	}

	if err = testDB.DeleteUser(context.Background(), id); err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}
}
