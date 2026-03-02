<script setup lang="ts">
import { VueDraggable } from 'vue-draggable-plus'
import { NBackTop, NButton, NButtonGroup, NDropdown, NModal, NSkeleton, NSpin, NEllipsis, useDialog, useMessage } from 'naive-ui'
import { computed, nextTick, onMounted, onUnmounted, ref, h } from 'vue'

const isMobile = ref(false)
function updateIsMobile() {
  isMobile.value = window.innerWidth <= 500
}
import { AppIcon, AppStarter, EditItem } from './components'
import { Clock, SearchBox, SystemMonitor } from '@/components/deskModule'
import { SvgIcon, SvgIconOnline } from '@/components/common'
import { deletes, edit, getListByGroupId, saveSort } from '@/api/panel/itemIcon'
import { getList as getGroupList } from '@/api/panel/itemIconGroup'

import { setTitle, updateLocalUserInfo } from '@/utils/cmn'
import { useAuthStore, usePanelState } from '@/store'
import { PanelPanelConfigStyleEnum, PanelStateNetworkModeEnum } from '@/enums'
import { VisitMode } from '@/enums/auth'
import { router } from '@/router'
import { t } from '@/locales'

interface ItemGroup extends Panel.ItemIconGroup {
  sortStatus?: boolean
  hoverStatus: boolean
  items?: Panel.ItemInfo[]
}

const ms = useMessage()
const dialog = useDialog()
const panelState = usePanelState()
const authStore = useAuthStore()

const scrollContainerRef = ref<HTMLElement | undefined>(undefined)
const groupRefs = ref<Record<number, HTMLElement | null>>({})
const groupDotMap = ref<Record<number, number>>({})

const editItemInfoShow = ref<boolean>(false)
const editItemInfoData = ref<Panel.ItemInfo | null>(null)
const windowShow = ref<boolean>(false)
const windowSrc = ref<string>('')
const windowTitle = ref<string>('')

const windowIframeRef = ref(null)
const windowIframeIsLoad = ref<boolean>(false)

const dropdownMenuX = ref(0)
const dropdownMenuY = ref(0)
const dropdownShow = ref(false)
const currentRightSelectItem = ref<Panel.ItemInfo | null>(null)
const currentAddItenIconGroupId = ref<number | undefined>()

const settingModalShow = ref(false)

const items = ref<ItemGroup[]>([])
const filterItems = ref<ItemGroup[]>([])
const currentGroupType = ref<'website' | 'webpage'>('website')

// 去掉分页逻辑代码，直接返回 items
// const webpagePagination = ref<Record<number, number>>({})
// const PAGE_SIZE = 20

// function getGroupPageItems(itemGroup: ItemGroup, groupId?: number) {
//   if (!itemGroup.items || currentGroupType.value !== 'webpage' || itemGroup.sortStatus) {
//     return itemGroup.items
//   }
//   const pid = groupId || 0
//   const currentPage = webpagePagination.value[pid] || 1
//   const start = (currentPage - 1) * PAGE_SIZE
//   return itemGroup.items.slice(start, start + PAGE_SIZE)
// }

const fixedElementStyle = computed(() => {
  const pos = panelState.panelConfig.floatingToolsPosition || 'right-bottom'
  const rightOffset = panelState.panelConfig.rightScrollBarShow === false ? '10px' : '24px'

  // 让悬浮按钮在“右上角”时也遵循用户配置的上边距/上内边距
  // 说明：原先内容区用 marginTop(%) + paddingTop(px) 控制顶部留白，这里用 vh + px 近似匹配视觉预期
  const marginTop = Number(panelState.panelConfig.marginTop ?? 0)
  const paddingTop = Number(panelState.panelConfig.paddingTop ?? 0)
  const topOffset = `calc(${marginTop}vh + ${paddingTop}px + 12px)`

  switch (pos) {
    case 'right-top':
      return { right: rightOffset, top: topOffset }
    case 'right-middle':
      return { right: rightOffset, top: '50%', transform: 'translateY(-50%)' }
    default:
      return { right: rightOffset, bottom: '50px' }
  }
})

function openPage(openMethod: number, url: string, title?: string) {
  switch (openMethod) {
    case 1:
      window.location.href = url
      break
    case 2:
      window.open(url)
      break
    case 3:
      windowShow.value = true
      windowSrc.value = url
      windowTitle.value = title || url
      windowIframeIsLoad.value = true
      break

    default:
      break
  }
}

function handleItemClick(itemGroupIndex: number, item: Panel.ItemInfo) {
  if (items.value[itemGroupIndex] && items.value[itemGroupIndex].sortStatus) {
    handleEditItem(item)
    return
  }

  let jumpUrl = ''

  if (item)
    jumpUrl = (panelState.networkMode === PanelStateNetworkModeEnum.lan ? item.lanUrl : item.url) as string
  if (item.lanUrl === '')
    jumpUrl = item.url

  openPage(item.openMethod, jumpUrl, item.title)
}

function handWindowIframeIdLoad(payload: Event) {
  windowIframeIsLoad.value = false
}

function getList() {
  // 获取组数据
  getGroupList<Common.ListResponse<ItemGroup[]>>(currentGroupType.value).then(({ code, data, msg }) => {
    if (code === 0)
      items.value = data.list
    for (let i = 0; i < data.list.length; i++) {
      const element = data.list[i]
      if (element.id)
        updateItemIconGroupByNet(i, element.id)
    }
    filterItems.value = items.value
    nextTick(() => recalcCatalogDots())
    // console.log(items)
  })
}

// 从后端获取组下面的图标
function updateItemIconGroupByNet(itemIconGroupIndex: number, itemIconGroupId: number) {
  getListByGroupId<Common.ListResponse<Panel.ItemInfo[]>>(itemIconGroupId).then((res) => {
    if (res.code === 0) {
      items.value[itemIconGroupIndex].items = res.data.list
      if (currentSearchKeyword.value) {
        itemFrontEndSearch(currentSearchKeyword.value)
      }
    }
    nextTick(() => recalcCatalogDots())
  })
}

