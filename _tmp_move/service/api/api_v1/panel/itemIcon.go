package panel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiData/panelApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/siteFavicon"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type ItemIcon struct {
}

func getPageTitle(targetURL string) string {
	parsed, err := url.Parse(targetURL)
	if err == nil {
		host := strings.ToLower(parsed.Host)
		// 头条：PC页常为JS渲染，优先走m站info接口
		if strings.Contains(host, "toutiao.com") {
			reArticle := regexp.MustCompile(`/article/(\d+)`)
			reA := regexp.MustCompile(`/a(\d+)`)
			id := ""
			if m := reArticle.FindStringSubmatch(parsed.Path); len(m) > 1 {
				id = m[1]
			} else if m := reA.FindStringSubmatch(parsed.Path); len(m) > 1 {
				id = m[1]
			}
			if id != "" {
				apiURL := "https://m.toutiao.com/i" + id + "/info/"
				client := &http.Client{Timeout: 8 * time.Second}
				req, _ := http.NewRequest(http.MethodGet, apiURL, nil)
				req.Header.Set("User-Agent", "Mozilla/5.0")
				if resp, err := client.Do(req); err == nil {
					defer resp.Body.Close()
					if body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024)); err == nil {
						var obj map[string]any
						if json.Unmarshal(body, &obj) == nil {
							if data, ok := obj["data"].(map[string]any); ok {
								if t, ok := data["title"].(string); ok && strings.TrimSpace(t) != "" {
									return strings.TrimSpace(t)
								}
							}
						}
					}
				}
			}
		}
	}

	client := &http.Client{Timeout: 6 * time.Second}
	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Sun-Panel)")

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 512*1024))
	if err != nil {
		return ""
	}

	re := regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)
	match := re.FindStringSubmatch(string(body))
	if len(match) < 2 {
		return ""
	}

	title := strings.TrimSpace(match[1])
	title = regexp.MustCompile(`\s+`).ReplaceAllString(title, " ")
	if len(title) > 200 {
		title = title[:200]
	}
	return title
}

