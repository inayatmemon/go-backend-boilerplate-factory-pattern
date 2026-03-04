package dependencies_models

import (
	brands_data_service "go_boilerplate_project/apps/service_one/layers/data/brands"
	products_data_service "go_boilerplate_project/apps/service_one/layers/data/products"
	brands_domain_service "go_boilerplate_project/apps/service_one/layers/domain/brands"
	products_domain_service "go_boilerplate_project/apps/service_one/layers/domain/products"
	brands_http_service "go_boilerplate_project/apps/service_one/layers/http/brands"
	products_http_service "go_boilerplate_project/apps/service_one/layers/http/products"
	serviceone_router "go_boilerplate_project/apps/service_one/router"
	application_middleware "go_boilerplate_project/middlewares/application"
	global_middleware "go_boilerplate_project/middlewares/global"
	env_models "go_boilerplate_project/models/env"
	context_repository "go_boilerplate_project/services/context"
	api_helpers_service "go_boilerplate_project/services/helpers/api"
	custom_helpers_service "go_boilerplate_project/services/helpers/custom"
	mongodb_helpers_service "go_boilerplate_project/services/helpers/db/mongodb"
	mysql_helpers_service "go_boilerplate_project/services/helpers/db/mysql"
	redis_helpers_service "go_boilerplate_project/services/helpers/db/redis"
	network_service "go_boilerplate_project/services/network"
	transactions_service "go_boilerplate_project/services/transactions"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type Databases struct {
	MongoDBClient *mongo.Client
	RedisClient   *redis.Client
	MySQLClient   *gorm.DB
}

type ServiceOneDependencies struct {
	Databases   *Databases
	Logger      *zap.SugaredLogger
	Environment *env_models.Environment
	Helpers     *Helpers
	Services    *Services
	Http        *Http
	Router      *Router
	Domain      *Domain
	Data        *Data
	Middleware  *Middleware
}

type Middleware struct {
	Global      global_middleware.Repository
	Application application_middleware.Repository
}

type Helpers struct {
	API           api_helpers_service.Repository
	DBHelpers     *DBHelpers
	CustomHelpers custom_helpers_service.Repository
}

type DBHelpers struct {
	MongoDB mongodb_helpers_service.Repository
	MySQL   mysql_helpers_service.Repository
	Redis   redis_helpers_service.Repository
}

type Services struct {
	Context      context_repository.Repository
	Transactions transactions_service.Repository
	Network      network_service.Repository
}

type Http struct {
	Brands   brands_http_service.Repository
	Products products_http_service.Repository
}

type Router struct {
	ServiceOne serviceone_router.Repository
}

type Domain struct {
	Brands   brands_domain_service.Repository
	Products products_domain_service.Repository
}

type Data struct {
	Brands   brands_data_service.Repository
	Products products_data_service.Repository
}