function handleRightMenuSelect(key: string | number) {
  dropdownShow.value = false
  // console.log(currentRightSelectItem, key)
  let jumpUrl = panelState.networkMode === PanelStateNetworkModeEnum.lan ? currentRightSelectItem.value?.lanUrl : currentRightSelectItem.value?.url
  if (currentRightSelectItem.value?.lanUrl === '')
    jumpUrl = currentRightSelectItem.value.url
  switch (key) {
    case 'newWindows':
      window.open(jumpUrl)
      break
    case 'copyUrl':
      if (currentRightSelectItem.value?.url) {
        copyUrlToClipboard(currentRightSelectItem.value.url)
      }
      break
    case 'copyLanUrl':
    case 'copyWanUrl':
      if (currentRightSelectItem.value?.lanUrl) {
        copyUrlToClipboard(currentRightSelectItem.value.lanUrl)
      }
      break
    case 'openWanUrl':
      if (currentRightSelectItem.value)
        openPage(currentRightSelectItem.value?.openMethod, currentRightSelectItem.value?.url, currentRightSelectItem.value?.title)
      break
    case 'openLanUrl':
      if (currentRightSelectItem.value && currentRightSelectItem.value.lanUrl)
        openPage(currentRightSelectItem.value?.openMethod, currentRightSelectItem.value.lanUrl, currentRightSelectItem.value?.title)
      break
    case 'edit':
      // 这里有个奇怪的问题，如果不使用{...}的方式 父组件的值会同步修改 标记一下
      handleEditItem({ ...currentRightSelectItem.value } as Panel.ItemInfo)
      break
    case 'delete':
      dialog.warning({
        title: t('common.warning'),
        content: t('common.deleteConfirmByName', { name: currentRightSelectItem.value?.title }),
        positiveText: t('common.confirm'),
        negativeText: t('common.cancel'),
        onPositiveClick: () => {
          const itemId = currentRightSelectItem.value?.id as number
          const groupId = currentRightSelectItem.value?.itemIconGroupId as number
          deletes([itemId]).then(({ code, msg }) => {
            if (code === 0) {
              ms.success(t('common.deleteSuccess'))
              removeItemLocal(itemId, groupId)
              refreshFilteredView()
              nextTick(() => recalcCatalogDots())
            }
            else {
              ms.error(`${t('common.deleteFail')}:${msg}`)
            }
          })
        },
      })

      break
    default:
      break
  }
}

function sortWebpageItemsInGroup(group: ItemGroup) {
  if (!group.items)
    return

  // 置顶在前；其余按创建时间从新到旧（无 createdAt 时用 id 兜底）
  group.items.sort((a: any, b: any) => {
    const ap = a?.pinned ? 1 : 0
    const bp = b?.pinned ? 1 : 0
    if (bp !== ap)
      return bp - ap

    const at = a?.createdAt ? new Date(a.createdAt).getTime() : 0
    const bt = b?.createdAt ? new Date(b.createdAt).getTime() : 0
    if (bt !== at)
      return bt - at

    return (b?.id ?? 0) - (a?.id ?? 0)
  })
}

function updatePinnedLocal(itemId: number, groupId: number, pinned: boolean) {
  const group = items.value.find(g => (g.id as number) === groupId)
  if (!group || !group.items)
    return

  const target = group.items.find(it => (it.id as number) === itemId)
  if (!target)
    return

  ;(target as any).pinned = pinned
  sortWebpageItemsInGroup(group)
}

function quickTogglePinWebpage(item: Panel.ItemInfo) {
  const itemId = item.id as number
  const groupId = item.itemIconGroupId as number
  const nextPinned = !item.pinned

  // 先本地更新，避免“置顶成功”后整页闪动
  updatePinnedLocal(itemId, groupId, nextPinned)

  const newItem = { ...item, pinned: nextPinned }
  edit(newItem).then(({ code, msg }) => {
    if (code === 0) {
      ms.success(nextPinned ? '置顶成功' : '已取消置顶')
      // 不再 getList()，只更新本地列表
      return
    }

    // 失败则回滚
    updatePinnedLocal(itemId, groupId, !nextPinned)
    ms.error(`操作失败：${msg}`)
  }).catch((e) => {
    updatePinnedLocal(itemId, groupId, !nextPinned)
    ms.error(`操作失败：${e?.message ?? e}`)
  })
}

function quickEditWebpage(item: Panel.ItemInfo) {
  handleEditItem({ ...item } as Panel.ItemInfo)
}

function quickDeleteWebpage(item: Panel.ItemInfo) {
  const itemId = item.id as number
  const groupId = item.itemIconGroupId as number

  dialog.warning({
    title: t('common.warning'),
    content: t('common.confirmDelete'),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: () => {
      deletes([itemId]).then(({ code, msg }) => {
        if (code === 0) {
          ms.success(t('common.deleteSuccess'))
          removeItemLocal(itemId, groupId)
          refreshFilteredView()
          nextTick(() => recalcCatalogDots())
        } else {
          ms.error(`${t('common.deleteFail')}:${msg}`)
        }
      })
    },
  })
}

function copyUrlToClipboard(url: string) {
  if (navigator.clipboard) {
    navigator.clipboard.writeText(url).then(() => {
      ms.success('链接已复制成功！')
    }).catch(() => {
      fallbackCopy(url)
    })
  } else {
    fallbackCopy(url)
  }
}

function fallbackCopy(text: string) {
  const textArea = document.createElement('textarea')
  textArea.value = text
  textArea.style.position = 'fixed'
  textArea.style.top = '-9999px'
  document.body.appendChild(textArea)
  textArea.focus()
  textArea.select()
  try {
    document.execCommand('copy')
    ms.success('链接已复制成功！')
  } catch (err) {
    ms.error('复制失败，可能是浏览器限制')
  }
  document.body.removeChild(textArea)
}

