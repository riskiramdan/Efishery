'use strict';

const _serializeSingleUser = (user) => {
  return {
    'id': user.id,
    'roleId': user.roleId,
    'name': user.name,
    'phone': user.phone,
  };
};

module.exports = class {

  serialize(data) {
    if (!data) {
      throw new Error('Expect data to be not undefined nor null');
    }
    if (Array.isArray(data)) {
      return data.map(_serializeSingleUser);
    }
    return _serializeSingleUser(data);
  }

};