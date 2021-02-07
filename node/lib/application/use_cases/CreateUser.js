'use strict';

const User = require('../../domain/User');

module.exports = async (roleId, name, phone, { userRepository }) => {
  const user = new User(null, roleId, name, phone, makepwd(4));
  return userRepository.persist(user);
};

function makepwd(length) {
  var result           = '';
  var characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  var charactersLength = characters.length;
  for ( var i = 0; i < length; i++ ) {
     result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
}