'use strict';

const jwt = require('jsonwebtoken');

const AccessTokenManager = require('../../application/security/AccessTokenManager');

const JWT_SECRET_KEY = Buffer.from("CJ59d3WFwke9r75q0bcYm6MDCwBxVqY3", 'utf8');

module.exports = class extends AccessTokenManager {

  generate(payload) {
    return jwt.sign(payload, JWT_SECRET_KEY, {expiresIn:"12h"});
  }

  decode(accessToken) {
    return jwt.verify(accessToken, JWT_SECRET_KEY);
  }

};