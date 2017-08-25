import {
  HttpEvent,
  HttpInterceptor,
  HttpHandler,
  HttpRequest
} from '@angular/common/http';
import {Observable} from 'rxjs/Observable';
import {AuthService} from './auth.service';
import {Injectable} from '@angular/core';

@Injectable()
export class AddHeaderInterceptor implements HttpInterceptor {

  constructor(private auth: AuthService) {
  }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {

    // if (!req.url.includes('/userinfo')) {
    //   return next.handle(req);
    // }

    // Clone the request to add the new header
    const clonedRequest = req.clone({headers: req.headers.set('Authorization', 'Bearer ' + this.auth.getAccessToken())});

    console.log(clonedRequest);

    // Pass the cloned request instead of the original request to the next handle
    return next.handle(clonedRequest);
  }
}
