import {Component, OnInit} from '@angular/core';
import {AuthService} from '../auth.service';
import {ActivatedRoute} from '@angular/router';
import {HttpClient, HttpHeaders} from '@angular/common/http';

@Component({
  selector: 'app-callback',
  templateUrl: './callback.component.html',
  styleUrls: ['./callback.component.css']
})
export class CallbackComponent implements OnInit {

  constructor(private route: ActivatedRoute, private auth: AuthService, private http: HttpClient) {
  }

  ngOnInit() {
    let code = this.route.snapshot.queryParams['code']

    let body = JSON.stringify({
      code: code,
      code_challenge: this.auth.getCodeChallenge(),
    });

    this.http.post('http://localhost:8080/auth', body).subscribe(
      (data) => {
        console.log(data);
      },
      (error) => {
        console.log(error);
      }
    );
  }
}