function handleContextMenu(e: MouseEvent, itemGroupIndex: number, item: Panel.ItemInfo) {
  if (items.value[itemGroupIndex] && items.value[itemGroupIndex].sortStatus)
    return

  e.preventDefault()
  currentRightSelectItem.value = item
  dropdownShow.value = false
  nextTick().then(() => {
    dropdownShow.value = true
    dropdownMenuX.value = e.clientX
    dropdownMenuY.value = e.clientY
  })
}

function onClickoutside() {
  // message.info('clickoutside')
  dropdownShow.value = false
}

function upsertItemLocal(saved: any) {
  const groupId = saved?.itemIconGroupId as number
  if (!groupId)
    return

  const group = items.value.find(g => (g.id as number) === groupId)
  if (!group)
    return

  if (!group.items)
    group.items = []

  const idx = group.items.findIndex(it => (it.id as number) === (saved.id as number))
  if (idx >= 0) {
    group.items.splice(idx, 1, { ...group.items[idx], ...saved })
  } else {
    // 新建：按创建时间“新到旧”逻辑插入（置顶在前）
    group.items.unshift(saved)
  }

  // 仅网页分组需要置顶/时间排序；网站图标保持原逻辑
  if (currentGroupType.value === 'webpage')
    sortWebpageItemsInGroup(group)
}

function removeItemLocal(itemId: number, groupId: number) {
  const group = items.value.find(g => (g.id as number) === groupId)
  if (!group || !group.items)
    return

  const idx = group.items.findIndex(it => (it.id as number) === itemId)
  if (idx >= 0)
    group.items.splice(idx, 1)
}

function refreshFilteredView() {
  if (currentSearchKeyword.value)
    itemFrontEndSearch(currentSearchKeyword.value)
  else
    filterItems.value = items.value
}

function handleEditSuccess(item: Panel.ItemInfo) {
  // 新建/修改成功：只本地更新对应分组列表，不全量 getList
  upsertItemLocal(item as any)
  refreshFilteredView()
  nextTick(() => recalcCatalogDots())
}

function handleChangeNetwork(mode: PanelStateNetworkModeEnum) {
  panelState.setNetworkMode(mode)
  if (mode === PanelStateNetworkModeEnum.lan)
    ms.success(t('panelHome.changeToLanModelSuccess'))

  else
    ms.success(t('panelHome.changeToWanModelSuccess'))
}

// 结束拖拽
// function handleEndDrag(event: any, itemIconGroup: Panel.ItemIconGroup) {
//   // console.log(event)
//   // console.log(items.value)
// }

function handleSaveSort(itemGroup: ItemGroup) {
  const saveItems: Common.SortItemRequest[] = []
  if (itemGroup.items) {
    for (let i = 0; i < itemGroup.items.length; i++) {
      const element = itemGroup.items[i]
      saveItems.push({
        id: element.id as number,
        sort: i + 1,
      })
    }

    saveSort({ itemIconGroupId: itemGroup.id as number, sortItems: saveItems }).then(({ code, msg }) => {
      if (code === 0) {
        ms.success(t('common.saveSuccess'))
        itemGroup.sortStatus = false
      }
      else {
        ms.error(`${t('common.saveFail')}:${msg}`)
      }
    })
  }
}

function getDropdownMenuOptions() {
  const dropdownMenuOptions: any[] = []

  const renderIconLabel = (label: string, iconName: string) => {
    return h(
      'div', 
      { class: 'flex items-center gap-2' }, 
      [
        h(SvgIcon, { icon: iconName, style: 'font-size: 16px; color: #666;' }),
        h('span', null, label)
      ]
    )
  }

  dropdownMenuOptions.push({
    label: () => renderIconLabel('打开链接', 'mdi:wan'), // 使用系统里确定有的常见图标代替
    key: 'newWindows',
  })

  dropdownMenuOptions.push({
    label: () => renderIconLabel('复制链接', 'ion-copy'),
    key: 'copyUrl',
  })

  if (currentRightSelectItem.value?.lanUrl && panelState.networkMode === PanelStateNetworkModeEnum.wan) {
    dropdownMenuOptions.push({
      label: () => renderIconLabel('打开内网链接', 'material-symbols:lan-outline-rounded'),
      key: 'openLanUrl',
    })
    
    dropdownMenuOptions.push({
      label: () => renderIconLabel('复制内网链接', 'ion-copy'),
      key: 'copyLanUrl',
    })
  }

  if (currentRightSelectItem.value?.lanUrl && panelState.networkMode === PanelStateNetworkModeEnum.lan) {
    dropdownMenuOptions.push({
      label: () => renderIconLabel('打开外网链接', 'mdi:wan'),
      key: 'openWanUrl',
    })

    dropdownMenuOptions.push({
      label: () => renderIconLabel('复制外网链接', 'ion-copy'),
      key: 'copyWanUrl',
    })
  }

  if (authStore.visitMode === VisitMode.VISIT_MODE_LOGIN) {
    dropdownMenuOptions.push({
      label: () => renderIconLabel(t('common.edit'), 'basil:edit-solid'),
      key: 'edit',
    }, {
      label: () => renderIconLabel(t('common.delete'), 'material-symbols:delete'),
      key: 'delete',
    })
  }

  return dropdownMenuOptions
}

onMounted(() => {
  updateIsMobile()
  window.addEventListener('resize', updateIsMobile)

  // 更新用户信息
  updateLocalUserInfo()
  getList()

  // 更新同步云端配置
  panelState.updatePanelConfigByCloud()

  // 设置标题
  if (panelState.panelConfig.logoText)
    setTitle(panelState.panelConfig.logoText)

  window.addEventListener('resize', recalcCatalogDots)
})

