.PHONY: mock

mock:
	mockgen -source=D:/gocode/src/basic-go/webook/internal/service/user.go -package=svcmocks -destination=D:/gocode/src/basic-go/webook/internal/service/mocks/user.mock.go
	mockgen -source=D:/gocode/src/basic-go/webook/internal/service/code.go -package=svcmocks -destination=D:/gocode/src/basic-go/webook/internal/service/mocks/code.mock.go
	mockgen -source=D:/gocode/src/basic-go/webook/internal/repository/user.go -package=repomocks -destination=D:/gocode/src/basic-go/webook/internal/repository/mocks/user.mock.go
	mockgen -source=D:/gocode/src/basic-go/webook/internal/repository/code.go -package=repomocks -destination=D:/gocode/src/basic-go/webook/internal/repository/mocks/code.mock.go
	go mod tidy

