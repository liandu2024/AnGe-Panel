package panel

import (
	"math"
	"strings"
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type ItemIconGroup struct {
}

type getListRequest struct {
	GroupType string `json:"groupType"`
}

func normalizeGroupType(groupType string) string {
	groupType = strings.TrimSpace(strings.ToLower(groupType))
	if groupType == "webpage" {
		return "webpage"
	}
	return "website"
}

func (a *ItemIconGroup) Edit(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	req := models.ItemIconGroup{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	req.UserId = userInfo.ID
	req.GroupType = normalizeGroupType(req.GroupType)

	if req.ID != 0 {
		// 修改
		updateField := []string{"Icon", "Title", "Description", "Sort", "GroupType", "UserId"}
		if req.Sort != 0 {
			updateField = append(updateField, "Sort")
		}
		global.Db.Model(&models.ItemIconGroup{}).
			Select(updateField).
			Where("id=?", req.ID).Updates(&req)
	} else {
		// 创建
		global.Db.Create(&req)
	}

	apiReturn.SuccessData(c, req)
}

func (a *ItemIconGroup) GetList(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	groups := []models.ItemIconGroup{}

	req := getListRequest{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil && err.Error() != "EOF" {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	groupType := normalizeGroupType(req.GroupType)

	err := global.Db.Transaction(func(tx *gorm.DB) error {
		// 兼容旧数据：历史分组没有类型时默认归为 website
		if err := tx.Model(&models.ItemIconGroup{}).
			Where("user_id=? AND (group_type='' OR group_type IS NULL)", userInfo.ID).
			Update("group_type", "website").Error; err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return err
		}

		query := tx.Order("sort ,created_at").Where("user_id=?", userInfo.ID)
		if req.GroupType != "" {
			query = query.Where("group_type=?", groupType)
		}
		if err := query.Find(&groups).Error; err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return err
		}

		// 分组为空时自动创建默认分组
		if len(groups) == 0 {
			if req.GroupType == "" {
				websiteGroup := models.ItemIconGroup{
					Title:     "网站",
					UserId:    userInfo.ID,
					Icon:      "material-symbols:language",
					Sort:      1,
					GroupType: "website",
				}
				if err := tx.Create(&websiteGroup).Error; err != nil {
					apiReturn.ErrorDatabase(c, err.Error())
					return err
				}

				webPageGroup := models.ItemIconGroup{
					Title:     "网页",
					UserId:    userInfo.ID,
					Icon:      "material-symbols:web-asset",
					Sort:      1,
					GroupType: "webpage",
				}
				if err := tx.Create(&webPageGroup).Error; err != nil {
					apiReturn.ErrorDatabase(c, err.Error())
					return err
				}

				// 将当前账号下所有无分组图标更新到默认“网站”组
				if err := tx.Model(&models.ItemIcon{}).
					Where("user_id=? AND (item_icon_group_id=0 OR item_icon_group_id IS NULL)", userInfo.ID).
					Update("item_icon_group_id", websiteGroup.ID).Error; err != nil {
					apiReturn.ErrorDatabase(c, err.Error())
					return err
				}

				groups = append(groups, websiteGroup, webPageGroup)
			} else {
				defaultGroup := models.ItemIconGroup{
					Title:     map[bool]string{true: "网页", false: "网站"}[groupType == "webpage"],
					UserId:    userInfo.ID,
					Icon:      map[bool]string{true: "material-symbols:web-asset", false: "material-symbols:language"}[groupType == "webpage"],
					Sort:      1,
					GroupType: groupType,
				}
				if err := tx.Create(&defaultGroup).Error; err != nil {
					apiReturn.ErrorDatabase(c, err.Error())
					return err
				}
				groups = append(groups, defaultGroup)
			}
		}

		// 返回 nil 提交事务
		return nil
	})

	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	} else {
		apiReturn.SuccessListData(c, groups, 0)
	}
}

func (a *ItemIconGroup) Deletes(c *gin.Context) {
	req := commonApiStructs.RequestDeleteIds[uint]{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	userInfo, _ := base.GetCurrentUserInfo(c)

	var count int64
	if err := global.Db.Model(&models.ItemIconGroup{}).Where(" user_id=?", userInfo.ID).Count(&count).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	} else {
		if math.Abs(float64(len(req.Ids))-float64(count)) < 1 {
			apiReturn.ErrorCode(c, 1201, "At least one must be retained", nil)
			return
		}

	}

	txErr := global.Db.Transaction(func(tx *gorm.DB) error {
		mitemIcon := models.ItemIcon{}
		if err := tx.Delete(&models.ItemIconGroup{}, "id in ? AND user_id=?", req.Ids, userInfo.ID).Error; err != nil {
			return err
		}

		if err := mitemIcon.DeleteByItemIconGroupIds(tx, userInfo.ID, req.Ids); err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		apiReturn.ErrorDatabase(c, txErr.Error())
		return
	}

	apiReturn.Success(c)
}

// 保存排序
func (a *ItemIconGroup) SaveSort(c *gin.Context) {
	req := commonApiStructs.SortRequest{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	transactionErr := global.Db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		for _, v := range req.SortItems {
			if err := tx.Model(&models.ItemIconGroup{}).Where("user_id=? AND id=?", userInfo.ID, v.Id).Update("sort", v.Sort).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
		}

		// 返回 nil 提交事务
		return nil
	})

	if transactionErr != nil {
		apiReturn.ErrorDatabase(c, transactionErr.Error())
		return
	}

	apiReturn.Success(c)
}