onUnmounted(() => {
  window.removeEventListener('resize', recalcCatalogDots)
  window.removeEventListener('resize', updateIsMobile)
})

// 前端搜索过滤
const currentSearchKeyword = ref('')

function itemFrontEndSearch(keyword?: string) {
  currentSearchKeyword.value = keyword ?? ''
  keyword = currentSearchKeyword.value.trim()
  if (keyword !== '' && panelState.panelConfig.searchBoxSearchIcon) {
    const filteredData = ref<ItemGroup[]>([])
    for (let i = 0; i < items.value.length; i++) {
      const element = items.value[i].items?.filter((item: Panel.ItemInfo) => {
        return (
          item.title.toLowerCase().includes(keyword?.toLowerCase() ?? '')
          || item.url.toLowerCase().includes(keyword?.toLowerCase() ?? '')
          || item.description?.toLowerCase().includes(keyword?.toLowerCase() ?? '')
        )
      })
      if (element && element.length > 0)
        filteredData.value.push({ ...items.value[i], items: element, hoverStatus: false })
    }
    filterItems.value = filteredData.value
  }
  else {
    filterItems.value = items.value
  }
}

function handleSetHoverStatus(groupIndex: number, hoverStatus: boolean) {
  if (items.value[groupIndex])
    items.value[groupIndex].hoverStatus = hoverStatus
}

function handleSetSortStatus(groupIndex: number, sortStatus: boolean) {
  if (items.value[groupIndex])
    items.value[groupIndex].sortStatus = sortStatus

  // 并未保存排序重新更新数据
  if (!sortStatus) {
    // 单独更新组
    if (items.value[groupIndex] && items.value[groupIndex].id)
      updateItemIconGroupByNet(groupIndex, items.value[groupIndex].id as number)
  }
}

function handleEditItem(item: Panel.ItemInfo) {
  editItemInfoData.value = item
  editItemInfoShow.value = true
  currentAddItenIconGroupId.value = undefined
}

function handleAddItem(itemIconGroupId?: number) {
  editItemInfoData.value = null
  editItemInfoShow.value = true
  if (itemIconGroupId)
    currentAddItenIconGroupId.value = itemIconGroupId
}

function handleSwitchGroupType(groupType: 'website' | 'webpage') {
  currentGroupType.value = groupType
  getList()
}

function setGroupRef(el: any, groupId?: number) {
  if (!groupId)
    return
  groupRefs.value[groupId] = el as HTMLElement | null
}

function scrollToGroup(groupId?: number) {
  if (!groupId || !scrollContainerRef.value)
    return
  const target = groupRefs.value[groupId]
  if (!target)
    return
  const container = scrollContainerRef.value
  const top = target.offsetTop - 20
  container.scrollTo({ top, behavior: 'smooth' })
}

function recalcCatalogDots() {
  if (!scrollContainerRef.value)
    return

  const container = scrollContainerRef.value
  const total = Math.max(1, container.scrollHeight)
  const nextMap: Record<number, number> = {}

  for (const key of Object.keys(groupRefs.value)) {
    const id = Number(key)
    const target = groupRefs.value[id]
    if (!target)
      continue
    const ratio = target.offsetTop / total
    nextMap[id] = Math.max(0, Math.min(100, ratio * 100))
  }

  groupDotMap.value = nextMap
}

function getGroupDotTop(groupId?: number) {
  if (!groupId)
    return '0%'
  const v = groupDotMap.value[groupId] ?? 0
  return `${v}%`
}

</script>

