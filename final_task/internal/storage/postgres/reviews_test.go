package postgres

import (
	"context"
	"pavlov061356/go-masters-homework/final_task/internal/models"
	"reflect"
	"testing"
)

func TestStorage_Reviews(t *testing.T) {
	service := models.Service{
		Name:        "test",
		Description: "test",
	}
	id, err := testDB.NewService(context.Background(), service)
	if err != nil {
		t.Fatal(err)
	}
	service.ID = id

	user := models.User{
		Email:    "test@mail.mail",
		Password: "test",
		Username: "test",
	}
	userID, err := testDB.NewUser(context.Background(), user)
	if err != nil {
		t.Fatal(err)
	}
	user.ID = userID

	newReview := models.Review{
		Content:    "test",
		Sentiment:  0,
		ReviewerID: user.ID,
		Score:      5,
		ServiceID:  service.ID,
	}

	id, err = testDB.NewReview(context.Background(), newReview)
	if err != nil {
		t.Fatal(err)
	}
	newReview.ID = id

	review, err := testDB.GetReview(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}
	if review.ID != id {
		t.Errorf("wrong id, got: %d, want: %d", review.ID, id)
	}

	newReview.CreatedAt = review.CreatedAt

	reviews, err := testDB.GetReviewsByService(context.Background(), service.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(reviews) != 1 {
		t.Errorf("wrong number of reviews, got: %d, want: %d", len(reviews), 1)
	}

	if !reflect.DeepEqual(newReview, reviews[0]) {
		t.Errorf("wrong review, got: %v, want: %v", reviews[0], newReview)
	}

	reviews, err = testDB.GetReviewsByUser(context.Background(), userID)
	if err != nil {
		t.Fatal(err)
	}

	if len(reviews) != 1 {
		t.Errorf("wrong number of reviews, got: %d, want: %d", len(reviews), 1)
	}

	if !reflect.DeepEqual(newReview, reviews[0]) {
		t.Errorf("wrong review, got: %v, want: %v", reviews[0], newReview)
	}

	newReview.Content = "new"
	err = testDB.UpdateReview(context.Background(), newReview)
	if err != nil {
		t.Fatal(err)
	}

	review, err = testDB.GetReview(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}
	if review.Content != "new" {
		t.Errorf("wrong review, got: %v, want: %v", review.Content, "new")
	}

	for i := range reviews {
		reviews[i].Sentiment = i%2 + 1
	}

	err = testDB.BatchUpdateReviewsSentiment(context.Background(), reviews)
	if err != nil {
		t.Fatal(err)
	}

	err = testDB.DeleteReview(context.Background(), newReview)
	if err != nil {
		t.Fatal(err)
	}
}

func getAvgScore(reviews []models.Review) float64 {
	var sum float64
	for _, review := range reviews {
		sum += float64(review.Score)
	}
	return sum / float64(len(reviews))
}

func TestStorage_AvgScore(t *testing.T) {
	service := models.Service{
		Name:        "test",
		Description: "test",
	}
	id, err := testDB.NewService(context.Background(), service)
	if err != nil {
		t.Fatal(err)
	}
	service.ID = id

	user := models.User{
		Email:    "test_mail@mail.mail",
		Password: "test",
		Username: "test",
	}
	userID, err := testDB.NewUser(context.Background(), user)
	if err != nil {
		t.Fatal(err)
	}
	user.ID = userID

	reviews := []models.Review{
		{
			Content:    "test",
			Sentiment:  0,
			ReviewerID: user.ID,
			Score:      5,
			ServiceID:  service.ID,
		},
		{
			Content:    "test",
			Sentiment:  0,
			ReviewerID: user.ID,
			Score:      4,
			ServiceID:  service.ID,
		},
	}

	for i, review := range reviews {
		id, err := testDB.NewReview(context.Background(), review)
		if err != nil {
			t.Fatal(err)
		}
		reviews[i].ID = id
	}

	service, err = testDB.GetService(context.Background(), service.ID)
	if err != nil {
		t.Fatal(err)
	}
	if service.AvgScore != getAvgScore(reviews) {
		t.Errorf("wrong avg score, got: %f, want: %f", service.AvgScore, getAvgScore(reviews))
	}

	reviews[0].Score = 1
	err = testDB.UpdateReview(context.Background(), reviews[0])
	if err != nil {
		t.Fatal(err)
	}

	service, err = testDB.GetService(context.Background(), service.ID)
	if err != nil {
		t.Fatal(err)
	}
	if service.AvgScore != getAvgScore(reviews) {
		t.Errorf("wrong avg score, got: %f, want: %f", service.AvgScore, getAvgScore(reviews))
	}

	testDB.DeleteReview(context.Background(), reviews[0])

	service, err = testDB.GetService(context.Background(), service.ID)
	if err != nil {
		t.Fatal(err)
	}

	if service.AvgScore != getAvgScore(reviews[1:]) {
		t.Errorf("wrong avg score, got: %f, want: %f", service.AvgScore, getAvgScore(reviews[1:]))
	}
}
