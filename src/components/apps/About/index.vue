<script setup lang="ts">
import { NDivider, NGradientText, NTag } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { get } from '@/api/system/about'
import srcSvglogo from '@/assets/logo.svg'
import srcGithub from '@/assets/about_image/github.png'
import srcTelegram from '@/assets/about_image/telegram.svg'
import srcMarket from '@/assets/about_image/market.svg'

interface Version {
  versionName: string
  versionCode: number
}

const versionName = ref('')
const frontVersion = import.meta.env.VITE_APP_VERSION || 'unknown'

onMounted(() => {
  get<Version>().then((res) => {
    if (res.code === 0)
      versionName.value = res.data.versionName
  })
})
</script>

<template>
  <div class="pt-5">
    <div class="flex flex-col items-center justify-center">
      <img :src="srcSvglogo" width="100" height="100" alt="">
      <div class="text-3xl font-semibold">
        AnGe-Panel
      </div>
      <div class="text-xl">
        <NGradientText type="info">
          <a href="https://github.com/liandu2024/AnGe-Panel/releases" class="font-semibold" :title="$t('apps.about.viewUpdateLog')" target="_blank">Ver: {{ versionName }}</a>
        </NGradientText>
      </div>
      <div class="mt-2">
        <a href="https://github.com/liandu2024/AnGe-Panel/releases" target="_blank" class="link">{{ $t('apps.about.checkUpdate') }}</a>
      </div>
    </div>

    <NDivider style="margin:10px 0">
      •
    </NDivider>
    <div class="flex flex-col items-center justify-center text-base">
      <div class="flex mt-[10px] flex-wrap justify-center">
        <div class="flex items-center mx-[10px]">
          <img class="w-[20px] h-[20px] mr-[5px]" :src="srcGithub" alt="">
          <a href="https://github.com/liandu2024/AnGe-Panel" target="_blank" class="link">Github</a>
        </div>
        <div class="flex items-center mx-[10px]">
          <img class="w-[20px] h-[20px] mr-[5px]" :src="srcTelegram" alt="">
          <a href="https://t.me/angeworld2024" target="_blank" class="link">安格视界TG群</a>
        </div>
        <div class="flex items-center mx-[10px]">
          <img class="w-[20px] h-[20px] mr-[5px]" :src="srcMarket" alt="">
          <a href="https://blog.angeworld.cc/market" target="_blank" class="link">安格超市</a>
        </div>
      </div>

      <div class="mt-5">
        <NTag :bordered="false" size="small">
          {{ $t("apps.about.frontVersionText") }}: FV-{{ frontVersion }}
        </NTag>
      </div>

      <div class="mt-4 text-sm text-gray-500">
        <a href="https://github.com/hslr-s/sun-panel" target="_blank" class="link">本项目基于 Sun-Panel 开源版 v1.3.0 二次开发</a>
      </div>
    </div>
  </div>
</template>

<style>
.link{
    color:rgb(0, 89, 255)
}
</style>
