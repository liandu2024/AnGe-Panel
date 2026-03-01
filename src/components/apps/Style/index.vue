<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { UploadFileInfo } from 'naive-ui'
import { NButton, NButtonGroup, NCard, NColorPicker, NGrid, NGridItem, NInput, NInputGroup, NPopconfirm, NSelect, NSpin, NSlider, NSwitch, NUpload, NUploadDragger, useMessage } from 'naive-ui'
import { getList } from '@/api/system/file'
import { RoundCardModal } from '@/components/common'
import { useAuthStore, usePanelState } from '@/store'
import { set as setUserConfig } from '@/api/panel/userConfig'
import { PanelPanelConfigStyleEnum } from '@/enums/panel'
import { t } from '@/locales'

const authStore = useAuthStore()
const panelState = usePanelState()
const ms = useMessage()
const showWallpaperInput = ref(false)

const wallpaperPickerShow = ref(false)
const wallpaperPickerLoading = ref(false)
const wallpaperList = ref<File.Info[]>([])
const wallpaperCountLabel = computed(() => wallpaperList.value.length)

async function loadWallpaperList() {
  wallpaperPickerLoading.value = true
  try {
    const { data } = await getList<Common.ListResponse<File.Info[]>>('wallpaper')
    wallpaperList.value = data.list || []
  }
  finally {
    wallpaperPickerLoading.value = false
  }
}

function openWallpaperPicker() {
  wallpaperPickerShow.value = true
  if (wallpaperList.value.length === 0)
    loadWallpaperList()
}

function handlePickWallpaper(src: string) {
  panelState.panelConfig.backgroundImageSrc = src
  setUserConfig({ panel: panelState.panelConfig })
  wallpaperPickerShow.value = false
}

const isSaveing = ref(false)

const iconTypeOptions = [
  {
    label: t('apps.baseSettings.detailIcon'),
    value: PanelPanelConfigStyleEnum.info,
  },
  {
    label: t('apps.baseSettings.smallIcon'),
    value: PanelPanelConfigStyleEnum.icon,
  },
]

const maxWidthUnitOption = [
  {
    label: 'px',
    value: 'px',
  },
  {
    label: '%',
    value: '%',
  },
]

const floatingToolsPositionOptions = [
  { label: '右上角', value: 'right-top' },
  { label: '右中间', value: 'right-middle' },
  { label: '右下角', value: 'right-bottom' },
]

watch(panelState.panelConfig, () => {
  if (!isSaveing.value) {
    isSaveing.value = true

    setTimeout(() => {
      panelState.recordState()// 本地记录
      isSaveing.value = false
      uploadCloud()
    }, 1000)
  }
})

function handleUploadBackgroundFinish({
  file,
  event,
}: {
  file: UploadFileInfo
  event?: ProgressEvent
}) {
  const res = JSON.parse((event?.target as XMLHttpRequest).response)
  panelState.panelConfig.backgroundImageSrc = res.data.imageUrl
  return file
}

function uploadCloud() {
  setUserConfig({ panel: panelState.panelConfig }).then((res) => {
    if (res.code === 0)
      ms.success(t('apps.baseSettings.configSaved'))
    else
      ms.error(t('apps.baseSettings.configFailed', { message: res.msg }))
  })
}

function resetPanelConfig() {
  panelState.resetPanelConfig()
  uploadCloud()
}
</script>

