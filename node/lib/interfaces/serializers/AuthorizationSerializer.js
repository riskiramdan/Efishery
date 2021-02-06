'use strict';

const _serializeAuth = (auth) => {
  return {
    'sessionId': auth.token,
    'claims' :{
      'exp' : auth.exp,
      'iat' : auth.iat,
      'name' : auth.name,
      'phone' : auth.phone,
      'roleId' : auth.roleId,
      'timestamp' :auth.timestamp,
    }
  };
};

module.exports = class {

  serialize(data) {
    if (!data) {
      throw new Error('Expect data to be not undefined nor null');
    }    
    if (Array.isArray(data)) {
      return data.map(_serializeAuth);
    }
    return _serializeAuth(data);
  }

};