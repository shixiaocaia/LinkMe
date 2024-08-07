// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"LinkMe/internal/api"
	"LinkMe/internal/domain/events/email"
	"LinkMe/internal/domain/events/post"
	"LinkMe/internal/domain/events/sms"
	"LinkMe/internal/repository"
	"LinkMe/internal/repository/cache"
	"LinkMe/internal/repository/dao"
	"LinkMe/internal/service"
	"LinkMe/ioc"
	"LinkMe/utils/jwt"
)

import (
	_ "github.com/google/wire"
)

// Injectors from wire.go:

func InitWebServer() *Cmd {
	db := ioc.InitDB()
	node := ioc.InitializeSnowflakeNode()
	logger := ioc.InitLogger()
	enforcer := ioc.InitCasbin(db)
	userDAO := dao.NewUserDAO(db, node, logger, enforcer)
	cmdable := ioc.InitRedis()
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache, logger)
	userService := service.NewUserService(userRepository, logger)
	handler := jwt.NewJWTHandler(cmdable)
	client := ioc.InitSaramaClient()
	syncProducer := ioc.InitSyncProducer(client)
	producer := sms.NewSaramaSyncProducer(syncProducer, logger)
	emailProducer := email.NewSaramaSyncProducer(syncProducer, logger)
	userHandler := api.NewUserHandler(userService, handler, logger, producer, emailProducer, enforcer)
	mongoClient := ioc.InitMongoDB()
	postDAO := dao.NewPostDAO(db, logger, mongoClient)
	postCache := cache.NewPostCache(cmdable, logger)
	postRepository := repository.NewPostRepository(postDAO, logger, postCache)
	interactiveDAO := dao.NewInteractiveDAO(db, logger)
	interactiveCache := cache.NewInteractiveCache(cmdable)
	interactiveRepository := repository.NewInteractiveRepository(interactiveDAO, logger, interactiveCache)
	interactiveService := service.NewInteractiveService(interactiveRepository, logger)
	checkDAO := dao.NewCheckDAO(db, logger)
	checkRepository := repository.NewCheckRepository(checkDAO, logger)
	historyCache := cache.NewHistoryCache(logger, cmdable)
	historyRepository := repository.NewHistoryRepository(logger, historyCache)
	activityDAO := dao.NewActivityDAO(db, logger)
	activityRepository := repository.NewActivityRepository(activityDAO)
	checkService := service.NewCheckService(checkRepository, postRepository, historyRepository, logger, activityRepository)
	postProducer := post.NewSaramaSyncProducer(syncProducer)
	postService := service.NewPostService(postRepository, logger, interactiveService, checkService, postProducer, historyRepository, checkRepository)
	postHandler := api.NewPostHandler(postService, interactiveService, enforcer, logger)
	historyService := service.NewHistoryService(historyRepository, logger)
	historyHandler := api.NewHistoryHandler(historyService)
	checkHandler := api.NewCheckHandler(checkService, logger, enforcer)
	v := ioc.InitMiddlewares(handler, logger, enforcer)
	permissionDAO := dao.NewPermissionDAO(enforcer, logger, db)
	permissionRepository := repository.NewPermissionRepository(logger, permissionDAO)
	permissionService := service.NewPermissionService(permissionRepository, logger)
	permissionHandler := api.NewPermissionHandler(permissionService, logger, enforcer)
	rankingRedisCache := cache.NewRankingRedisCache(cmdable)
	rankingLocalCache := cache.NewRankingLocalCache()
	rankingRepository := repository.NewRankingCache(rankingRedisCache, rankingLocalCache)
	rankingService := service.NewRankingService(interactiveService, postRepository, rankingRepository, logger)
	rankingHandler := api.NewRakingHandler(rankingService)
	plateDAO := dao.NewPlateDAO(logger, db)
	plateRepository := repository.NewPlateRepository(logger, plateDAO)
	plateService := service.NewPlateService(logger, plateRepository)
	plateHandler := api.NewPlateHandler(plateService, logger, enforcer)
	activityService := service.NewActivityService(activityRepository)
	activityHandler := api.NewActivityHandler(activityService, enforcer, logger)
	engine := ioc.InitWebServer(userHandler, postHandler, historyHandler, checkHandler, v, permissionHandler, rankingHandler, plateHandler, activityHandler)
	cron := ioc.InitRanking(logger, rankingService)
	interactiveReadEventConsumer := post.NewInteractiveReadEventConsumer(interactiveRepository, client, logger)
	smsDAO := dao.NewSmsDAO(db, logger)
	smsCache := cache.NewSMSCache(cmdable)
	tencentSms := ioc.InitSms()
	smsRepository := repository.NewSmsRepository(smsDAO, smsCache, logger, tencentSms)
	smsConsumer := sms.NewSMSConsumer(smsRepository, client, logger, smsCache)
	emailCache := cache.NewEmailCache(cmdable)
	emailRepository := repository.NewEmailRepository(emailCache, logger)
	emailConsumer := email.NewEmailConsumer(emailRepository, client, logger)
	v2 := ioc.InitConsumers(interactiveReadEventConsumer, smsConsumer, emailConsumer)
	cmd := &Cmd{
		server:   engine,
		Cron:     cron,
		consumer: v2,
	}
	return cmd
}
