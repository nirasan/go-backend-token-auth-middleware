import { Injectable } from '@angular/core';

const CODE_CHALLENGE_KEY = 'code_challenge';
const ACCESS_TOKEN_KEY = 'access_token';

@Injectable()
export class AuthService {

  constructor() { }

  setCodeChallenge(codeChallenge: string) {
    localStorage.setItem(CODE_CHALLENGE_KEY, codeChallenge);
  }

  getCodeChallenge(): string {
    return localStorage.getItem(CODE_CHALLENGE_KEY);
  }

  setAccessToken(accessToken: string) {
    localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
  }

  getAccessToken(): string {
    return localStorage.getItem(ACCESS_TOKEN_KEY);
  }
}
