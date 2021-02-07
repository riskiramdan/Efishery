'use strict';

module.exports = (accessToken, path, { accessTokenManager }) => {
  const decoded = accessTokenManager.decode(accessToken);

  if (!decoded) {
    throw new Error('Invalid access token');
  }
  if (path === "/aggregate") {
    if (decoded.roleId !== 1) {
      throw new Error('Restricted Area, Admin Only !');
    }
  }

  return { uid: decoded.uid };
};