<template>
  <div class="w-full h-full sun-main">
    <div
      class="cover wallpaper" :style="{
        filter: `blur(${panelState.panelConfig.backgroundBlur}px)`,
        background: `url(${panelState.panelConfig.backgroundImageSrc}) no-repeat`,
        backgroundSize: 'cover',
        backgroundPosition: 'center',
      }"
    />
    <div class="mask" :style="{ backgroundColor: `rgba(0,0,0,${panelState.panelConfig.backgroundMaskNumber})` }" />
    <div
      ref="scrollContainerRef"
      class="absolute w-full h-full overflow-y-auto overflow-x-hidden"
      :class="panelState.panelConfig.rightScrollBarShow === false ? 'hide-native-scrollbar' : 'show-native-scrollbar'"
    >
      <div
        class="p-2.5 mx-auto"
        :style="{
          marginTop: `${panelState.panelConfig.marginTop}%`,
          marginBottom: `${panelState.panelConfig.marginBottom}%`,
          paddingTop: `${panelState.panelConfig.paddingTop || 0}px`,
          paddingBottom: `${panelState.panelConfig.paddingBottom || 0}px`,
          maxWidth: (panelState.panelConfig.maxWidth ?? '1200') + panelState.panelConfig.maxWidthUnit,
        }"
      >
        <!-- 头 -->
        <div class="mx-[auto] w-[80%]">
          <div v-if="panelState.panelConfig.topHeaderShow" class="flex mx-[auto] items-center justify-center text-white">
            <div class="logo">
              <span class="text-2xl md:text-6xl font-bold text-shadow">
                {{ panelState.panelConfig.logoText }}
              </span>
            </div>
            <div class="divider text-base lg:text-2xl mx-[10px]">
              |
            </div>
            <div class="text-shadow">
              <Clock :hide-second="!panelState.panelConfig.clockShowSecond" />
            </div>
          </div>
          <div
            v-if="panelState.panelConfig.searchBoxShow"
            class="flex items-center justify-between mt-[20px] mx-auto sm:w-full lg:w-[80%]"
            :style="{ marginBottom: `${panelState.panelConfig.searchToBelowGapPx ?? 30}px` }"
          >
            <div class="flex-1">
              <SearchBox @itemSearch="itemFrontEndSearch" />
            </div>
            <!-- 移动端：登录按钮（右对齐、更大） -->
            <div v-if="isMobile" class="ml-2">
              <NButton
                color="#2a2a2a6b"
                class="!h-10 !w-10"
                circle
                @click="router.push('/login')"
              >
                <template #icon>
                  <SvgIcon class="text-2xl text-white" icon="material-symbols:account-circle" />
                </template>
              </NButton>
            </div>
          </div>
        </div>

        <!-- 应用盒子 -->
        <div :style="{ marginLeft: `${panelState.panelConfig.marginX}px`, marginRight: `${panelState.panelConfig.marginX}px` }">
          <div class="flex justify-center text-white items-center" :style="{ marginTop: `${panelState.panelConfig.searchBoxShow ? 0 : 20}px`, marginBottom: '0px' }">
            <div
              class="cursor-pointer px-4 py-1 rounded-l-full hover:text-[#fef08a] transition-colors"
              :class="currentGroupType === 'website' ? 'bg-white/30' : 'bg-white/10'"
              @click="handleSwitchGroupType('website')"
            >
              网站导航
            </div>
            <div class="px-2">|</div>
            <div
              class="cursor-pointer px-4 py-1 rounded-r-full hover:text-[#fef08a] transition-colors"
              :class="currentGroupType === 'webpage' ? 'bg-white/30' : 'bg-white/10'"
              @click="handleSwitchGroupType('webpage')"
            >
              网页收藏
            </div>
          </div>
          <!-- 系统监控状态 -->
          <div
            v-if="panelState.panelConfig.systemMonitorShow
              && ((panelState.panelConfig.systemMonitorPublicVisitModeShow && authStore.visitMode === VisitMode.VISIT_MODE_PUBLIC)
                || authStore.visitMode === VisitMode.VISIT_MODE_LOGIN)"
            class="flex mx-auto"
          >
            <SystemMonitor
              :allow-edit="authStore.visitMode === VisitMode.VISIT_MODE_LOGIN"
              :show-title="panelState.panelConfig.systemMonitorShowTitle"
            />
          </div>

          <!-- 组纵向排列 -->
          <div
            v-for="(itemGroup, itemGroupIndex) in filterItems" :key="itemGroupIndex"
            :ref="(el) => setGroupRef(el, itemGroup.id)"
            class="item-list"
            :style="{
              marginTop: `${itemGroupIndex === 0 ? (panelState.panelConfig.tabsToFirstGroupGapPx ?? 40) : (panelState.panelConfig.groupToGroupGapPx ?? 40)}px`,
            }"
            :class="itemGroup.sortStatus ? 'shadow-2xl border shadow-[0_0_30px_10px_rgba(0,0,0,0.3)]  p-[10px] rounded-2xl' : ''"
            @mouseenter="handleSetHoverStatus(itemGroupIndex, true)"
            @mouseleave="handleSetHoverStatus(itemGroupIndex, false)"
          >
            <!-- 分组标题 -->
            <div class="text-white text-xl font-extrabold ml-[10px] flex items-center" :style="{ marginBottom: '20px' }">
              <span class="group-title text-shadow flex items-center gap-2">
                {{ itemGroup.title }}
                <span v-if="itemGroup.items" class="text-xs bg-white/20 px-2 py-0.5 rounded-full font-normal">
                  {{ itemGroup.items.length }}
                </span>
              </span>
              <div
                v-if="authStore.visitMode === VisitMode.VISIT_MODE_LOGIN"
                class="group-buttons ml-2 delay-100 transition-opacity flex"
                :class="itemGroup.hoverStatus ? 'opacity-100' : 'opacity-0'"
              >
                <span class="mr-2 cursor-pointer" :title="t('common.add')" @click="handleAddItem(itemGroup.id)">
                  <SvgIcon class="text-white font-xl" icon="typcn:plus" />
                </span>
                <span class="mr-2 cursor-pointer " :title="t('common.sort')" @click="handleSetSortStatus(itemGroupIndex, !itemGroup.sortStatus)">
                  <SvgIcon class="text-white font-xl" icon="ri:drag-drop-line" />
                </span>
              </div>
            </div>

            <!-- 网页分组：默认列表（每行一个） -->
            <div v-if="currentGroupType === 'webpage'">
              <div v-if="itemGroup.items">
                <VueDraggable
                  v-model="itemGroup.items" item-key="sort" :animation="300"
                  :class="[
                    'w-full webpage-list-container show-native-scrollbar',
                    (itemGroup.items && itemGroup.items.length > 10) ? 'webpage-list-fixed' : ''
                  ]"
                  :style="{ gap: `${panelState.panelConfig.webpageListGapPx ?? 10}px` }"
                  filter=".not-drag"
                  :disabled="!itemGroup.sortStatus"
                >
                  <div
                    v-for="item, index in (itemGroup.items || [])"
                    :key="index"
                    class="group relative w-full h-12 px-3 rounded-lg flex justify-between items-center hover:text-[#fef08a] transition-colors shrink-0"
                    :class="[
                      itemGroup.sortStatus ? 'cursor-move' : 'cursor-pointer',
                      item.pinned ? 'text-sky-300 bg-black/20' : 'text-white bg-black/20'
                    ]"
                    :title="item.description"
                    @click="handleItemClick(itemGroupIndex, item)"
                    @contextmenu="(e) => handleContextMenu(e, itemGroupIndex, item)"
                  >
                    <!-- 通过 Vue 指令判断文字是否截断再显示原生 title（此处使用封装一个只在溢出时显示 title 的逻辑较繁琐，妥协为外层不加 title，或通过 CSS tooltip 实现，这里用最简单的 css hover + ellipsis 判断比较难，咱们改用自定义组件或者根据字符长度估算，这里用简单的 title 移除，靠用户自行判断，但用户需要“全标题显示” -->
                      <NEllipsis 
                        class="webpage-title flex-1" 
                        :tooltip="{ 
                          content: item.title, 
                          placement: 'bottom',
                          style: { maxWidth: `calc(${(panelState.panelConfig.maxWidth ?? '1200') + (panelState.panelConfig.maxWidthUnit ?? 'px')} * 0.9)`, whiteSpace: 'normal', wordBreak: 'break-all' }
                        }" 
                        :line-clamp="1"
                      >
                      {{ item.title }}
                    </NEllipsis>
                    
                    <!-- 悬浮操作按钮（移动端不显示） -->
                    <div v-if="!itemGroup.sortStatus && !isMobile" class="opacity-0 group-hover:opacity-100 flex items-center gap-2 transition-opacity ml-2">
                      <div 
                        class="p-1 rounded bg-black/20 hover:bg-black/40 cursor-pointer flex items-center justify-center text-white"
                        :title="item.pinned ? '取消置顶' : '置顶'"
                        @click.stop="quickTogglePinWebpage(item)"
                      >
                        <SvgIconOnline class="text-sm" :icon="item.pinned ? 'mdi:pin-off' : 'mdi:pin'" />
                      </div>
                      <div 
                        class="p-1 rounded bg-black/20 hover:bg-black/40 cursor-pointer flex items-center justify-center text-white"
                        title="修改"
                        @click.stop="quickEditWebpage(item)"
                      >
                        <SvgIcon class="text-sm" icon="basil:edit-solid" />
                      </div>
                      <div 
                        class="p-1 rounded bg-black/20 hover:bg-black/40 cursor-pointer flex items-center justify-center text-white"
                        title="删除"
                        @click.stop="quickDeleteWebpage(item)"
                      >
                        <SvgIcon class="text-sm" icon="material-symbols:delete" />
                      </div>
                    </div>
                  </div>

                  <div v-if="itemGroup.items.length === 0" class="not-drag text-white/70 py-2 px-3">
                    {{ $t('common.add') }}
                  </div>
                </VueDraggable>
              </div>
            </div>

            <!-- 网站分组：保留原图标样式 -->
            <div v-else>
              <!-- 详情图标 -->
              <div v-if="panelState.panelConfig.iconStyle === PanelPanelConfigStyleEnum.info">
                <div v-if="itemGroup.items">
                  <VueDraggable
                    v-model="itemGroup.items" item-key="sort" :animation="300"
                    class="icon-info-box"
                    filter=".not-drag"
                    :disabled="!itemGroup.sortStatus"
                  >
                    <div v-for="item, index in itemGroup.items" :key="index" :title="item.description" @contextmenu="(e) => handleContextMenu(e, itemGroupIndex, item)">
                      <AppIcon
                        :class="itemGroup.sortStatus ? 'cursor-move' : 'cursor-pointer'"
                        :item-info="item"
                        :icon-text-color="panelState.panelConfig.iconTextColor"
                        :icon-text-info-hide-description="panelState.panelConfig.iconTextInfoHideDescription || false"
                        :icon-text-icon-hide-title="panelState.panelConfig.iconTextIconHideTitle || false"
                        :style="0"
                        @click="handleItemClick(itemGroupIndex, item)"
                      />
                    </div>

                    <div v-if="itemGroup.items.length === 0" class="not-drag">
                      <AppIcon
                        :class="itemGroup.sortStatus ? 'cursor-move' : 'cursor-pointer'"
                        :item-info="{ icon: { itemType: 3, text: 'subway:add' }, title: t('common.add'), url: '', openMethod: 0 }"
                        :icon-text-color="panelState.panelConfig.iconTextColor"
                        :icon-text-info-hide-description="panelState.panelConfig.iconTextInfoHideDescription || false"
                        :icon-text-icon-hide-title="panelState.panelConfig.iconTextIconHideTitle || false"
                        :style="0"
                        @click="handleAddItem(itemGroup.id)"
                      />
                    </div>
                  </VueDraggable>
                </div>
              </div>

              <!-- APP图标宫型盒子 -->
              <div v-if="panelState.panelConfig.iconStyle === PanelPanelConfigStyleEnum.icon">
                <div v-if="itemGroup.items">
                  <VueDraggable
                    v-model="itemGroup.items" item-key="sort" :animation="300"
                    class="icon-small-box"

                    filter=".not-drag"
                    :disabled="!itemGroup.sortStatus"
                  >
                    <div v-for="item, index in itemGroup.items" :key="index" :title="item.description" @contextmenu="(e) => handleContextMenu(e, itemGroupIndex, item)">
                      <AppIcon
                        :class="itemGroup.sortStatus ? 'cursor-move' : 'cursor-pointer'"
                        :item-info="item"
                        :icon-text-color="panelState.panelConfig.iconTextColor"
                        :icon-text-info-hide-description="!panelState.panelConfig.iconTextInfoHideDescription"
                        :icon-text-icon-hide-title="panelState.panelConfig.iconTextIconHideTitle || false"
                        :style="1"
                        @click="handleItemClick(itemGroupIndex, item)"
                      />
                    </div>

                    <div v-if="itemGroup.items.length === 0" class="not-drag">
                      <AppIcon
                        class="cursor-pointer"
                        :item-info="{ icon: { itemType: 3, text: 'subway:add' }, title: $t('common.add'), url: '', openMethod: 0 }"
                        :icon-text-color="panelState.panelConfig.iconTextColor"
                        :icon-text-info-hide-description="!panelState.panelConfig.iconTextInfoHideDescription"
                        :icon-text-icon-hide-title="panelState.panelConfig.iconTextIconHideTitle || false"
                        :style="1"
                        @click="handleAddItem(itemGroup.id)"
                      />
                    </div>
                  </vuedraggable>
                </div>
              </div>
            </div>

            <!-- 编辑栏 -->
            <div v-if="itemGroup.sortStatus" class="flex mt-[10px]">
              <div>
                <NButton color="#2a2a2a6b" @click="handleSaveSort(itemGroup)">
                  <template #icon>
                    <SvgIcon class="text-white font-xl" icon="material-symbols:save" />
                  </template>
                  <div>
                    {{ $t('common.saveSort') }}
                  </div>
                </NButton>
              </div>
            </div>
          </div>
        </div>
        <div class="mt-5 footer" v-html="panelState.panelConfig.footerHtml" />
      </div>
    </div>

    <!-- 左侧分组目录 -->
    <div v-if="panelState.panelConfig.leftCatalogShow !== false && !isMobile" class="left-catalog">
      <div class="left-catalog-track">
        <div
          v-for="group in filterItems"
          :key="`catalog-${group.id}`"
          class="catalog-dot"
          :style="{ top: getGroupDotTop(group.id) }"
          @click="scrollToGroup(group.id)"
        >
          <span class="catalog-label">{{ group.title }}</span>
        </div>
      </div>
    </div>

    <!-- 右键菜单 -->
    <NDropdown
      placement="bottom-start" trigger="manual" :x="dropdownMenuX" :y="dropdownMenuY"
      :options="getDropdownMenuOptions()" :show="dropdownShow" :on-clickoutside="onClickoutside" @select="handleRightMenuSelect"
    />

    <!-- 悬浮按钮 -->
    <div v-if="!isMobile" class="fixed-element" :style="fixedElementStyle">
      <div class="fixed-tools-shell">
        <div class="fixed-tools-wrapper">
          <NButtonGroup vertical>
        <!-- 网络模式切换按钮组 -->
        <NButton
          v-if="panelState.networkMode === PanelStateNetworkModeEnum.lan && panelState.panelConfig.netModeChangeButtonShow" color="#2a2a2a6b"
          :title="t('panelHome.changeToWanModel')" @click="handleChangeNetwork(PanelStateNetworkModeEnum.wan)"
        >
          <template #icon>
            <SvgIcon class="text-white font-xl" icon="material-symbols:lan-outline-rounded" />
          </template>
        </NButton>

        <NButton
          v-if="panelState.networkMode === PanelStateNetworkModeEnum.wan && panelState.panelConfig.netModeChangeButtonShow" color="#2a2a2a6b"
          :title="t('panelHome.changeToLanModel')" @click="handleChangeNetwork(PanelStateNetworkModeEnum.lan)"
        >
          <template #icon>
            <SvgIcon class="text-white font-xl" icon="mdi:wan" />
          </template>
        </NButton>

        <NButton v-if="authStore.visitMode === VisitMode.VISIT_MODE_LOGIN" color="#2a2a2a6b" @click="settingModalShow = !settingModalShow">
          <template #icon>
            <SvgIcon class="text-white font-xl" icon="majesticons-applications" />
          </template>
        </NButton>

        <NButton v-if="authStore.visitMode === VisitMode.VISIT_MODE_PUBLIC" color="#2a2a2a6b" :title="$t('panelHome.goToLogin')" @click="router.push('/login')">
          <template #icon>
            <SvgIcon class="text-white font-xl" icon="material-symbols:account-circle" />
          </template>
        </NButton>
          </NButtonGroup>
        </div>
      </div>

      <AppStarter v-model:visible="settingModalShow" />
      <!-- <Setting v-model:visible="settingModalShow" /> -->
    </div>

    <NBackTop
      :listen-to="() => scrollContainerRef"
      :right="panelState.panelConfig.rightScrollBarShow === false ? 10 : 24"
      :bottom="Number(panelState.panelConfig.paddingBottom ?? 10)"
      style="background-color:transparent;border: none;box-shadow: none;"
    >
      <div class="backtop-shell">
        <div class="backtop-wrapper">
          <NButton color="#2a2a2a6b">
          <template #icon>
            <SvgIcon class="text-white font-xl" icon="icon-park-outline:to-top" />
          </template>
          </NButton>
        </div>
      </div>
    </NBackTop>

    <EditItem
      v-model:visible="editItemInfoShow"
      :item-info="editItemInfoData"
      :item-group-id="currentAddItenIconGroupId"
      :group-type="currentGroupType"
      @done="handleEditSuccess"
    />

    <!-- 弹窗 -->
    <NModal
      v-model:show="windowShow" :mask-closable="false" preset="card"
      style="max-width: 1000px;height: 600px;border-radius: 1rem;" :bordered="true" size="small" role="dialog"
      aria-modal="true"
    >
      <template #header>
        <div class="flex items-center">
          <span class="mr-[20px]">
            {{ windowTitle }}
          </span>

          <NSpin v-if="windowIframeIsLoad" size="small" />
        </div>
      </template>
      <div class="w-full h-full rounded-2xl overflow-hidden border dark:border-zinc-700">
        <div v-if="windowIframeIsLoad" class="flex flex-col p-5">
          <NSkeleton height="50px" width="100%" class="rounded-lg" />
          <NSkeleton height="180px" width="100%" class="mt-[20px] rounded-lg" />
          <NSkeleton height="180px" width="100%" class="mt-[20px] rounded-lg" />
        </div>
        <iframe
          v-show="!windowIframeIsLoad" id="windowIframeId" ref="windowIframeRef" :src="windowSrc"
          class="w-full h-full" frameborder="0" @load="handWindowIframeIdLoad"
        />
      </div>
    </NModal>
  </div>
