import { Component } from '@angular/core';
import { Router } from '@angular/router';

import { UserLevels } from '../../openapi/const';
import { RobotApi, LlmodelCreateParam } from '../../openapi/wrobot';


@Component({
    selector: 'page-llmodel-create',
    templateUrl: 'create.html'
})
export class LLModelCreateComponent {

    public userLevels = UserLevels;

    public formdata: LlmodelCreateParam = {
        mid: '',
        family: 'Gemini',
        provider: 'google',
        model: 'gemini-pro',
        secret: '',
    };

    constructor(private router: Router) { }

    public createLLModel() {
        if (this.formdata.level) {
            this.formdata.level = +this.formdata.level;
        }
        RobotApi.llmodelCreate(this.formdata).then(() => {
            this.router.navigate(['llmodel/list']);
        });
    }

}
