package integration

import (
	"bytes"
	"encoding/json"
	"github.com/czh0913/gocode/basic-go/webook/internal/integration/startup"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/dao"
	ijwt "github.com/czh0913/gocode/basic-go/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (s *ArticleTestSuite) SetupSuite() {
	// 在所有测试执行之前，初始化内容
	// 使用之前，初始化 Handler 和路由，以及文章功能所需要的 service
	s.db = startup.InitTestDB()
	s.server = gin.Default()

	s.server.Use(func(ctx *gin.Context) {
		ctx.Set("claims", &ijwt.UserClaims{
			Uid: 123,
		})
	})
	artHdl := startup.InitArticleHandler()
	artHdl.RegisterRoutes(s.server)
}

func (s *ArticleTestSuite) TearDownTest() {
	s.db.Exec("TRUNCATE TABLE articles")
}

// 每一个测试都会执行
func (s *ArticleTestSuite) TestEdit() {
	t := s.T()
	testCases := []struct {
		name string
		art  Article

		// 准备数据，这里我能做什么 ：
		before func(t *testing.T)
		// 验证数据
		after func(t *testing.T)

		wantCode int
		wantRes  Result[int64]
	}{
		{
			name: "新建帖子--保存成功",
			before: func(t *testing.T) {
				// 新建帖子之前不需要做什么
			},
			after: func(t *testing.T) {
				// 在帖子保存之后，需要检查 数据库是否存有新建的帖子
				var art dao.Article
				err := s.db.Where("id=?", 1).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Ctime = 0
				art.Utime = 0
				// 用户id 自定义为123 ，帖子id 是1
				assert.Equal(t, dao.Article{
					ID:       1,
					Title:    "标题",
					Content:  "内容",
					AuthorID: 123,
					Ctime:    0,
					Utime:    0,
				}, art)
			},
			art: Article{
				Title:   "标题",
				Content: "内容",
			},
			wantCode: 200,
			wantRes: Result[int64]{
				Data: 1,
				Msg:  "OK",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/articles/edit", bytes.NewBuffer([]byte(reqBody)))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			// 将构造的请求发送进 server 进入测试逻辑
			s.server.ServeHTTP(recorder, req)

			code := recorder.Code
			// 反序列化为结果
			var result Result[int64]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, code)
			assert.Equal(t, tc.wantRes, result)
			tc.after(t)
		})
	}
}

func (s *ArticleTestSuite) TestA() {
	s.T().Log("这里是测试")
}

func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}

type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
