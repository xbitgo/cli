package tpls

import (
	"bytes"
	"text/template"
)

const repoTpl = `
package repo

import (
	"context"

	"{{.ProjectName}}/apps/demo/domain/entity"
	"github.com/xbitgo/components/dtx"
	"github.com/xbitgo/components/filterx"
)

type {{.EntityName}}Repo interface {
	Get(ctx context.Context, id int64) (*entity.{{.EntityName}}, error)
	Query(ctx context.Context, query filterx.FilteringList, pg *filterx.Page) (entity.{{.EntityName}}List, int, error)
	Create(ctx context.Context, input *entity.{{.EntityName}}) (*entity.{{.EntityName}}, error)
	UpdateById(ctx context.Context, updates dtx.SetItemList, id int64) error
	DeleteById(ctx context.Context, id int64) error
}
`

const repoImplTpl = `package repo_impl

import (
	"context"

	"{{.ProjectName}}/apps/demo/domain/entity"
	"{{.ProjectName}}/apps/demo/repo_impl/converter"
	"{{.ProjectName}}/apps/demo/repo_impl/dao"
	"{{.ProjectName}}/apps/demo/repo_impl/do"

	"github.com/xbitgo/components/database"
	"github.com/xbitgo/components/dtx"
	"github.com/xbitgo/components/filterx"
)

// @IMPL[repo.{{.EntityName}}epo] @DI
type {{.EntityName}}Impl struct {
	DB  *database.Database ` + "`" + `di:"conf.DB"` + "`" + `
	Dao *dao.{{.EntityName}}Dao
}

func New{{.EntityName}}Impl() *{{.EntityName}}Impl {
	return &{{.EntityName}}Impl{
		Dao: dao.New{{.EntityName}}Dao(),
	}
}

func (impl *{{.EntityName}}Impl) Get(ctx context.Context, id int64) (*entity.{{.EntityName}}, error) {
	session := impl.DB.NewSession(ctx)
	_do, err := impl.Dao.GetById(session, id)
	if err != nil {
		return nil, err
	}
	return converter.To{{.EntityName}}Entity(_do), nil
}

func (impl *{{.EntityName}}Impl) Query(ctx context.Context, query filterx.FilteringList, pg *filterx.Page) (entity.{{.EntityName}}List, int, error) {
	session := impl.DB.NewSession(ctx)
	session, err := query.GormOption(session)
	if err != nil {
		return nil, 0, err
	}
	session, noCount := filterx.PageGormOption(session, pg)
	var (
		doList do.{{.EntityName}}DoList
		count  int
	)
	if noCount {
		doList, err = impl.Dao.FindAll(session)
	} else {
		doList, count, err = impl.Dao.FindPage(session)
	}
	if err != nil {
		return nil, 0, err
	}
	return converter.To{{.EntityName}}List(doList), count, nil
}

func (impl *{{.EntityName}}Impl) Create(ctx context.Context, input *entity.{{.EntityName}}) (*entity.{{.EntityName}}, error) {
	session := impl.DB.NewSession(ctx)
	_do := converter.From{{.EntityName}}Entity(input)
	err := impl.Dao.Create(session, _do)
	if err != nil {
		return nil, err
	}
	output := converter.To{{.EntityName}}Entity(_do)
	// 注册回滚操作
	dtx.Register(ctx, func() error {
		return impl.rollbackCreate(ctx, _do.ID)
	}, nil)
	return output, err
}

func (impl *{{.EntityName}}Impl) rollbackCreate(ctx context.Context, id int64) error {
	session := impl.DB.NewSession(ctx)
	session = session.Where("id = ?", id)
	// 执行补偿回滚
	return impl.Dao.Delete(session)
}

func (impl *{{.EntityName}}Impl) UpdateById(ctx context.Context, updates dtx.SetItemList, id int64) error {
	_updates, err := updates.ToGormUpdates()
	if err != nil {
		return err
	}
	// 查询老的数据
	_old, err := impl.Get(ctx, id)
	if err != nil {
		return err
	}
	session := impl.DB.NewSession(ctx)
	err = impl.Dao.UpdateById(session, _updates, id)
	if err != nil {
		return err
	}
	// 注册回滚操作
	dtx.Register(ctx, func() error {
		return impl.rollbackUpdateById(ctx, updates, id, _old)
	}, nil)
	return err
}

func (impl *{{.EntityName}}Impl) rollbackUpdateById(ctx context.Context, updates dtx.SetItemList, id int64, _old *entity.{{.EntityName}}) error {
	_oldMap := _old.ToTagMap("db")
	sets := make(map[string]interface{})
	for _, update := range updates {
		if update.Operator == dtx.SET {
			sets[update.Field] = _oldMap[update.Field]
		}
	}
	session := impl.DB.NewSession(ctx)
	// 生成rollback updates
	_updates, err := updates.RollbackGormUpdates(sets)
	if err != nil {
		return err
	}
	// 执行补偿回滚
	err = impl.Dao.UpdateById(session, _updates, id)
	return err
}

func (impl *{{.EntityName}}Impl) DeleteById(ctx context.Context, id int64) error {
	// 查询老的数据
	_old, err := impl.Get(ctx, id)
	if err != nil {
		return err
	}
	session := impl.DB.NewSession(ctx)
	err = impl.Dao.DeleteById(session, id)
	if err != nil {
		return err
	}
	// 注册回滚操作
	dtx.Register(ctx, func() error {
		return impl.rollbackDeleteById(ctx, _old)
	}, nil)
	return err
}

func (impl *{{.EntityName}}Impl) rollbackDeleteById(ctx context.Context, _old *entity.{{.EntityName}}) error {
	session := impl.DB.NewSession(ctx)
	_do := converter.From{{.EntityName}}Entity(_old)
	// 执行补偿回滚
	err := impl.Dao.Create(session, _do)
	if err != nil {
		return err
	}
	return err
}
`

type Repo struct {
	ProjectName string
	EntityName  string
}

func (s *Repo) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("repl").Parse(repoTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Repo) ExecuteImpl() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("repl.impl").Parse(repoImplTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
