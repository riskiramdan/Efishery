'use strict';

const constants = require('./constants');
const environment = require('./environment');
const JwtAccessTokenManager = require('../security/JwtAccessTokenManager');
const UserSerializer = require('../../interfaces/serializers/UserSerializer');
const AuthSerializer = require('../../interfaces/serializers/AuthorizationSerializer');

function buildBeans() {

  const beans = {
    accessTokenManager: new JwtAccessTokenManager(),
    userSerializer: new UserSerializer(),
    authSerializer: new AuthSerializer(),
  };

  if (environment.database.dialect === constants.SUPPORTED_DATABASE.IN_MEMORY) {
    const UserRepositoryInMemory = require('../repositories/UserRepositoryInMemory');
    beans.userRepository = new UserRepositoryInMemory();
  } else if (environment.database.dialect === constants.SUPPORTED_DATABASE.MONGO) {
    const UserRepositoryMongo = require('../repositories/UserRepositoryMongo');
    beans.userRepository = new UserRepositoryMongo();
  } else if (environment.database.dialect === constants.SUPPORTED_DATABASE.POSTGRES) {
    const UserRepositoryPostgreSQL= require('../repositories/UserRepositoryPostgreSQL');
    beans.userRepository = new UserRepositoryPostgreSQL();
  } else {
    const UserRepositorySQLite= require('../repositories/UserRepositorySQLite');
    beans.userRepository = new UserRepositorySQLite();
  }

  return beans;
}

module.exports = buildBeans();
