'use strict';

const Boom = require('@hapi/boom');
const ExtractAccessToken = require('../../application/use_cases/ExtractAccessToken');
const ListPrices = require('../../application/use_cases/ListPrices');
const VerifyAccessToken = require('../../application/use_cases/VerifyAccessToken');

module.exports = {

  async getListDataPrice(request) {

    // Context
    const serviceLocator = request.server.app.serviceLocator;

    // Treatment
    try {
      const prices = await ListPrices(serviceLocator);
      // Output
      return {"prices" : prices, "total" : prices.length};
    } catch (err) {
      console.log(err)
      return Boom.unauthorized('Bad credentials');
    }
  },

  async extractAcccessToken(request) {
    // Context
    const serviceLocator = request.server.app.serviceLocator;
    // Input
    const token = request.payload['token'];
    // Treatment
    try {
      const auth = await ExtractAccessToken(token, serviceLocator);
      // Output
      return auth;
    } catch (err) {
      console.log(err)
      return Boom.unauthorized('Bad credentials');
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

    // Treatment
    try {
      const { uid } = VerifyAccessToken(accessToken, serviceLocator);

      // Output
      return h.authenticated({
        credentials: { uid },
        artifacts: { accessToken: accessToken }
      });
    } catch (err) {
      return Boom.unauthorized('Bad credentials');
    }
  },

};