import {Component, OnInit} from '@angular/core';
import {AuthService} from '../auth.service';
import {HttpClient, HttpHeaders} from '@angular/common/http';

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.css']
})
export class UserComponent implements OnInit {

  name = '';

  constructor(private http: HttpClient,
              private auth: AuthService) {
  }

  ngOnInit() {
    this.http.get('http://localhost:8080/userinfo', {
      headers: new HttpHeaders().set('Authorization', 'Bearer ' + this.auth.getAccessToken())
    }).subscribe(
      data => {
        this.name = data['name'];
        console.log(data);
      },
      error => {
        console.log(error);
      }
    );
  }
}
