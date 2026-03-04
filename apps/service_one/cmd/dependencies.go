package main

import (
	env_service_one "go_boilerplate_project/apps/service_one/env"
	brands_data_service "go_boilerplate_project/apps/service_one/layers/data/brands"
	products_data_service "go_boilerplate_project/apps/service_one/layers/data/products"
	brands_domain_service "go_boilerplate_project/apps/service_one/layers/domain/brands"
	products_domain_service "go_boilerplate_project/apps/service_one/layers/domain/products"
	brands_http_service "go_boilerplate_project/apps/service_one/layers/http/brands"
	products_http_service "go_boilerplate_project/apps/service_one/layers/http/products"
	serviceone_middleware "go_boilerplate_project/apps/service_one/middlewares"
	serviceone_router "go_boilerplate_project/apps/service_one/router"
	global_middleware "go_boilerplate_project/middlewares/global"
	dependencies_models "go_boilerplate_project/models/dependencies"
	env_models "go_boilerplate_project/models/env"
	context_repository "go_boilerplate_project/services/context"
	api_helpers_service "go_boilerplate_project/services/helpers/api"
	custom_helpers_service "go_boilerplate_project/services/helpers/custom"
	mongodb_helpers_service "go_boilerplate_project/services/helpers/db/mongodb"
	mysql_helpers_service "go_boilerplate_project/services/helpers/db/mysql"
	redis_helpers_service "go_boilerplate_project/services/helpers/db/redis"
	logger_service "go_boilerplate_project/services/logger"
	network_service "go_boilerplate_project/services/network"
	transactions_service "go_boilerplate_project/services/transactions"
	"go_boilerplate_project/storage/mongodb"
	"go_boilerplate_project/storage/mysql"
	"go_boilerplate_project/storage/redis"
	"log"

	"go.uber.org/zap"
)

func initService() {
	env, err := env_service_one.LoadEnv()
	if err != nil {
		log.Fatalf("Failed to load environment for service one: %v", err)
	}
	l, err := logger_service.InitLogger(env.Logger)
	if err != nil {
		log.Fatalf("Failed to initialize logger for service one: %v", err)
	}

	initDependencies(env, l)
}

func initDatabases(env *env_models.Environment, logger *zap.SugaredLogger) (*dependencies_models.Databases, error) {
	mongoDBClient, err := mongodb.InitMongoDB(env.Databases.MongoDB, logger)
	if err != nil {
		return nil, err
	}
	redisClient, err := redis.InitRedis(env.Databases.Redis, logger)
	if err != nil {
		return nil, err
	}
	mysqlClient, err := mysql.InitMySQL(env.Databases.MySQL, logger)
	if err != nil {
		return nil, err
	}
	return &dependencies_models.Databases{
		MongoDBClient: mongoDBClient,
		RedisClient:   redisClient,
		MySQLClient:   mysqlClient,
	}, nil
}

