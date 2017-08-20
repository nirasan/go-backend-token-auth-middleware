import {Component, OnInit} from '@angular/core';
import {environment} from '../../environments/environment';

const oauth2Endpoint = 'https://accounts.google.com/o/oauth2/v2/auth';

@Component({
  selector: 'app-top',
  templateUrl: './top.component.html',
  styleUrls: ['./top.component.css']
})
export class TopComponent implements OnInit {


  constructor() { }


  ngOnInit() {
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
