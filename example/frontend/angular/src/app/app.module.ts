import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {HttpClientModule, HTTP_INTERCEPTORS} from '@angular/common/http';

import {AuthService} from './auth.service';

import {AppComponent} from './app.component';
import {TopComponent} from './top/top.component';
import {UserComponent} from './user/user.component';
import {CallbackComponent} from './callback/callback.component';
import {AuthGuard} from './auth.guard';
import {AddHeaderInterceptor} from './add-header.interceptor';

const routes: Routes = [
  {path: '', redirectTo: 'top', pathMatch: 'full'},
  {path: 'top', component: TopComponent},
  {path: 'callback', component: CallbackComponent},
  {path: 'user', component: UserComponent, canActivate: [AuthGuard]}
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
    HttpClientModule
  ],
  providers: [
    AuthService,
    AuthGuard,
    {
      provide: HTTP_INTERCEPTORS,
      useClass: AddHeaderInterceptor,
      multi: true
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
