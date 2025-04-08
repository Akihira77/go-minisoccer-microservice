package seeders

import (
	"user-service/constants"
	"user-service/domain/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RunUserSeeder(db *gorm.DB) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	users := []models.User{
		{
			UUID:        uuid.New(),
			Name:        "Administrator",
			Username:    "admin",
			Password:    string(hashPass),
			PhoneNumber: "081285942567",
			Email:       "admin@gmail.com",
			RoleID:      constants.Admin,
		},
	}

	logrus.Info("Seeder user start")
	for _, user := range users {
		err := db.FirstOrCreate(&user, models.User{Username: user.Username}).Error
		if err != nil {
			logrus.Errorf("failed to seed user: %v", err)
			panic(err)
		}

		logrus.Info("user %s successfully seeded", user.ID)
	}
	logrus.Info("Seeder user finish")
}
