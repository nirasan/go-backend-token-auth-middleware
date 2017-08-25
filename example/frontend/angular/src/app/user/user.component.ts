import {Component, OnInit} from '@angular/core';
import {HttpClient, HttpHeaders, HttpParams} from '@angular/common/http';
import {AuthService} from '../auth.service';

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
    this.http.get('http://localhost:8080/userinfo').subscribe(
      data => {
        this.name = data['name'];
        console.log(data);
      },
      error => {
        console.log(error);
      }
    )
    ;
  }
}
