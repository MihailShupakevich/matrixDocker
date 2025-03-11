package repository

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"matrixDocker/internal/domain"
	"os"
	"testing"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	var err error
	db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestFindUser(t *testing.T) {
	testUser := domain.User{ID: 1, Username: "John Doe", Age: 30}
	db.Create(&testUser)
	repo := NewUserRepository(db)
	ctx := context.Background()
	user, err := repo.FindUser(ctx, 1)
	require.NoError(t, err)
	require.EqualValues(t, testUser, user)
}

func TestFindAllUsers(t *testing.T) {
	testUsers := []domain.User{
		{ID: 1, Username: "Alice", Age: 25},
		{ID: 2, Username: "Bob", Age: 30},
	}
	for _, user := range testUsers {
		db.Create(&user)
	}
	repo := NewUserRepository(db)
	ctx := context.Background()
	users, err := repo.FindAllUsers(ctx)
	require.NoError(t, err)
	require.ElementsMatch(t, testUsers, users)
}

func TestCreateUser(t *testing.T) {
	newUser := domain.User{Username: "New User", Age: 20}
	repo := NewUserRepository(db)
	ctx := context.Background()
	createdUser, err := repo.CreateUser(ctx, newUser)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.Equal(t, newUser.Username, createdUser.Username)
	require.Equal(t, newUser.Age, createdUser.Age)
}

func TestUpdateUser(t *testing.T) {
	testUser := domain.User{ID: 1, Username: "John Doe", Age: 30}
	db.Create(&testUser)
	updatedUser := domain.User{ID: 1, Username: "John Doe", Age: 35}
	repo := NewUserRepository(db)
	ctx := context.Background()
	updated, err := repo.UpdateUser(ctx, testUser.ID, updatedUser)
	require.NoError(t, err)
	require.EqualValues(t, updatedUser, updated)
}

func TestDeleteUser(t *testing.T) {
	testUser := domain.User{ID: 1, Username: "John Doe", Age: 30}
	db.Create(&testUser)
	repo := NewUserRepository(db)
	ctx := context.Background()
	msg, err := repo.DeleteUser(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, "User successfully deleted", msg)
}
