import {Component, OnInit} from '@angular/core';
import {AuthService} from '../auth.service';
import {ActivatedRoute, Router} from '@angular/router';
import {HttpClient} from '@angular/common/http';

@Component({
  selector: 'app-callback',
  templateUrl: './callback.component.html',
  styleUrls: ['./callback.component.css']
})
export class CallbackComponent implements OnInit {

  constructor(private route: ActivatedRoute,
              private auth: AuthService,
              private http: HttpClient,
              private router: Router) {
  }

  ngOnInit() {
    const code = this.route.snapshot.queryParams['code'];

    const body = JSON.stringify({
      code: code,
      code_challenge: this.auth.getCodeChallenge()
    });

    this.http.post('http://localhost:8080/auth', body).subscribe(
      (data) => {
        console.log(data);
        this.auth.setAccessToken(data['access_token']);
        this.router.navigate(['user']);
      },
      (error) => {
        console.log(error);
      }
    );
  }
}