<template>
  <div class="bg-slate-200 dark:bg-zinc-900 rounded-[10px] p-[8px] overflow-auto">
    <NCard style="border-radius:10px" size="small">
      <div class="text-slate-500 mb-[5px] font-bold">
        LOGO
      </div>

      <div>
        <div class="flex items-center mt-[5px] mb-[5px]">
          <span class="mr-[10px]">显示顶部标题/时间</span>
          <NSwitch v-model:value="panelState.panelConfig.topHeaderShow" />
        </div>
        <div>
          {{ $t('apps.baseSettings.textContent') }}
        </div>
        <div class="flex items-center mt-[5px]">
          <NInput v-model:value="panelState.panelConfig.logoText" type="text" show-count :maxlength="20" placeholder="请输入文字" />
        </div>
      </div>
    </NCard>

    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <div class="text-slate-500 mb-[5px] font-bold">
        {{ $t('apps.baseSettings.clock') }}
      </div>
      <div class="flex items-center mt-[5px]">
        <span class="mr-[10px]">{{ $t('apps.baseSettings.clockSecondShow') }}</span>
        <NSwitch v-model:value="panelState.panelConfig.clockShowSecond" />
      </div>
    </NCard>

    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <div class="text-slate-500 mb-[5px] font-bold">
        {{ $t('apps.baseSettings.searchBar') }}
      </div>
      <div class="flex items-center mt-[5px]">
        <span class="mr-[10px]">{{ $t('common.show') }}</span>
        <NSwitch v-model:value="panelState.panelConfig.searchBoxShow" />
      </div>
      <div v-if="panelState.panelConfig.searchBoxShow" class="flex items-center mt-[5px]">
        <span class="mr-[10px]">{{ $t('apps.baseSettings.searchBarSearchItem') }}</span>
        <NSwitch v-model:value="panelState.panelConfig.searchBoxSearchIcon" />
      </div>
    </NCard>

    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <div class="text-slate-500 mb-[5px] font-bold">
        {{ $t('apps.baseSettings.systemMonitorStatus') }}
      </div>
      <div class="flex items-center mt-[5px]">
        <span class="mr-[10px]">{{ $t('common.show') }}</span>
        <NSwitch v-model:value="panelState.panelConfig.systemMonitorShow" />
      </div>
      <div v-if="panelState.panelConfig.systemMonitorShow" class="flex items-center mt-[5px]">
        <span class="mr-[10px]">{{ $t('apps.baseSettings.showTitle') }}</span>
        <NSwitch v-model:value="panelState.panelConfig.systemMonitorShowTitle" />
      </div>
      <div v-if="panelState.panelConfig.systemMonitorShow" class="flex items-center mt-[5px]">
        <span class="mr-[10px]">{{ $t('apps.baseSettings.publicVisitModeShow') }}</span>
        <NSwitch v-model:value="panelState.panelConfig.systemMonitorPublicVisitModeShow" />
      </div>
    </NCard>

    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <div class="text-slate-500 mb-[5px] font-bold">
        {{ $t('common.icon') }}
      </div>
      <div class="mt-[5px]">
        <div>
          {{ $t('common.style') }}
        </div>
        <div class="flex items-center mt-[5px]">
          <NSelect v-model:value="panelState.panelConfig.iconStyle" :options="iconTypeOptions" />
        </div>
      </div>

      <div v-if="panelState.panelConfig.iconStyle === PanelPanelConfigStyleEnum.info" class="mt-[5px]">
        <div>
          {{ $t('apps.baseSettings.hideDescription') }}
        </div>
        <div class="flex items-center mt-[5px]">
          <NSwitch v-model:value="panelState.panelConfig.iconTextInfoHideDescription" />
        </div>
      </div>

      <div v-if="panelState.panelConfig.iconStyle === PanelPanelConfigStyleEnum.icon" class="mt-[5px]">
        <div>
          {{ $t('apps.baseSettings.hideTitle') }}
        </div>
        <div class="flex items-center mt-[5px]">
          <NSwitch v-model:value="panelState.panelConfig.iconTextIconHideTitle" />
        </div>
      </div>

      <div class="mt-[5px]">
        <div>
          {{ $t('common.textColor') }}
        </div>
        <div class="flex items-center mt-[5px]">
          <NColorPicker
            v-model:value="panelState.panelConfig.iconTextColor"
            :show-alpha="false"
            size="small"
            :modes="['hex']"
            :swatches="[
              '#000000',
              '#ffffff',
              '#18A058',
              '#2080F0',
              '#F0A020',
            ]"
          />
        </div>
      </div>
    </NCard>
    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <div class="text-slate-500 mb-[5px] font-bold">
        {{ $t('apps.baseSettings.wallpaper') }}
      </div>
      <NUpload
        action="/api/file/uploadImg?fileType=wallpaper"
        :show-file-list="false"
        name="imgfile"
        :headers="{
          token: authStore.token as string,
        }"
        :directory-dnd="true"
        @finish="handleUploadBackgroundFinish"
      >
        <NUploadDragger style="width: 100%;">
          <div
            class="h-[200px] w-full border bg-slate-100 flex justify-center items-center cursor-pointer rounded-[10px]"
            :style="{ background: `url(${panelState.panelConfig.backgroundImageSrc}) no-repeat`, backgroundSize: 'cover' }"
          >
            <div class="text-shadow text-white">
              {{ $t('apps.baseSettings.uploadOrDragText') }}
            </div>
          </div>
        </NUploadDragger>
      </NUpload>

      <div class="flex justify-end mt-2">
        <NButton size="small" quaternary type="info" @click="openWallpaperPicker">
          选择已上传壁纸（{{ wallpaperCountLabel }}）
        </NButton>
      </div>

      <RoundCardModal v-model:show="wallpaperPickerShow" style="max-width: 900px;" size="small" title="选择已上传壁纸">
        <div class="min-h-[180px]">
          <div class="flex items-center justify-between mb-2">
            <NButtonGroup size="small">
              <NButton @click="loadWallpaperList" :loading="wallpaperPickerLoading">
                刷新
              </NButton>
            </NButtonGroup>
          </div>

          <NSpin v-show="wallpaperPickerLoading" size="small" />
          <div v-if="!wallpaperPickerLoading && wallpaperList.length === 0" class="text-slate-500 text-sm">
            暂无已上传壁纸
          </div>

          <NGrid v-else cols="2 500:3 800:4" :x-gap="10" :y-gap="10">
            <NGridItem v-for="(it, idx) in wallpaperList" :key="`wall-pick-${idx}`">
              <div class="wall-pick-card" @click="handlePickWallpaper(it.src)">
                <div class="wall-pick-thumb" :style="{ background: `url(${it.src}) center / cover no-repeat` }" />
                <div class="wall-pick-name">{{ it.fileName }}</div>
              </div>
            </NGridItem>
          </NGrid>
        </div>
      </RoundCardModal>

      <div class="flex items-center mt-[5px]">
        <span class="mr-[10px]">{{ $t('apps.baseSettings.customImageAddress') }}</span>
        <NSwitch v-model:value="showWallpaperInput" />
      </div>
      <div v-if="showWallpaperInput" class="mt-1">
        <NInput v-model:value="panelState.panelConfig.backgroundImageSrc" type="text" size="small" clearable />
      </div>

      <div class="flex items-center mt-[10px]">
        <span class="mr-[10px]">{{ $t('apps.baseSettings.vague') }}</span>
        <NSlider v-model:value="panelState.panelConfig.backgroundBlur" class="max-w-[200px]" :step="2" :max="20" />
      </div>

      <div class="flex items-center mt-[10px]">
        <span class="mr-[10px]">{{ $t('apps.baseSettings.mask') }}</span>
        <NSlider v-model:value="panelState.panelConfig.backgroundMaskNumber" class="max-w-[200px]" :step="0.1" :max="1" />
      </div>
    </NCard>

    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <div class="text-slate-500 mb-[5px] font-bold">
        {{ $t('apps.baseSettings.contentArea') }}
      </div>

      <NGrid cols="2">
        <NGridItem span="12 400:12">
          <div class="flex items-center mt-[5px]">
            <span class="mr-[10px]">{{ $t('apps.baseSettings.netModeChangeButtonShow') }}</span>
            <NSwitch v-model:value="panelState.panelConfig.netModeChangeButtonShow" />
          </div>
        </NGridItem>

        <NGridItem span="12 400:12">
          <div class="flex items-center mt-[10px]">
            <span class="mr-[10px]">悬浮按钮位置</span>
            <NSelect
              v-model:value="panelState.panelConfig.floatingToolsPosition"
              :options="floatingToolsPositionOptions"
              :style="{ width: '140px' }"
              size="small"
            />
          </div>
        </NGridItem>

        <NGridItem span="12 400:12">
          <div class="flex items-center mt-[10px]">
            <span class="mr-[10px]">左侧分组目录</span>
            <NSwitch v-model:value="panelState.panelConfig.leftCatalogShow" />
          </div>
        </NGridItem>

        <NGridItem span="12 400:12">
          <div class="flex items-center mt-[10px]">
            <span class="mr-[10px]">右侧滚动进度条</span>
            <NSwitch v-model:value="panelState.panelConfig.rightScrollBarShow" />
          </div>
        </NGridItem>

        <NGridItem span="12 400:12">
          <div class="flex items-center mt-[10px]">
            <span class="mr-[10px]">{{ $t('apps.baseSettings.maxWidth') }}</span>
            <div class="flex">
              <NInputGroup>
                <NInput v-model:value="panelState.panelConfig.maxWidth" size="small" type="number" :maxlength="10" :style="{ width: '100px' }" placeholder="1200" />
                <NSelect v-model:value="panelState.panelConfig.maxWidthUnit" :style="{ width: '80px' }" :options="maxWidthUnitOption" size="small" />
              </NInputGroup>
            </div>
          </div>
        </NGridItem>
        <NGridItem span="12 400:12">
          <div class="flex items-center mt-[10px]">
            <span class="mr-[10px]">{{ $t('apps.baseSettings.leftRightMargin') }}</span>
            <NSlider v-model:value="panelState.panelConfig.marginX" class="max-w-[200px]" :step="1" :max="100" />
          </div>
        </NGridItem>
        <NGridItem span="12 400:12">
          <div class="flex items-center mt-[10px]">
            <span class="mr-[10px]">{{ $t('apps.baseSettings.topMargin') }} (%)</span>
            <NSlider v-model:value="panelState.panelConfig.marginTop" class="max-w-[200px]" :step="1" :max="50" />
          </div>
        </NGridItem>
        <NGridItem span="12 400:6">
          <div class="flex items-center mt-[10px]">
            <span class="mr-[10px]">{{ $t('apps.baseSettings.bottomMargin') }} (%)</span>
            <NSlider v-model:value="panelState.panelConfig.marginBottom" class="max-w-[200px]" :step="1" :max="50" />
          </div>
        </NGridItem>

        <NGridItem span="12 400:12">
          <div class="flex items-center mt-[10px]">
            <span class="mr-[10px]">上内边距 (px)</span>
            <NSlider v-model:value="panelState.panelConfig.paddingTop" class="max-w-[200px]" :step="2" :max="200" />
          </div>
        </NGridItem>

        <NGridItem span="12 400:12">
          <div class="flex items-center mt-[10px]">
            <span class="mr-[10px]">下内边距 (px)</span>
            <NSlider v-model:value="panelState.panelConfig.paddingBottom" class="max-w-[200px]" :step="2" :max="200" />
          </div>
        </NGridItem>

        <!-- Spacing -->
        <NGridItem span="12 400:12">
          <div class="flex items-center mt-[10px]">
            <span class="mr-[10px]">搜索框到下面间距 (px)</span>
            <NSlider v-model:value="panelState.panelConfig.searchToBelowGapPx" class="max-w-[200px]" :step="1" :max="100" />
          </div>
        </NGridItem>

        <NGridItem span="12 400:12">
          <div class="flex items-center mt-[10px]">
            <span class="mr-[10px]">网站/网页到第1分组间距 (px)</span>
            <NSlider v-model:value="panelState.panelConfig.tabsToFirstGroupGapPx" class="max-w-[200px]" :step="1" :max="120" />
          </div>
        </NGridItem>

        <NGridItem span="12 400:12">
          <div class="flex items-center mt-[10px]">
            <span class="mr-[10px]">分组与分组间距 (px)</span>
            <NSlider v-model:value="panelState.panelConfig.groupToGroupGapPx" class="max-w-[200px]" :step="1" :max="120" />
          </div>
        </NGridItem>

        <NGridItem span="12 400:12">
          <div class="flex items-center mt-[10px]">
            <span class="mr-[10px]">网页收藏列表间距 (px)</span>
            <NSlider v-model:value="panelState.panelConfig.webpageListGapPx" class="max-w-[200px]" :step="1" :max="30" />
          </div>
        </NGridItem>
      </NGrid>
    </NCard>

    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <div class="text-slate-500 mb-[5px] font-bold">
        {{ $t('apps.baseSettings.customFooter') }}
      </div>

      <NInput
        v-model:value="panelState.panelConfig.footerHtml"
        type="textarea"
        clearable
      />
    </NCard>

    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <NPopconfirm
        @positive-click="resetPanelConfig"
      >
        <template #trigger>
          <NButton size="small" quaternary type="error">
            {{ $t('common.reset') }}
          </NButton>
        </template>
        {{ $t('apps.baseSettings.resetWarnText') }}
      </NPopconfirm>

      <NButton size="small" quaternary type="success" class="ml-[10px]" @click="uploadCloud">
        {{ $t('common.save') }}
      </NButton>
    </NCard>
  </div>
</template>

<style scoped>
.text-shadow{
  text-shadow: 0px 0px 5px gray;
}

.wall-pick-card {
  border: 1px solid rgba(148, 163, 184, 0.35);
  border-radius: 10px;
  overflow: hidden;
  cursor: pointer;
  background: rgba(255, 255, 255, 0.7);
}

.wall-pick-thumb {
  height: 120px;
  width: 100%;
}

.wall-pick-name {
  padding: 6px 8px;
  font-size: 12px;
  color: #475569;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.wall-pick-card {
  border: 1px solid rgba(148, 163, 184, 0.35);
  border-radius: 10px;
  overflow: hidden;
  cursor: pointer;
  background: rgba(255, 255, 255, 0.7);
}

.wall-pick-thumb {
  height: 140px;
  width: 100%;
}

.wall-pick-name {
  padding: 6px 8px;
  font-size: 12px;
  color: #475569;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