</template>

<style>
body,
html {
  overflow: hidden;
  background-color: rgb(54, 54, 54);
}
</style>

<style scoped>
.mask {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.sun-main {
  user-select: none;
}

.cover {
  position: absolute;
  width: 100%;
  height: 100%;
  overflow: hidden;
  /* background: url(@/assets/start_sky.jpg) no-repeat; */

  transform: scale(1.05);
}

.text-shadow {
  text-shadow: 2px 2px 50px rgb(0, 0, 0);
}

.app-icon-text-shadow {
  text-shadow: 2px 2px 5px rgb(0, 0, 0);
}

.fixed-element {
  position: fixed;
  z-index: 20;
}

.fixed-tools-shell,
.backtop-shell {
  border-radius: var(--sp-ui-radius);
  /* 用 shell 做阴影 + 圆角，避免出现“矩形阴影尖角” */
  box-shadow: 0 0 10px 2px rgba(0, 0, 0, 0.2);
  background: transparent;
}

.fixed-tools-wrapper,
.backtop-wrapper {
  border-radius: var(--sp-ui-radius);
  overflow: hidden;
  background: transparent;
}

/* 让圆角由外层容器控制，避免 NaiveUI 默认小圆角和用户风格不一致 */
:global(.fixed-tools-wrapper .n-button),
:global(.backtop-wrapper .n-button) {
  border-radius: 0 !important;
}

/* 避免 NButtonGroup/NButton 自带的背景在圆角外“露底色” */
:global(.fixed-tools-wrapper .n-button-group),
:global(.backtop-wrapper .n-button-group) {
  background: transparent !important;
}

.icon-info-box {
  width: 100%;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 18px;

}

.icon-small-box {
  width: 100%;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(75px, 1fr));
  gap: 18px;

}

