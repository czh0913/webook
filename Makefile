.PHONY: mock

mock:
	mockgen -source=D:/gocode/src/basic-go/webook/internal/service/user.go -package=svcmocks -destination=D:/gocode/src/basic-go/webook/internal/service/mocks/user.mock.go
	mockgen -source=D:/gocode/src/basic-go/webook/internal/service/code.go -package=svcmocks -destination=D:/gocode/src/basic-go/webook/internal/service/mocks/code.mock.go
	mockgen -source=D:/gocode/src/basic-go/webook/internal/repository/user.go -package=repomocks -destination=D:/gocode/src/basic-go/webook/internal/repository/mocks/user.mock.go
	mockgen -source=D:/gocode/src/basic-go/webook/internal/repository/code.go -package=repomocks -destination=D:/gocode/src/basic-go/webook/internal/repository/mocks/code.mock.go
	mockgen -source=D:/gocode/src/basic-go/webook/internal/repository/dao/user.go -package=daomocks -destination=D:/gocode/src/basic-go/webook/internal/repository/dao/mocks/user.mock.go
	mockgen -source=D:/gocode/src/basic-go/webook/internal/repository/cache/user.go -package=cachemocks -destination=D:/gocode/src/basic-go/webook/internal/repository/cache/mocks/user.mock.go
	mockgen -source=D:/gocode/src/basic-go/webook/internal/repository/cache/code.go -package=cachemocks -destination=D:/gocode/src/basic-go/webook/internal/repository/cache/mocks/code.mock.go
	mockgen -package=redismocks -destination=D:/gocode/src/basic-go/webook/internal/repository/cache/redismocks/cmdable.mock.go github.com/redis/go-redis/v9 Cmdable

	go mod tidy

