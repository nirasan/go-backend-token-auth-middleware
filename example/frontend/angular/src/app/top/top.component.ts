import {Component, OnInit} from '@angular/core';
import {environment} from '../../environments/environment';
import {HttpClient} from '@angular/common/http';
import {AuthService} from '../auth.service';

const backendAuthStartUrl = 'http://localhost:8080/auth/start';
const oauth2Endpoint = 'https://accounts.google.com/o/oauth2/v2/auth';

interface AuthStartResponse {
  code_challenge: string
}

@Component({
  selector: 'app-top',
  templateUrl: './top.component.html',
  styleUrls: ['./top.component.css']
})
export class TopComponent implements OnInit {

  constructor(private http: HttpClient, private auth: AuthService) {
  }

  ngOnInit() {
    this.http.get<AuthStartResponse>(backendAuthStartUrl).subscribe(
      (data) => {
        this.auth.setCodeChallenge(data.code_challenge);
      },
      (error) => {
        console.log(error);
      }
    );
  }

  login() {
    // Create <form> element to submit parameters to OAuth 2.0 endpoint.
    var form = document.createElement('form');
    form.setAttribute('method', 'GET'); // Send as a GET request.
    form.setAttribute('action', oauth2Endpoint);

    // Parameters to pass to OAuth 2.0 endpoint.
    var params = {
      'client_id': environment['CLIENT_ID'],
      'redirect_uri': 'http://localhost:4200/callback',
      'response_type': 'code',
      'scope': 'profile openid',
      'code_challenge': this.auth.getCodeChallenge(),
      'code_challenge_method': 'S256',
    };

    // Add form parameters as hidden input values.
    for (var p in params) {
      var input = document.createElement('input');
      input.setAttribute('type', 'hidden');
      input.setAttribute('name', p);
      input.setAttribute('value', params[p]);
      form.appendChild(input);
    }

    // Add form to page and submit it to open the OAuth 2.0 endpoint.
    document.body.appendChild(form);
    form.submit();
  }

}
