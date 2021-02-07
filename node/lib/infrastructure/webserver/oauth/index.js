'use strict';

const AuthorizationController = require('../../../interfaces/controllers/AuthorizationController');
const ConcurrenctyController = require('../../../interfaces/controllers/ConcurrenctyController');

module.exports = {
  name: 'oauth',
  version: '1.0.0',
  register: (server) => {

    server.auth.scheme('oauth', require('./scheme'));

    server.auth.strategy('oauth-jwt', 'oauth');

    server.route([{
      method: 'POST',
      path: '/oauth/token',
      handler: AuthorizationController.getAccessToken,
      options: {
        description: 'Return an OAuth 2 access token',
        tags: ['api'],
      },
    },
    {
      method:'POST',
      path:'/oauth/verify',
      handler: AuthorizationController.extractAcccessToken,
      options: {
        tags : ['api']
      },
    },
    {
      method:'GET',
      path:'/prices',
      handler: ConcurrenctyController.getListDataPrice,
      config: {
        auth:'oauth-jwt',
        tags : ['api']
      },
    },
    {
      method:'GET',
      path:'/conversion',
      handler: ConcurrenctyController.getConversionPrice,
      config: {
        auth:'oauth-jwt',
        tags : ['api']
      },
    }
  ]);
  }
};