@media (max-width: 500px) {
  .icon-info-box{
    grid-template-columns: repeat(4, 1fr);
    gap: 6px;
  }
  .icon-small-box{
    grid-template-columns: repeat(4, 1fr);
    gap: 6px;
  }

  /* 移动端左侧目录不要贴着内容 */
  .left-catalog {
    left: 4px;
  }
  .p-2\.5 {
    padding-left: 10px !important;
    padding-right: 10px !important;
  }

  /* 移动端图标整体进一步缩小 */
  .icon-info-box > div,
  .icon-small-box > div {
    padding: 0 !important;
    overflow: hidden;
    box-sizing: border-box;
  }
  .icon-info-box,
  .icon-small-box {
    overflow-x: hidden;
  }
  .icon-info-box :deep(.icon),
  .icon-small-box :deep(.icon) {
    width: 24px !important;
    height: 24px !important;
  }
  .icon-info-box :deep(.text),
  .icon-small-box :deep(.text) {
    font-size: 8px !important;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

.left-catalog {
  position: fixed;
  left: 8px;
  top: 12px;
  bottom: 12px;
  z-index: 20;
  display: flex;
  align-items: stretch;
}

.left-catalog-track {
  height: 100%;
  width: var(--sp-scrollbar-width, 14px);
  position: relative;
  box-sizing: border-box; /* 保证含边框后的实际宽度仍是设定值 */
  /* 恢复之前：不加背景、不加中线；只保留轻微灰色包边 */
  border: 1px solid rgba(148, 163, 184, 0.35);
  border-radius: var(--sp-ui-radius);
}

.catalog-dot {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: #ffffff;
  border: 2px solid rgba(255, 255, 255, 0.9);
  position: absolute;
  left: 2px;
  transform: translateY(-50%);
  cursor: pointer;
}

.catalog-label {
  position: absolute;
  left: 18px;
  top: 50%;
  transform: translateY(-50%);
  background: rgba(0, 0, 0, 0.35);
  color: #fff;
  border-radius: var(--sp-ui-radius); /* 跟全局一致，不要半圆角 */
  padding: 6px 12px;
  white-space: nowrap;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.2s, color 0.2s;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

/* 去掉名称前面的点 */
.catalog-label::before {
  content: none;
}

.catalog-dot:hover .catalog-label {
  opacity: 1;
  color: #fef08a;
}

/* .webpage-title 的相关样式可以去掉，交由 NEllipsis 处理 */
.webpage-title {
  text-align: left;
  line-height: 1.5;
  word-break: break-all;
}

.webpage-list-container {
  /* 默认：按内容自然高度（用于少于等于10条的分组，避免把分组之间撑出大空白） */
  overflow-x: hidden;
  padding-right: 4px; /* 防止滚动条紧贴文字 */
  gap: 10px; /* 网页列表条目间距 */
  display: flex;
  flex-direction: column;
}

.webpage-list-fixed {
  /* 只有当条目 > 10 时才固定为10行高度，并在容器内滚动 */
  height: calc(3rem * 10 + 0.25rem * 9);
  max-height: calc(3rem * 10 + 0.25rem * 9);
  overflow-y: auto;
}

/* 覆盖局部样式，使滚动条颜色在深色背景下协调 */
.webpage-list-container::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2) !important;
}
.webpage-list-container::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.4) !important;
}

:global(:root) {
  --sp-scrollbar-width: 14px;
  --sp-ui-radius: 10px;
}

:global(.show-native-scrollbar) {
  /* Firefox: thin 会强制变窄，导致和左侧目录条宽度不一致 */
  scrollbar-width: auto;
  scrollbar-color: rgba(148, 163, 184, 0.9) transparent;
}

:global(.show-native-scrollbar::-webkit-scrollbar) {
  width: var(--sp-scrollbar-width);
}

:global(.show-native-scrollbar::-webkit-scrollbar-track) {
  background: transparent;
  margin: 48px 0;
}

:global(.show-native-scrollbar::-webkit-scrollbar-thumb) {
  background: rgba(148, 163, 184, 0.92);
  border-radius: 999px;
  border: 3px solid transparent;
  background-clip: content-box;
}

:global(.show-native-scrollbar::-webkit-scrollbar-thumb:hover) {
  background: rgba(203, 213, 225, 0.98);
  background-clip: content-box;
}

:global(.hide-native-scrollbar) {
  scrollbar-width: none;
}

:global(.hide-native-scrollbar::-webkit-scrollbar) {
  width: 0;
  height: 0;
}
</style>
