'use strict';

module.exports = async (phone, password, { userRepository, accessTokenManager }) => {
  const user = await userRepository.getByPhone(phone);
  if (!user || user.password !== password) {
    throw new Error('Bad credentials');
  }
  var today = new Date();
  const token = accessTokenManager.generate({ 
    uid: user.id,
    exp : Math.floor(Date.now() / 1000) + (60 * 60),
    iat: Math.floor(Date.now() / 1000) - 30,
    name: user.name,
    phone: user.phone,
    roleId : user.roleId,
    timestamp: Date.now(),
  })
  user.token = token
  user.tokenExpiredAt = today.setHours(today.getHours() + 1)

  await userRepository.merge(user);

  const auth = {
    'sessionId' :token,
    'claims': {
      'exp':Math.floor(Date.now() / 1000) + (60 * 60),
      'iat':Math.floor(Date.now() / 1000) - 30,
      'name':user.name,
      'phone' :user.phone,
      'roleId' : user.roleId,
      timestamp :Date.now()
    }
  }

  return auth;
};