func initDependencies(env *env_models.Environment, logger *zap.SugaredLogger) {
	databases, err := initDatabases(env, logger)
	if err != nil {
		log.Fatalf("Failed to initialize databases for service one: %v", err)
	}
	dependencies := &dependencies_models.ServiceOneDependencies{
		Databases:   databases,
		Logger:      logger,
		Environment: env,
		Helpers: &dependencies_models.Helpers{
			API:           nil,
			DBHelpers:     nil,
			CustomHelpers: nil,
		},
		Services: &dependencies_models.Services{
			Context:      nil,
			Transactions: nil,
			Network:      nil,
		},
		Http: &dependencies_models.Http{
			Brands:   nil,
			Products: nil,
		},
		Router: &dependencies_models.Router{
			ServiceOne: nil,
		},
		Domain: &dependencies_models.Domain{
			Brands:   nil,
			Products: nil,
		},
		Data: &dependencies_models.Data{
			Brands:   nil,
			Products: nil,
		},
		Middleware: nil,
	}

	err = initContextService(dependencies)
	if err != nil {
		log.Fatalf("Failed to initialize context service for service one: %v", err)
	}

	err = initHelpers(dependencies)
	if err != nil {
		log.Fatalf("Failed to initialize helpers for service one: %v", err)
	}

	err = initMiddlewares(dependencies)
	if err != nil {
		log.Fatalf("Failed to initialize middlewares for service one: %v", err)
	}

	err = initServices(dependencies)
	if err != nil {
		log.Fatalf("Failed to initialize services for service one: %v", err)
	}

	err = initDataLayers(dependencies)
	if err != nil {
		log.Fatalf("Failed to initialize data layers for service one: %v", err)
	}

	err = initDomainLayers(dependencies)
	if err != nil {
		log.Fatalf("Failed to initialize domain layers for service one: %v", err)
	}

	err = initHttpLayers(dependencies)
	if err != nil {
		log.Fatalf("Failed to initialize http layers for service one: %v", err)
	}

	err = initRouter(dependencies)
	if err != nil {
		log.Fatalf("Failed to initialize router for service one: %v", err)
	}

}

func initHelpers(d *dependencies_models.ServiceOneDependencies) error {
	apiHelpers := api_helpers_service.InitService(api_helpers_service.Input{
		Logger: d.Logger,
	})
	d.Helpers.API = apiHelpers

	customHelpers := custom_helpers_service.InitService(custom_helpers_service.Input{
		Logger: d.Logger,
	})
	d.Helpers.CustomHelpers = customHelpers

	mongodbHelpers := mongodb_helpers_service.InitService(mongodb_helpers_service.Input{
		Logger: d.Logger,
		Client: &mongodb_helpers_service.Client{
			MongoDBClient: &mongodb_helpers_service.MongoClient{
				Client:   d.Databases.MongoDBClient,
				Database: d.Environment.Databases.MongoDB.Database,
			},
		},
		Services: &mongodb_helpers_service.Services{
			Context: d.Services.Context,
		},
	})

	mysqlHelpers := mysql_helpers_service.InitService(mysql_helpers_service.Input{
		Logger: d.Logger,
		Client: &mysql_helpers_service.Client{
			MySQLClient: d.Databases.MySQLClient,
		},
		Services: &mysql_helpers_service.Services{
			Context: d.Services.Context,
		},
		Env: d.Environment.Databases.MySQL,
	})

	redisHelpers := redis_helpers_service.InitService(redis_helpers_service.Input{
		Logger: d.Logger,
		Client: &redis_helpers_service.Client{
			RedisClient: d.Databases.RedisClient,
		},
		Services: &redis_helpers_service.Services{
			Context: d.Services.Context,
		},
	})

	d.Helpers.DBHelpers = &dependencies_models.DBHelpers{
		MongoDB: mongodbHelpers,
		MySQL:   mysqlHelpers,
		Redis:   redisHelpers,
	}
	return nil
}

func initMiddlewares(d *dependencies_models.ServiceOneDependencies) error {
	global := global_middleware.InitService(global_middleware.Input{
		Logger: d.Logger,
	})
	application := serviceone_middleware.InitService(serviceone_middleware.Input{
		Logger:     d.Logger,
		AppName:    d.Environment.App.AppName,
		AppVersion: d.Environment.App.AppVersion,
	})
	d.Middleware = &dependencies_models.Middleware{
		Global:      global,
		Application: application,
	}
	return nil
}

func initContextService(d *dependencies_models.ServiceOneDependencies) error {
	context := context_repository.InitService(context_repository.Input{})
	d.Services.Context = context
	return nil
}

func initServices(d *dependencies_models.ServiceOneDependencies) error {
	transactions := transactions_service.InitService(transactions_service.Input{
		Helpers: &transactions_service.Helpers{
			MongoDB: d.Helpers.DBHelpers.MongoDB,
			MySQL:   d.Helpers.DBHelpers.MySQL,
		},
		Services: &transactions_service.Services{
			Context: d.Services.Context,
		},
		Logger: d.Logger,
	})
	d.Services.Transactions = transactions

	network := network_service.InitService(network_service.Input{
		Logger: d.Logger,
	})
	d.Services.Network = network

	return nil
}