func (a *ItemIcon) Edit(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	req := models.ItemIcon{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if req.ItemIconGroupId == 0 {
		// apiReturn.Error(c, "Group is mandatory")
		apiReturn.ErrorParamFomat(c, "Group is mandatory")
		return
	}

	req.UserId = userInfo.ID

	// 增加重复性检查
	if req.Url != "" {
		var existItem models.ItemIcon
		checkDb := global.Db.Model(&models.ItemIcon{}).Where("url = ? AND user_id = ?", req.Url, req.UserId)
		if req.ID != 0 {
			checkDb = checkDb.Where("id != ?", req.ID)
		}
		if err := checkDb.First(&existItem).Error; err == nil {
			apiReturn.Error(c, "已有重复链接：["+existItem.Title+"]，请修改！")
			return
		}
	}

	// json转字符串
	if j, err := json.Marshal(req.Icon); err == nil {
		req.IconJson = string(j)
	}

	if req.ID != 0 {
		// 修改
		updateField := []string{"IconJson", "Icon", "Title", "Url", "LanUrl", "Description", "OpenMethod", "GroupId", "UserId", "ItemIconGroupId", "Pinned"}
		if req.Sort != 0 {
			updateField = append(updateField, "Sort")
		}
		global.Db.Model(&models.ItemIcon{}).
			Select(updateField).
			Where("id=?", req.ID).Updates(&req)
	} else {
		req.Sort = 9999
		// 创建
		global.Db.Create(&req)
	}

	apiReturn.SuccessData(c, req)
}

// 添加多个图标
func (a *ItemIcon) AddMultiple(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	// type Request
	req := []models.ItemIcon{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	for i := 0; i < len(req); i++ {
		if req[i].ItemIconGroupId == 0 {
			apiReturn.ErrorParamFomat(c, "Group is mandatory")
			return
		}
		req[i].UserId = userInfo.ID
		// json转字符串
		if j, err := json.Marshal(req[i].Icon); err == nil {
			req[i].IconJson = string(j)
		}
	}

	global.Db.Create(&req)

	apiReturn.SuccessData(c, req)
}

// // 获取详情
// func (a *ItemIcon) GetInfo(c *gin.Context) {
// 	req := systemApiStructs.AiDrawGetInfoReq{}

// 	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
// 		apiReturn.ErrorParamFomat(c, err.Error())
// 		return
// 	}

// 	userInfo, _ := base.GetCurrentUserInfo(c)

// 	aiDraw := models.AiDraw{}
// 	aiDraw.ID = req.ID
// 	if err := aiDraw.GetInfo(global.Db); err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			apiReturn.Error(c, "不存在记录")
// 			return
// 		}
// 		apiReturn.ErrorDatabase(c, err.Error())
// 		return
// 	}

// 	if userInfo.ID != aiDraw.UserID {
// 		apiReturn.ErrorNoAccess(c)
// 		return
// 	}

// 	apiReturn.SuccessData(c, aiDraw)
// }

func (a *ItemIcon) GetListByGroupId(c *gin.Context) {
	req := models.ItemIcon{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	itemIcons := []models.ItemIcon{}

	// 默认排序：置顶在前，然后按创建时间从新到旧
	if err := global.Db.Order("pinned DESC, created_at DESC").Find(&itemIcons, "item_icon_group_id = ? AND user_id=?", req.ItemIconGroupId, userInfo.ID).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	for k, v := range itemIcons {
		json.Unmarshal([]byte(v.IconJson), &itemIcons[k].Icon)
	}

	apiReturn.SuccessListData(c, itemIcons, 0)
}

func (a *ItemIcon) Deletes(c *gin.Context) {
	req := commonApiStructs.RequestDeleteIds[uint]{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	if err := global.Db.Delete(&models.ItemIcon{}, "id in ? AND user_id=?", req.Ids, userInfo.ID).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

// 保存排序
func (a *ItemIcon) SaveSort(c *gin.Context) {
	req := panelApiStructs.ItemIconSaveSortRequest{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	transactionErr := global.Db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		for _, v := range req.SortItems {
			if err := tx.Model(&models.ItemIcon{}).Where("user_id=? AND id=? AND item_icon_group_id=?", userInfo.ID, v.Id, req.ItemIconGroupId).Update("sort", v.Sort).Error; err != nil {
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

// 支持获取并直接下载对方网站图标到服务器
func (a *ItemIcon) GetSiteFavicon(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	req := panelApiStructs.ItemIconGetSiteFaviconReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	resp := panelApiStructs.ItemIconGetSiteFaviconResp{}
	resp.PageTitle = getPageTitle(req.Url)
	fullUrl := ""
	if iconUrl, err := siteFavicon.GetOneFaviconURL(req.Url); err == nil {
		fullUrl = iconUrl
	}

	parsedURL, err := url.Parse(req.Url)
	if err != nil {
		apiReturn.Error(c, "acquisition failed:"+err.Error())
		return
	}

	if fullUrl != "" {
		protocol := parsedURL.Scheme
		global.Logger.Debug("protocol:", protocol)
		global.Logger.Debug("fullUrl:", fullUrl)

		// 如果URL以双斜杠（//）开头，则使用当前页面协议
		if strings.HasPrefix(fullUrl, "//") {
			fullUrl = protocol + "://" + fullUrl[2:]
		} else if !strings.HasPrefix(fullUrl, "http://") && !strings.HasPrefix(fullUrl, "https://") {
			// 如果URL既不以http://开头也不以https://开头，则默认为http协议
			fullUrl = "http://" + fullUrl
		}
		global.Logger.Debug("fullUrl:", fullUrl)
		// 去除图标的get参数
		{
			parsedIcoURL, err := url.Parse(fullUrl)
			if err == nil {
				fullUrl = parsedIcoURL.Scheme + "://" + parsedIcoURL.Host + parsedIcoURL.Path
			}
		}
		global.Logger.Debug("fullUrl:", fullUrl)

		// 生成保存目录
		configUpload := global.Config.GetValueString("base", "source_path")
		savePath := fmt.Sprintf("%s/%d/%d/%d/", configUpload, time.Now().Year(), time.Now().Month(), time.Now().Day())
		isExist, _ := cmn.PathExists(savePath)
		if !isExist {
			os.MkdirAll(savePath, os.ModePerm)
		}

		// 下载
		var imgInfo *os.File
		{
			var err error
			if imgInfo, err = siteFavicon.DownloadImage(fullUrl, savePath, 1024*1024); err == nil {
				// 保存到数据库
				ext := path.Ext(fullUrl)
				mFile := models.File{}
				if _, err := mFile.AddFile(userInfo.ID, parsedURL.Host, ext, imgInfo.Name(), "icon"); err == nil {
					resp.IconUrl = imgInfo.Name()[1:]
				}
			}
		}
	}

	if resp.IconUrl == "" && strings.TrimSpace(resp.PageTitle) == "" {
		apiReturn.Error(c, "acquisition failed")
		return
	}

	apiReturn.SuccessData(c, resp)
}
