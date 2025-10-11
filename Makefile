.PHONY: mock

mock:
	@echo ">>> Generating mocks..."
	mockgen -source=./webook/internal/service/user.go -package=svcmocks -destination=./webook/internal/service/mocks/user.mock.go
	mockgen -source=./webook/internal/service/code.go -package=svcmocks -destination=./webook/internal/service/mocks/code.mock.go
	mockgen -source=./webook/internal/service/article.go -package=svcmocks -destination=./webook/internal/service/mocks/article.mock.go
	mockgen -source=./webook/internal/repository/user.go -package=repomocks -destination=./webook/internal/repository/mocks/user.mock.go
	mockgen -source=./webook/internal/repository/code.go -package=repomocks -destination=./webook/internal/repository/mocks/code.mock.go
	mockgen -source=./webook/internal/repository/article/article.go -package=artrepomocks -destination=./webook/internal/repository/article/mocks/ariticle.mock.go
	mockgen -source=./webook/internal/repository/article/article_auther.go -package=artrepomocks -destination=./webook/internal/repository/article/mocks/ariticle_auther.mock.go
	mockgen -source=./webook/internal/repository/article/article_reader.go -package=artrepomocks -destination=./webook/internal/repository/article/mocks/ariticle_reader.mock.go
	mockgen -source=./webook/internal/repository/dao/user.go -package=daomocks -destination=./webook/internal/repository/dao/mocks/user.mock.go
	mockgen -source=./webook/internal/repository/cache/user.go -package=cachemocks -destination=./webook/internal/repository/cache/mocks/user.mock.go
	mockgen -source=./webook/internal/repository/cache/code.go -package=cachemocks -destination=./webook/internal/repository/cache/mocks/code.mock.go
	mockgen -package=redismocks -destination=./webook/internal/repository/cache/redismocks/cmdable.mock.go github.com/redis/go-redis/v9 Cmdable
	@echo ">>> Running go mod tidy..."
	go mod tidy
	@echo "✅ All mocks generated!"
