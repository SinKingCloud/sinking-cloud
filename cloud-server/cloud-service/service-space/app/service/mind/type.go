package mind

import (
	"errors"
	"service-space/app/constant"
	"service-space/app/enum/mind_type"
	"service-space/app/model"
	mindDataRepository "service-space/app/repository/mind_data"
	"strconv"
	"strings"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// CreateType 创建自定义导图分类
func (s *service) CreateType(userId int64, name string) (*Type, error) {
	if userId <= 0 {
		return nil, errors.New("用户信息错误")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("分类名称不能为空")
	}
	settings := s.configService.GetSystemConfig(constant.MindGroup, false)
	typeMax, _ := strconv.ParseInt(strings.TrimSpace(settings[constant.MindTypeMax]), 10, 64)
	if typeMax != 0 && typeMax < 10 {
		typeMax = 10
	}
	if typeMax > 10000 {
		typeMax = 10000
	}
	lockKey := constant.LockMindTypeOperate + strconv.FormatInt(userId, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeMindTypeOperate) {
		return nil, errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)
	data := &model.MindType{
		Id:     str.GetSnowWorkIns().GetId(),
		UserId: userId,
		Name:   name,
	}
	err := s.database.Transaction(func(tx *gorm.DB) error {
		info, err := s.mindTypeRepository.FindByNameAndUserId(name, userId, tx)
		if err != nil {
			return errors.Join(errors.New("获取导图分类失败"), err)
		}
		if info != nil && info.Id > 0 {
			return errors.New("分类名称已存在")
		}
		if typeMax > 0 {
			total, err := s.mindTypeRepository.CountByUserId(userId, tx)
			if err != nil {
				return errors.Join(errors.New("获取导图分类数量失败"), err)
			}
			if total >= typeMax {
				return errors.New("自定义分类数量已达到上限")
			}
		}
		if err = s.mindTypeRepository.Create(data, tx); err != nil {
			return errors.Join(errors.New("创建导图分类失败"), err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &Type{
		Id:         data.Id,
		Name:       data.Name,
		CreateTime: data.CreateTime,
		UpdateTime: data.UpdateTime,
	}, nil
}

// DeleteType 删除自定义导图分类
func (s *service) DeleteType(userId int64, id int64) error {
	if userId <= 0 || id <= 0 {
		return errors.New("参数错误")
	}
	if _, ok := mind_type.Map()[int(id)]; ok {
		return errors.New("预设分类不能删除")
	}
	lockKey := constant.LockMindTypeOperate + strconv.FormatInt(userId, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeMindTypeOperate) {
		return errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)
	return s.database.Transaction(func(tx *gorm.DB) error {
		info, err := s.mindTypeRepository.FindByIdAndUserId(id, userId, tx)
		if err != nil {
			return errors.Join(errors.New("获取导图分类失败"), err)
		}
		if info == nil || info.Id <= 0 {
			return errors.New("导图分类不存在")
		}
		if err = s.mindDataRepository.UpdateByTypeIdAndUserId(id, userId, &mindDataRepository.UpdateMindData{TypeId: int64(mind_type.All)}, tx); err != nil {
			return errors.Join(errors.New("更新导图分类失败"), err)
		}
		if err = s.mindTypeRepository.DeleteByIdAndUserId(id, userId, tx); err != nil {
			return errors.Join(errors.New("删除导图分类失败"), err)
		}
		return nil
	})
}

// SelectTypes 查询自定义导图分类列表
func (s *service) SelectTypes(userId int64, query *page.Query) (*page.Result[*Type], error) {
	if userId <= 0 {
		return nil, errors.New("用户信息错误")
	}
	if query == nil {
		return nil, errors.New("分页参数错误")
	}
	data, err := s.mindTypeRepository.SelectByUserId(userId, query)
	if err != nil {
		return nil, errors.Join(errors.New("查询导图分类失败"), err)
	}
	list := make([]*Type, 0, len(data.List))
	for _, item := range data.List {
		if item == nil {
			list = append(list, nil)
			continue
		}
		list = append(list, &Type{
			Id:         item.Id,
			Name:       item.Name,
			CreateTime: item.CreateTime,
			UpdateTime: item.UpdateTime,
		})
	}
	return &page.Result[*Type]{
		Total:            data.Total,
		Page:             data.Page,
		PageSize:         data.PageSize,
		List:             list,
		HasMore:          data.HasMore,
		NextCursorId:     data.NextCursorId,
		NextCursorLastId: data.NextCursorLastId,
	}, nil
}

// UpdateType 更新自定义导图分类
func (s *service) UpdateType(userId int64, id int64, name string) error {
	if userId <= 0 || id <= 0 {
		return errors.New("参数错误")
	}
	if _, ok := mind_type.Map()[int(id)]; ok {
		return errors.New("预设分类不能修改")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("分类名称不能为空")
	}
	lockKey := constant.LockMindTypeOperate + strconv.FormatInt(userId, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeMindTypeOperate) {
		return errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)
	info, err := s.mindTypeRepository.FindByIdAndUserId(id, userId)
	if err != nil {
		return errors.Join(errors.New("获取导图分类失败"), err)
	}
	if info == nil || info.Id <= 0 {
		return errors.New("导图分类不存在")
	}
	if err = s.mindTypeRepository.UpdateByIdAndUserId(id, userId, name); err != nil {
		return errors.Join(errors.New("更新导图分类失败"), err)
	}
	return nil
}
