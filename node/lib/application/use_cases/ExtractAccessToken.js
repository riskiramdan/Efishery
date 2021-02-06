'use strict';

module.exports = async (token, { userRepository, accessTokenManager  }) => {
  const user = await userRepository.getByToken(token);
  if (!user) {
    throw new Error('Invalid access token');
  }

  const decoded = accessTokenManager.decode(token);
  if (!decoded) {
    throw new Error('Invalid access token');
  }
  
  return user;
};
