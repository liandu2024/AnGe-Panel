declare namespace Panel {

    interface Info extends ItemInfo {

    }

    interface ItemInfo extends Common.InfoBase {
        icon: ItemIcon |null
        title: string
        url: string
        sort?: number
        pinned?: boolean
        lanUrl?: string
        description?: string
        openMethod: number
        itemIconGroupId ?:number
    }

    interface ItemIconGroup extends Common.InfoBase {
        icon?: string
        title?: string
        sort?:number
        groupType?: 'website' | 'webpage'
    }

    interface ItemIcon {
        itemType: number
        src ?: string
        text ?: string
        // bgColor ?: string
        backgroundColor ?: string
    }

    interface State {
        rightSiderCollapsed: boolean
        leftSiderCollapsed: boolean
        networkMode:PanelStateNetworkModeEnum | null
        panelConfig:panelConfig
    }

    interface panelConfig{
        backgroundImageSrc?:string
        backgroundBlur?:number
        backgroundMaskNumber?:number
        iconStyle?:PanelPanelConfigStyleEnum
        iconTextColor?:string
        iconTextInfoHideDescription?:boolean
        iconTextIconHideTitle?:boolean
        logoText?:string
        logoImageSrc?:string
        topHeaderShow?:boolean
        clockShowSecond?:boolean
        clockColor?:string
        searchBoxShow?:boolean
        searchBoxSearchIcon?:boolean
        marginTop?:number
        marginBottom?:number
        maxWidth?:number
        maxWidthUnit:string
        marginX?:number
        paddingTop?:number
        paddingBottom?:number
        footerHtml?:string
        systemMonitorShow?:boolean
        systemMonitorShowTitle?:boolean
        systemMonitorPublicVisitModeShow?:boolean
        netModeChangeButtonShow?:boolean
        floatingToolsPosition?: 'right-top' | 'right-middle' | 'right-bottom'
        leftCatalogShow?: boolean
        rightScrollBarShow?: boolean

        // Spacing (px)
        searchToBelowGapPx?: number
        tabsToFirstGroupGapPx?: number
        groupToGroupGapPx?: number
        webpageListGapPx?: number
    }

    interface userConfig{
        panel:panelConfig
        searchEngine?:any
    }

    interface ItemIconSortRequest{
        sortItems:Common.SortItemRequest[]
        itemIconGroupId:number
    }
}

