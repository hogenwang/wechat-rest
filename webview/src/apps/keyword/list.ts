import { Component } from '@angular/core';

import { KeywordTargets, KeywordLevels } from '../../openapi/const';
import { RobotApi, TablesKeyword } from '../../openapi/wrobot';
import { WrestApi, WcfrestContactPayload } from '../../openapi/wcfrest';


@Component({
    selector: 'page-keyword-list',
    templateUrl: 'list.html'
})
export class KeywordListComponent {

    public keywordTargets = KeywordTargets;
    public keywordLevels = KeywordLevels;

    public wcfChatrooms: Record<string, WcfrestContactPayload> = {};

    public keywords: Array<TablesKeyword> = [];

    constructor() {
        this.getKeywords();
        this.getWcfChatrooms();
    }

    public getKeywords() {
        RobotApi.keywordList({}).then((data) => {
            this.keywords = data || [];
        });
    }

    public deleteKeyword(item: TablesKeyword) {
        RobotApi.keywordDelete({ rd: item.rd }).then(() => {
            this.getKeywords();
        });
    }

    public getWcfChatrooms() {
        WrestApi.chatrooms().then((data) => {
            data.forEach((item) => this.wcfChatrooms[item.wxid] = item);
        });
    }

}
