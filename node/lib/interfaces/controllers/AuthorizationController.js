'use strict';

const Boom = require('@hapi/boom');
const ExtractAccessToken = require('../../application/use_cases/ExtractAccessToken');
const GetAccessToken = require('../../application/use_cases/GetAccessToken');
const VerifyAccessToken = require('../../application/use_cases/VerifyAccessToken');

module.exports = {

  async getAccessToken(request) {

    // Context
    const serviceLocator = request.server.app.serviceLocator;
    // Input
    const phone = request.payload['phone'];
    const password = request.payload['password'];

    // Treatment
    try {
      const auth = await GetAccessToken(phone, password, serviceLocator);
      // Output
      return auth;
    } catch (err) {
      return Boom.unauthorized('Bad credentials');
    }
  },

  async extractAcccessToken(request) {
    // Context
    const serviceLocator = request.server.app.serviceLocator;
    // Input
    let token = request.payload['token']
    // Treatment
    try {
      const auth = await ExtractAccessToken(token, serviceLocator);
      // Output
      return auth;
    } catch (err) {
      return Boom.unauthorized(err);
    }
  },

  verifyAccessToken(request, h) {

    // Context
    const serviceLocator = request.server.app.serviceLocator;

    // Input
    const authorizationHeader = request.headers.authorization;
    if (!authorizationHeader || !authorizationHeader.startsWith('Bearer ')) {
      throw Boom.badRequest('Missing or wrong Authorization request header', 'oauth');
    }
    const accessToken = authorizationHeader.replace(/Bearer/gi, '').replace(/ /g, '');
    const path = h.request.route.path

    // Treatment
    try {
      const { uid } = VerifyAccessToken(accessToken,path, serviceLocator);

      // Output
      return h.authenticated({
        credentials: { uid },
        artifacts: { accessToken: accessToken }
      });
    } catch (err) {
      return Boom.unauthorized(err);
    }
  },

};