func initDataLayers(d *dependencies_models.ServiceOneDependencies) error {
	brandsData := brands_data_service.InitService(brands_data_service.Input{
		Helpers: &brands_data_service.Helpers{
			MySQL:   d.Helpers.DBHelpers.MySQL,
			Redis:   d.Helpers.DBHelpers.Redis,
			MongoDB: d.Helpers.DBHelpers.MongoDB,
			Custom:  d.Helpers.CustomHelpers,
		},
		Services: &brands_data_service.Services{
			Context: d.Services.Context,
		},
		Logger: d.Logger,
	})
	d.Data.Brands = brandsData

	productsData := products_data_service.InitService(products_data_service.Input{
		Helpers: &products_data_service.Helpers{
			MySQL:   d.Helpers.DBHelpers.MySQL,
			Redis:   d.Helpers.DBHelpers.Redis,
			MongoDB: d.Helpers.DBHelpers.MongoDB,
			Custom:  d.Helpers.CustomHelpers,
		},
		Services: &products_data_service.Services{
			Context: d.Services.Context,
		},
		Logger: d.Logger,
	})
	d.Data.Products = productsData

	return nil
}

func initDomainLayers(d *dependencies_models.ServiceOneDependencies) error {
	brandsDomain := brands_domain_service.InitService(brands_domain_service.Input{
		Data: &brands_domain_service.Data{
			Brands:   d.Data.Brands,
			Products: d.Data.Products,
		},
		Services: &brands_domain_service.Services{
			Context:      d.Services.Context,
			Transactions: d.Services.Transactions,
		},
		Logger: d.Logger,
	})
	d.Domain.Brands = brandsDomain

	productsDomain := products_domain_service.InitService(products_domain_service.Input{
		Data: &products_domain_service.Data{
			Products: d.Data.Products,
			Brands:   d.Data.Brands,
		},
		Services: &products_domain_service.Services{
			Context: d.Services.Context,
		},
		Logger: d.Logger,
	})
	d.Domain.Products = productsDomain

	return nil
}

func initHttpLayers(d *dependencies_models.ServiceOneDependencies) error {
	brandsHttp := brands_http_service.InitService(brands_http_service.Input{
		Domain: &brands_http_service.Domain{
			Brands: d.Domain.Brands,
		},
		Logger: d.Logger,
		Services: &brands_http_service.Services{
			Context: d.Services.Context,
		},
		Helpers: &brands_http_service.Helpers{
			API: d.Helpers.API,
		},
	})
	d.Http.Brands = brandsHttp

	productsHttp := products_http_service.InitService(products_http_service.Input{
		Domain: &products_http_service.Domain{
			Products: d.Domain.Products,
		},
		Logger: d.Logger,
		Services: &products_http_service.Services{
			Context: d.Services.Context,
		},
		Helpers: &products_http_service.Helpers{
			API: d.Helpers.API,
		},
	})
	d.Http.Products = productsHttp

	return nil
}

func initRouter(d *dependencies_models.ServiceOneDependencies) error {
	router := serviceone_router.InitService(serviceone_router.Input{
		Logger: d.Logger,
		Services: &serviceone_router.Services{
			Context: d.Services.Context,
		},
		Helpers: &serviceone_router.Helpers{
			API: d.Helpers.API,
		},
		Env: d.Environment,
		Http: &serviceone_router.Http{
			Brands:   d.Http.Brands,
			Products: d.Http.Products,
		},
		Middleware: &serviceone_router.Middleware{
			Global:      d.Middleware.Global,
			Application: d.Middleware.Application,
		},
	})
	d.Router.ServiceOne = router

	router.ConfigureRouter()
	router.SetupRoutes()
	router.Run()

	return nil
}
