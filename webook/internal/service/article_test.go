package service

import (
	"context"
	"errors"
	"github.com/czh0913/gocode/basic-go/webook/internal/domain"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/article"
	artrepomocks "github.com/czh0913/gocode/basic-go/webook/internal/repository/article/mocks"
	"github.com/czh0913/gocode/basic-go/webook/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_articleService_Publish(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (article.ArticleAutherRepository, article.ArticleReaderRepository)
		article domain.Article

		wantErr error
		wantId  int64
	}{
		{
			name: "新建发表成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAutherRepository, article.ArticleReaderRepository) {
				auther := artrepomocks.NewMockArticleAutherRepository(ctrl)
				auther.EXPECT().Creat(gomock.Any(), domain.Article{
					Title:   "标题",
					Content: "内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				reader := artrepomocks.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      1,
					Title:   "标题",
					Content: "内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)

				return auther, reader
			},
			article: domain.Article{
				Title:   "标题",
				Content: "内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId: 1,
		},
		{
			name: "修改并发表成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAutherRepository, article.ArticleReaderRepository) {
				author := artrepomocks.NewMockArticleAutherRepository(ctrl)
				author.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "标题",
					Content: "内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				reader := artrepomocks.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "标题",
					Content: "内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(2), nil)

				return author, reader
			},
			article: domain.Article{
				Id:      2,
				Title:   "标题",
				Content: "内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId: 2,
		},
		{
			name: "保存到制作库失败",
			mock: func(ctrl *gomock.Controller) (article.ArticleAutherRepository, article.ArticleReaderRepository) {
				auther := artrepomocks.NewMockArticleAutherRepository(ctrl)
				auther.EXPECT().Creat(gomock.Any(), domain.Article{
					Title:   "标题",
					Content: "内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("更新制作库失败"))
				reader := artrepomocks.NewMockArticleReaderRepository(ctrl)
				//reader.EXPECT().Save(gomock.Any(), domain.Article{
				//	Id:      1,
				//	Title:   "标题",
				//	Content: "内容",
				//	Author: domain.Author{
				//		Id: 123,
				//	},
				//}).Return()

				return auther, reader
			},
			article: domain.Article{
				Title:   "标题",
				Content: "内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId:  0,
			wantErr: errors.New("更新制作库失败"),
		},
		{
			name: "制作库更新成功，线上库保存重试成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAutherRepository, article.ArticleReaderRepository) {
				author := artrepomocks.NewMockArticleAutherRepository(ctrl)
				author.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "标题",
					Content: "内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				reader := artrepomocks.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "标题",
					Content: "内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("Save 失败"))
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "标题",
					Content: "内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(2), nil)

				return author, reader
			},
			article: domain.Article{
				Id:      2,
				Title:   "标题",
				Content: "内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId: 2,
		},
		{
			name: "制作库更新成功，线上库保存重试全部失败",
			mock: func(ctrl *gomock.Controller) (article.ArticleAutherRepository, article.ArticleReaderRepository) {
				author := artrepomocks.NewMockArticleAutherRepository(ctrl)
				author.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "标题",
					Content: "内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				reader := artrepomocks.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "标题",
					Content: "内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Times(3).Return(int64(0), errors.New("Save 失败"))

				return author, reader
			},
			article: domain.Article{
				Id:      2,
				Title:   "标题",
				Content: "内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId:  0,
			wantErr: errors.New("Save 失败"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			auther, reader := tc.mock(ctrl)
			svc := NewArticleService(auther, reader, &logger.NopLogger{})

			id, err := svc.Publish(context.Background(), tc.article)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantId, id)

		})
	}
}
