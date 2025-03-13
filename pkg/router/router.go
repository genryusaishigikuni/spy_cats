package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/genryusaishigikuni/spy_cats/internal/cat"
	"github.com/genryusaishigikuni/spy_cats/internal/mission"
	"github.com/genryusaishigikuni/spy_cats/internal/note"
	"github.com/genryusaishigikuni/spy_cats/internal/target"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 1) Repositories
	catRepo := cat.NewRepository(db)
	missionRepo := mission.NewRepository(db)
	targetRepo := target.NewRepository(db)
	noteRepo := note.NewRepository(db)

	// 2) Services
	catService := cat.NewService(catRepo)
	// Pass *all* required repos to mission.NewService
	missionService := mission.NewService(missionRepo, catRepo, targetRepo)
	targetService := target.NewService(targetRepo)
	// Pass the note repo + target repo to note.NewService
	noteService := note.NewService(noteRepo, targetRepo)

	// 3) Handlers
	catHandler := cat.NewHandler(catService)
	missionHandler := mission.NewHandler(missionService)
	targetHandler := target.NewHandler(targetService)
	noteHandler := note.NewHandler(noteService)

	// 4) Register routes
	catHandler.RegisterRoutes(r)
	missionHandler.RegisterRoutes(r)
	targetHandler.RegisterRoutes(r)
	noteHandler.RegisterRoutes(r)

	return r
}
