import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';

import {AppComponent} from './app.component';
import {TopComponent} from './top/top.component';
import {UserComponent} from './user/user.component';
import { CallbackComponent } from './callback/callback.component';

const routes: Routes = [
  {path: '', redirectTo: 'top', pathMatch: 'full'},
  {path: 'top', component: TopComponent},
  {path: 'callback', component: CallbackComponent},
  {path: 'user', component: UserComponent},
];

@NgModule({
  declarations: [
    AppComponent,
    TopComponent,
    UserComponent,
    CallbackComponent
  ],
  imports: [
    BrowserModule,
    RouterModule.forRoot(routes),
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {
}
